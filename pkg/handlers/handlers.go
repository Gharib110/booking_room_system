package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DapperBlondie/booking_system/pkg/config"
	"github.com/DapperBlondie/booking_system/pkg/driver"
	"github.com/DapperBlondie/booking_system/pkg/models"
	"github.com/DapperBlondie/booking_system/pkg/renderer"
	"github.com/DapperBlondie/booking_system/pkg/repository"
	"github.com/DapperBlondie/booking_system/pkg/repository/dbrepo"
	"github.com/DapperBlondie/booking_system/validation"
	"github.com/go-chi/chi"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Repository struct {
	AppConf *config.AppConfig
	DB      repository.DatabaseRepo
}

type JsonResponse struct {
	OK      int    `json:"ok"`
	Message string `json:"message"`
}

type BookRoomResponse struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	OK        int    `json:"ok"`
}

// Repo a variable we can share our wide configuration with handlers
var Repo *Repository

// NewRepo make a repository for our handlers, Such as constructor for our struct
func NewRepo(ac *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		AppConf: ac,
		DB:      dbrepo.NewPostgresRepo(db.SQL, ac),
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
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	layout := "2006-01-02"
	start_date, err := time.Parse(layout, start)
	if err != nil {
		log.Println("Error in parsing the start_date : " + err.Error())
	}
	end_date, err := time.Parse(layout, end)
	if err != nil {
		log.Println("Error in parsing the end_date : " + err.Error())
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(start_date, end_date)
	if err != nil {
		_, err := fmt.Fprint(w, "We got an error during search in our database for availability")
		if err != nil {
			return
		}
		return
	}

	if len(rooms) != 0 {
		for _, room := range rooms {
			fmt.Println("ROOM : ", room.ID, room.RoomName)
		}
		data := make(map[string]interface{})
		data["rooms"] = rooms
		renderer.RenderByCacheTemplates(&w, r, "choose_room.page.tmpl", &models.TemplateData{
			Data: data,
		})

		res := models.Reservations{
			StartDate: start_date,
			EndDate:   end_date,
		}
		repo.AppConf.Session.Put(r.Context(), "reservation", res)
	} else {
		repo.AppConf.Session.Put(r.Context(), "error", "NO Availability Exists")
		http.Redirect(w, r, "/Availability", http.StatusSeeOther)
		return
	}

	_, err = w.Write([]byte("Unexpected trouble happened !"))
}

// JSONAvailability used for getting the information for start date and end date
func (repo *Repository) JSONAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Something went wrong")
		log.Println("Error in parsing form in JSONAvailablity : " + err.Error() + "\n")
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))

	layout := "2006-08-12"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	available, _ := repo.DB.SearchAvailabilityByDateByRoomID(startDate, endDate, roomId)

	resp := &JsonResponse{
		OK:      available,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		fmt.Fprint(w, "We can not give you the JSON file :(")
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		log.Println(err)
	}
}

//BookRoom used for getting the booking information from URL query params
func (repo *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	sd := r.URL.Query().Get("start_date")
	ed := r.URL.Query().Get("end_date")

	roomResp := &BookRoomResponse{
		StartDate: sd,
		EndDate:   ed,
		OK:        200,
	}

	roomRespB, err := json.MarshalIndent(roomResp, "", "   ")
	if err != nil {
		log.Println("could not be able to marshal the BookRoomResponse : " + err.Error() + "\n")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(roomRespB)
	if err != nil {
		log.Println("Error occurred in writing response into the response writer : " + err.Error() + "\n")
		return
	}
}

// Reservation for handling the Reserve operation
func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := repo.AppConf.Session.Get(r.Context(), "reservation").(models.Reservations)
	if !ok {
		log.Println("An error occurred in Reservation handler : " + errors.New("not be able to retrieve reservation data\nfrom the session").Error())
		return
	}

	room, err := repo.DB.GetRoomByID(res.RoomID)
	if err != nil {
		log.Println("An error occurred in Reservation handler : " + err.Error() + "\n")
		return
	}

	res.Room.RoomName = room.RoomName

	data := make(map[string]interface{})
	data["reservation"] = res

	repo.AppConf.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-11-12")
	ed := res.StartDate.Format("2006-11-12")

	strMap := make(map[string]string)
	strMap["start_date"] = sd
	strMap["end_date"] = ed

	renderer.RenderByCacheTemplates(&w, r, "make_reservation.page.tmpl", &models.TemplateData{
		Form:      validation.New(nil),
		Data:      data,
		StringMap: strMap,
	})
}

// PostReservation handling the post request that send from // make_reservation.page.tmpl reservation form
func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := repo.AppConf.Session.Get(r.Context(), "reservation").(models.Reservations)
	if !ok {
		log.Println("An error occurred in PostReservation : " +
			errors.New("we can not pullout the session from browser").Error())
	}
	err := r.ParseForm()
	if err != nil {
		log.Println("We have some errors in parsing reservation from data")
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	start_date, err := time.Parse(layout, sd)
	if err != nil {
		log.Println("Error in parsing the start_date : " + err.Error())
	}
	end_date, err := time.Parse(layout, ed)
	if err != nil {
		log.Println("Error in parsing the end_date : " + err.Error())
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		log.Println("Error in parsing the room_id : " + err.Error())
	}

	reservation := &models.ReservationData{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	reservationDB := &models.Reservations{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: res.StartDate,
		EndDate:   res.EndDate,
		RoomID:    roomID,
	}

	form := validation.New(r.PostForm)
	form.RequiredField("first_name", "last_name", "phone", "email")
	form.MinLength("first_name", 3, r)
	//form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		repo.AppConf.Session.Put(r.Context(), "reservation", reservation)

		renderer.RenderByCacheTemplates(&w, r, "make_reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	newReservationID, err := repo.DB.InsertReservation(*reservationDB)
	if err != nil {
		log.Fatal("We can not be able to add reservation data into the database : " + err.Error())
	}

	restrictionData := models.RoomsRestrictions{
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		StartDate:     start_date,
		EndDate:       end_date,
	}

	err = repo.DB.InsertRestrictionRoom(restrictionData)
	if err != nil {
		log.Println("We can not be able to add restriction data into the database : " + err.Error())
	}

	emailMsg := &models.MailData{
		To:      res.Email,
		From:    "alirezagharib@dapper.me",
		Subject: "Reservation Confirmation",
		Content: "FirstName : " + reservationDB.FirstName + "\n" + "LastName : " + reservationDB.LastName + "\n" + "E-mail : " + reservationDB.Email + "\n",
	}

	repo.AppConf.MailChan <- *emailMsg

	repo.AppConf.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/Reservation_Summary", http.StatusSeeOther)
}

// ReservationSummary for getting the PostReservation data and show them to the user
func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.AppConf.Session.Get(r.Context(), "reservation").(models.ReservationData)
	if !ok {
		log.Println("We have error in finding the reservation data in session.")
		repo.AppConf.Session.Put(r.Context(), "error", "We can not find your reservation data.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	renderer.RenderByCacheTemplates(&w, r, "reservation_summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// ChooseRoom for handle choosing a room from available rooms
func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_, err = fmt.Fprintf(w, "There is no any id Comeback to the homepage")
		return
	}

	resData, ok := repo.AppConf.Session.Get(r.Context(), "reservation").(models.Reservations)
	if !ok {
		_, err := fmt.Fprintf(w, "We have some troble for getting the reservation data in browser cookie")
		if err != nil {
			return
		}
		return
	}

	resData.RoomID = roomID
	repo.AppConf.Session.Put(r.Context(), "reservation", resData)
	http.Redirect(w, r, "/Reserve", http.StatusSeeOther)
	return
}

func (repo *Repository) LoginHandler(w http.ResponseWriter, r *http.Request)  {
	renderer.RenderByCacheTemplates(&w, r, "login.page.tmpl",
		&models.TemplateData{Form: validation.New(nil)})
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
