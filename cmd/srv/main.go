package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/swtch1/jokesontap"
	"github.com/swtch1/jokesontap/cli"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// buildVersion is the application version and should be populated at build time by build ldflags
	// this default message should be overwritten
	buildVersion string = "unset: please file an issue"

	// defaultNamesUrl is the default API URL used to get new random names
	defaultNamesUrl string = "https://uinames.com/api/?amount=500"
	// defaultJokesUrl is the default API URL used to get new random jokes
	defaultJokesUrl string = "http://api.icndb.com/jokes/random"
	// defaultNameChanSize is the default size of the channel used to store names
	// so that names can be eagerly retrieved from the API
	defaultNameChanSize int64 = 10000
)

func main() {
	// TODO: only allow GET methods on the server
	// TODO: translate quotes like '&quot;' to real chars.
	// TODO: set appropriate headers in the server response
	// TODO: solve this error:
	// ERRO[0118]/mnt/c/important/code/jokesontap/name.go:136 github.com/swtch1/jokesontap.(*BudgetNameReq).RequestOften() unable to get names from names client
	//        error="unable to unmarshal names API response body: invalid character '<' looking for beginning of value"

	cli.Init(buildVersion)
	jokesontap.InitLogger(os.Stderr, cli.LogLevel, cli.LogFormat, cli.PrettyPrintJsonLogs)

	namesUrl, err := url.Parse(defaultNamesUrl)
	if err != nil {
		log.WithError(err).Fatal("unable to parse default names URL, please submit an issue")
	}
	nameClient := jokesontap.NewNameClient(*namesUrl)
	namesChan := make(chan jokesontap.Name, defaultNameChanSize)

	budgetReq := jokesontap.BudgetNameReq{
		MinDiff:    time.Second * 59,
		NameClient: nameClient,
		NameChan:   namesChan,
	}
	go budgetReq.RequestOften()

	jokesUrl, err := url.Parse(defaultJokesUrl)
	if err != nil {
		log.WithError(err).Fatalf("unable to parse jokes URL '%s', please file an issue", defaultJokesUrl)
	}
	jokeClient := jokesontap.NewJokeClient(*jokesUrl)

	go HandleInterrupt()
	log.Infof("starting server on port %d", cli.Port)
	srv := &jokesontap.Server{
		Port:       cli.Port,
		Names:      namesChan,
		JokeClient: jokeClient,
	}
	log.Fatal(srv.ListenAndServe())
}

// HandleInterrupt will immediately terminate the server if it detects an interrupt signal.
func HandleInterrupt() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("interrupt: stopping server...")
		os.Exit(1)
	}()
}
