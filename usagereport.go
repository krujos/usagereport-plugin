package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/usagereport-plugin/apihelper"
)

//UsageReportCmd the plugin
type UsageReportCmd struct {
	apiHelper apihelper.CFAPIHelper
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
	running   bool
}

// contains CLI flag values
type flagVal struct {
	OrgName string
}

func ParseFlags(args []string) flagVal {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)

	// Create flags
	orgName := flagSet.String("o", "", "-o orgName")
	err := flagSet.Parse(args[1:])
	if err != nil {

	}

	return flagVal{
		OrgName: string(*orgName),
	}
}

//GetMetadata returns metatada
func (cmd *UsageReportCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "usage-report",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 3,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "usage-report",
				HelpText: "Report AI and memory usage for orgs and spaces",
				UsageDetails: plugin.Usage{
					Usage: "cf usage-report [-o orgName]",
					Options: map[string]string{
						"o": "organization",
					},
				},
			},
		},
	}
}

//UsageReportCommand doer
func (cmd *UsageReportCmd) UsageReportCommand(args []string) {
	flagVals := ParseFlags(args)

	fmt.Println("Gathering usage information")

	totalApps := 0
	totalInstances := 0

	var orgs []org
	var err error

	if flagVals.OrgName != "" {
		org, err := cmd.getOrg(flagVals.OrgName)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		orgs = append(orgs, *org)
	} else {
		orgs, err = cmd.getOrgs()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	for _, org := range orgs {
		appsPerOrg, instancesPerOrg := cmd.printOrg(org)
		totalApps += appsPerOrg
		totalInstances += instancesPerOrg
	}

	fmt.Printf("You are running %d apps in %d org(s), with a total of %d instances.\n",
		totalApps, len(orgs), totalInstances)
}

func (cmd *UsageReportCmd) printOrg(o org) (int, int) {
	totalApps := 0
	totalInstances := 0

	fmt.Printf("Org %s is consuming %d MB of %d MB.\n", o.name, o.memoryUsage, o.memoryQuota)
	for _, space := range o.spaces {
		consumed := 0
		instances := 0
		runningApps := 0
		runningInstances := 0
		for _, app := range space.apps {
			if app.running {
				consumed += int(app.instances * app.ram)
				runningApps++
				runningInstances += app.instances
			}
			instances += int(app.instances)
		}
		fmt.Printf("\tSpace %s is consuming %d MB memory (%d%%) of org quota.\n",
			space.name, consumed, (100 * consumed / o.memoryQuota))
		fmt.Printf("\t\t%d apps: %d running %d stopped\n", len(space.apps),
			runningApps, len(space.apps)-runningApps)
		fmt.Printf("\t\t%d instances: %d running, %d stopped\n", instances,
			runningInstances, instances-runningInstances)
		totalInstances += instances
		totalApps += len(space.apps)
	}

	return totalApps, totalInstances
}

func (cmd *UsageReportCmd) getOrgs() ([]org, error) {
	rawOrgs, err := cmd.apiHelper.GetOrgs()
	if nil != err {
		return nil, err
	}

	var orgs = []org{}

	for _, o := range rawOrgs {
		orgDetails, err := cmd.getOrgDetails(o)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, *orgDetails)
	}
	return orgs, nil
}

func (cmd *UsageReportCmd) getOrg(name string) (*org, error) {
	rawOrg, err := cmd.apiHelper.GetOrg(name)
	if nil != err {
		return nil, err
	}

	return cmd.getOrgDetails(rawOrg)
}

func (cmd *UsageReportCmd) getOrgDetails(o apihelper.Organization) (*org, error) {
	usage, err := cmd.apiHelper.GetOrgMemoryUsage(o)
	if nil != err {
		return nil, err
	}
	quota, err := cmd.apiHelper.GetQuotaMemoryLimit(o.QuotaURL)
	if nil != err {
		return nil, err
	}
	spaces, err := cmd.getSpaces(o.SpacesURL)
	if nil != err {
		return nil, err
	}

	return &org{
		name:        o.Name,
		memoryQuota: int(quota),
		memoryUsage: int(usage),
		spaces:      spaces,
	}, nil
}

func (cmd *UsageReportCmd) getSpaces(spaceURL string) ([]space, error) {
	rawSpaces, err := cmd.apiHelper.GetOrgSpaces(spaceURL)
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
	rawApps, err := cmd.apiHelper.GetSpaceApps(appsURL)
	if nil != err {
		return nil, err
	}
	var apps = []app{}
	for _, a := range rawApps {
		apps = append(apps, app{
			instances: int(a.Instances),
			ram:       int(a.RAM),
			running:   a.Running,
		})
	}
	return apps, nil
}

//Run runs the plugin
func (cmd *UsageReportCmd) Run(cli plugin.CliConnection, args []string) {
	if args[0] == "usage-report" {
		cmd.apiHelper = apihelper.New(cli)
		cmd.UsageReportCommand(args)
	}
}

func main() {
	plugin.Start(new(UsageReportCmd))
}
