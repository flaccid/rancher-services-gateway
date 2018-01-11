package main

import (
	"os"

	"github.com/urfave/cli"
	log "github.com/Sirupsen/logrus"
	gw "github.com/flaccid/rancher-services-gateway/discover"
	ui "github.com/flaccid/rancher-services-gateway/ui"
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
			Name:  "rancher-url",
			Value: "http://localhost:8080/",
			Usage: "full url of the rancher server",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:  "rancher-access-key",
			Usage: "rancher access Key",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:  "rancher-secret-key",
			Usage: "rancher secret Key",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:  "lb-id",
			Usage: "public rancher load balancer id",
		},
		cli.StringFlag{
			Name:  "router-service-label",
			Value: "services_router",
			Usage: "label used to identify the load balancer serviced used for routing",
			EnvVar: "ROUTER_SERVICE_TAG",
		},
		cli.BoolFlag{
			Name: "ui,u",
			Usage: "run the basic ui",
		},
		cli.BoolFlag{
			Name: "debug,d",
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
		gw.Discover(c.String("rancher-url"), c.String("rancher-access-key"), c.String("rancher-secret-key"), "")
	}

  return nil
}
