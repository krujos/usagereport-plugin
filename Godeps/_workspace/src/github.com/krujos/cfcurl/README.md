cfcurl
======

A library for calling the cf curl command. This is intended to be used when
writing CF CLI plugins. Sometimes it's easier to interact with the api directly.
During the implementation of the
[scaleover](https://github.com/krujos/scaleover-plugin/) we had a conversation
about accessing the CF api through the CLI framework. The conclusion was that
isn't such a great idea to use the internals of the CLI, so this feels like the
rght way to do it to me.

#Usage
[Here's an example plugin](https://github.com/krujos/cfcurl-testplugin) that makes use of it, it prints the contents of the marshaled JSON. 

```
//Run a command
func (cmd *TestCmd) Run(cliConnection plugin.CliConnection, args []string) {
        out, _ := cfcurl.Curl(cliConnection, "/v2/apps")
        fmt.Println(out)

        out, _ = cfcurl.CurlDepricated(cliConnection, "/v2/domains")
        fmt.Println(out)
}

```

#API
The package offers two methods `Curl` and `CurlDepricated`.

* `Curl` calls "current" (or experimental) API specified by the path argument and returns a `map[string]interface{}`. If you call a deprecated API with this method it will `panic`.
* `CurlDepricated` will let you call a "Endpoint deprecated" API, and returns `map[string]interface{}` representing the JSON. You can call a current API with no issue, but the converse will panic. 

See the [tests](https://github.com/krujos/cfcurl/blob/master/cfcurl_test.go) for more details, but you can handle the `map[string]interface{}` pretty intuitively. 

For instances, to read the total number of results returned from `/v2/apps` you would use the following code:

```
appsJSON, _ := Curl(fakeCliConnection, "/v2/apps")
fmt.Println(appsJSON["total_results"]))
```