package cmd

import (
  fmt   "fmt"
  http  "net/http"
  log   "github.com/sirupsen/logrus"
  mux   "github.com/gorilla/mux"
  cobra "github.com/spf13/cobra"
  viper "github.com/spf13/viper"
  internal "github.com/gpenaud/alterconso-mailer/internal"
)

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "runs mail server",
	Long:  `this command runs mail server and it won't turn off till you manually do it'`,
	Run:   runServer,
}

//adds serverCmd to rootCmd
func init() {
	rootCmd.AddCommand(serverCmd)
}

//runs the server and also does the calculations and send result to client
func runServer(cmd *cobra.Command, args []string) {
  router := mux.NewRouter()
  router.HandleFunc("/send", internal.EmailController).Methods("POST")
  log.Info("Routes are initialized")

  server_address :=
    fmt.Sprintf("%s:%s", viper.GetString("server_host"), viper.GetString("server_port"))

  server_message := fmt.Sprintf(
  `

START INFOS
-----------
Listening alterconso-mailer on %s:%s...

BUILD INFOS
-----------
time: @TODO
release: @TODO
commit: @TODO

`,
    viper.GetString("server_host"),
    viper.GetString("server_port"),
  )

  fmt.Println(server_address)

  httpServer := &http.Server{
		Addr:    server_address,
		Handler: router,
	}

  go func() {
    log.Info(server_message)
    httpServer.ListenAndServe()
  }()

  return
}
