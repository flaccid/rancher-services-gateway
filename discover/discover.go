package discover

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/Jeffail/gabs"

	log "github.com/Sirupsen/logrus"
	r "github.com/flaccid/rancher-services-gateway/rancher"
	rancher "github.com/rancher/go-rancher/v2"
)

// NOTE:
// not using go-rancher for parts, because of fail :(
// https://forums.rancher.com/t/using-go-rancher-to-update-lb-service-rules/8041

func Discover(rancherUrl string, rancherAccessKey string, rancherSecretKey string, lbId string, dry bool) {
	rancherClient := r.CreateClient(rancherUrl, rancherAccessKey, rancherSecretKey)
	loadBalancerServices := r.ListRancherLoadBalancerServices(rancherClient)

	var servicesGateway *rancher.LoadBalancerService
	var targetLBs []*rancher.LoadBalancerService

	// get the lb
	if len(lbId) > 0 {
		lB, err := rancherClient.LoadBalancerService.ById(lbId)
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("lb ", lB)
	} else {
		for _, s := range loadBalancerServices {
			log.Debug("processing ", s.Uuid)
			log.Debug("data ", s)
			log.Debug("labels: ", s.LaunchConfig.Labels)

			for k, v := range s.LaunchConfig.Labels {
				// this lb is a services gateway
				if k == "services_gateway" && v == "true" {
					servicesGateway = s
					log.Debug(reflect.TypeOf(s))
					log.Debug(reflect.TypeOf(servicesGateway))
					log.Debug(s)
					break
				}

				// this lb is a service target
				if k == "dns_target" {
					log.Debug("found a target service ", s)
					targetLBs = append(targetLBs, s)
				}
			}
		}
	}

	// by now we should have the services gateway resource
	if servicesGateway == nil {
		log.Fatalf("no services gateway found")
	}
	log.WithFields(log.Fields{
		"name":         servicesGateway.Name,
		"link":         servicesGateway.Links["self"],
		"target_count": len(targetLBs),
	}).Info("services gateway")
	log.Debug("target LBs", targetLBs)

	// get the gw data first
	client := &http.Client{}
	req, err := http.NewRequest("GET", servicesGateway.Links["self"], nil)
	req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal(err2)
	}
	gwData, err := gabs.ParseJSON(bodyBytes)

	var portRules []rancher.PortRule

	// get the default-website service, we assume its the first returned atm
	defaultWebsite := r.GetRancherServiceByName(rancherClient, "default-website")[0]

	// add default-website portRule (hard-coded for now)
	portRule := rancher.PortRule{
		Resource: rancher.Resource{Type: "portRule"},
		Protocol: "https",
		Hostname: "",
		Path:     "",
		Priority: 1,
		// only https frontend supported currently
		SourcePort: 443,
		// target port is the source port of the downstream lb
		TargetPort: 8080,
		ServiceId:  defaultWebsite.Id,
	}
	portRules = append(portRules, portRule)

	// for each target
	for i, t := range targetLBs {
		var dnsAlias string
		var dnsTarget string

		log.Debug(i, t)

		// get the dns alias and target
		for k, v := range t.LaunchConfig.Labels {
			log.Debug(k)
			if k == "dns_alias" {
				log.Debug(v)
				log.Debug(reflect.TypeOf(v))
				dnsAlias = fmt.Sprint(v)
			}
			if k == "dns_target" {
				log.Debug(v)
				log.Debug(reflect.TypeOf(v))
				dnsTarget = fmt.Sprint(v)
			}
		}

		// create port rule
		portRule := rancher.PortRule{
			Resource: rancher.Resource{Type: "portRule"},
			Protocol: "https",
			Hostname: dnsAlias,
			Path:     "",
			Priority: int64(i + 1),
			// only https frontend supported currently
			SourcePort: 443,
			// target port is the source port of the downstream lb
			TargetPort: t.LbConfig.PortRules[0].SourcePort,
			ServiceId:  t.Id,
		}

		log.Debug("portRule: ", portRule)

		log.WithFields(log.Fields{
			"target_service": portRule.ServiceId,
			"dns_alias":      portRule.Hostname,
			"dns_target":     dnsTarget,
			"target_port":    portRule.TargetPort,
		}).Info("target ", portRule.Priority)

		portRules = append(portRules, portRule)
	}

	log.Debug("portRules", portRules)
	gabs.Consume(gwData)
	gwData.SetP(portRules, "lbConfig.portRules")

	// update the lb, with port rule
	if !dry {
		log.Info("update load balancer")
		req, err := http.NewRequest("PUT", servicesGateway.Links["self"], strings.NewReader(gwData.String()))
		req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
		req.Header.Set("Content-Type", "application/json")
		req.ContentLength = int64(len(gwData.String()))
		response, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Debug("the calculated length is:", len(string(contents)), "for the url:")
			log.Debug("   ", response.StatusCode)
			hdr := response.Header
			for key, value := range hdr {
				log.Debug("   ", key, ":", value)
			}
			log.Debug(string(contents))
		}
	} else {
		log.Info("dry mode, skipping update")
	}
}
