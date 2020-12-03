package commands

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

// TODO: extract to file
const resultAsJson = "{\n" +
	"  \"artifacts\": [\n" +
	"    {\n" +
	"      \"general\": {\n" +
	"        \"component_id\": \"org.apache.httpcomponents:httpclient:4.5.9\",\n" +
	"        \"name\": \"org.apache.httpcomponents:httpclient\",\n" +
	"        \"pkg_type\": \"maven\"\n" +
	"      },\n" +
	"      \"issues\": [\n" +
	"        {\n" +
	"          \"components\": [\n" +
	"            {\n" +
	"              \"component_id\": \"org.apache.httpcomponents:httpclient\",\n" +
	"              \"fixed_versions\": [\n" +
	"                \"[4.5.13]\"\n" +
	"              ]\n" +
	"            },\n" +
	"            {\n" +
	"              \"component_id\": \"org.apache.httpcomponents.client5:httpclient5\",\n" +
	"              \"fixed_versions\": [\n" +
	"                \"[5.0.3]\"\n" +
	"              ]\n" +
	"            }\n" +
	"          ],\n" +
	"          \"created\": \"2020-10-11T00:00:00.429Z\",\n" +
	"          \"cves\": [\n" +
	"            {\n" +
	"              \"cvss_v2\": \"5.0/AV:N/AC:L/Au:N/C:N/I:P/A:N\"\n" +
	"            }\n" +
	"          ],\n" +
	"          \"description\": \"Apache HttpComponents HttpClient contains a flaw that is triggered as malformed authority components in request URIs are improperly accepted. This may allow a remote attacker to bypass intended restrictions.\",\n" +
	"          \"issue_type\": \"security\",\n" +
	"          \"provider\": \"JFrog\",\n" +
	"          \"severity\": \"Medium\",\n" +
	"          \"summary\": \"Apache HttpComponents HttpClient Request URI Authority Component Handling Acceptance\"\n" +
	"        }\n" +
	"      ],\n" +
	"      \"licenses\": [\n" +
	"        {\n" +
	"          \"components\": [\n" +
	"            \"gav://org.apache.httpcomponents:httpclient:4.5.9\"\n" +
	"          ],\n" +
	"          \"full_name\": \"The Apache Software License, Version 2.0\",\n" +
	"          \"more_info_url\": [\n" +
	"            \"http://www.apache.org/licenses/LICENSE-2.0\",\n" +
	"            \"https://spdx.org/licenses/Apache-2.0.html\",\n" +
	"            \"https://spdx.org/licenses/Apache-2.0\",\n" +
	"            \"http://www.opensource.org/licenses/apache2.0.php\",\n" +
	"            \"http://www.opensource.org/licenses/Apache-2.0\"\n" +
	"          ],\n" +
	"          \"name\": \"Apache-2.0\"\n" +
	"        }\n" +
	"      ]\n" +
	"    }\n" +
	"  ]\n" +
	"}"

func Test_resultPrinter_print(t *testing.T) {
	printer, err := newPrinter()
	require.NoError(t, err)

	result := &ComponentSummaryResult{}
	err = json.Unmarshal([]byte(resultAsJson), result)
	require.NoError(t, err)

	bufferString := bytes.NewBufferString("")
	err = printer.print(*result, bufferString)
	require.NoError(t, err)
	expected :=
		`| component                 |    issues (high/medium/low) |     min fix version   |      licences
| golang.org/x/net:1.8.2    |    2/3/4                    |     1.12              |       Unknown license`
	require.Equal(t, expected, bufferString.String())
}
