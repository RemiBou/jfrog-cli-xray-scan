package commands

import (
	"bufio"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io"
	"os"
	"strings"
)

const componentFlagKey = "component"
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
	}
}

type scanConfiguration struct {
	component string
}

type xrayScanner func(comps []component) (*ComponentSummaryResult, error)

type cmdContext interface {
	GetStringFlagValue(flagName string) string
}

// Reads from stdin or from the argument and converts to a channel
func scanCmd(c cmdContext, stdin io.Reader, scanner xrayScanner) error {
	var conf = &scanConfiguration{
		component: c.GetStringFlagValue(componentFlagKey),
	}
	lines := make(chan string)
	if conf.component == "" {
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
		go func() {
			defer close(lines)
			lines <- conf.component
		}()
	}
	return scan(lines, scanner)
}

// Central method where everything is orchestrated:
// - reads lines from the channel
// - tries to parse to an Xray component (maven, go...)
// - once buffer size is reached or no more lines are given, sends to Xray
// - prints back a summary result
func scan(lines <-chan string, scanner xrayScanner) error {
	printer, err := newPrinter(os.Stdout)
	if err != nil {
		return err
	}
	buffer := make([]component, 0, scanBufferSize)
	for line := range lines {
		comp, ok := parse(line)
		if ok {
			log.Debug("Detected component %v", comp)
			buffer = append(buffer, comp)
		}
		if len(buffer) == scanBufferSize {
			log.Debug("Component buffer is full")
			err := callScanPrintResult(scanner, buffer, printer)
			buffer = make([]component, 0, scanBufferSize)
			if err != nil {
				return err
			}
		}
	}
	if len(buffer) > 0 {
		err := callScanPrintResult(scanner, buffer, printer)
		if err != nil {
			return err
		}
	}
	return nil
}

func callScanPrintResult(scanner xrayScanner, buffer []component, printer *resultPrinter) error {
	result, err := scanner(buffer)
	if err != nil {
		return err
	}
	err = printer.print(*result)
	if err != nil {
		return err
	}
	return nil
}
