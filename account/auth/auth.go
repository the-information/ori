package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/oauth2/jws"
	"time"
)

// AllScope is the auth scope set when a token is valid for all
// of a user's roles.
var AllScope = "__ALL"

// All errors returned by package auth have type Error, so they
// can be differentiated in a type switch as follows:
//
//   err := doSomethingInvolvingJWTs()
//   switch t := err.(type) {
//   case auth.Error:
//     fmt.Println("Authentication error: ", err)
//   default:
//     fmt.Println("Generic error: ", err)
//   }
type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ExpiredJWTError       = Error("JWT has expired")
	InvalidJWTError       = Error("Not a valid JWT")
	BadSignatureError     = Error("Signatures don't match")
	InvalidHeaderError    = Error("Header isn't type JWT")
	InvalidAlgorithmError = Error("Algorithm isn't HS256")

	// SuperClaimSet is a special jws.ClaimSet returned when
	// the JWT supplied to a Decode call is actually just the
	// auth secret itself. A Super user can perform literally
	// any action it is possible to perform.
	SuperClaimSet = &jws.ClaimSet{
		Sub: "_super",
		Exp: time.Now().AddDate(10, 0, 0).Unix(),
	}

	// NobodyClaimSet is a special jws.ClaimSet returned when
	// the JWT supplied to a Decode call is the empty string.
	NobodyClaimSet *jws.ClaimSet = &jws.ClaimSet{
		Sub: "_nobody",
		Exp: time.Now().AddDate(10, 0, 0).Unix(),
	}

	separator []byte = []byte(".")
)

// Encode converts claimSet into a JWT and signs it with HMAC-256
// using secret.
func Encode(claimSet *jws.ClaimSet, secret []byte) ([]byte, error) {

	result, err := jws.EncodeWithSigner(&jws.Header{
		Typ:       "JWT",
		Algorithm: "HS256",
	}, claimSet, func(data []byte) ([]byte, error) {

		sig := make([]byte, 0, 32)
		mac := hmac.New(sha256.New, secret)
		mac.Write(data)
		return mac.Sum(sig), nil

	})

	return []byte(result), err
}

// Decode checks jwt's signature against secret. If it matches
// and jwt has not expired, Decode returns a jws.ClaimSet containing
// jwt's claims.
//
// There are two special cases:
//
// - If jwt is equal to secret, SuperClaimSet is returned.
//
// - If jwt is the empty string, NobodyClaimSet is returned.
func Decode(jwt []byte, secret []byte) (*jws.ClaimSet, error) {

	var theirSignature [32]byte
	var ourSignature [32]byte

	sepCount := bytes.Count(jwt, separator)

	if sepCount == 0 && hmac.Equal(jwt, secret) {
		// in the special case where the "JWT" is actually just the
		// auth secret itself, the claim is authorized as the SuperClaimSet
		return SuperClaimSet, nil
	} else if len(jwt) == 0 {
		// in the special case where the JWT is nothing, the claim is
		// authorized as the NobodyClaimSet
		return NobodyClaimSet, nil
	} else if sepCount != 2 {
		return nil, InvalidJWTError
	}

	firstSeparatorIndex := bytes.Index(jwt, separator)
	secondSeparatorIndex := bytes.Index(jwt[firstSeparatorIndex+1:], separator)

	header := jwt[0:firstSeparatorIndex]
	payload := jwt[firstSeparatorIndex+1 : firstSeparatorIndex+1+secondSeparatorIndex]
	sig := jwt[firstSeparatorIndex+1+secondSeparatorIndex+1:]

	var decodedHeader jws.Header

	if err := readHeader(header, &decodedHeader); err != nil {
		return nil, err
	} else if decodedHeader.Algorithm != "HS256" {
		return nil, InvalidAlgorithmError
	} else if decodedHeader.Typ != "JWT" {
		return nil, InvalidHeaderError
	}

	_, err := base64.RawURLEncoding.Decode(theirSignature[0:], sig)
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write(header)
	mac.Write(separator)
	mac.Write(payload)
	mac.Sum(ourSignature[:0])

	if !hmac.Equal(theirSignature[0:], ourSignature[0:]) {
		return nil, BadSignatureError
	}

	// decode the claim set and ensure it hasn't expired
	now := time.Now().Unix()
	claimSet, err := jws.Decode(string(jwt))
	if err != nil {
		return nil, err
	} else if claimSet.Exp <= now {
		return nil, ExpiredJWTError
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(string(payload))
	if err != nil {
		return nil, err
	}

	// TODO(goldibex): This operation is likely to be quite slow due to reflection
	// Ideally, what we want is to accept an interface{} to decode the private claims
	// of the JWT into for the client
	claimSet.PrivateClaims = map[string]interface{}{}
	err = json.Unmarshal(payloadBytes, &claimSet.PrivateClaims)
	if err != nil {
		return nil, err
	}

	delete(claimSet.PrivateClaims, "iss")
	delete(claimSet.PrivateClaims, "sub")
	delete(claimSet.PrivateClaims, "iat")
	delete(claimSet.PrivateClaims, "exp")
	delete(claimSet.PrivateClaims, "aud")
	delete(claimSet.PrivateClaims, "scope")

	return claimSet, nil

}

func readHeader(headPart []byte, header *jws.Header) error {

	jsonHeader, err := base64.RawURLEncoding.DecodeString(string(headPart))
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonHeader, header)

}
