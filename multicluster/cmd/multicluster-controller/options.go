// Copyright 2021 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	mcsv1alpha1 "antrea.io/antrea/multicluster/apis/multicluster/v1alpha1"
)

type Options struct {
	// The path of configuration file.
	configFile     string
	SelfSignedCert bool
	options        ctrl.Options
	// The Service ClusterIP range used in the member cluster.
	ServiceCIDR string
	// The precedence about which IP (private or public one) of Node is preferred to
	// be used as tunnel endpoint. If not specified, private IP will be chosen.
	GatewayIPPrecedence mcsv1alpha1.Precedence
}

func newOptions() *Options {
	return &Options{
		SelfSignedCert: true,
	}
}

func (o *Options) complete(args []string) error {
	var err error
	o.setDefaults()
	options := ctrl.Options{Scheme: scheme}
	ctrlConfig := &mcsv1alpha1.MultiClusterConfig{}
	if len(o.configFile) > 0 {
		klog.InfoS("Loading config", "file", o.configFile)
		options, err = options.AndFrom(ctrl.ConfigFile().AtPath(o.configFile).OfKind(ctrlConfig))
		if err != nil {
			klog.ErrorS(err, "Failed to load options")
			return fmt.Errorf("failed to load options from configuration file %s", o.configFile)
		}
		o.options = options
		if ctrlConfig.ServiceCIDR != "" {
			if _, _, err := net.ParseCIDR(ctrlConfig.ServiceCIDR); err != nil {
				return fmt.Errorf("failed to parse serviceCIDR, invalid CIDR string %s", ctrlConfig.ServiceCIDR)
			}
		}
		o.ServiceCIDR = ctrlConfig.ServiceCIDR
		o.GatewayIPPrecedence = ctrlConfig.GatewayIPPrecedence
		klog.InfoS("Using config from file", "config", o.options)
	} else {
		klog.InfoS("Using default config", "config", o.options)
	}
	return nil
}

// addFlags adds flags to fs and binds them to options.
func (o *Options) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.configFile, "config", o.configFile, "The path to the configuration file")
}

func (o *Options) setDefaults() {
	o.options = ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     "0",
		Port:                   9443,
		HealthProbeBindAddress: ":8080",
		LeaderElection:         false,
	}
}
