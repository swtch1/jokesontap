package main

import (
	"fmt"
	"github.com/swtch1/jokesontap"
	"github.com/swtch1/jokesontap/cli"
	"os"
	"os/signal"
	"syscall"
)

var (
	// buildVersion is the application version and should be populated at build time by build ldflags
	// this default message should be overwritten
	buildVersion string = "unset: please file an issue"

	// defaultNamesUrl is the default  // FIXME: update the docs
	defaultNamesUrl     string = "https://uinames.com/api/?amount=500"
	defaultBaseJokesUrl string = "http://api.icndb.com/jokes/random"
)

func init() {
	cli.Init(buildVersion)
	jokesontap.SetLogger(os.Stderr, cli.LogLevel, cli.LogFormat, cli.PrettyPrintJsonLogs)
}

func main() {
	//namesUrl, err := url.Parse(defaultNamesUrl)
	//if err != nil {
	//	log.WithError(err).Fatal("unable to parse default names URL, please submit an issue")
	//}
	//jokeUrl, err := url.Parse(defaultBaseJokesUrl)
	//if err != nil {
	//	log.WithError(err).Fatal("unable to parse default jokes URL, please submit an issue")
	//}
	//

	//jc := jokesontap.NewJokeClient(*jokeUrl)
	//joke, err := jc.JokeWithCustomName("john", "smit")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(joke)
	//nc := jokesontap.NewNameClient(defaultNamesUrl)

	//names, err := nc.Names()
	//if err != nil {
	//	panic(err)
	//}
	//for _, name := range names {
	//	fmt.Println(name.Name)
	//}

	//go HandleInterrupt()
	//srv := jokesontap.Server{Port: cli.Port}
	//log.Infof("starting server on port %d", cli.Port)
	//log.Fatal(srv.ListenAndServe())
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
