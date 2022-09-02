package main

import (
  cli      "github.com/urfave/cli/v2"
  context  "context"
  internal "github.com/gpenaud/alterconso-mailer/internal"
  log      "github.com/sirupsen/logrus"
  os       "os"
  signal   "os/signal"
  syscall  "syscall"
)

// -------------------------------------------------------------------------- //
// 1. Application Initialization
// -------------------------------------------------------------------------- //

var a internal.Application

func init() {
  registerConfiguration(&a)
  registerVersion(&a)
  a.Initialize()
}

// -------------------------------------------------------------------------- //
// 2. Application Configuration
// -------------------------------------------------------------------------- //

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str { return true }
	}

	return false
}

func registerConfiguration(a *internal.Application) {
  a.Config = &internal.Configuration{}

  app := &cli.App{
    Action: func(c *cli.Context) error {
      return nil
    },
    Flags: []cli.Flag{
      &cli.StringFlag{Name: "server_host", Value: "", Usage: "API server host `HOST`", Destination: &a.Config.ServerHost, EnvVars: []string{"ALTERCONSO_MAILER_SERVER_HOST"}},
      &cli.StringFlag{Name: "server_port", Value: "8010", Usage: "API server port `PORT`", Destination: &a.Config.ServerPort, EnvVars: []string{"ALTERCONSO_MAILER_SERVER_PORT"}},
      &cli.StringFlag{Name: "smtp_host", Value: "8010", Usage: "API server port `PORT`", Destination: &a.Config.SmtpHost, EnvVars: []string{"ALTERCONSO_MAILER_SMTP_HOST"}},
      &cli.StringFlag{Name: "smtp_port", Value: "8010", Usage: "API server port `PORT`", Destination: &a.Config.SmtpPort, EnvVars: []string{"ALTERCONSO_MAILER_SMTP_PORT"}},
      &cli.StringFlag{Name: "smtp_username", Value: "8010", Usage: "API server port `PORT`", Destination: &a.Config.SmtpUsername, EnvVars: []string{"ALTERCONSO_MAILER_SMTP_USERNAME"}},
      &cli.StringFlag{Name: "smtp_password", Value: "8010", Usage: "API server port `PORT`", Destination: &a.Config.SmtpPassword, EnvVars: []string{"ALTERCONSO_MAILER_SMTP_PASSWORD"}},
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}

// -------------------------------------------------------------------------- //
// 3. Application Version
// -------------------------------------------------------------------------- //

var BuildTime = "unset"
var Commit 		= "unset"
var Release 	= "unset"

func registerVersion(a *internal.Application) {
  a.Version = &internal.Version{BuildTime, Commit, Release}
}

// -------------------------------------------------------------------------- //
// 4. Main function
// -------------------------------------------------------------------------- //

func main() {
  c := make(chan os.Signal, 1) // creation of a channel of type os.Signal
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM) // add 2 signals to the channel
	ctx, cancel := context.WithCancel(context.Background())

  go func() {
		<-c
    log.Warn("received a system call")
		cancel() // linked with ctx, cancel
	}()

  a.Run(ctx)
}
