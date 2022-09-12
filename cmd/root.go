package cmd

import (
	"github.com/spf13/cobra"
)

type Version struct {
  BuildTime string
  Commit    string
  Release   string
}

var (
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
