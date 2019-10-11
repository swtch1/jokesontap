package main

import (
	"fmt"
	"github.com/swtch1/jokesontap"
	"github.com/swtch1/jokesontap/cli"
	"os"
)

var (
	// buildVersion is the application version and should be populated at build time by build ldflags
	// this default message should be overwritten
	buildVersion string = "unset: please file an issue"
)

func init() {
	cli.Init(buildVersion)
	jokesontap.SetLogger(os.Stderr, cli.LogLevel, cli.LogFormat, cli.PrettyPrintJsonLogs)
}

func main() {
	fmt.Println("go")
}
