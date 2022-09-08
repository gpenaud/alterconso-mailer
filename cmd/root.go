package cmd

import (

)
import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Version struct {
  BuildTime string
  Commit    string
  Release   string
}

var (
	// Used for flags.
	configFile  string
	secretsFile string

	rootCmd = &cobra.Command{
		Use:   "mailer",
		Short: "A tiny mail server developed with golang",
		Long: `Replace the embedded mail server but wrongly developped
within alterconso web application.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
  cobra.OnInitialize(initSecrets)

	rootCmd.PersistentFlags().StringVar(&configFile, "config-file", "", "config file (default is $CWD/config.yaml)")
  rootCmd.PersistentFlags().StringVar(&secretsFile, "secrets-file", "", "secrets file (default is $CWD/secrets.yaml)")

	rootCmd.Flags().StringP("server-host", "H", viper.GetString("server_host"), "Address to run Application server on")
	viper.BindPFlag("server_host", rootCmd.Flags().Lookup("server-host"))

	rootCmd.Flags().Int("server-port", viper.GetInt("server_port"), "Port to run Application server on")
	viper.BindPFlag("server_port", rootCmd.Flags().Lookup("server-port"))
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

func initSecrets() {
	if secretsFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigFile(configFile)
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

  viper_secrets := viper.New()
  viper_secrets.SetConfigName("secrets")
  viper_secrets.AddConfigPath(".")

  err := viper_secrets.ReadInConfig()
  if err != nil {
    panic(fmt.Errorf("fatal error config file: %w", err))
  }

  viper.Set("smtp_username", viper_secrets.Get("smtp_username"))
  viper.Set("smtp_password", viper_secrets.Get("smtp_password"))
}
