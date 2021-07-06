package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/DapperBlondie/booking_system/pkg/config"
	"github.com/DapperBlondie/booking_system/pkg/driver"
	"github.com/DapperBlondie/booking_system/pkg/handlers"
	"github.com/DapperBlondie/booking_system/pkg/models"
	"github.com/DapperBlondie/booking_system/pkg/renderer"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const (
	ConnHost         = "localhost:"
	ConnPort         = "8080"
	PostgresDBString = "host=localhost port=5720 dbname=postgres user=postgres password=alireza1380##"
)

var appConfig *config.AppConfig
var session *scs.SessionManager

func main() {
	dbConn, err := run()
	defer func(SQL *sql.DB) {
		err := SQL.Close()
		if err != nil {
			log.Fatal(err.Error() + " ,We can not be able to close connection to postgreSQL :(")
		}
	}(dbConn.SQL)

	defer close(appConfig.MailChan)

	log.Println("Dummy Email server configured ...")
	listenToMail()

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

func run() (*driver.DB, error) {

	appConfig = new(config.AppConfig)
	appConfig.IsProduction = false

	gob.Register(models.ReservationData{})
	gob.Register(models.Rooms{})
	gob.Register(models.Users{})
	gob.Register(models.Restrictions{})
	gob.Register(models.Reservations{})

	mailChan := make(chan models.MailData)
	appConfig.MailChan = mailChan

	session = scs.New()
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Secure = appConfig.IsProduction

	appConfig.Session = session

	dbConn, err := driver.ConnectSQL(PostgresDBString)
	if err != nil {
		log.Fatal(err.Error() + " ,During connecting to the database !")
		return nil, err
	}

	tmplCache, err := renderer.CreateCacheTemplates()
	if err != nil {
		log.Println("This is an error about CreateCacheTemplates.")
		return nil, err
	}

	appConfig.TemplateCache = tmplCache
	appConfig.UseCache = true

	renderer.NewAppConfig(appConfig)
	repo := handlers.NewRepo(appConfig, dbConn)
	handlers.NewHandlers(repo)

	return dbConn, nil
}
