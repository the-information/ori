package test

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"golang.org/x/oauth2/jws"
	"math/rand"
)

// JWT generates a JWT from claimSet and secret. It panics if it encounters an error.
func JWT(claimSet *jws.ClaimSet, secret string) string {

	bytes, err := jws.EncodeWithSigner(&jws.Header{Algorithm: "HS256", Typ: "JWT"}, claimSet, func(data []byte) ([]byte, error) {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(data)
		return mac.Sum(nil), nil
	})
	if err != nil {
		panic(err)
	}

	return string(bytes)

}

// ConsumableJWT generates a consumable JWT from claimset and secret, with u uses available.
// It panics if it encounters an error.
func ConsumableJWT(claimSet *jws.ClaimSet, secret string, uses int) string {

	var newClaimSet jws.ClaimSet
	newClaimSet = *claimSet

	if newClaimSet.PrivateClaims == nil {
		newClaimSet.PrivateClaims = map[string]interface{}{}
	}

	newClaimSet.PrivateClaims["u"] = uses
	newClaimSet.PrivateClaims["jti"] = fmt.Sprintf("j%x", rand.Int63()%1000000000)

	return JWT(&newClaimSet, secret)

}
