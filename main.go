package main

// Import our dependencies. We'll use the standard HTTP library as well as the gorilla router for this app
import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jimareed/slideshow-data/data"
	"github.com/rs/cors"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

/* user info */
type UserInfo struct {
	Sub   string `json:"sub"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserCache struct {
	Email string
	Token string
}

type UpdateData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ResourceId  string `json:"resourceId"`
}

var d = data.Data{}
var usersCache = []UserCache{}

func main() {

	var filePath = "."

	input := flag.String("input", "", "path to source")
	help := flag.Bool("help", false, "help")

	flag.Parse()

	if *help {
		log.Fatal("usage: slideshow-data [-input <path>][-help]")
	}

	if *input != "" {
		filePath = *input
	}

	log.Print("reading model and policy from ", filePath)

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := os.Getenv("DATA_API_ID")
			domain := os.Getenv("DATA_DOMAIN")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := "https://" + domain + "/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	d = data.Init(filePath+"/model.conf", filePath+"/policy.csv")

	r := mux.NewRouter()

	r.Handle("/data", jwtMiddleware.Handler(GetDataHandler)).Methods("GET")
	r.Handle("/data", jwtMiddleware.Handler(NewDataHandler)).Methods("POST")
	r.Handle("/data/{id}", jwtMiddleware.Handler(UpdateDataHandler)).Methods("PUT")
	r.Handle("/data/{id}", jwtMiddleware.Handler(DeleteDataHandler)).Methods("DELETE")

	// For dev only - Set up CORS so React client can consume our API
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	log.Printf("server started\n")

	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(r)))
}

var GetDataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	token := authHeaderParts[1]

	userId, err := getUserEmail(token)

	if err != nil {
		log.Printf("error: %v\n", err)
	}

	filteredData := d.ReadData(userId)

	payload, _ := json.Marshal(filteredData)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

var UpdateDataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	token := authHeaderParts[1]

	userId, _ := getUserEmail(token)

	i, _ := strconv.Atoi(id)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var updateData UpdateData

	json.Unmarshal(reqBody, &updateData)

	d.UpdateData(userId, i, updateData.Name, updateData.Description)

	filteredData := d.ReadData(userId)

	payload, _ := json.Marshal(filteredData)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

var DeleteDataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	token := authHeaderParts[1]

	userId, _ := getUserEmail(token)

	i, _ := strconv.Atoi(id)

	_ = d.DeleteData(userId, i)

	filteredData := d.ReadData(userId)

	payload, _ := json.Marshal(filteredData)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

var NewDataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	token := authHeaderParts[1]

	userId, _ := getUserEmail(token)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var updateData UpdateData

	json.Unmarshal(reqBody, &updateData)

	s, err := d.NewData(userId, updateData.ResourceId)

	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		payload, _ := json.Marshal(s)
		w.Write([]byte(payload))
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
})

func getPemCert(token *jwt.Token) (string, error) {

	cert := ""
	domain := os.Getenv("DATA_DOMAIN")
	resp, err := http.Get("https://" + domain + "/.well-known/jwks.json")

	if err != nil {
		log.Printf("error: %v\n", err)
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func getUserEmail(token string) (string, error) {

	var userInfo UserInfo

	userEmail := lookupUser(token)

	if userEmail != "" {
		return userEmail, nil
	}

	domain := os.Getenv("DATA_DOMAIN")

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://"+domain+"/userinfo", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error: %v\n", err)
		return "", err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		log.Printf("error: %v\n", err)
		return "", err
	}

	newUser := UserCache{}
	newUser.Email = userInfo.Email
	newUser.Token = token

	usersCache = append(usersCache, newUser)

	return userInfo.Email, nil
}

func lookupUser(token string) string {
	for _, u := range usersCache {
		if u.Token == token {
			return u.Email
		}
	}

	return ""
}
