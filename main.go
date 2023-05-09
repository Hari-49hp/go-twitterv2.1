package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func main() {
	goth.UseProviders(
		twitter.New(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("BASE_URL")+"/auth/twitter/callback"),
		// twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),
	)

	m := make(map[string]string)
	m["twitter"] = "Twitter"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	providerIndex := &ProviderIndex{
		Providers:    keys,
		ProvidersMap: m,
	}

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)
	})

	p.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
			t, _ := template.ParseFiles("templates/success.html")
			t.Execute(res, gothUser)
		} else {
			gothic.BeginAuthHandler(res, req)
		}
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, providerIndex)
	})
	log.Println("listening on https://8fd3-103-99-148-171.ngrok-free.app")
	log.Fatal(http.ListenAndServe(":8080", p))
}
