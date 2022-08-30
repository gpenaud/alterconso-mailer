package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// Mail ....
type Mail struct {
	Sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

// SMTPServer ...
type SMTPServer struct {
	Host      string
	Port      string
	TLSConfig *tls.Config
}

// ServerName ...
func (s *SMTPServer) ServerName() string {
	return s.Host + ":" + s.Port
}

// EmailController controls email operations
func EmailController(subject, body, sender, to, cc, bcc string) {

	toArr := strings.Split(to, ";")
	ccArr := strings.Split(cc, ";")
	bccArr := strings.Split(bcc, ";")
	request := Mail{
		Sender: sender,
		To: toArr,
		Cc: ccArr,
		Bcc: bccArr,
		Subject: subject,
		Body: body,
	}
	fmt.Println("request: ", request)

	send(request)
}

// BuildMessage ...
func (mail *Mail) BuildMessage() string {
	header := ""
	header += fmt.Sprintf("From: %s\r\n", mail.Sender)
	if len(mail.To) > 0 {
		header += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	}
	if len(mail.Cc) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ";"))
	}

	header += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	header += "\r\n" + mail.Body

	return header
}

func send(mailObject Mail) {

	mail := Mail{}
	mail.Sender = mailObject.Sender
	mail.To = mailObject.To
	mail.Cc = mailObject.Cc
	mail.Bcc = mailObject.Bcc
	mail.Subject = mailObject.Subject
	mail.Body = mailObject.Body

	messageBody := mail.BuildMessage()

	smtpServer := SMTPServer{Host: "ssl0.ovh.net", Port: "465"}
	smtpServer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.Host,
	}

	auth := smtp.PlainAuth("", mail.Sender, "elisabeth", smtpServer.Host)

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), smtpServer.TLSConfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.Host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.Sender); err != nil {
		log.Panic(err)
	}

	receivers := mail.To

	if len(mail.Cc) > 0 && mail.Cc[0] != "" {
		receivers = append(receivers, mail.Cc...)
	}

	if len(mail.Bcc) > 0 && mail.Bcc[0] != "" {
		receivers = append(receivers, mail.Bcc...)
	}

	for _, k := range receivers {
		log.Println("sending to: ", k)
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	wr, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = wr.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = wr.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")
}

func main() {
	EmailController("mail de test", "Envoyé par mailer.go", "alterconso@leportail.org", "alterconso@leportail.org", "", "guillaume.penaud@gmail.com")
}
