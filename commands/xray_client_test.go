//+build itest

package commands

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// This test assume an Xray is provisioned with a valid user. For example on JFrog cloud free tier.
// Then test should be launched with the following env vars
// XRAY_URL = "https://myarti.jfrog.io/xray"
// ARTI_USERNAME = "myuser"
// ARTI_PASSWORD = "pazz"

func Test_xrayClient_scan_single_component(t *testing.T) {
	err, xray := createXrayClient(t)
	require.NoError(t, err)

	result, err := xray.scan([]component{"go://golang.org/x/net:1.8.2"})

	require.NotNil(t, result)
	assert.Equal(t, 1, len(result.Artifacts))
	assert.True(t, contains(result, criteria{
		id:       "golang.org/x/net:1.8.2",
		nbIssues: 1,
		license:  "Unknown",
	}))
	assert.Equal(t, result.Artifacts[0].Issues[0].Severity, "Medium")
}

func Test_xrayClient_scan_few_components(t *testing.T) {
	err, xray := createXrayClient(t)
	require.NoError(t, err)

	result, err := xray.scan([]component{
		"go://golang.org/x/net:1.8.2",
		"gav://org.apache.thrift:libthrift:0.13.0",
		"gav://org.codehaus.plexus:plexus-utils:3.0.16",
		"gav://org.flywaydb:flyway-core:4.1.1",
	})

	require.NotNil(t, result)

	assert.True(t, contains(result, criteria{
		id:       "golang.org/x/net:1.8.2",
		nbIssues: 1,
		license:  "Unknown",
	}))

	assert.True(t, contains(result, criteria{
		id:       "org.apache.thrift:libthrift:0.13.0",
		nbIssues: 0,
		license:  "Apache-2.0",
	}))

	assert.True(t, contains(result, criteria{
		id:       "org.flywaydb:flyway-core:4.1.1",
		nbIssues: 0,
		license:  "Apache-2.0",
	}))

	assert.True(t, contains(result, criteria{
		id:       "org.codehaus.plexus:plexus-utils:3.0.16",
		nbIssues: 2,
		license:  "Apache-2.0",
	}))
}

type criteria struct {
	id       string
	nbIssues int
	license  string
}

func contains(result *ComponentSummaryResult, crit criteria) bool {
	for _, art := range result.Artifacts {
		// There might be more in the future !
		if art.General.ComponentID == crit.id &&
			len(art.Issues) >= crit.nbIssues &&
			art.Licenses[0].Name == crit.license {
			return true
		}
	}
	return false
}

func createXrayClient(t *testing.T) (error, *xrayClient) {
	details := httputils.HttpClientDetails{
		User:   os.Getenv("ARTI_USERNAME"),
		ApiKey: os.Getenv("ARTI_PASSWORD"),
	}
	artifactoryDetails := auth.NewArtifactoryDetails()
	client, err := httpclient.ArtifactoryClientBuilder().
		SetServiceDetails(&artifactoryDetails).Build()
	require.NoError(t, err)

	xray := newXrayClient(os.Getenv("XRAY_URL"), client, details)
	return err, xray
}
