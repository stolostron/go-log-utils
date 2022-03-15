// Copyright (c) 2022 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

/*

Package zaputil provides some helpful functions for consistently setting up a zap logger, allowing
for command-line configuration and some preferences we've developed. It is built to be easily
adjusted and expanded if there are any additional preferences on specific projects.

The `go-log-utils` repository is part of the `open-cluster-management` community. For more
information, visit: [open-cluster-management.io](https://open-cluster-management.io).

Usual setup to be configurable via command-line, and send klog messages (from other packages your
project might use) in the same format:

	import (
		"flag"
		"fmt"

		ctrl "sigs.k8s.io/controller-runtime"
		"k8s.io/klog/v2"
		"github.com/go-logr/zapr"
		"github.com/stolostron/go-log-utils/zaputil"
	)

	func main() {
		zflags := zaputil.NewFlagConfig()
		zflags.Bind(flag.CommandLine)

		// ... define your custom flags here...

		flag.Parse()

		ctrlZap, err := zflags.BuildForCtrl()
		if err != nil {
			panic(fmt.Sprintf("Failed to build zap logger for controller: %v", err))
		}

		ctrl.SetLogger(zapr.NewLogger(ctrlZap))

		klogZap, err := zaputil.BuildForKlog(zflags.GetConfig(), flag.CommandLine)
		if err != nil {
			panic(fmt.Sprintf("Failed to build zap logger for klog: %v", err))
		}

		klog.SetLogger(zapr.NewLogger(klogZap).WithName("klog"))

		// ...
	}

*/
package zaputil
