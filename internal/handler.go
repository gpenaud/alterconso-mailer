package internal

import (
  fmt     "fmt"
  http    "net/http"
  json    "encoding/json"
  log     "github.com/sirupsen/logrus"
  // mux     "github.com/gorilla/mux"
  // "github.com/gpenaud/alterconso-mailer/internal/user"
  "io/ioutil"
  // "html"
)

var handlerLog *log.Entry

func init() {
  fmt.Println("init")
}

// -------------------------------------------------------------------------- //
// Common functions for handlers

func respondHTTPCodeOnly(w http.ResponseWriter, code int) {
  w.WriteHeader(code)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
  handlerLog.Error(message)
  respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  response, err := json.Marshal(payload)
  if err != nil {
    fmt.Printf("Error: %s", err)
    return;
  }
  // response := payload
  log.Info(fmt.Sprintf("Payload: %s", payload))
  log.Info(fmt.Sprintf("JSON response: %s", response))

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)

  fmt.Println(string(response))
}

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

func (a *Application) emailController(writer http.ResponseWriter, request *http.Request) {
  defer request.Body.Close()

  requestBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("Oh no! There was an error:", err)
		return
	}

  var params MailParameters
  json.Unmarshal(requestBytes, &params)

  to  := []string { params.FromEmail }
  cc  := []string {}
  bcc := params.MailList()

  mail_request := Mail {
    Sender: params.FromEmail,
    To: to,
    Cc: cc,
    Bcc: bcc,
    Subject: params.Subject,
    Body: params.Body,
  }

  fmt.Println("request: ", mail_request)
	send(mail_request)
}
