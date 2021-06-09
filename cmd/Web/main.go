package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DapperBlondie/booking_system/pkg/config"
	"github.com/DapperBlondie/booking_system/pkg/handlers"
	"github.com/DapperBlondie/booking_system/pkg/renderer"
	"github.com/alexedwards/scs/v2"
)

const (
	ConnHost = "localhost:"
	ConnPort = "8080"
)

var appConfig *config.AppConfig
var session *scs.SessionManager

func main() {

	appConfig = new(config.AppConfig)
	appConfig.IsProduction = false

	session = scs.New()
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Secure = appConfig.IsProduction

	appConfig.Session = session

	tmplCache, err := renderer.CreateCacheTemplates()
	if err != nil {

		log.Fatal("We have problem with creating cache for our templates : ", err.Error())
	}

	appConfig.TemplateCache = tmplCache
	appConfig.UseCache = true

	renderer.NewAppConfig(appConfig)
	repo := handlers.NewRepo(appConfig)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:         ConnHost + ConnPort,
		Handler:      chiRoutes(appConfig),
		ReadTimeout:  10 * time.Microsecond,
		WriteTimeout: 10 * time.Millisecond,
	}

	fmt.Println("Basic Server is listening on this address : ", ConnHost+ConnPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("This error occurred during serving : ", err.Error())
	}

	return
}
