# UsageReport Plugin
This CF CLI Plugin to shows memory consumption and application instances for each org and space you have permission to access.

[![wercker status](https://app.wercker.com/status/8881b5530809e3636080d2df6433aada/s/master "wercker status")](https://app.wercker.com/project/bykey/8881b5530809e3636080d2df6433aada)


#Usage

```
➜  usagereport-plugin git:(master) ✗ cf usage-report
Gathering usage information
Org platform-eng is consuming 53400 MB of 204800 MB.
	Space CFbook is consuming 128 MB memory (0%) of org quota.
		1 apps: 1 running 0 stopped
		1 instances: 1 running, 0 stopped
Org krujos is consuming 512 MB of 10240 MB.
	Space development is consuming 0 MB memory (0%) of org quota.
		4 apps: 0 running 4 stopped
		4 instances: 0 running, 4 stopped
	Space production is consuming 512 MB memory (5%) of org quota.
		1 apps: 1 running 0 stopped
		2 instances: 2 running, 0 stopped
Org pcfp is consuming 7296 MB of 102400 MB.
	Space development is consuming 0 MB memory (0%) of org quota.
		0 apps: 0 running 0 stopped
		0 instances: 0 running, 0 stopped
	Space docs-staging is consuming 512 MB memory (0%) of org quota.
		2 apps: 1 running 1 stopped
		4 instances: 2 running, 2 stopped
	Space docs-prod is consuming 512 MB memory (0%) of org quota.
		3 apps: 1 running 2 stopped
		5 instances: 2 running, 3 stopped
	Space guillermo-playground is consuming 2560 MB memory (2%) of org quota.
		1 apps: 1 running 0 stopped
		5 instances: 5 running, 0 stopped
	Space haydon-playground is consuming 1024 MB memory (1%) of org quota.
		1 apps: 1 running 0 stopped
		1 instances: 1 running, 0 stopped
	Space jkruck-playground is consuming 128 MB memory (0%) of org quota.
		1 apps: 1 running 0 stopped
		1 instances: 1 running, 0 stopped
	Space rsalas-dev is consuming 0 MB memory (0%) of org quota.
		0 apps: 0 running 0 stopped
		0 instances: 0 running, 0 stopped
	Space shekel-dev is consuming 1536 MB memory (1%) of org quota.
		3 apps: 3 running 0 stopped
		3 instances: 3 running, 0 stopped
	Space shekel-qa is consuming 0 MB memory (0%) of org quota.
		0 apps: 0 running 0 stopped
		0 instances: 0 running, 0 stopped
	Space hd-playground is consuming 0 MB memory (0%) of org quota.
		0 apps: 0 running 0 stopped
		0 instances: 0 running, 0 stopped
	Space dwallraff-dev is consuming 1024 MB memory (1%) of org quota.
		1 apps: 1 running 0 stopped
		1 instances: 1 running, 0 stopped
You are running 18 apps in all orgs, with a total of 27 instances.
```

##Installation
#####Install from CLI 
  ```
  $ cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
  $ cf install-plugin 'Usage Report' -r CF-Community
  ```


#####Install from Source (need to have [Go](http://golang.org/dl/) installed)
  ```
  $ go get github.com/cloudfoundry/cli
  $ go get github.com/krujos/usagereport-plugin
  $ cd $GOPATH/src/github.com/krujos/usagereport-plugin
  $ go build
  $ cf install-plugin usagereport-plugin
  ```
