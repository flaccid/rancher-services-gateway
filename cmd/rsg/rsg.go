package main

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	gw "github.com/flaccid/rancher-services-gateway/discover"
	ui "github.com/flaccid/rancher-services-gateway/ui"
	"github.com/urfave/cli"
)

var (
	VERSION = "v0.0.0-dev"
)

func beforeApp(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "rsg"
	app.Version = VERSION
	app.Usage = "rsg"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "rancher-url",
			Value:  "http://localhost:8080/",
			Usage:  "full url of the rancher server",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "rancher-access-key",
			Usage:  "rancher access Key",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "rancher-secret-key",
			Usage:  "rancher secret Key",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:  "lb-id",
			Usage: "public rancher load balancer id",
		},
		cli.StringFlag{
			Name:   "router-service-label",
			Value:  "services_router",
			Usage:  "label used to identify the load balancer serviced used for routing",
			EnvVar: "ROUTER_SERVICE_TAG",
		},
		cli.IntFlag{
			Name:   "poll-interval,t",
			Usage:  "polling interval in seconds",
			EnvVar: "POLL_INTERVAL",
			Value:  0,
		},
		cli.BoolFlag{
			Name:  "dry",
			Usage: "run in dry mode",
		},
		cli.BoolFlag{
			Name:  "ui,u",
			Usage: "run the basic ui",
		},
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "run in debug mode",
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	log.Info("starting up")

	if c.Bool("ui") {
		ui.Run(c.String("rancher-url"), c.String("rancher-access-key"), c.String("rancher-secret-key"))
	} else {
		if c.Int("poll-interval") > 0 {
			for {
				gw.Discover(c.String("rancher-url"), c.String("rancher-access-key"), c.String("rancher-secret-key"), "", c.Bool("dry"))
				log.Debug("sleeping ", c.Int("poll-interval"), " second(s)")
				time.Sleep(time.Duration(c.Int("poll-interval")) * time.Second)
			}
		} else {
			gw.Discover(c.String("rancher-url"), c.String("rancher-access-key"), c.String("rancher-secret-key"), "", c.Bool("dry"))
		}
	}

	return nil
}
