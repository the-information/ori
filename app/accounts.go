package app

func createAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	var accountCreationRequest models.AccountCreationRequest
	if err := api.ReadJSON(r, &accountCreationRequest); err != nil {
		api.WriteJSON(w, err)
		return
	}

}

func readAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}
