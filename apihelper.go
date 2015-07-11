package main

import (
	"errors"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/cfcurl"
)

//APIHelper to wrap cf curl results
type APIHelper struct{}

func (api *APIHelper) getOrgs(cli plugin.CliConnection) ([]organization, error) {
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

func (api *APIHelper) getQuotaMemoryLimit(cli plugin.CliConnection, quotaURL string) (float64, error) {
	quotaJSON, err := cfcurl.Curl(cli, quotaURL)
	if nil != err {
		return 0, err
	}
	return quotaJSON["entity"].(map[string]interface{})["memory_limit"].(float64), nil
}

func (api *APIHelper) getOrgMemoryUsage(cli plugin.CliConnection, org organization) (float64, error) {
	usageJSON, err := cfcurl.Curl(cli, org.url+"/memory_usage")
	if nil != err {
		return 0, err
	}
	return usageJSON["memory_usage_in_mb"].(float64), nil
}
