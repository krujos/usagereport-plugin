package main

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/usagereport-plugin/apihelper"
)

//UsageReportCmd the plugin
type UsageReportCmd struct {
	apiHelper apihelper.CFAPIHelper
	cli       plugin.CliConnection
}

type Org struct {
	name        string
	memoryQuota float64
	memoryUsage float64
}

//GetMetadata returns metatada
func (cmd *UsageReportCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "usage-report",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "usage-report",
				HelpText: "Report AI and memory usage for orgs and spaces",
				UsageDetails: plugin.Usage{
					Usage: "cf usage-report",
				},
			},
		},
	}
}

//UsageReportCommand doer
func (cmd *UsageReportCmd) UsageReportCommand(cli plugin.CliConnection, args []string) {
	//Do the things
}

func (cmd *UsageReportCmd) getOrgs() ([]Org, error) {
	rawOrgs, err := cmd.apiHelper.GetOrgs(cmd.cli)
	if nil != err {
		return nil, err
	}

	var orgs = []Org{}

	for _, org := range rawOrgs {

		usage, err := cmd.apiHelper.GetOrgMemoryUsage(cmd.cli, org)
		if nil != err {
			return nil, err
		}

		quota, err := cmd.apiHelper.GetQuotaMemoryLimit(cmd.cli, org.QuotaURL)
		if nil != err {
			return nil, err
		}

		orgs = append(orgs, Org{
			name:        org.Name,
			memoryQuota: quota,
			memoryUsage: usage,
		})
	}
	return orgs, nil
}

//Run runs the plugin
func (cmd *UsageReportCmd) Run(cli plugin.CliConnection, args []string) {
	if args[0] == "usage-report" {
		cmd.UsageReportCommand(cli, args)
	}
}

func main() {
	plugin.Start(new(UsageReportCmd))
}
