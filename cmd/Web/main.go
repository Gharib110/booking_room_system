package main

import (
	"encoding/gob"
	"fmt"
	"github.com/DapperBlondie/booking_system/pkg/models"
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
	err := run()
	if err != nil {
		log.Fatal("This error occurred in the run method : ", err)
	}
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

func run() error {

	appConfig = new(config.AppConfig)
	appConfig.IsProduction = false

	gob.Register(models.ReservationData{})
	session = scs.New()
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Secure = appConfig.IsProduction

	appConfig.Session = session

	tmplCache, err := renderer.CreateCacheTemplates()
	if err != nil {
		log.Println("This is an error about CreateCacheTemplates.")
		return err
	}

	appConfig.TemplateCache = tmplCache
	appConfig.UseCache = true

	renderer.NewAppConfig(appConfig)
	repo := handlers.NewRepo(appConfig)
	handlers.NewHandlers(repo)

	return nil
}