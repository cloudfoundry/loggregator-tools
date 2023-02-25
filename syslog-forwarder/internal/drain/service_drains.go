package drain

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"code.cloudfoundry.org/loggregator-tools/syslog-forwarder/internal/cloudcontroller"
)

type ServiceDrainLister struct {
	c                 cloudcontroller.Curler
	appNameBatchLimit int
}

func NewServiceDrainLister(c cloudcontroller.Curler, opts ...ServiceDrainListerOption) *ServiceDrainLister {
	dl := &ServiceDrainLister{
		c:                 c,
		appNameBatchLimit: 100,
	}

	for _, o := range opts {
		o(dl)
	}

	return dl
}

type ServiceDrainListerOption func(l *ServiceDrainLister)

func WithServiceDrainAppBatchLimit(limit int) ServiceDrainListerOption {
	return func(l *ServiceDrainLister) {
		l.appNameBatchLimit = limit
	}
}

type Drain struct {
	Name     string
	Guid     string
	Apps     []string
	AppGuids []string
	Type     string
	DrainURL string
	UseAgent bool
}

func (l *ServiceDrainLister) Drains(spaceGuid string) ([]Drain, error) {
	instances, err := l.fetchServiceInstances(
		fmt.Sprintf("/v2/user_provided_service_instances?q=space_guid:%s", spaceGuid),
	)
	if err != nil {
		return nil, err
	}

	var appGuids []string
	var drains []Drain
	for _, s := range instances {
		if s.Entity.SyslogDrainURL == "" {
			continue
		}

		apps, err := l.fetchApps(s.Entity.ServiceBindingsURL)
		if err != nil {
			return nil, err
		}
		appGuids = append(appGuids, apps...)

		drains = append(drains, s.drain(apps))
	}

	appNames, err := l.fetchBatchAppNames(appGuids)
	if err != nil {
		return nil, err
	}

	var namedDrains []Drain
	for _, d := range drains {
		var names, guids []string
		for _, guid := range d.Apps {
			names = append(names, appNames[guid])
			guids = append(guids, guid)
		}
		d.Apps = uniqueStringSlice(names)
		d.AppGuids = uniqueStringSlice(guids)
		namedDrains = append(namedDrains, d)
	}

	return namedDrains, nil
}

func (l *ServiceDrainLister) fetchServiceInstances(url string) ([]userProvidedServiceInstance, error) {
	instances := []userProvidedServiceInstance{}
	for url != "" {
		resp, err := l.c.Curl(url, "GET", "")
		if err != nil {
			return nil, err
		}

		var services userProvidedServiceInstancesResponse
		err = json.Unmarshal(resp, &services)
		if err != nil {
			return nil, err
		}

		instances = append(instances, services.Resources...)

		url = services.NextURL
	}
	return instances, nil
}

func (l *ServiceDrainLister) fetchApps(url string) ([]string, error) {
	var apps []string
	for url != "" {
		resp, err := l.c.Curl(url, "GET", "")
		if err != nil {
			return nil, err
		}

		var serviceBindingsResponse serviceBindingsResponse
		err = json.Unmarshal(resp, &serviceBindingsResponse)
		if err != nil {
			return nil, err
		}

		for _, r := range serviceBindingsResponse.Resources {
			apps = append(apps, r.Entity.AppGuid)
		}

		url = serviceBindingsResponse.NextURL
	}

	return apps, nil
}

func (l *ServiceDrainLister) fetchBatchAppNames(guids []string) (map[string]string, error) {
	guids = uniqueStringSlice(guids)

	allAppNames := make(map[string]string)
	for i := 0; i < len(guids); i += l.appNameBatchLimit {
		end := i + l.appNameBatchLimit

		if end > len(guids) {
			end = len(guids)
		}

		appNames, err := l.fetchAppNames(guids[i:end])
		if err != nil {
			return nil, err
		}

		for k, v := range appNames {
			allAppNames[k] = v
		}
	}

	return allAppNames, nil
}

func (l *ServiceDrainLister) fetchAppNames(guids []string) (map[string]string, error) {
	if len(guids) == 0 {
		return nil, nil
	}

	params := url.Values{
		"guids": {strings.Join(guids, ",")},
	}

	url := "/v3/apps?" + params.Encode()
	apps := make(map[string]string)
	for url != "" {
		resp, err := l.c.Curl(url, "GET", "")
		if err != nil {
			return nil, err
		}

		var appsResp appsResponse
		err = json.Unmarshal(resp, &appsResp)
		if err != nil {
			return nil, err
		}

		for _, a := range appsResp.Apps {
			apps[a.Guid] = a.Name
		}
		url = appsResp.Pagination.Next
	}

	return apps, nil
}

type userProvidedServiceInstancesResponse struct {
	NextURL   string                        `json:"next_url"`
	Resources []userProvidedServiceInstance `json:"resources"`
}

type userProvidedServiceInstance struct {
	MetaData struct {
		Guid string `json:"guid"`
	} `json:"metadata"`
	Entity struct {
		Name               string `json:"name"`
		ServiceBindingsURL string `json:"service_bindings_url"`
		SyslogDrainURL     string `json:"syslog_drain_url"`
	} `json:"entity"`
}

func (s userProvidedServiceInstance) drainType() string {
	uri, err := url.Parse(s.Entity.SyslogDrainURL)
	if err != nil {
		return ""
	}

	drainTypes := uri.Query()["drain-type"]
	if len(drainTypes) == 0 {
		return "logs"
	}
	return drainTypes[0]
}

func (s userProvidedServiceInstance) useAgent() bool {
	instanceURL, err := url.Parse(s.Entity.SyslogDrainURL)
	if err != nil {
		return false
	}

	return strings.HasSuffix(instanceURL.Scheme, "-v3")
}

func (s userProvidedServiceInstance) drain(apps []string) Drain {
	return Drain{
		Name:     s.Entity.Name,
		Guid:     s.MetaData.Guid,
		Apps:     apps,
		Type:     s.drainType(),
		DrainURL: s.Entity.SyslogDrainURL,
		UseAgent: s.useAgent(),
	}
}

type serviceBindingsResponse struct {
	NextURL   string           `json:"next_url"`
	Resources []serviceBinding `json:"resources"`
}

type serviceBinding struct {
	Entity struct {
		AppGuid string `json:"app_guid"`
		AppUrl  string `json:"app_url"`
	} `json:"entity"`
}

type appsResponse struct {
	Apps       []appData `json:"resources"`
	Pagination struct {
		Next string `json:"next"`
	} `json:"pagination"`
}

type appData struct {
	Name string `json:"name"`
	Guid string `json:"guid"`
}

func uniqueStringSlice(str []string) []string {
	var results []string
	for _, s := range str {
		results = appendIfMissing(results, s)
	}

	return results
}

func appendIfMissing(data []string, element string) []string {
	for _, elem := range data {
		if elem == element {
			return data
		}
	}

	return append(data, element)
}
