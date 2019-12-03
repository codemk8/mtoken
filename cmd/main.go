package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codemk8/mtoken/pkg/token"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2/jwt"
)

var ip = flag.String("addr", "127.0.0.1:8001", "Serving host and port")
var userAuthService = flag.String("user_service", "http://localhost:8000/v1/user/auth", "User microservice downstream endpoint")
var apiRoot = flag.String("api_root", "/v1", "api root path")
var appName = flag.String("app_name", "yourapp", "app name")
var privKeyFile = flag.String("priv_key", "", "private key file path")
var pubKeyFile = flag.String("pub_key", "", "public key file path")
var signer token.Signer
var verifier token.Verifier

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))
	return nil
}

func genToken(username string) (*string, error) {
	expiry := jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(10)))
	claims := jwt.Claims{Issuer: *appName,
		Subject: username,
		Expiry:  expiry}
	jwt, err := signer.Sign(&claims)
	if err != nil {
		glog.Errorf("Internal error %v", err)
	}
	return jwt, err
}

func issueHandler(w http.ResponseWriter, r *http.Request) {
	username, passwd, authOK := r.BasicAuth()
	if authOK == false {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				claims, err := verifier.Verify(&bearerToken[1])
				if err != nil {
					glog.Errorf("Invalid token %v", err)
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}
				if (int64)(*claims.Expiry) < time.Now().Unix() {
					glog.Errorf("Token expired")
					http.Error(w, "Token expired", http.StatusUnauthorized)
					return
				}
				jwt, err := genToken(claims.Subject)
				if err != nil {
					glog.Errorf("Internal error %v", err)
					http.Error(w, "Internal Error", http.StatusInternalServerError)
					return
				}
				glog.Info("generated jwt.")
				w.Write([]byte(*jwt))
				return
			}
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	client := &http.Client{CheckRedirect: redirectPolicyFunc}
	req, err := http.NewRequest("GET", *userAuthService, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(username, passwd))
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorf("Internal error %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}
	// Generate a token
	// expiry := jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)) // add a month
	jwt, err := genToken(username)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	//http.SetCookie(w, &http.Cookie{
	//	Name:    "token",
	//	Value:   tokenString,
	//	Expires: expirationTime,
	//})
	glog.Infof("generated jwt %s.", *jwt)
	w.Write([]byte(*jwt))
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("authorization")
	if authorizationHeader != "" {
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) == 2 {
			claims, err := verifier.Verify(&bearerToken[1])
			if err != nil {
				glog.Errorf("Invalid token %v", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			if (int64)(*claims.Expiry) < time.Now().Unix() {
				glog.Errorf("Token expired")
				http.Error(w, "Token expired", http.StatusUnauthorized)
			}
			claimBytes, err := json.Marshal(claims)
			if err != nil {
				glog.Warningf("Error marshal claims: %v", err)
				http.Error(w, "Internal error.", http.StatusInternalServerError)
			}
			w.Write(claimBytes)
			return
		}
	}
	http.Error(w, "Invalid input", http.StatusUnauthorized)
	return
}

func main() {
	flag.Parse()

	priv, pub := token.GenerateRsaKeyPairIfNotExist(*privKeyFile, *pubKeyFile, true)
	var err error
	signer, err = token.NewSigner(priv, pub)
	if err != nil {
		log.Fatalf("Error creating the token signer: %v", err)
	}
	verifier, err = token.NewVerifier(pub)
	if err != nil {
		log.Fatalf("Error creating the token verifier: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc(*apiRoot+"/token/issue", issueHandler).Methods("GET")
	r.HandleFunc(*apiRoot+"/token/auth", authHandler).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    *ip,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	glog.Errorf("Running server on %s%s/token/issue(auth).\n", *ip, *apiRoot)
	log.Fatal(srv.ListenAndServe())
}
