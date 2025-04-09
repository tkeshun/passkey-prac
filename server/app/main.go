package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
)

var webAuthn *webauthn.WebAuthn
var userStore = map[string]*User{}

func main() {
	var err error
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Go WebAuthn Example",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
	})
	if err != nil {
		log.Fatal(err)
	}
	userStore["tShun"] = &User{
		ID:   "tShun", // 本番では UUID 推奨
		Name: "tShun",
	}

	http.Handle("/", http.FileServer(http.Dir("./pages")))

	http.HandleFunc("/webauthn/register/begin", beginRegistration)
	http.HandleFunc("/webauthn/register/finish", finishRegistration)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func beginRegistration(w http.ResponseWriter, r *http.Request) {
	log.Println("request begin register")
	var req struct {
		Username string `json:"username"`
		webauthn.Credential
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("invalid request:" + err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user := getOrCreateUser(req.Username)

	options, sessionData, err := webAuthn.BeginRegistration(
		user,
	)
	if err != nil {
		log.Panicln("error starting registration  " + err.Error())
		http.Error(w, "error starting registration", http.StatusInternalServerError)
		return
	}
	log.Println(options.Response.Challenge)

	user.RegistrationSession = sessionData
	log.Println(sessionData)

	writeJSON(w, options)
}

func finishRegistration(w http.ResponseWriter, r *http.Request) {
	user, ok := userStore["tShun"]
	if !ok {
		log.Println("user not found")
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	fmt.Println(user)

	log.Println(*user.RegistrationSession)
	cred, err := webAuthn.FinishRegistration(user, *user.RegistrationSession, r)
	if err != nil {
		log.Println("registration failed: " + err.Error())
		http.Error(w, "registration failed", http.StatusInternalServerError)
		return
	}

	user.Credentials = append(user.Credentials, *cred)
	w.WriteHeader(http.StatusOK)
}

func getOrCreateUser(username string) *User {
	if user, exists := userStore[username]; exists {
		return user
	}

	user := &User{
		ID:   username, // 本番では UUID 推奨
		Name: username,
	}
	userStore[username] = user
	return user
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
