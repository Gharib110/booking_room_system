package config

import (
	"github.com/DapperBlondie/booking_system/pkg/models"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLogger    *log.Logger
	IsProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
