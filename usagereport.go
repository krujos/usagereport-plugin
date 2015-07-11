package main

import (
	"errors"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/cfcurl"
)

//UsageReportCmd the plugin
type UsageReportCmd struct {
}

type organization struct {
	url       string
	name      string
	quotaURL  string
	spacesURL string
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

//Run runs the plugin
func (cmd *UsageReportCmd) Run(cli plugin.CliConnection, args []string) {
	if args[0] == "usage-report" {
		cmd.UsageReportCommand(cli, args)
	}
}

func (cmd *UsageReportCmd) getOrgs(cli plugin.CliConnection) ([]organization, error) {
	orgsJSON, err := cfcurl.Curl(cli, "/v2/organizations")

	if nil != err {
		//TODO Swollow this?
		return nil, errors.New("Failed to get orgs!")
	}
	orgs := []organization{}
	for _, o := range orgsJSON["resources"].([]interface{}) {
		theOrg := o.(map[string]interface{})
		entity := theOrg["entity"].(map[string]interface{})
		metadata := theOrg["metadata"].(map[string]interface{})
		orgs = append(orgs,
			organization{
				name:      entity["name"].(string),
				url:       metadata["url"].(string),
				quotaURL:  entity["quota_definition_url"].(string),
				spacesURL: entity["spaces_url"].(string),
			})
	}
	return orgs, nil
}

func (cmd *UsageReportCmd) getQuotaMemoryLimit(cli plugin.CliConnection, quotaURL string) (float64, error) {
	quotaJSON, err := cfcurl.Curl(cli, quotaURL)
	if nil != err {
		return 0, err
	}
	return quotaJSON["entity"].(map[string]interface{})["memory_limit"].(float64), nil
}

func (cmd *UsageReportCmd) getOrgMemoryUsage(cli plugin.CliConnection, org organization) (float64, error) {
	usageJSON, err := cfcurl.Curl(cli, org.url+"/memory_usage")
	if nil != err {
		return 0, err
	}
	return usageJSON["memory_usage_in_mb"].(float64), nil
}

func main() {
	plugin.Start(new(UsageReportCmd))
}
