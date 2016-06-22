package app

import (
	"github.com/the-information/api2"
	"golang.org/x/net/context"
	"net/http"
)

func readConfig(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	var conf api.Config
	if err := api.GetConfig(ctx, &conf); err != nil {
		api.WriteJSON(w, err)
	} else {
		api.WriteJSON(w, &conf)
	}

}

func updateConfig(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	var conf api.Config
	if err := api.GetConfig(ctx, &conf); err != nil {
		api.WriteJSON(w, err)
	} else if err := api.ReadJSON(r, &conf); err != nil {
		api.WriteJSON(w, err)
	} else if err := api.SaveConfig(ctx, &conf); err != nil {
		api.WriteJSON(w, err)
	} else {
		api.WriteJSON(w, &conf)
	}

}
