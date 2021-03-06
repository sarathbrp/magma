/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"errors"
	"fmt"
	"os"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/service/client"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/cloud/go/services/dispatcher/gw_client_apis/service303"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	cmdMetrics := &cobra.Command{
		Use:   "metrics <service> [--gateway-service (--hwid=<hardware-id> | --network=<network-id> --gateway=<gateway-id>)]",
		Short: "Get service metrics",
		Args:  validateMetricsArgs,
		Run:   metricsCmd,
	}

	rootCmd.AddCommand(cmdMetrics)
}

func validateMetricsArgs(cmd *cobra.Command, args []string) error {
	if err := validateGlobalFlags(); err != nil {
		return err
	}
	if err := setHwIdFlag(); err != nil {
		return err
	}
	if len(args) != 1 {
		return errors.New("requires 1 arg")
	}
	if !isGatewayServiceQuery && !isValidService(args[0], services) {
		return fmt.Errorf("service %s is invalid, needs to match one of %v", args[0], services)
	}
	if isGatewayServiceQuery && !isValidGwService(gateway_registry.GwServiceType(args[0]), gwServices) {
		return fmt.Errorf("service %s is invalid, needs to match one of %v", args[0], gwServices)
	}
	return nil
}

func metricsCmd(cmd *cobra.Command, args []string) {
	err := getMetrics(args[0])
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}

func getMetrics(service string) error {
	metrics, err := getMetricsOrGwMetrics(service)
	if err != nil {
		return fmt.Errorf("Failed to GetMetrics for %s: %s", service, err)
	}
	if metrics == nil {
		return fmt.Errorf("Metrics is nil")
	}
	for _, mt := range metrics.Family {
		fmt.Printf("%v:\n", *mt.Name)
		for _, metric := range mt.Metric {
			fmt.Printf("       %v\n", metric)
		}
		fmt.Println()
	}
	return nil
}

func getMetricsOrGwMetrics(service string) (*protos.MetricsContainer, error) {
	if isGatewayServiceQuery {
		return service303.GWService303GetMetrics(gateway_registry.GwServiceType(service), hwId)
	} else {
		return client.Service303GetMetrics(service)
	}
}
