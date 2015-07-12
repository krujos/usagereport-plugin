# UsageReport Plugin
This CF CLI Plugin to shows memory consumption and application instances for each org and space you have permission to access. 

[![wercker status](https://app.wercker.com/status/f5f8d90193968cce6f5d60583d85be3c/s "wercker status")](https://app.wercker.com/project/bykey/f5f8d90193968cce6f5d60583d85be3c)


#Usage

```
➜  usagereport-plugin git:(master) ✗ cf usage-report                                                                                                                               $
Gathering usage information
Org platform-eng is using 46232 MB of 204800 MB.
        Space CFbook is using 128 MB memory (0%) of org quota. 1 apps 1 instances.
Org krujos is using 512 MB of 10240 MB.
        Space development is using 3200 MB memory (31%) of org quota. 4 apps 4 instances.
        Space production is using 512 MB memory (5%) of org quota. 1 apps 2 instances.
Org pcfp is using 5120 MB of 102400 MB.
        Space development is using 0 MB memory (0%) of org quota. 0 apps 0 instances.
        Space docs-staging is using 1024 MB memory (1%) of org quota. 2 apps 4 instances.
        Space docs-prod is using 2048 MB memory (2%) of org quota. 3 apps 5 instances.
        Space guillermo-playground is using 512 MB memory (0%) of org quota. 1 apps 1 instances.
        Space haydon-playground is using 1024 MB memory (1%) of org quota. 1 apps 1 instances.
        Space jkruck-playground is using 0 MB memory (0%) of org quota. 0 apps 0 instances.
        Space rsalas-dev is using 0 MB memory (0%) of org quota. 0 apps 0 instances.
        Space shekel-dev is using 1536 MB memory (1%) of org quota. 3 apps 3 instances.
        Space shekel-qa is using 0 MB memory (0%) of org quota. 0 apps 0 instances.
        Space hd-playground is using 0 MB memory (0%) of org quota. 0 apps 0 instances.
        Space dwallraff-dev is using 1024 MB memory (1%) of org quota. 1 apps 1 instances.
You are running 17 apps in all orgs, with a total of 22 instances.
```

##Installation 
#####Install from CLI (this will not work until [#34](https://github.com/cloudfoundry-incubator/cli-plugin-repo/pull/34) is merged)

  ```
  $ cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
  $ cf install-plugin usagereport -r CF-Community
  ```
  
  
#####Install from Source (need to have [Go](http://golang.org/dl/) installed)
  ```
  $ go get github.com/cloudfoundry/cli
  $ go get github.com/krujos/usagereport-plugin
  $ cd $GOPATH/src/github.com/krujos/usagereport-plugin
  $ go build
  $ cf install-plugin usagereport-plugin
  ```
