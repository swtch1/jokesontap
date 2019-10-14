package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var appName = "jokesontap"

var (
	// cmd is the root of our CLI
	cmd = &cobra.Command{
		Use:   appName,
		Short: appName,
		Long:  fmt.Sprintf("%s is a user specific joke retrieval server built with a concentration on resiliency and high throughput.", appName),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	Help                bool
	Version             bool
	Port                int32
	LogLevel            string
	LogFormat           string
	PrettyPrintJsonLogs bool
)

// Init performs setup for the application CLI commands and flags, setting application version as provided.
func Init(version string) {
	cmd.PersistentFlags().BoolVarP(&Help, "help", "h", false, "Display this help and exit.")
	cmd.PersistentFlags().BoolVar(&Version, "version", false, "Print the application version and exit.")
	cmd.PersistentFlags().Int32VarP(&Port, "port", "p", 5000, "Port which the server will listen on.")
	cmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "info", "Log level should be one of trace, debug, info, warn, error, fatal.")
	cmd.PersistentFlags().StringVar(&LogFormat, "log-format", "text", "Log format should be one of text, json.")
	cmd.PersistentFlags().BoolVar(&PrettyPrintJsonLogs, "pretty-json", false, "If writing JSON logs, pretty print those logs.")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if Help {
		os.Exit(0)
	}

	// handle the version manually since the built in version options for Cobra do not exit after printing
	if Version {
		fmt.Printf("%s version %s\n", appName, version)
		os.Exit(0)
	}
}
