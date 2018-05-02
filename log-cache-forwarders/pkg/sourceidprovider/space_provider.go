package sourceidprovider

import (
	"encoding/json"
	"fmt"
	"log"
)

type Curler interface {
	Get(url string) []byte
}

type SpaceProvider struct {
	c               Curler
	includeServices bool
	spaceGuid       string
	apiURL          string
}

type SpaceOption func(*SpaceProvider)

func WithSpaceServiceInstances() SpaceOption {
	return func(s *SpaceProvider) {
		s.includeServices = true
	}
}

func Space(c Curler, apiURL, spaceGuid string, opts ...SpaceOption) *SpaceProvider {
	sp := &SpaceProvider{
		c:         c,
		spaceGuid: spaceGuid,
		apiURL:    apiURL,
	}

	for _, o := range opts {
		o(sp)
	}

	return sp
}

func (s *SpaceProvider) SourceIDs() []string {
	sourceIDs := s.guidsFor("apps")

	if s.includeServices {
		sourceIDs = append(sourceIDs, s.guidsFor("service_instances")...)
	}

	return sourceIDs
}

func (s *SpaceProvider) guidsFor(resource string) []string {
	resp := s.c.Get(fmt.Sprintf("%s/v3/%s?space_guids=%s", s.apiURL, resource, s.spaceGuid))

	var capiResources response
	err := json.Unmarshal(resp, &capiResources)
	if err != nil {
		log.Printf("error getting app info from CAPI: %s", err)
		return nil
	}

	var guids []string
	for _, resource := range capiResources.Resources {
		guids = append(guids, resource.Guid)
	}
	return guids
}

type response struct {
	Resources []struct {
		Guid string `json:"guid"`
	} `json:"resources"`
}
