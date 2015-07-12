package cfcurl

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

func callAndValidateCLI(cli plugin.CliConnection, path string) ([]string, error) {
	output, err := cli.CliCommandWithoutTerminalOutput("curl", path)

	if nil != err {
		return nil, err
	}

	if nil == output || 0 == len(output) {
		return nil, errors.New("CF API returned no output")
	}

	return output, nil
}

func parseOutput(output []string) (map[string]interface{}, error) {
	if nil == output || 0 == len(output) {
		return nil, errors.New("CF API returned no output")
	}

	data := strings.Join(output, "\n")

	if 0 == len(data) || "" == data {
		return nil, errors.New("Failed to join output")
	}

	var f interface{}
	err := json.Unmarshal([]byte(data), &f)
	return f.(map[string]interface{}), err
}

// Curl calls cf curl  and return the resulting json. This method will panic if
// the api is depricated
func Curl(cli plugin.CliConnection, path string) (map[string]interface{}, error) {
	output, err := cli.CliCommandWithoutTerminalOutput("curl", path)

	if nil != err {
		return nil, err
	}

	return parseOutput(output)
}

// CurlDepricated calls cf curl and return the resulting json, even if the api is depricated
func CurlDepricated(cli plugin.CliConnection, path string) (map[string]interface{}, error) {
	output, err := callAndValidateCLI(cli, path)
	if nil != err {
		return nil, err
	}

	if strings.Contains(output[len(output)-1], "Endpoint deprecated") {
		output = output[:len(output)-1]
	}

	return parseOutput(output)
}
