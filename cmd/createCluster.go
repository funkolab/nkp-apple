/*
Copyright © 2025 Christophe Jauffret <reg-github@geo6.net>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// createClusterCmd represents the createCluster command
var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create a Kubernetes cluster, one of [aks, aws, azure, eks, gcp, nutanix, preprovisioned, vsphere]",
	Run: func(cmd *cobra.Command, args []string) {
		selfManaged := false

		// check for --self-managed flag
		for i, arg := range args {
			if arg == "--self-managed" {
				selfManaged = true
				args = append(args[:i], args[i+1:]...)
				break
			}
		}

		// Search cluster-name in args
		clusterName := ""
		for i, arg := range args {
			if (arg == "--cluster-name" || arg == "-c") && i+1 < len(args) {
				clusterName = args[i+1]
				break
			}
		}

		// failed if no cluster-name provided
		if clusterName == "" {
			cobra.CheckErr(fmt.Errorf("cluster-name is required"))
		}

		// create bootstrap cluster if self-managed
		if selfManaged {
			err := createBootstrap()
			cobra.CheckErr(err)

		}

		// create target cluster
		cmdArgs := append([]string{"create", "cluster"}, args...)
		myCmd := exec.Command("nkp", cmdArgs...)
		if err := runCommand(myCmd, true); err != nil {
			cobra.CheckErr(err)
		}

		// retrieve kubeconfig for target cluster
		myCmd = exec.Command("nkp", "get", "kubeconfig", "-c", clusterName)
		output, err := myCmd.Output()
		if err != nil {
			cobra.CheckErr(fmt.Errorf("failed to get kubeconfig: %w", err))
		}

		kubeconfigPath := filepath.Join(".", fmt.Sprintf("%s.conf", clusterName))
		if err := os.WriteFile(kubeconfigPath, output, 0600); err != nil {
			cobra.CheckErr(err)
		}

		// Prepare self-managed cluster
		if selfManaged {

			// Install CAPI components on target cluster
			myCmd = exec.Command("nkp", "create", "capi-components", "--kubeconfig", kubeconfigPath)
			if err := runCommand(myCmd, true); err != nil {
				cobra.CheckErr(err)
			}

			// Move CAPI objects to target cluster
			myCmd = exec.Command("nkp", "move", "capi-resources", "--to-kubeconfig", kubeconfigPath)
			if err := runCommand(myCmd, true); err != nil {
				cobra.CheckErr(err)
			}

			// Delete bootstrap cluster
			err := deleteBootstrap()
			cobra.CheckErr(err)
		}

	},
}

func init() {
	createCmd.AddCommand(createClusterCmd)

	createClusterCmd.Flags().SetInterspersed(false)
}
