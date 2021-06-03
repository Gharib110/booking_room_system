package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

// WriteToConsole provide a functionality for write the every request route to the console
func WriteToConsole(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Request : ", r.URL)
		next.ServeHTTP(w, r)
	})
}

// CSRFTokenGenerator generate a CSRF token for us in the templateData structure
func CSRFTokenGenerator(next http.Handler) http.Handler {

	csrfToken := nosurf.New(next)
	csrfToken.SetBaseCookie(http.Cookie{
		Path:     "/",
		Secure:   appConfig.IsProduction,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfToken
}

// SessionLoad load the session for us and store them in the browser for each request
// SessionLoad and use this middleware in the main.chiRoutes function
func SessionLoad(next http.Handler) http.Handler {

	return session.LoadAndSave(next)
}
