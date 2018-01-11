package discover

import (
	"fmt"

  "reflect"
	log "github.com/Sirupsen/logrus"
	r "github.com/flaccid/rancher-services-gateway/rancher"
	rancher "github.com/rancher/go-rancher/v2"
)

func Discover(rancherUrl string, rancherAccessKey string, rancherSecretKey string, lbId string) {
	rancherClient := r.CreateClient(rancherUrl, rancherAccessKey, rancherSecretKey)

	var servicesRouter *rancher.LoadBalancerService
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
				// this lb is a service router
				if k == "services_router" && v == "true" {
					servicesRouter = s
					log.Debug(reflect.TypeOf(s))
					log.Debug(reflect.TypeOf(servicesRouter))
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

	// by now we should have the service router resource!
	//if servicesRouter ? {
	//	log.Errorf("no services router! found")
	//	os.Exit(2)
	//}
	log.Info("using services router, ", servicesRouter.Name)

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

	log.Info("servicesRouter LbConfig", servicesRouter.LbConfig)
	log.Info("servicesRouter LinkedServices", servicesRouter.LinkedServices)

	// create port rule client
	// this is just a dummy rule atm
	portRule := rancher.PortRule{
		Protocol: "https",
		Hostname: "flaccid-is-learning",
		Path: "",
		Priority: 10,
		SourcePort: 443,
		TargetPort: 80,
		ServiceId: "1s1128",
	}
	log.Info("new port rule ", portRule)

	// The strategy allows you to upgrade
 	strategy := &rancher.ServiceUpgrade{
		InServiceStrategy: &rancher.InServiceUpgradeStrategy{
 		  BatchSize: 1,
 		  StartFirst: true,
 		  LaunchConfig: servicesRouter.LaunchConfig,
 	  },
  }

	// Here we can see the new rule added to the current list of rules
	log.Info("updating the services router")
	servicesRouter.LbConfig.PortRules = append(servicesRouter.LbConfig.PortRules, portRule)
	log.Infof("New service Rules: %+v", servicesRouter.LbConfig.PortRules)
	update, err := rancherClient.LoadBalancerService.ActionUpgrade(servicesRouter, strategy)
	if err != nil {
		panic(err)
	} else {
		// Notice here it isn't present in the returned object
		log.Infof("update complete: +%v", update.LbConfig.PortRules)
	}
}
