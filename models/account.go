package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	flag           int
	CreatedAt      time.Time       `json:"createdAt"`
	LastUpdatedAt  time.Time       `json:"lastUpdatedAt"`
	Email          string          `json:"email"`
	Roles          []string        `json:"roles"`
	SecurePassword []byte          `json:"-" datastore:",noindex"`
	Type           string          `json:"type"`
	Settings       AccountSettings `json:"settings"`
}

type AccountSettings struct {
	Email struct {
		Articles    bool `json:"articles"`
		Digests     bool `json:"digests"`
		Contributor bool `json:"contributor"`
		Marketing   bool `json:"marketing"`
	} `json:"email"`
}

const (
	super = iota + 1
	nobody
)

var SuperAccount Account = Account{
	flag:  super,
	Email: "super@",
	Roles: []string{},
}

var NobodyAccount Account = Account{
	flag:  nobody,
	Email: "nobody@",
	Roles: []string{},
}

func (a *Account) HasRole(role string) bool {

	for _, existingRole := range a.Roles {
		if role == existingRole {
			return true
		}
	}
	return false

}

func (a *Account) Super() bool {
	return a.flag == super
}

func (a *Account) Nobody() bool {
	return a.flag == nobody
}

func (a *Account) CheckPassword(proposedPassword string) error {
	return bcrypt.CompareHashAndPassword(a.SecurePassword, []byte(proposedPassword))
}

func (a *Account) SetPassword(plaintextPassword string) (err error) {
	a.SecurePassword, err = bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	return err
}

type AccountCreationRequest struct {
	Email              string `json:"email" valid:"email"`
	Name               string `json:"name"`
	Password           string `json:"password" valid:"length(6|)"`
	Type               string `json:"type"`
	PaymentMethodNonce string `json:"paymentMethodNonce"`
	OfferCode          string `json:"offerCode"`
}
