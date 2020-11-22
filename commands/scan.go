package commands

import (
	"github.com/jfrog/jfrog-cli-core/plugins/components"
)

const componentFlagKey = "component"

func GetScanCommand() components.Command {
	return components.Command{
		Name:        "scan",
		Description: "Scan the given components, either with a single arg or stdin",
		Aliases:     []string{"s"},
		Flags:       getScanFlags(),
		Action: func(c *components.Context) error {
			return scanCmd(c)
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

func scanCmd(c *components.Context) error {
	var conf = &scanArgs{
		component: c.GetStringFlagValue(componentFlagKey),
	}
	if conf.component == "" {
		conf.stdin = true
	}
	// TODO: handle parsing of stdin, might unify both single string and stdin parsing
	comp := parse(conf.component)
	client, err := newXrayClient()
	if err != nil {
		return err
	}
	result, err := client.scan([]component{comp})
	if err != nil {
		return err
	}
	printResult(*result)
	return nil
}
