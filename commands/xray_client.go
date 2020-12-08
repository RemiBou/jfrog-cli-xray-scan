package commands

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	utils2 "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"time"
)

type xrayClient struct {
	client         *httpclient.ArtifactoryHttpClient
	cliHttpDetails httputils.HttpClientDetails
	xrayUrl        string
}

func newXrayClient(xrayUrl string, client *httpclient.ArtifactoryHttpClient, cliHttpDetails httputils.HttpClientDetails) *xrayClient {
	return &xrayClient{client: client, xrayUrl: xrayUrl, cliHttpDetails: cliHttpDetails}
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

	requestFullUrl, err := utils2.BuildArtifactoryUrl(x.xrayUrl, "/api/v1/summary/component", make(map[string]string))
	utils2.SetContentType("application/json", &x.cliHttpDetails.Headers)
	log.Debug("Calling xray /api/v1/summary/component with components: ", components)
	_, body, err := x.client.SendPost(requestFullUrl, data, &x.cliHttpDetails)
	if err != nil {
		return nil, err
	}
	result := &ComponentSummaryResult{}
	log.Debug("Xray response: ", string(body))
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, errors.Wrap(err, "Failed parsing Xray response")
	}
	return result, nil
}

type ComponentSummaryRequest struct {
	ComponentDetails []ComponentSummaryRequestDetails `json:"component_details"`
}
type ComponentSummaryRequestDetails struct {
	ComponentID string `json:"component_id"`
}
type ComponentArtifact struct {
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
}

type ComponentSummaryResult struct {
	Artifacts []ComponentArtifact `json:"artifacts"`
}
