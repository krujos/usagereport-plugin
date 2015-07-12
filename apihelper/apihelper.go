package apihelper

import (
	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/cfcurl"
)

//Organization representation
type Organization struct {
	url       string
	name      string
	quotaURL  string
	spacesURL string
}

//CFAPIHelper to wrap cf curl results
type CFAPIHelper interface {
	GetOrgs(plugin.CliConnection) ([]Organization, error)
	GetQuotaMemoryLimit(plugin.CliConnection, string) (float64, error)
	GetOrgMemoryUsage(plugin.CliConnection, Organization) (float64, error)
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
				name:      entity["name"].(string),
				url:       metadata["url"].(string),
				quotaURL:  entity["quota_definition_url"].(string),
				spacesURL: entity["spaces_url"].(string),
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

//GetOrgMemoryUsage returns teh amount of memory (in MB) that the org is consuming
func (api *APIHelper) GetOrgMemoryUsage(cli plugin.CliConnection, org Organization) (float64, error) {
	usageJSON, err := cfcurl.Curl(cli, org.url+"/memory_usage")
	if nil != err {
		return 0, err
	}
	return usageJSON["memory_usage_in_mb"].(float64), nil
}
