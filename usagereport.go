package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/usagereport-plugin/apihelper"
)

//UsageReportCmd the plugin
type UsageReportCmd struct {
	apiHelper apihelper.CFAPIHelper
	cli       plugin.CliConnection
}

type org struct {
	name        string
	memoryQuota int
	memoryUsage int
	spaces      []space
}

type space struct {
	apps []app
	name string
}

type app struct {
	ram       int
	instances int
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
func (cmd *UsageReportCmd) UsageReportCommand(args []string) {
	fmt.Println("Gathering usage information")

	if nil == cmd.cli {
		fmt.Println("ERROR: CLI Connection is nil!")
		os.Exit(1)
	}

	orgs, err := cmd.getOrgs()
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	totalApps := 0
	totalInstances := 0
	for _, org := range orgs {
		fmt.Printf("Org %s is using %d MB of %d MB.\n", org.name, org.memoryUsage, org.memoryQuota)
		for _, space := range org.spaces {
			consumed := 0
			instances := 0
			for _, app := range space.apps {
				consumed += int(app.instances * app.ram)
				instances += int(app.instances)
			}
			fmt.Printf("\tSpace %s is using %d MB memory (%d%%) of org quota. %d apps %d instances.\n",
				space.name, consumed, (100 * consumed / org.memoryQuota), len(space.apps), instances)
			totalInstances += instances
			totalApps += len(space.apps)
		}
	}
	fmt.Printf("You are running %d apps in all orgs, with a total of %d instances.\n",
		totalApps, totalInstances)
}

func (cmd *UsageReportCmd) getOrgs() ([]org, error) {
	rawOrgs, err := cmd.apiHelper.GetOrgs(cmd.cli)
	if nil != err {
		return nil, err
	}

	var orgs = []org{}

	for _, o := range rawOrgs {
		usage, err := cmd.apiHelper.GetOrgMemoryUsage(cmd.cli, o)
		if nil != err {
			return nil, err
		}
		quota, err := cmd.apiHelper.GetQuotaMemoryLimit(cmd.cli, o.QuotaURL)
		if nil != err {
			return nil, err
		}
		spaces, err := cmd.getSpaces(o.SpacesURL)
		if nil != err {
			return nil, err
		}

		orgs = append(orgs, org{
			name:        o.Name,
			memoryQuota: int(quota),
			memoryUsage: int(usage),
			spaces:      spaces,
		})
	}
	return orgs, nil
}

func (cmd *UsageReportCmd) getSpaces(spaceURL string) ([]space, error) {
	rawSpaces, err := cmd.apiHelper.GetOrgSpaces(cmd.cli, spaceURL)
	if nil != err {
		return nil, err
	}
	var spaces = []space{}
	for _, s := range rawSpaces {
		apps, err := cmd.getApps(s.AppsURL)
		if nil != err {
			return nil, err
		}
		spaces = append(spaces,
			space{
				apps: apps,
				name: s.Name,
			},
		)
	}
	return spaces, nil
}

func (cmd *UsageReportCmd) getApps(appsURL string) ([]app, error) {
	rawApps, err := cmd.apiHelper.GetSpaceApps(cmd.cli, appsURL)
	if nil != err {
		return nil, err
	}
	var apps = []app{}
	for _, a := range rawApps {
		apps = append(apps, app{
			instances: int(a.Instances),
			ram:       int(a.RAM),
		})
	}
	return apps, nil
}

//Run runs the plugin
func (cmd *UsageReportCmd) Run(cli plugin.CliConnection, args []string) {
	if args[0] == "usage-report" {
		cmd.apiHelper = &apihelper.APIHelper{}
		cmd.cli = cli
		cmd.UsageReportCommand(args)
	}
}

func main() {
	plugin.Start(new(UsageReportCmd))
}
