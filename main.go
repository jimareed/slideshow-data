package main

// Import our dependencies. We'll use the standard HTTP library as well as the gorilla router for this app
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

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

/* Data type */
type Data struct {
	Id          int
	ResourceId  string
	Name        string
	Description string
}

var data = []Data{
	Data{Id: 1, ResourceId: "default", Name: "Slideshow", Description: "Overview"},
	Data{Id: 2, ResourceId: "instructions", Name: "Instructions", Description: "Steps to use"},
	Data{Id: 3, ResourceId: "emotional-intelligence", Name: "Emotional Intelligence", Description: "Sample slideshow"},
}

var NULL_DATA = Data{0, "", "", ""}

func Get(resourceId string) (Data, error) {

	for _, s := range data {
		if s.ResourceId == resourceId {
			return s, nil
		}
	}

	return NULL_DATA, fmt.Errorf("Get error: invalid id %s", resourceId)
}

func Duplicate(resourceId string) (Data, error) {

	for _, s := range data {
		if s.ResourceId == resourceId {
			s1 := Data{Id: len(data) + 1, ResourceId: "2345", Name: "copy of " + s.Name, Description: s.Description}
			data = append(data, s1)
			return s1, nil
		}
	}

	return NULL_DATA, fmt.Errorf("Duplicate error: invalid id %s", resourceId)
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

	// For dev only - Set up CORS so React client can consume our API
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	fmt.Printf("server started\n")

	http.ListenAndServe(":8080", corsWrapper.Handler(r))
}

var DataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

var DuplicateDataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	s, err := Duplicate("default")

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
