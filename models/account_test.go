package models

import (
	"testing"
)

func TestAccountSuper(t *testing.T) {

	blank := &Account{}
	if blank.Super() {
		t.Errorf("Expected a blank account not to be super, but it was")
	}
	if !SuperAccount.Super() {
		t.Errorf("Expected SuperAccount to be super, but it wasn't")
	}
}

func TestAccountNobody(t *testing.T) {

	blank := &Account{}
	if blank.Nobody() {
		t.Errorf("Expected a blank account not to be nobody, but it was")
	}
	if !NobodyAccount.Nobody() {
		t.Errorf("Expected NobodyAccount to be nobody, but it wasn't")
	}

}
