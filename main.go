package main

// Import our dependencies. We'll use the standard HTTP library as well as the gorilla router for this app
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
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

func main() {

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

	r := mux.NewRouter()

	r.Handle("/data", jwtMiddleware.Handler(DataHandler)).Methods("GET")
	r.Handle("/data", jwtMiddleware.Handler(DuplicateDataHandler)).Methods("POST")
	r.Handle("/data/{id}", jwtMiddleware.Handler(PutDataHandler)).Methods("PUT")
	r.Handle("/data/{id}", jwtMiddleware.Handler(DeleteDataHandler)).Methods("DELETE")

	// For dev only - Set up CORS so React client can consume our API
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	fmt.Printf("server started\n")

	http.ListenAndServe(":8080", corsWrapper.Handler(r))
}

var DataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	userId := ""
	filteredData := readData(userId)

	payload, _ := json.Marshal(filteredData)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

var PutDataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	token := authHeaderParts[1]

	userId, _ := getUserEmail(token)

	log.Printf("put user=%s id=%s", userId, id)

	filteredData := readData(userId)

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

	_ = deleteData(userId, i)

	filteredData := readData(userId)

	payload, _ := json.Marshal(filteredData)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

var DuplicateDataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	userId := ""
	resourceId := "default"

	s, err := duplicateData(userId, resourceId)

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

	domain := os.Getenv("DATA_DOMAIN")

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://"+domain+"/userinfo", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return "", err
	}

	return userInfo.Email, nil
}
