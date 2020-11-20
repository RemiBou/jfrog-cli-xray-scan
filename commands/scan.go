package commands

import (
	"encoding/json"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	utils2 "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"time"
)

const currentProjectFlagKey = "current-project"
const componentFlagKey = "component"

func GetScanCommand() components.Command {
	return components.Command{
		Name:        "scan",
		Description: "Scan the current project dependency and display security issues",
		Aliases:     []string{"s"},
		Flags:       getScanFlags(),
		Action: func(c *components.Context) error {
			return scanCmd(c)
		},
	}
}

func getScanFlags() []components.Flag {
	return []components.Flag{
		components.BoolFlag{
			Name:         currentProjectFlagKey,
			Description:  "Scan the project in the current folder",
			DefaultValue: true,
		},
		components.StringFlag{
			Name:         componentFlagKey,
			Description:  "Display results for a specific component.",
			DefaultValue: "",
		},
	}
}

type scanConfiguration struct {
	scanCurrentProject bool
	component          string
}

func scanCmd(c *components.Context) error {
	var conf = new(scanConfiguration)
	conf.scanCurrentProject = c.GetBoolFlagValue(currentProjectFlagKey)
	conf.component = c.GetStringFlagValue(componentFlagKey)
	if conf.component != "" {
		conf.scanCurrentProject = false
	}
	res := ComponentSummaryRequest{
		ComponentDetails: []ComponentSummaryRequestDetails{
			{
				"gav://org.apache.httpcomponents:httpclient:4.5.9",
			},
		},
	}
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	confArti, err := config.GetDefaultArtifactoryConf()
	if err != nil {
		return err
	}
	servicesManager, err := utils.CreateServiceManager(confArti, false)
	client := servicesManager.Client()
	if err != nil {
		return err
	}
	serviceDetails := servicesManager.GetConfig().GetServiceDetails()
	details := serviceDetails.CreateHttpClientDetails()
	requestFullUrl, err := utils2.BuildArtifactoryUrl(serviceDetails.GetUrl(), "/xray/api/v1/summary/component", make(map[string]string))
	_, body, err := client.SendPost(requestFullUrl, data, &details)
	if err != nil {
		return err
	}
	log.Output(string(body))
	return nil
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
