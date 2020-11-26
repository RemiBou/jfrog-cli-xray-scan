package commands

import (
	"bufio"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"io"
	"os"
)

const componentFlagKey = "component"

type xrayScanner struct {
	xrayClient xrayClientInterface
	stdin      io.Reader
}

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
			s := xrayScanner{
				xrayClient: client,
				stdin:      os.Stdin,
			}
			return s.scanCmd(c)
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

type scanArgs struct {
	component string
	stdin     bool
}

type ContextInterface interface {
	GetStringFlagValue(flagName string) string
}

func (s *xrayScanner) scanCmd(c ContextInterface) error {
	var conf = &scanArgs{
		component: c.GetStringFlagValue(componentFlagKey),
	}
	var comps []component
	if conf.component == "" {
		conf.stdin = true
		comps = s.readStdIn()
	} else {
		comps = []component{component(conf.component)}
	}

	result, err := s.xrayClient.scan(comps)
	if err != nil {
		return err
	}
	printResult(*result)
	return nil
}

func (s *xrayScanner) readStdIn() []component {
	scanner := bufio.NewScanner(s.stdin)
	res := make([]component, 0)
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			component := parse(text)
			if component != "" {
				res = append(res, component)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return res
}
