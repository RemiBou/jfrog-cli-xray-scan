package commands

import (
	"encoding/json"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	utils2 "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"strings"
	"time"
)

type xrayClient struct {
	client         *httpclient.ArtifactoryHttpClient
	serviceDetails auth.ServiceDetails
	xrayUrl        string
}

type xrayClientInterface interface {
	scan(comps []component) (*ComponentSummaryResult, error)
}

// TODO: probably hard to itest like this
func newXrayClient() (xrayClientInterface, error) {
	confArti, err := config.GetDefaultArtifactoryConf()
	if err != nil {
		return nil, err
	}
	servicesManager, err := utils.CreateServiceManager(confArti, false)
	if err != nil {
		return nil, err
	}
	serviceDetails := servicesManager.GetConfig().GetServiceDetails()
	xrayUrl := strings.Replace(serviceDetails.GetUrl(), "/artifactory", "/xray", 1)

	return &xrayClient{client: servicesManager.Client(), xrayUrl: xrayUrl, serviceDetails: serviceDetails}, nil
}

func (x *xrayClient) scan(comps []component) (*ComponentSummaryResult, error) {
	var components []ComponentSummaryRequestDetails
	for _, comp := range comps {
		components = append(components, ComponentSummaryRequestDetails{
			ComponentID: comp.toString(),
		})
	}
	request := ComponentSummaryRequest{
		ComponentDetails: components,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	details := x.serviceDetails.CreateHttpClientDetails()
	requestFullUrl, err := utils2.BuildArtifactoryUrl(x.xrayUrl, "/api/v1/summary/component", make(map[string]string))
	utils2.SetContentType("application/json", &details.Headers)
	_, body, err := x.client.SendPost(requestFullUrl, data, &details)
	if err != nil {
		return nil, err
	}
	result := &ComponentSummaryResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type ComponentSummaryRequest struct {
	ComponentDetails []ComponentSummaryRequestDetails `json:"component_details"`
}
type ComponentSummaryRequestDetails struct {
	ComponentID string `json:"component_id"`
}
type ComponentSummaryResult struct {
	Artifacts []struct {
		General struct {
			ComponentID string `json:"component_id"`
			Name        string `json:"name"`
			PkgType     string `json:"pkg_type"`
		} `json:"general"`
		Issues []struct {
			Components []struct {
				ComponentID   string   `json:"component_id"`
				FixedVersions []string `json:"fixed_versions"`
			} `json:"components"`
			Created time.Time `json:"created"`
			Cves    []struct {
				CvssV2 string `json:"cvss_v2"`
			} `json:"cves"`
			Description string `json:"description"`
			IssueType   string `json:"issue_type"`
			Provider    string `json:"provider"`
			Severity    string `json:"severity"`
			Summary     string `json:"summary"`
		} `json:"issues"`
		Licenses []struct {
			Components  []string `json:"components"`
			FullName    string   `json:"full_name"`
			MoreInfoURL []string `json:"more_info_url"`
			Name        string   `json:"name"`
		} `json:"licenses"`
	} `json:"artifacts"`
}
