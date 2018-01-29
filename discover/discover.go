package discover

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	r "github.com/flaccid/rancher-services-gateway/rancher"
	rancher "github.com/rancher/go-rancher/v2"
	"reflect"
)

type PortRuleClient struct {
	rancherClient *rancher.RancherClient
}

type PortRules struct {
	Rules []*rancher.PortRule
}

func Discover(rancherUrl string, rancherAccessKey string, rancherSecretKey string, lbId string) {
	rancherClient := r.CreateClient(rancherUrl, rancherAccessKey, rancherSecretKey)

	var servicesGateway *rancher.LoadBalancerService
	var targetLBs []*rancher.LoadBalancerService

	loadBalancerServices := r.ListRancherLoadBalancerServices(rancherClient)

	// get the lb
	if len(lbId) > 0 {
		loadBalancerService, err := rancherClient.LoadBalancerService.ById(lbId)
		if err != nil {
			panic(err)
		}

		log.Debug(loadBalancerService)
		log.Debug(loadBalancerService.LaunchConfig.Labels)
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
				if k == "platform_identifier" {
					log.Debug("found a target service ", s)
					targetLBs = append(targetLBs, s)
				}
			}
		}
	}

	// by now we should have the services gateway resource
	// TODO: error out if not
	// if servicesGateway ? {
	//	log.Errorf("no services gateway found")
	//	os.Exit(2)
	// }
	log.Info("using services gateway, ", servicesGateway.Name)

	log.Info("found ", len(targetLBs), " target service(s)")
	log.Debug("target LBs", targetLBs)

	var requestHost string

	// for each target
	for i, t := range targetLBs {
		var platformId string
		var platformDnsDomain string

		log.Debug(i, t)

		// construct the request host
		for k, v := range t.LaunchConfig.Labels {
			log.Debug(k)
			if k == "platform_identifier" {
				log.Debug(v)
				log.Debug(reflect.TypeOf(v))
				platformId = fmt.Sprint(v)
			}
			if k == "platform_dns_domain" {
				log.Debug(v)
				log.Debug(reflect.TypeOf(v))
				platformDnsDomain = fmt.Sprint(v)
			}
		}
		requestHost = fmt.Sprintf(platformId + "." + platformDnsDomain)

		log.Info("request host: ", requestHost)
	}

	log.Info("servicesGateway LbConfig", servicesGateway.LbConfig)
	log.Info("servicesGateway LinkedServices", servicesGateway.LinkedServices)

	// create port rule
	// this is just a dummy rule atm
	portRule := rancher.PortRule{
		Resource:   rancher.Resource{Type: "portRule"},
		Protocol:   "https",
		Hostname:   "my.host.com",
		Path:       "",
		Priority:   3,
		SourcePort: 443,
		TargetPort: 80,
		ServiceId:  "1s1292",
	}
	log.Info("new port rule ", portRule)

	// The strategy allows you to upgrade
	strategy := &rancher.ServiceUpgrade{
		InServiceStrategy: &rancher.InServiceUpgradeStrategy{
			BatchSize:    1,
			StartFirst:   true,
			LaunchConfig: servicesGateway.LaunchConfig,
		},
	}
	log.Debug("upgrade strategy", strategy)

	log.Debug("current port rules: ", servicesGateway.LbConfig.PortRules)

	// show each
	for p, k := range servicesGateway.LbConfig.PortRules {
		log.Debug(reflect.TypeOf(k))
		log.Debug("port rule - ", p, k)
	}
	// log.Debug(reflect.TypeOf(portRule))
	// show what we are adding
	log.Debug("port rule + ", portRule)

	// append the additional port rule
	servicesGateway.LbConfig.PortRules = append(servicesGateway.LbConfig.PortRules, portRule)
	// override instead
	// servicesGateway.LbConfig.PortRules[1] = portRule

	log.Debug("new port rules: ", servicesGateway.LbConfig.PortRules)

	// carry out the update
	update, err := rancherClient.LoadBalancerService.ActionUpdate(servicesGateway)
	if err != nil {
		panic(err)
	}
	log.Debug(update)
}
