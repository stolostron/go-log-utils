// Copyright (c) 2022 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

/*
Package zaputil provides some helpful functions for consistently setting up a zap logger, allowing
for command-line configuration and some preferences we've developed. It is built to be easily
adjusted and expanded if there are any additional preferences on specific projects.

The `go-log-utils` repository is part of the `open-cluster-management` community. For more
information, visit: [open-cluster-management.io](https://open-cluster-management.io).

This example setup allows for command-line configuration, easy setup of the controller-runtime
logger, and will make klog messages (from other packages your project might use) appear in the same
format as your zap logs:

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
		klog.InitFlags(flag.CommandLine)

		// ... define your custom flags here...

		flag.Parse()

		ctrlZap, err := zflags.BuildForCtrl()
		if err != nil {
			panic(fmt.Sprintf("Failed to build zap logger for controller: %v", err))
		}

		ctrl.SetLogger(zapr.NewLogger(ctrlZap))
		setupLog := ctrl.Log.WithName("setup")

		klogZap, err := zaputil.BuildForKlog(zflags.GetConfig(), flag.CommandLine)
		if err != nil {
			setupLog.Error(err, "Failed to build zap logger for klog, those logs will not go through zap")
		} else {
			klog.SetLogger(zapr.NewLogger(klogZap).WithName("klog"))
		}

		// ... the rest of your main function ...
	}

Note: for consistent log annotations with the filename, line number, and function name, call the
`WithName` or `WithValues` method in your function - do not just use a global variable for the
logger directly, or the "caller" value will be off by one.
*/
package zaputil
