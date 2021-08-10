package main

import (
	"github.com/DapperBlondie/booking_system/pkg/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func listenToMail() {
	go func() {
		for {
			msg := <-appConfig.MailChan
			sendMsg(msg)
		}
	}()
}

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
