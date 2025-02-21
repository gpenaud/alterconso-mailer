package tools

func RespondHTTPCodeOnly(w http.ResponseWriter, code int) {
  w.WriteHeader(code)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
  log.Error(message)
  respondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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
