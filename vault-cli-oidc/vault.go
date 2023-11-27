package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/skratchdot/open-golang/open"
)

func LoginExample(vaultURL string) {
	// 1. Setup an HTTP server to recieve the oidc callback
	httpServer := &http.Server{Addr: "127.0.0.1:8250"}
	http.HandleFunc("/", handleCallback)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("failed to start callback webserver: %s", err.Error())
		}
	}()
	defer httpServer.Shutdown(context.TODO())

	ctx := context.Background()

	// 2.  prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress(vaultURL),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Get the authorization URL. In our case this will be google.
	r, err := client.Auth.JwtOidcRequestAuthorizationUrl(ctx,
		schema.JwtOidcRequestAuthorizationUrlRequest{
			RedirectUri: "http://localhost:8250/oidc/callback",
		},
		vault.WithMountPath("oidc"),
	)
	if err != nil {
		fmt.Println("error JwtOidcRequestAuthorizationUrl:", err.Error())
		os.Exit(42)
	}

	// 4a. Open the authroziation url
	if u, ok := r.Data["auth_url"].(string); ok {
		fmt.Printf("response: %s\n\n", u)
		open.Run(u)
	}

	// 4b. Wait for CallBackURL to be defined
	for {
		if CallBackURL != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 5. Send callback info to vault and get our ResponseAuth
	r, err = client.Auth.JwtOidcCallback(ctx,
		"",
		CallBackURL.Query().Get("code"),
		CallBackURL.Query().Get("state"),
		vault.WithMountPath("oidc"),
	)
	if err != nil {
		fmt.Printf("error JwtOidcCallback: %s\n\n", err.Error())
		os.Exit(42)
	}

	// testing: peak inside the reponse
	// fmt.Printf("JwtOidcCallback response: %+v\n", r)
	// fmt.Printf("Auth.EntityID: %+v\n", r.Auth.EntityID)
	// fmt.Printf("Auth.ClientToken: %+v\n", r.Auth.ClientToken)
	// fmt.Printf("Auth.Accessor: %+v\n", r.Auth.Accessor)
	// fmt.Printf("Auth.Policies: %+v\n", r.Auth.Policies)
	// fmt.Printf("Auth.TokenPolicies: %+v\n", r.Auth.TokenPolicies)
	// fmt.Printf("Auth.IdentityPolicies: %+v\n", r.Auth.IdentityPolicies)
	// fmt.Printf("Auth.LeaseDuration: %+v\n", r.Auth.LeaseDuration)
	// fmt.Printf("Auth.Renewable: %+v\n", r.Auth.Renewable)
	// fmt.Printf("Auth.Metadata: %+v\n", r.Auth.Metadata)

	// 6. Set client token to the token returned in the ResponseAuth
	if err := client.SetToken(r.Auth.ClientToken); err != nil {
		fmt.Printf("client.SetToken error: %s\n", err.Error())
		os.Exit(84)
	}

	// 7. Get secrets (this secret is for testing and not sensitive)
	resp, err := client.Secrets.KvV2Read(
		ctx,
		"developers/test",
		vault.WithMountPath("secret"),
	)
	fmt.Printf("resp: %+v", resp)
}
