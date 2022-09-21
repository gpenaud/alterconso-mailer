package internal

import (
	"crypto/tls"
	"fmt"
	log     "github.com/sirupsen/logrus"
	"net/smtp"
	"strings"
	// viper   "github.com/spf13/viper"
)

type SMTPServer struct {
	Host      string
	Port      string
	Password  string
	TLSConfig *tls.Config
}

// Mail ....
type Mail struct {
	Sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

// ServerName ...
func (s *SMTPServer) ServerName() string {
	return s.Host + ":" + s.Port
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
	header += "MIME-Version: 1.0;\n"
	header += "Content-Type: text/html; charset=\"utf-8\";\n\n"

	header += "\r\n" + mail.Body

	return header
}

func send(smtpServer SMTPServer, mail Mail) {
	messageBody := mail.BuildMessage()
	smtpServer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.Host,
	}

	fmt.Println("Connection parameters:\n")
	fmt.Println(smtpServer.ServerName())
	fmt.Println(smtpServer.TLSConfig)
	fmt.Println(mail.Sender)
	fmt.Println(smtpServer.Password)
	fmt.Println(smtpServer.Host)

	auth := smtp.PlainAuth("", mail.Sender, smtpServer.Password, smtpServer.Host)
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

	// log.Println(messageBody)
	log.Println("Mail sent successfully")
}
