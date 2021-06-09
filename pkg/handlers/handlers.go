package handlers

import (
	"errors"
	"fmt"
	"github.com/DapperBlondie/booking_system/pkg/config"
	"github.com/DapperBlondie/booking_system/pkg/models"
	"github.com/DapperBlondie/booking_system/pkg/renderer"
	"log"
	"math"
	"net/http"
)

type Repository struct {
	AppConf *config.AppConfig
}

// Repo a variable we can share our wide configuration with handlers
var Repo *Repository

// NewRepo make a repository for our handlers
// NewRepo Such as constructor for our struct
func NewRepo(ac *config.AppConfig) *Repository {

	return &Repository{

		AppConf: ac,
	}
}

// NewHandlers assign the repo to internal Repo variable
func NewHandlers(repo *Repository) {

	Repo = repo
}

// HomePg handle function  for HomePage
func (repo *Repository) HomePg(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	repo.AppConf.Session.Put(r.Context(), "remote-ip", remoteIP)
	renderer.RenderByCacheTemplates(&w, r, "home.page.tmpl", &models.TemplateData{})
}

// About handle the about page
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	renderer.RenderByCacheTemplates(&w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: map[string]string{"test": "Hello, I am here !",
			"remote_ip": repo.AppConf.Session.GetString(r.Context(), "remote-ip")},
	})
}

// Generals for handling the generals rooms
func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	renderer.RenderByCacheTemplates(&w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors for handling the majors rooms
func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	renderer.RenderByCacheTemplates(&w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability for handling the availability page
func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	renderer.RenderByCacheTemplates(&w, r, "availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability for handling the availability page
func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Your Request posted"))
	if err != nil {
		fmt.Fprintf(w, "We get an error")
	}
}

// Reservation for handling the Reserve operation
func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	renderer.RenderByCacheTemplates(&w, r, "make_reservation.page.tmpl", &models.TemplateData{})
}

// AdditionPg handle /About
func (repo *Repository) AdditionPg(w http.ResponseWriter, r *http.Request) {

	result, err := addValues(12, 45)
	if err == nil {

		_, err := fmt.Fprintf(w, fmt.Sprintf("The result is : %d\n", result))
		if err != nil {
			_, _ = fmt.Fprintf(w, "Some Error Occurred !")
			log.Print("This error occurred : ", err)
		}
	}
}

// DivisionPg handle /Division
func (repo *Repository) DivisionPg(w http.ResponseWriter, r *http.Request) {

	result, err := divisionValues(23, 0)
	if err != nil {

		_, err := fmt.Fprintf(w, fmt.Sprintf("The result is %f with this msg: %s\n", result, err.Error()))
		if err != nil {

			_, _ = fmt.Fprintf(w, fmt.Sprintf("This error occurred : %s\n", err.Error()))
		}
	} else {

		_, err = fmt.Fprintf(w, fmt.Sprintf("The result is %f\n", result))
		if err != nil {

			_, _ = fmt.Fprintf(w, fmt.Sprintf("This error occurred : %s\n", err.Error()))
		}
	}
}

// addValues add two integer value and return error and result
func addValues(a, b int) (int, error) {

	return a + b, nil
}

// divisionValues divide two float32 value and return error and result
func divisionValues(a, b float32) (float32, error) {

	if b == 0 && a == 0 {

		return float32(math.NaN()), errors.New("0/0 is not valid you gotten NAN")
	}

	if b == 0 {

		if a > 0 {

			return float32(math.Inf(1)), errors.New("a / 0 is Positive Infinity")
		}

		if a < 0 {

			return float32(math.Inf(-1)), errors.New("a / 0 is Negative Infinity")
		}
	}

	return a / b, nil
}
