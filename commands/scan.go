package commands

import (
	"bufio"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"os"
)

const componentFlagKey = "component"

func GetScanCommand() components.Command {
	return components.Command{
		Name:        "scan",
		Description: "Scan the given components, either with a single arg or stdin",
		Aliases:     []string{"s"},
		Flags:       getScanFlags(),
		Action: func(c *components.Context) error {
			client, err := newXrayClient()
			if err != nil {
				return err
			}
			return scanCmd(c, client.scan)
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

func scanCmd(c *components.Context, scanner xrayScanner) error {
	var conf = &scanConfiguration{
		component: c.GetStringFlagValue(componentFlagKey),
	}
	lines := make(chan string)
	if conf.component == "" {
		go func() {
			stdinScanner := bufio.NewScanner(os.Stdin)
			defer close(lines)
			for stdinScanner.Scan() {
				lines <- stdinScanner.Text()
			}
			if err := stdinScanner.Err(); err != nil {
				panic(fmt.Sprintf("Could not read from stdin: %+v", err))
			}
		}()
	} else {
		lines <- conf.component
	}

	return scan(lines, scanner)
}

func scan(lines <-chan string, scanner xrayScanner) error {
	for line := range lines {
		comp := parse(line)
		// TODO: introduce buffering somewhere to avoid hammering xray
		result, err := scanner([]component{comp})
		if err != nil {
			return err
		}
		printResult(*result)
	}
	return nil
}
