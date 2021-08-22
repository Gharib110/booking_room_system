# BookingManagementSystem
An application for manage a website that rent rooms for any type of meetings .
that have two type such as generals or majors rooms for any type of porpuses.
I use html template library for generating and parsing HTML templates.

***

## Rendering HTML Templates
I use native golang libraries for creating and caching all templates for every pages that My app have.<br>
Creating cache from Templates:<br>

- First I get all my templates by regex then I create a template for each of them and parse them.<br>
- Second I get all layouts that I designed for applying on every page that I have with ParseGlob
and parse with every template that I created.<br>
- Third I return my templates as a ``` map[string]*template.Template ```
 ```go
 
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
 
 ```

***
## Rendering HTML Templates
 I implement two function for rendering a html template : <br>

- First I get information and structures that I need to such as ```http.ResponseWriter```
- Second pull out the target template from my template cache then execute it with default function that I use for 
using some information And also CSRF token for preventing cross-site-request-forgery

```go
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
```

 ```go
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
```
***
## Reservation E-mail
 I use a goroutine that I run it at the beginning of the application
 to get the ``` MailData ``` through a channel then pass it to a function that I used
 it for sending email to a specific mail.<br>

```go
//MailData using for sending mail management structure
type MailData struct {
To      string
From    string
Subject string
Content string
}

func listenToMail() {
	go func() {
		for {
			msg := <-appConfig.MailChan
			sendMsg(msg)
		}
	}()
}
```
<br>

```go
func sendMsg(msg models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.SendTimeout = time.Second * 10
	server.ConnectTimeout = time.Second * 10

	client, err := server.Connect()
	if err != nil {
		log.Println("Some error in getting the client from smtp server : " + err.Error() + "\n")
		return
	}

	mailTemplate, err := ioutil.ReadFile("./html_sources/basic_email_template.html")
	if err != nil {
		log.Println("Error in reading the html template : " + err.Error() + "\n")
	}

	mailStr := string(mailTemplate)
	mailStr = strings.Replace(mailStr, "[%name%]", msg.To, 2)
	mailStr = strings.Replace(mailStr, "[%subject%]", msg.Subject, 2)
	mailStr = strings.Replace(mailStr, "[%content%]", msg.Content, 2)

	email := mail.NewMSG()
	email.SetSubject(msg.Subject)
	email.SetFrom(msg.From)
	email.SetFrom(msg.To)
	email.SetBody(mail.TextHTML, mailStr)

	err = email.Send(client)
	if err != nil {
		log.Println("Error in sending mail : " + err.Error() + "\n")
		return
	}
}
```
<br>

***
## Special Thanks 
 I should say thank you to these legend developers for creating and
 maintaining fabulous packages.

```go
        github.com/alexedwards/scs/v2 v2.4.0
	github.com/asaskevich/govalidator
	github.com/go-chi/chi v1.5.4
	github.com/jackc/pgconn v1.8.1
	github.com/jackc/pgx/v4 v4.11.0
	github.com/justinas/nosurf v1.1.1
	github.com/xhit/go-simple-mail/v2 v2.9.1
```