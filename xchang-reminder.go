package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/MikeAlbertFleetSolutions/xchango"
	"github.com/gen2brain/beeep"
	"gopkg.in/yaml.v2"
)

// configuration holds the common application configuration
type configuration struct {
	Domain       string
	Username     string
	Password     string
	MaxFetchSize int
	ExchangeURL  string
	Reminder     int
	Icon         string
}

var (
	// config holds base application configuration
	config configuration

	// connection to exchange
	xchang *xchango.ExchangeClient
)

func init() {
	// show file & location, date & time
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// command line app
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage of %s\n", os.Args[0])
		flag.PrintDefaults()
	}

	// get config file from command line
	configFile := flag.String("config", "", "System configuration file")
	flag.Parse()

	if len(*configFile) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// read config
	// #nosec G304
	bytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// create exchange client
	xchangconfig := xchango.ExchangeConfig{
		ExchangeUser: xchango.ExchangeUser{
			Domain:   config.Domain,
			Username: config.Username,
			Password: config.Password,
		},
		MaxFetchSize: config.MaxFetchSize,
		ExchangeURL:  config.ExchangeURL,
	}

	xchang, err = xchango.NewExchangeClient(xchangconfig)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func main() {
	cal, err := xchang.GetCalendar()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	appointments, err := xchang.GetAppointments(cal)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// determine icon to use for notify message
	icon := config.Icon
	if _, err = os.Stat(config.Icon); os.IsNotExist(err) {
		icon = ""
	}

	// notify of events within the reminder period
	start := time.Now()
	end := start.Add(time.Minute * time.Duration(config.Reminder))
	for _, appointment := range appointments {
		if inTimeSpan(start, end, appointment.Start.In(time.Local)) {
			err := beeep.Alert("Calendar Event", fmt.Sprintf("<b>%s</b>", appointment.Subject), icon)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		}
	}
}
