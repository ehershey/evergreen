package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/cloud"
	"github.com/evergreen-ci/evergreen/cloud/providers"
	"github.com/evergreen-ci/evergreen/db"
	"github.com/evergreen-ci/evergreen/model/distro"
	_ "github.com/evergreen-ci/evergreen/plugin/config"
	"github.com/evergreen-ci/evergreen/util"
	"github.com/tychoish/grip"
	"github.com/tychoish/grip/slogger"
)

var (
	// requestTimeout is the duration to wait until killing
	// active requests and stopping the server.
	requestTimeout = 10 * time.Second
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s handles communication with running tasks and command line tools.\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Supported flags are:\n")
		flag.PrintDefaults()
	}
}

func main() {
	go util.DumpStackOnSIGQUIT(os.Stdout)
	settings := evergreen.GetSettingsOrExit()
	evergreen.SetLogger("/tmp/ernie.log")
	grip.SetName("launch-spot-instances")

	db.SetGlobalSessionProvider(db.SessionFactoryFromConfig(settings))

	distroId := "debian81-build"
	d, err := distro.FindOne(distro.ById(distroId))
	if err != nil {
		evergreen.Logger.Logf(slogger.ERROR, "Failed to find distro '%v': %v", distroId, err)
		os.Exit(1)
	}

	cloudManager, err := providers.GetCloudManager(d.Provider, settings)
	if err != nil {
		evergreen.Logger.Errorf(slogger.ERROR, "Error getting cloud manager for distro: %v", err)
		os.Exit(1)
	}

	hostOptions := cloud.HostOptions{
		UserName: evergreen.User,
		UserHost: false,
	}
	for i := 0; i < 100; i++ {
		newHost, err := cloudManager.SpawnInstance(d, hostOptions)
		if err != nil {
			evergreen.Logger.Errorf(slogger.ERROR, "Error spawning instance #: %d (%v),",
				i, err)
			continue
		}
		evergreen.Logger.Errorf(slogger.ERROR, "Success spawning instance #: %d", i)
		evergreen.Logger.Errorf(slogger.ERROR, "Host ID: %v", newHost.Id)
	}
	os.Exit(0)
}
