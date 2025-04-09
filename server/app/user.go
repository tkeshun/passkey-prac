package main

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	ID                  string
	Name                string
	Credentials         []webauthn.Credential
	RegistrationSession *webauthn.SessionData
}

func (u *User) WebAuthnID() []byte {
	return []byte(u.ID)
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	return u.Name
}

func (u *User) WebAuthnIcon() string {
	return ""
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
