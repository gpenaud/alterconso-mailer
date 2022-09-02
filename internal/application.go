package internal

import (
  context       "context"
  fmt           "fmt"
  http          "net/http"
  log           "github.com/sirupsen/logrus"
  mux           "github.com/gorilla/mux"
  time          "time"
)

// -------------------------------------------------------------------------- //
// 2. Application Declarative Configuration
// -------------------------------------------------------------------------- //

type Configuration struct {
  ServerHost   string
  ServerPort   string
  SmtpHost     string
  SmtpPort     string
  SmtpUsername string
  SmtpPassword string
}

type Version struct {
  BuildTime string
  Commit    string
  Release   string
}

type Application struct {
  Config  *Configuration
  Router  *mux.Router
  Version *Version
}

// -------------------------------------------------------------------------- //
// 4. Router setup
// -------------------------------------------------------------------------- //

func (a *Application) initializeRoutes() {
  // application user-related routes
  a.Router.HandleFunc("/send", a.emailController).Methods("GET")
}

// -------------------------------------------------------------------------- //
// 5. Application Setup
// -------------------------------------------------------------------------- //

func (a *Application) Initialize() {
  a.Router = mux.NewRouter()
  a.initializeRoutes()
  log.Info("application is initialized")
}

func (a *Application) Run(ctx context.Context) {
  server_address :=
    fmt.Sprintf("%s:%s", a.Config.ServerHost, a.Config.ServerPort)

  server_message :=
    fmt.Sprintf(
  `

START INFOS
-----------
Listening alterconso-mailer on %s:%s...

BUILD INFOS
-----------
time: %s
release: %s
commit: %s

`,
      a.Config.ServerHost,
      a.Config.ServerPort,
      a.Version.BuildTime,
      a.Version.Release,
      a.Version.Commit,
    )

  // ---------------------------------------------------------------------------
  // 5.2. manage healthchecks and healthchecks server
  // ---------------------------------------------------------------------------

  httpServer := &http.Server{
		Addr:    server_address,
		Handler: a.Router,
	}

  go func() {
    // we keep this log on standard format
    log.Info(server_message)
    httpServer.ListenAndServe()
  }()

  // ---------------------------------------------------------------------------
  // 5.4. manage server shutdown
  // ---------------------------------------------------------------------------

  <-ctx.Done()
  log.Info("application server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

  var err error

	if err = httpServer.Shutdown(ctxShutdown); err != nil {
    log.Fatal("application server shutdown failed")
	}

  log.Info("application server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
