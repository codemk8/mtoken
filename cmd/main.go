package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/codemk8/mtoken/pkg/token"
	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2/jwt"
)

var ip = flag.String("addr", "127.0.0.1:8001", "Serving host and port")
var userAuthService = flag.String("user_service", "http://localhost:8000/v1/user/auth", "User microservice downstream endpoint")
var apiRoot = flag.String("api_root", "/v1", "api root path")
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

func issueHandler(w http.ResponseWriter, r *http.Request) {
	username, passwd, authOK := r.BasicAuth()
	if authOK == false {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	client := &http.Client{CheckRedirect: redirectPolicyFunc}
	req, err := http.NewRequest("GET", *userAuthService, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(username, passwd))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Internal error %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}
	// Generate a token
	expiry := jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)) // add a month
	claims := jwt.Claims{Issuer: username,
		Subject: "yourapp",
		Expiry:  expiry}
	jwt, err := signer.Sign(&claims)
	if err != nil {
		fmt.Printf("Internal error %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(*jwt))
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
	r.HandleFunc(*apiRoot+"/token/issue", issueHandler).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    *ip,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Running server on %s/token/login.\n", *ip)
	log.Fatal(srv.ListenAndServe())
}
