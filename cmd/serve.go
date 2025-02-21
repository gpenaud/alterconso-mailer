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

var (
  configFile string

  serverCmd = &cobra.Command{
  	Use:   "serve",
  	Short: "runs mail server",
  	Long:  `this command runs mail server and it won't turn off till you manually do it'`,
  	Run:   runServer,
  }
)

func init() {
	cobra.OnInitialize(initConfig)
	serverCmd.Flags().StringVar(&configFile, "config-file", "", "config file (default is $CWD/config.yaml)")
	serverCmd.Flags().StringP("server-host", "H", "", "Address to run Application server on")
	viper.BindPFlag("server_host", serverCmd.Flags().Lookup("server-host"))
  viper.SetDefault("server_host", viper.GetString("server_host"))

	serverCmd.Flags().IntP("server-port", "P", 0, "Port to run Application server on")
	viper.BindPFlag("server_port", serverCmd.Flags().Lookup("server-port"))
  viper.SetDefault("server_port", viper.GetString("server_port"))

  serverCmd.Flags().StringP("smtp-username", "u", "", "SMTP account username")
  viper.BindPFlag("smtp_username", serverCmd.Flags().Lookup("smtp-username"))
  viper.SetDefault("smtp_username", viper.GetString("smtp_username"))

  fmt.Println(viper.GetString("smtp_password"))

  serverCmd.Flags().StringP("smtp-password", "p", "", "SMTP account password")
  viper.BindPFlag("smtp_password", serverCmd.Flags().Lookup("smtp-password"))
  viper.SetDefault("smtp_password", viper.GetString("smtp_password"))

	rootCmd.AddCommand(serverCmd)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigFile(configFile)
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

  err := viper.ReadInConfig()
  if err != nil {
    panic(fmt.Errorf("fatal error config file: %w", err))
  }
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
