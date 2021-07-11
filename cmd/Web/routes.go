package main

import (
	"github.com/DapperBlondie/booking_system/pkg/config"
	"github.com/DapperBlondie/booking_system/pkg/handlers"
	"github.com/bmizerany/pat"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

// routes manage out Application Routes for us using pat pkg
func routes(appConfig *config.AppConfig) http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Repo.HomePg))
	mux.Get("/About", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/Addition", http.HandlerFunc(handlers.Repo.AdditionPg))
	mux.Get("/Division", http.HandlerFunc(handlers.Repo.DivisionPg))

	return mux
}

// chiRoutes use chi pkg for creating a request multiplexer for us
func chiRoutes(appConfig *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Use(middleware.Recoverer)
	mux.Use(CSRFTokenGenerator)
	mux.Use(SessionLoad)
	mux.Use(AuthChecker)
	mux.Use(WriteToConsole)

	mux.Get("/", handlers.Repo.HomePg)
	mux.Get("/About", handlers.Repo.About)

	mux.Get("/Addition", handlers.Repo.AdditionPg)
	mux.Get("/Division", handlers.Repo.DivisionPg)

	mux.Get("/Majors", handlers.Repo.Majors)
	mux.Get("/Generals", handlers.Repo.Generals)

	mux.Get("/Availability", handlers.Repo.Availability)
	mux.Post("/Availability", handlers.Repo.PostAvailability)
	mux.Post("/Availability-json", handlers.Repo.JSONAvailability)

	mux.Get("/Reserve", handlers.Repo.Reservation)
	mux.Post("/Reserve", handlers.Repo.PostReservation)
	mux.Get("/Reservation_Summary", handlers.Repo.ReservationSummary)
	mux.Get("/choose_room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book_room", handlers.Repo.BookRoom)
	mux.Get("/User/Login", handlers.Repo.LoginHandler)
	mux.Post("/User/Login", handlers.Repo.PostLoginHandler)
	mux.Get("/User/Logout", handlers.Repo.LogoutHandler)

	mux.Route("/admin", func(router chi.Router) {
		router.Use(AuthChecker)
		router.Get("/dashboard", handlers.Repo.AdminDashboard)
	})

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
