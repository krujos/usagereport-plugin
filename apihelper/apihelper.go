package apihelper

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/cfcurl"
)

//Organization representation
type Organization struct {
	URL       string
	Name      string
	QuotaURL  string
	SpacesURL string
}

//Space representation
type Space struct {
	Name    string
	AppsURL string
}

//App representation
type App struct {
	Instances float64
	RAM       float64
}

//CFAPIHelper to wrap cf curl results
type CFAPIHelper interface {
	GetOrgs(plugin.CliConnection) ([]Organization, error)
	GetQuotaMemoryLimit(plugin.CliConnection, string) (float64, error)
	GetOrgMemoryUsage(plugin.CliConnection, Organization) (float64, error)
	GetOrgSpaces(plugin.CliConnection, string) ([]Space, error)
	GetSpaceApps(plugin.CliConnection, string) ([]App, error)
}

//APIHelper implementation
type APIHelper struct{}

//GetOrgs returns a struct that represents critical fields in the JSON
func (api *APIHelper) GetOrgs(cli plugin.CliConnection) ([]Organization, error) {
	orgsJSON, err := cfcurl.Curl(cli, "/v2/organizations")
	if nil != err {
		return nil, err
	}

	orgs := []Organization{}
	for _, o := range orgsJSON["resources"].([]interface{}) {
		theOrg := o.(map[string]interface{})
		entity := theOrg["entity"].(map[string]interface{})
		metadata := theOrg["metadata"].(map[string]interface{})
		orgs = append(orgs,
			Organization{
				Name:      entity["name"].(string),
				URL:       metadata["url"].(string),
				QuotaURL:  entity["quota_definition_url"].(string),
				SpacesURL: entity["spaces_url"].(string),
			})
	}
	return orgs, nil
}

//GetQuotaMemoryLimit retruns the amount of memory (in MB) that the org is allowed
func (api *APIHelper) GetQuotaMemoryLimit(cli plugin.CliConnection, quotaURL string) (float64, error) {
	quotaJSON, err := cfcurl.Curl(cli, quotaURL)
	if nil != err {
		return 0, err
	}
	return quotaJSON["entity"].(map[string]interface{})["memory_limit"].(float64), nil
}

//GetOrgMemoryUsage returns the amount of memory (in MB) that the org is consuming
func (api *APIHelper) GetOrgMemoryUsage(cli plugin.CliConnection, org Organization) (float64, error) {
	usageJSON, err := cfcurl.Curl(cli, org.URL+"/memory_usage")
	if nil != err {
		return 0, err
	}
	return usageJSON["memory_usage_in_mb"].(float64), nil
}

//GetOrgSpaces returns the spaces in an org.
func (api *APIHelper) GetOrgSpaces(cli plugin.CliConnection, spacesURL string) ([]Space, error) {
	spacesJSON, err := cfcurl.Curl(cli, spacesURL)
	if nil != err {
		return nil, err
	}
	spaces := []Space{}
	for _, s := range spacesJSON["resources"].([]interface{}) {
		theSpace := s.(map[string]interface{})
		entity := theSpace["entity"].(map[string]interface{})
		spaces = append(spaces,
			Space{
				AppsURL: entity["apps_url"].(string),
				Name:    entity["name"].(string),
			})
	}
	return spaces, nil
}

//GetSpaceApps returns the apps in a space
func (api *APIHelper) GetSpaceApps(cli plugin.CliConnection, spaceURL string) ([]App, error) {
	appsJSON, err := cfcurl.Curl(cli, spaceURL)
	if nil != err {
		return nil, err
	}
	apps := []App{}
	for _, a := range appsJSON["resources"].([]interface{}) {
		theApp := a.(map[string]interface{})
		entity := theApp["entity"].(map[string]interface{})
		apps = append(apps,
			App{
				Instances: entity["instances"].(float64),
				RAM:       entity["memory"].(float64),
			})
	}
	return apps, nil
}
