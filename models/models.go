package models

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Org struct {
	Name        string
	MemoryQuota int
	MemoryUsage int
	Spaces      []Space
}

type Space struct {
	Apps []App
	Name string
}

type App struct {
	Ram       int
	Instances int
	Running   bool
}

type Report struct {
	Orgs []Org
}

func (org *Org) InstancesCount() int {
	instancesCount := 0
	for _, space := range org.Spaces {
		instancesCount += space.InstancesCount()
	}
	return instancesCount
}

func (org *Org) AppsCount() int {
	appsCount := 0
	for _, space := range org.Spaces {
		appsCount += len(space.Apps)
	}
	return appsCount
}

func (space *Space) ConsumedMemory() int {
	consumed := 0
	for _, app := range space.Apps {
		if app.Running {
			consumed += int(app.Instances * app.Ram)
		}
	}
	return consumed
}

func (space *Space) RunningAppsCount() int {
	runningAppsCount := 0
	for _, app := range space.Apps {
		if app.Running {
			runningAppsCount++
		}
	}
	return runningAppsCount
}

func (space *Space) InstancesCount() int {
	instancesCount := 0
	for _, app := range space.Apps {
		instancesCount += int(app.Instances)
	}
	return instancesCount
}

func (space *Space) RunningInstancesCount() int {
	runningInstancesCount := 0
	for _, app := range space.Apps {
		if app.Running {
			runningInstancesCount += app.Instances
		}
	}
	return runningInstancesCount
}

func (report *Report) String() string {
	var response bytes.Buffer

	totalApps := 0
	totalAppsRunning := 0
	totalInstances := 0
	totalInstancesRunning := 0

	for _, org := range report.Orgs {
		response.WriteString(fmt.Sprintf("Org %s is consuming %d MB of %d MB.\n",
			org.Name, org.MemoryUsage, org.MemoryQuota))

		for _, space := range org.Spaces {
			spaceRunningAppsCount := space.RunningAppsCount()
			spaceInstancesCount := space.InstancesCount()
			spaceRunningInstancesCount := space.RunningInstancesCount()
			spaceConsumedMemory := space.ConsumedMemory()

			spaceStoppedInstancesCount := spaceInstancesCount - spaceRunningInstancesCount

			response.WriteString(
				fmt.Sprintf("\tSpace %s is consuming %d MB memory (%d%%) of org quota.\n",
					space.Name, spaceConsumedMemory, (100 * spaceConsumedMemory / org.MemoryQuota)))
			response.WriteString(
				fmt.Sprintf("\t\t%d apps: %d running, %d stopped\n", len(space.Apps),
					spaceRunningAppsCount, len(space.Apps)-spaceRunningAppsCount))
			response.WriteString(
				fmt.Sprintf("\t\t%d instances: %d running, %d stopped\n", spaceInstancesCount,
					spaceRunningInstancesCount, spaceStoppedInstancesCount))

			totalInstancesRunning += spaceRunningInstancesCount
			totalAppsRunning += spaceRunningAppsCount
		}

		totalApps += org.AppsCount()
		totalInstances += org.InstancesCount()
	}

	response.WriteString(
		fmt.Sprintf("You are running %d apps in %d org(s), with a total of %d instances.\n",
			totalApps, len(report.Orgs), totalInstances))
	response.WriteString(
		fmt.Sprintf("\t\t%d apps: %d running, %d stopped\n",
			totalApps, totalAppsRunning, totalApps-totalAppsRunning))
	response.WriteString(
		fmt.Sprintf("\t\t%d instances: %d running, %d stopped\n",
			totalInstances, totalInstancesRunning, totalInstances-totalInstancesRunning))

	return response.String()
}

func (report *Report) CSV() string {
	var rows = [][]string{}
	var csv bytes.Buffer

	var headers = []string{"OrgName", "SpaceName", "SpaceMemoryUsed", "OrgMemoryQuota", "AppsDeployed", "AppsRunning", "AppInstancesDeployed", "AppInstancesRunning"}

	rows = append(rows, headers)

	for _, org := range report.Orgs {
		for _, space := range org.Spaces {
			appsDeployed := len(space.Apps)

			spaceResult := []string{
				org.Name,
				space.Name,
				strconv.Itoa(space.ConsumedMemory()),
				strconv.Itoa(org.MemoryQuota),
				strconv.Itoa(appsDeployed),
				strconv.Itoa(space.RunningAppsCount()),
				strconv.Itoa(space.InstancesCount()),
				strconv.Itoa(space.RunningInstancesCount()),
			}

			rows = append(rows, spaceResult)
		}
	}

	for i := range rows {
		csv.WriteString(strings.Join(rows[i], ", "))
		csv.WriteString("\n")
	}

	return csv.String()
}
