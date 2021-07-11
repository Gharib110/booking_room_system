package models

import "github.com/DapperBlondie/booking_system/validation"

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSFRToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *validation.Form
	IsAuthenticated int
}
