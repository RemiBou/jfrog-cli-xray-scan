package commands

import (
	"bufio"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-cli-core/utils/coreutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io"
	"os"
	"strings"
)

const componentFlagKey = "component"
const displayNoIssuesKey = "display-no-issues"
const failKey = "fail"
const scanBufferSize = 50

func GetScanCommand() components.Command {
	return components.Command{
		Name:        "scan",
		Description: "Scan the given components, either with a single arg or stdin",
		Aliases:     []string{"s"},
		Flags:       getScanFlags(),
		Action: func(c *components.Context) error {
			confArti, err := config.GetDefaultArtifactoryConf()
			if err != nil {
				return err
			}
			servicesManager, err := utils.CreateServiceManager(confArti, false)
			if err != nil {
				return err
			}
			serviceDetails := servicesManager.GetConfig().GetServiceDetails()
			xrayUrl := strings.Replace(serviceDetails.GetUrl(), "/artifactory", "/xray", 1)
			client := newXrayClient(xrayUrl, servicesManager.Client(), serviceDetails.CreateHttpClientDetails())
			return scanCmd(c, os.Stdin, client.scan)
		},
	}
}

func getScanFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         componentFlagKey,
			Description:  "Display results for a specific component.",
			DefaultValue: "",
		},
		components.BoolFlag{
			Name:         failKey,
			Description:  "Return non zero exit code if at least one issue is found",
			DefaultValue: true,
		},
		components.BoolFlag{
			Name:         displayNoIssuesKey,
			Description:  "Display component with no issues.",
			DefaultValue: false,
		},
	}
}

type scanConfiguration struct {
	component       string
	displayNoIssues bool
	fail            bool
}

type xrayScanner func(comps []component) (*ComponentSummaryResult, error)

type cmdContext interface {
	GetStringFlagValue(flagName string) string
	GetBoolFlagValue(flagName string) bool
}

// Reads from stdin or from the argument and converts to a channel
func scanCmd(c cmdContext, stdin io.Reader, scanner xrayScanner) error {
	var conf = scanConfiguration{
		component:       c.GetStringFlagValue(componentFlagKey),
		displayNoIssues: c.GetBoolFlagValue(displayNoIssuesKey),
		fail:            c.GetBoolFlagValue(failKey),
	}
	lines := make(chan string)
	if conf.component == "" {
		log.Debug("No explicit component, scanning stdin")
		go func() {
			stdinScanner := bufio.NewScanner(stdin)
			defer close(lines)
			for stdinScanner.Scan() {
				lines <- stdinScanner.Text()
			}
			if err := stdinScanner.Err(); err != nil {
				panic(fmt.Sprintf("Could not read from stdin: %+v", err))
			}
		}()
	} else {
		log.Debug("Explicit component received ", conf.component)
		go func() {
			defer close(lines)
			lines <- conf.component
		}()
	}
	return scan(lines, scanner, conf)
}

// Central method where everything is orchestrated:
// - reads lines from the channel
// - tries to parse to an Xray component (maven, go...)
// - once buffer size is reached or no more lines are given, sends to Xray
// - prints back a summary result
func scan(lines <-chan string, scanner xrayScanner, conf scanConfiguration) error {
	printer := newPrinter(os.Stdout, printerConfig{printNoIssues: conf.displayNoIssues})
	var err error
	var result error
	var any bool
	buffer := make([]component, 0, scanBufferSize)
	alreadyDones := make(map[component]bool)
	for line := range lines {
		comp, ok := parse(line)
		if ok {
			log.Debug("Detected component ", comp)
			if _, alreadyDone := alreadyDones[comp]; !alreadyDone {
				buffer = append(buffer, comp)
				alreadyDones[comp] = true
			}
		}
		if len(buffer) == scanBufferSize {
			log.Debug("Component buffer is full")
			any, err = callScanPrintResult(scanner, buffer, printer)

			buffer = make([]component, 0, scanBufferSize)
			if err != nil {
				return err
			}
		}
	}
	if len(buffer) > 0 {
		any, err = callScanPrintResult(scanner, buffer, printer)
		if err != nil {
			return err
		}
	}
	printer.flush()
	if any && conf.fail {
		result = coreutils.CliError{
			ExitCode: coreutils.ExitCodeVulnerableBuild,
			ErrorMsg: "Xray vulnerabilities found",
		}
	}
	return result
}

func callScanPrintResult(scanner xrayScanner, buffer []component, printer *resultPrinter) (bool, error) {
	log.Debug("Scanning ", len(buffer), " components")
	result, err := scanner(buffer)
	if err != nil {
		return false, err
	}
	printer.print(*result)
	for _, r := range result.Artifacts {
		if len(r.Issues) > 0 {
			return true, nil
		}
	}
	return false, nil
}
