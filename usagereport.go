package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/usagereport-plugin/apihelper"
	"github.com/krujos/usagereport-plugin/models"
)

//UsageReportCmd the plugin
type UsageReportCmd struct {
	apiHelper apihelper.CFAPIHelper
}

// contains CLI flag values
type flagVal struct {
	OrgName      string
	Format       string
	ExcludedOrgs []string
}

type multipleFlagValues []string

func (i *multipleFlagValues) String() string {
	return "multiple flag values"
}

func (i *multipleFlagValues) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func ParseFlags(args []string) flagVal {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)

	// Create flags
	orgName := flagSet.String("o", "", "-o orgName")
	format := flagSet.String("f", "format", "-f <csv>")

	var orgsToExclude multipleFlagValues
	flagSet.Var(&orgsToExclude, "e", "-e orgName to exclude (multiple)")

	err := flagSet.Parse(args[1:])
	if err != nil {

	}

	return flagVal{
		OrgName:      string(*orgName),
		Format:       string(*format),
		ExcludedOrgs: orgsToExclude,
	}
}

//GetMetadata returns metatada
func (cmd *UsageReportCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "usage-report",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 4,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "usage-report",
				HelpText: "Report AI and memory usage for orgs and spaces",
				UsageDetails: plugin.Usage{
					Usage: "cf usage-report [-o orgName] [-f <csv>] [-e orgNameToExclude]",
					Options: map[string]string{
						"o": "organization",
						"f": "format",
						"e": "organization(s) to exclude",
					},
				},
			},
		},
	}
}

//UsageReportCommand doer
func (cmd *UsageReportCmd) UsageReportCommand(args []string) {
	flagVals := ParseFlags(args)

	var orgs []models.Org
	var err error
	var report models.Report

	if flagVals.OrgName != "" {
		org, err := cmd.getOrg(flagVals.OrgName)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		orgs = append(orgs, org)
	} else {
		orgs, err = cmd.getOrgs()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	report.Orgs = cmd.filterOrgs(orgs, flagVals.ExcludedOrgs)

	if flagVals.Format == "csv" {
		fmt.Println(report.CSV())
	} else {
		fmt.Println(report.String())
	}
}

func (cmd *UsageReportCmd) filterOrgs(orgs []models.Org, excludedOrgs []string) []models.Org {
	var filtered = []models.Org{}

	for _, o := range orgs {
		if contains(excludedOrgs, o.Name) {

		} else {
			filtered = append(filtered, o)
		}
	}

	return filtered
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if strings.TrimSpace(x) == strings.TrimSpace(n) {
			return true
		}
	}
	return false
}

func (cmd *UsageReportCmd) getOrgs() ([]models.Org, error) {
	rawOrgs, err := cmd.apiHelper.GetOrgs()
	if nil != err {
		return nil, err
	}

	var orgs = []models.Org{}

	for _, o := range rawOrgs {
		orgDetails, err := cmd.getOrgDetails(o)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, orgDetails)
	}

	return orgs, nil
}

func (cmd *UsageReportCmd) getOrg(name string) (models.Org, error) {
	rawOrg, err := cmd.apiHelper.GetOrg(name)
	if nil != err {
		return models.Org{}, err
	}

	return cmd.getOrgDetails(rawOrg)
}

func (cmd *UsageReportCmd) getOrgDetails(o apihelper.Organization) (models.Org, error) {
	usage, err := cmd.apiHelper.GetOrgMemoryUsage(o)
	if nil != err {
		return models.Org{}, err
	}
	quota, err := cmd.apiHelper.GetQuotaMemoryLimit(o.QuotaURL)
	if nil != err {
		return models.Org{}, err
	}
	spaces, err := cmd.getSpaces(o.SpacesURL)
	if nil != err {
		return models.Org{}, err
	}

	return models.Org{
		Name:        o.Name,
		MemoryQuota: int(quota),
		MemoryUsage: int(usage),
		Spaces:      spaces,
	}, nil
}

func (cmd *UsageReportCmd) getSpaces(spaceURL string) ([]models.Space, error) {
	rawSpaces, err := cmd.apiHelper.GetOrgSpaces(spaceURL)
	if nil != err {
		return nil, err
	}
	var spaces = []models.Space{}
	for _, s := range rawSpaces {
		apps, err := cmd.getApps(s.AppsURL)
		if nil != err {
			return nil, err
		}
		spaces = append(spaces,
			models.Space{
				Apps: apps,
				Name: s.Name,
			},
		)
	}
	return spaces, nil
}

func (cmd *UsageReportCmd) getApps(appsURL string) ([]models.App, error) {
	rawApps, err := cmd.apiHelper.GetSpaceApps(appsURL)
	if nil != err {
		return nil, err
	}
	var apps = []models.App{}
	for _, a := range rawApps {
		apps = append(apps, models.App{
			Instances: int(a.Instances),
			Ram:       int(a.RAM),
			Running:   a.Running,
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
