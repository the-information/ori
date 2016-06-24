package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2/jws"
	"math/rand"
	"testing"
	"time"
)

var secret []byte = []byte("wat")

func makeFakeJWT(fakeHeader, fakePayload, secret string) (jwt []byte, sub string) {

	sub = fmt.Sprintf("%d", rand.Uint32())
	if len(fakeHeader) == 0 {
		fakeHeader = `{"typ":"JWT","alg":"HS256"}`
	}

	if len(fakePayload) == 0 {
		fakePayload = fmt.Sprintf(`{"sub":"%s","PrivateClaims":{"admin":true},"exp":%d,"nbf":%d}`, sub, time.Now().AddDate(1, 0, 0).Unix(), time.Now().AddDate(-1, 0, 0).Unix())
	}

	if len(secret) == 0 {
		secret = "wat"
	}

	fakeHeaderBase64 := base64.RawURLEncoding.EncodeToString([]byte(fakeHeader))
	fakePayloadBase64 := base64.RawURLEncoding.EncodeToString([]byte(fakePayload))

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fakeHeaderBase64))
	mac.Write([]byte("."))
	mac.Write([]byte(fakePayloadBase64))
	fakeSig := mac.Sum(nil)
	fakeSigBase64 := base64.RawURLEncoding.EncodeToString([]byte(fakeSig))
	jwt = []byte(fmt.Sprintf("%s.%s.%s", fakeHeaderBase64, fakePayloadBase64, fakeSigBase64))
	return

}

func TestInvalidJWT(t *testing.T) {

	if _, err := Decode([]byte("a.b"), secret); err != InvalidJWTError {
		t.Errorf("Should have gotten InvalidJWTError, but got %s", err)
	}

	if _, err := Decode([]byte("a.b.c"), secret); err == nil {
		t.Errorf("Should have gotten an error with a nonsense JWT")
	}

	jwt, _ := makeFakeJWT("", `INVALIDJSON`, "")
	if _, err := Decode(jwt, secret); err == nil {
		t.Errorf("Should have gotten an error with invalid claim body")
	}

}

func TestInvalidHeader(t *testing.T) {

	jwt, _ := makeFakeJWT(`INVALIDJSON`, "", "")
	if _, err := Decode(jwt, secret); err == nil {
		t.Errorf("Should have gotten an error with invalid header JSON")
	}

	jwt, _ = makeFakeJWT(`{"typ":"wat", "alg":"HS256"}`, "", "")
	if _, err := Decode(jwt, secret); err != InvalidHeaderError {
		t.Errorf("Should have gotten InvalidHeaderError, but got %s", err)
	}

	jwt, _ = makeFakeJWT(`{"typ":"JWT", "alg":"none"}`, "", "")
	if _, err := Decode(jwt, secret); err != InvalidAlgorithmError {
		t.Errorf("Should have gotten InvalidAlgorithmError, but got %s", err)
	}

}

func TestExpiredJWT(t *testing.T) {

	jwt, _ := makeFakeJWT("", `{"exp": 0}`, "")
	if _, err := Decode(jwt, secret); err != ExpiredJWTError {
		t.Errorf("Expected ExpiredJWTError, but got %s", err)
	}

}

func TestInvalidSignature(t *testing.T) {

	jwt, _ := makeFakeJWT("", "", "wrong")
	if _, err := Decode(jwt, secret); err != BadSignatureError {
		t.Errorf("Expected BadSignatureError, but got %s", err)
	}

}

func TestValidSignature(t *testing.T) {

	jwt, sub := makeFakeJWT("", "", "")

	// make an account for this sub
	if claimSet, err := Decode(jwt, []byte("wat")); err != nil {
		t.Errorf("Expected no error from Decode, but got %s", err)
	} else if claimSet.Sub != sub {
		t.Errorf("Expected the sub to be %s, but got %s", sub, claimSet.Sub)
	}

}

func TestSuper(t *testing.T) {

	if claimSet, err := Decode([]byte("wat"), []byte("wat")); err != nil {
		t.Errorf("Expected no error from Decode, but got %s", err)
	} else if claimSet != SuperClaimSet {
		t.Errorf("Expected the claimset to be SuperClaimSet, but got %+v", claimSet)
	}

}

func TestEncode(t *testing.T) {

	_, sub := makeFakeJWT("", "", "")

	claimSet := &jws.ClaimSet{
		Sub: sub,
	}

	if result, err := Encode(claimSet, secret); err != nil {
		t.Errorf("Expected no error from Encode, but got %s", err)
	} else if newClaimSet, err := Decode(result, secret); err != nil {
		t.Errorf("Encoded no error using Decode on Encode-d JWT, but got %s", err)
	} else if newClaimSet.Sub != claimSet.Sub {
		t.Errorf("Bad sub on new token: wanted %s, got %s", claimSet.Sub, newClaimSet.Sub)
	}

}
