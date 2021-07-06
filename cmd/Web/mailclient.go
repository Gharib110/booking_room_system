package main

import (
	"github.com/DapperBlondie/booking_system/pkg/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
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

	email := mail.NewMSG()
	email.SetSubject(msg.Subject)
	email.SetFrom(msg.From)
	email.SetFrom(msg.To)
	email.SetBody(mail.TextHTML, string(msg.Content))

	err = email.Send(client)
	if err != nil {
		log.Println("Error in sending mail : " + err.Error() + "\n")
		return
	}
}
