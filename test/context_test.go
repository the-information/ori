package test

import (
	"github.com/the-information/ori/account"
	"github.com/the-information/ori/account/auth"
	"github.com/the-information/ori/config"
	"golang.org/x/net/context"
	"testing"
)

type FakeConf struct {
	Foo string
}

func TestBlessContext(t *testing.T) {

	conf := FakeConf{"bar"}
	conf2 := FakeConf{}

	acct := account.Account{}

	ctx := context.Background()

	// without decoration, we should get an error

	if err := config.Get(ctx, &conf); err != config.ErrNotInConfigContext {
		t.Errorf("Got unexpected error %s from config.Get undecorated", err)
	}

	if err := auth.GetAccount(ctx, &acct); err != auth.ErrNotInAuthContext {
		t.Errorf("Got unexpected error %s from auth.GetAccount undecorated", err)
	}

	// with decoration, we should not get an error

	ctx, err := BlessContext(ctx, &conf, &account.Super)
	if err != nil {
		t.Fatalf("Unexpected error from BlessContext: %s", err)
	}

	if err := config.Get(ctx, &conf2); err != nil {
		t.Errorf("Got unexpected error %s", err)
	} else if conf2.Foo != "bar" {
		t.Errorf("Config wasn't blessed properly")
	}

	if err := auth.GetAccount(ctx, &acct); err != nil {
		t.Errorf("Got unexpected error %s", err)
	} else if !acct.Super() {
		t.Errorf("Config wasn't blessed properly")
	}

}
