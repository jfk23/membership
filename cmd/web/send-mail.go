package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/jfk23/gobookings/internal/model"
	"github.com/xhit/go-simple-mail/v2"
)

func ListenForMail() {
	go func() {
		for {
			msg := <- appConfig.MailChan
			sendMail(msg)
		}

	}()
}

func sendMail(m model.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10*time.Second
	server.SendTimeout = 10*time.Second

	client, err := server.Connect()
	if err != nil {
		appConfig.ErrorLog.Println(err)
	}

	message := mail.NewMSG()
	message.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		message.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-template/%s", m.Template))
		if err != nil {
			appConfig.ErrorLog.Println(err)
		}
		dataString := string(data)
		msgToSend := strings.Replace(dataString, "[%body%]", m.Content, 1)
		message.SetBody(mail.TextHTML, msgToSend)
	}
	
	err = message.Send(client)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("email sent succesfully!")
	}


}