package main

import (
	"fmt"
	"github.com/DapperBlondie/booking_system/pkg/handlers"
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
	csrfTokenHandler := nosurf.New(next)
	csrfTokenHandler.SetBaseCookie(http.Cookie {
		Path:     "/",
		Secure:   appConfig.IsProduction,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfTokenHandler
}

// SessionLoad load the session for us and store them in the browser for each request
// SessionLoad and use this middleware in the main.chiRoutes function
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// AuthChecker for checking the user is authenticated or not
func AuthChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !handlers.Repo.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log-in first !")
			http.Redirect(w, r, "/User/Login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
