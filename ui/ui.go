package ui

import (
	"fmt"
	"html/template"
	"net/http"

	rancher "github.com/flaccid/rancher-services-gateway/rancher"
	log "github.com/Sirupsen/logrus"
)

const (
	listenPort = 8080
)

type Service struct {
	DnsName string
	ServiceName string
	StackId string
	StackUrl string
	State string
	Url string
	Uuid string
}

type PageData struct {
	Title    string
	Services []Service
}

func render(w http.ResponseWriter, r *http.Request, rancherUrl string, rancherAccessKey string, rancherSecretKey string) {
	rancherClient := rancher.CreateClient(rancherUrl, rancherAccessKey, rancherSecretKey)
	loadBalancerServices := rancher.ListRancherLoadBalancerServices(rancherClient)

	var services []Service

	for _, s := range loadBalancerServices {
		for k, v := range s.LaunchConfig.Labels {
			if k == "dns_name" {
				log.Debug(s)
				service := Service{
										DnsName: fmt.Sprintf("%v", v),
										ServiceName: s.Name,
										StackId: s.StackId,
										StackUrl: fmt.Sprintf("%v/env/foo/apps/stacks/bar", rancherUrl),
										State: s.State,
										Url: fmt.Sprintf("https://%v/", v),
										Uuid: s.Uuid,
									}
				services = append(services, service)
			}
		}
	}

	// expects template in cwd
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		panic(err)
	}
	pageData := PageData{
		Title: "Rancher Services Gateway",
		Services: services,
	}
	log.Debug(pageData)
	tmpl.Execute(w, pageData)
}

func Run(rancherUrl string, rancherAccessKey string, rancherSecretKey string) {
	log.Info("ui listening for requests on :" + fmt.Sprintf("%v", listenPort))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, rancherUrl, rancherAccessKey, rancherSecretKey)
	})

	http.ListenAndServe(":" + fmt.Sprintf("%v", listenPort), nil)
}
