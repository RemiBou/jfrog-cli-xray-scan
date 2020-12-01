package commands

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testData = "testData/"

//TODO: replace by itest that will use free tier instance
func xrayMockServer(file string) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/v1/summary/component", func(writer http.ResponseWriter, request *http.Request) {
		file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", testData, file))
		if err != nil {
			panic(fmt.Sprintf("Cannot read test file %v", file))
		}
		_, _ = writer.Write(file)
	})

	srv := httptest.NewServer(handler)
	return srv
}

func Test_xrayClient_scan(t *testing.T) {
	tests := []struct {
		response string
		wantErr  bool
		wantRes  bool
	}{
		{response: "one_component.json", wantRes: true, wantErr: false},
		{response: "few_components.json", wantRes: true, wantErr: false},
		{response: "unrecognized_response.json", wantRes: true, wantErr: false},
		{response: "bad_response", wantRes: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.response, func(t *testing.T) {

			srv := xrayMockServer(tt.response)
			defer srv.Close()

			details := httputils.HttpClientDetails{
				User:   "admin",
				ApiKey: "XXX",
			}
			artifactoryDetails := auth.NewArtifactoryDetails()
			client, err := httpclient.ArtifactoryClientBuilder().
				SetServiceDetails(&artifactoryDetails).Build()
			require.NoError(t, err)

			xray, err := newXrayClient(srv.URL, client, details)
			require.NoError(t, err)

			result, err := xray.scan([]component{"does_not_matter_for_this_test"})

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantRes {
				file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", testData, tt.response))
				expected := &ComponentSummaryResult{}
				err = json.Unmarshal(file, expected)
				require.NoError(t, err)
				assert.Equal(t, result, expected)
			} else {
				require.Nil(t, result)
			}
		})
	}
}
