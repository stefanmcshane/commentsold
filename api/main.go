package main

import (
	"context"
	"log"
	"net/http"

	"github.com/stefanmcshane/commentsold/auth"
	"github.com/stefanmcshane/commentsold/db"
)

func main() {
	ctx := context.Background()

	// Initialise DB
	csvdb, err := db.NewCSVDatabase(ctx, inventoryData)
	if err != nil {
		log.Fatal(err)
	}

	// Setup JWT Signing cert
	signingKey := auth.JWTToken{
		SigningKey: "SuperSecureKey",
	}

	mux := http.NewServeMux()

	signinHandler := http.HandlerFunc(signin)
	mux.Handle("/api/signin", authenticateUser(signingKey, signinHandler))

	// DB endpoints
	h := NewDBHandler(csvdb)
	mux.Handle("/api/inventory", checkTokenValidity(signingKey, http.HandlerFunc(h.listInventory)))
	mux.Handle("/api/inventory/", checkTokenValidity(signingKey, http.HandlerFunc(h.manageInventoryItemWithID)))

	log.Println("Listening on :3000...")
	err = http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
