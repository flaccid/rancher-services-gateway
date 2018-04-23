package rancher

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	rancher "github.com/rancher/go-rancher/v2"
)

type Client struct {
	client *rancher.RancherClient
}

var withoutPagination *rancher.ListOpts

func CreateClient(url, accessKey, secretKey string) *rancher.RancherClient {
	client, err := rancher.NewRancherClient(&rancher.ClientOpts{
		Url:       url,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Timeout:   time.Second * 5,
	})

	if err != nil {
		log.Errorf("Failed to create a client for rancher, error: %s", err)
		os.Exit(1)
	}

	return client
}

// not used
func GetServicesRouter(client *rancher.RancherClient) *rancher.LoadBalancerService {
	var servicesRouter *rancher.LoadBalancerService

	loadBalancers := ListRancherLoadBalancerServices(client)

	servicesRouter = loadBalancers[0]

	return servicesRouter
}

func ListRancherLoadBalancerServices(client *rancher.RancherClient) []*rancher.LoadBalancerService {
	var servicesList []*rancher.LoadBalancerService

	services, err := client.LoadBalancerService.List(withoutPagination)

	if err != nil {
		log.Errorf("cannot get services: %+v", err)
	}

	for k := range services.Data {
		servicesList = append(servicesList, &services.Data[k])
	}

	return servicesList
}

func GetRancherServiceByName(client *rancher.RancherClient, name string) []*rancher.Service {
	var servicesList []*rancher.Service

	opts := &rancher.ListOpts{
		Filters: map[string]interface{}{
			"name": name,
		},
	}

	services, err := client.Service.List(opts)

	if err != nil {
		log.Errorf("cannot get services: %+v", err)
	}

	for k := range services.Data {
		servicesList = append(servicesList, &services.Data[k])
	}

	return servicesList
}
