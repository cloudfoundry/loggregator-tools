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
	c         Curler
	spaceGuid string
	apiURL    string
}

func Space(c Curler, apiURL, spaceGuid string) *SpaceProvider {
	return &SpaceProvider{
		c:         c,
		spaceGuid: spaceGuid,
		apiURL:    apiURL,
	}
}

func (s *SpaceProvider) SourceIDs() []string {
	resp := s.c.Get(fmt.Sprintf("%s/v3/apps?space_guids=%s", s.apiURL, s.spaceGuid))

	var cApps response
	err := json.Unmarshal(resp, &cApps)
	if err != nil {
		log.Printf("error getting app info from CAPI: %s", err)
		return nil
	}

	var apps []string
	for _, app := range cApps.Resources {
		apps = append(apps, app.Guid)
	}
	return apps
}

type response struct {
	Resources []struct {
		Guid string `json:"guid"`
	} `json:"resources"`
}
