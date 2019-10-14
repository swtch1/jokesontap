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
	defaultNameChanSize int64 = 100000
)

func main() {
	cli.Init(buildVersion)
	jokesontap.InitLogger(os.Stderr, cli.LogLevel, cli.LogFormat, cli.PrettyPrintJsonLogs)

	namesUrl, err := url.Parse(defaultNamesUrl)
	if err != nil {
		log.WithError(err).Fatal("unable to parse default names URL, please submit an issue")
	}
	nameClient := jokesontap.NewNameClient(*namesUrl)
	namesChan := make(chan jokesontap.Name, defaultNameChanSize)

	// NOTE: the size of the budget array has been shortened to 6 rather than the API specified 7 requests per minute as
	// real world testing showed that rate limit errors were still being seen at 7 requests per every 65 seconds.
	// TODO: re-evaluate the names API at regular intervals to determine the optimal request rate
	budgetReq := jokesontap.BudgetNameReq{
		// allow buffer to avoid getting rate limited from names API
		// ref: http://www.icndb.com/api/
		MinDiff:    time.Second * 61,
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
