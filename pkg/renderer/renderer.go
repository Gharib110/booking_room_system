package renderer

import (
	"bytes"
	"fmt"
	"github.com/DapperBlondie/booking_system/pkg/config"
	"github.com/DapperBlondie/booking_system/pkg/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var appConfig *config.AppConfig

// NewAppConfig for setting the appConfig in the renderer.go
func NewAppConfig(config *config.AppConfig) {
	appConfig = config
}

// AddDefaultData use for adding default data to every single *.page.tmpl for using it in them
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = appConfig.Session.PopString(r.Context(), "flash")
	td.Warning = appConfig.Session.PopString(r.Context(), "warning")
	td.Error = appConfig.Session.PopString(r.Context(), "error")

	if appConfig.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	} else {
		td.IsAuthenticated = 0
	}

	td.CSFRToken = nosurf.Token(r)
	return td
}

// RenderTemplate render a go template
func RenderTemplate(w *http.ResponseWriter, tmpl string, tmplData *models.TemplateData) {
	parsedTmpl, err := template.ParseFiles("./templates/" + tmpl)
	if err != nil || parsedTmpl == nil {

		_, _ = fmt.Fprintf(*w, fmt.Sprintf("We have some problem with parsing template :(("))
		log.Println("This problem occurred during parsing the template : ", err.Error())
	} else {

		err = parsedTmpl.Execute(*w, nil)
		if err != nil {

			_, _ = fmt.Fprintf(*w, fmt.Sprintf("We have some problem with executing template :(("))
			log.Println("This problem occurred during executing the template : ", err)
		}
	}
}

// RenderByCacheTemplates render a go template by templatesCache
func RenderByCacheTemplates(w *http.ResponseWriter, r *http.Request, tmpl string, tmplData *models.TemplateData) {
	cacheTmpls := appConfig.TemplateCache

	Tmpl, ok := cacheTmpls[tmpl]
	if ok == false {

		_, _ = fmt.Fprintf(*w, "We do not be able to find your page !")
		log.Println("The tmpl did not found : ", tmpl)
	} else {

		tmplBuff := new(bytes.Buffer)
		err := Tmpl.Execute(tmplBuff, AddDefaultData(tmplData, r))
		if err != nil {

			log.Println("Error occurred during Executing the template !")
			_, _ = fmt.Fprintf(*w, "Error occurred during loading the page ! => \n"+err.Error())
		} else {

			_, err := tmplBuff.WriteTo(*w)
			if err != nil {

				log.Println("Error occurred during writing the buffer to responseWriter !")
				_, _ = fmt.Fprintf(*w, "Error occurred during loading the page !")
			}
		}
	}
}

// CreateCacheTemplates for rendering and caching the templates
func CreateCacheTemplates() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {

		log.Println("We have some error with finding the pages : ", err.Error())
		return myCache, err
	}

	for _, page := range pages {

		pageName := filepath.Base(page)
		parsedTmpl, err := template.New(pageName).Funcs(functions).ParseFiles(page)
		if err != nil {

			log.Println("We have some problem during parsing the templates : ", err.Error())
			return myCache, err
		}

		layouts, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {

			log.Println("We have some problems with finding matches layouts : ", err.Error())
			return myCache, err
		}

		if len(layouts) > 0 {

			parsedTmpl, err = parsedTmpl.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {

				log.Println("We have some problems with parsing layouts with pages : ", err.Error())
				return nil, err
			}

			myCache[pageName] = parsedTmpl
		}
	}

	return myCache, nil
}
