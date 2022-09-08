package internal

import (
  fmt     "fmt"
  http    "net/http"
  json    "encoding/json"
  // log     "github.com/sirupsen/logrus"
  // mux     "github.com/gorilla/mux"
  // "github.com/gpenaud/alterconso-mailer/internal/user"
  "io/ioutil"
  // "html"
  viper   "github.com/spf13/viper"
)

type To struct {
  Name string  `json:"name"`
  Email string `json:"email"`
}

func (m MailParameters) MailList() []string {
  var list []string
  for _, to := range m.To {
    list = append(list, to.Email)
  }

  return list
}

type MailParameters struct {
  Subject string    `json:"subject"`
  Body string       `json:"body"`
  FromName string   `json:"from_name"`
  FromEmail string  `json:"from_email"`
  To []To           `json:"to"`
  Headers struct {} `json:"headers"`
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func EmailController(writer http.ResponseWriter, request *http.Request) {
  defer request.Body.Close()

  requestBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Oh no! There was an error:", err)
		return
	}

  smtpServer := SMTPServer{
    Host: viper.GetString("smtp_host"),
    Port: viper.GetString("smtp_port"),
    Password: viper.GetString("smtp_password"),
  }

  var params MailParameters
  json.Unmarshal(requestBytes, &params)

  to  := []string {}
  cc  := []string {}

  bcc_chunks := chunkSlice(params.MailList(), 10)
  from_chunk := []string { params.FromEmail }
  bcc_chunks = append(bcc_chunks, from_chunk)

  for _, chunk := range chunks {
    mail_request := Mail {
      Sender: params.FromEmail,
      To: to,
      Cc: cc,
      Bcc: bcc_chunks,
      Subject: params.Subject,
      Body: params.Body,
    }

    fmt.Println("request: ", mail_request)
    send(smtpServer, mail_request)
  }
}
