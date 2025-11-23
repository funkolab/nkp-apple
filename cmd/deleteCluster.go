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

	"github.com/spf13/cobra"
)

// deleteClusterCmd represents the deleteCluster command
var deleteClusterCmd = &cobra.Command{
	Use:                "cluster",
	Short:              "Delete a Kubernetes cluster",
	DisableFlagParsing: true,
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
			cobra.CheckErr(fmt.Errorf("required flag \"cluster-name\" not set"))
		}

		// Search kubeconfig in args
		kubeconfig := ""
		for i, arg := range args {
			if arg == "--kubeconfig" && i+1 < len(args) {
				kubeconfig = args[i+1]
				break
			}
		}

		// failed if no kubeconfig provided
		if kubeconfig == "" {
			cobra.CheckErr(fmt.Errorf("required flag \"kubeconfig\" not set"))
		}

		// create bootstrap cluster if self-managed and move CAPI resources
		if selfManaged {
			err := createBootstrap()
			cobra.CheckErr(err)

			// move CAPI resources from target to bootstrap
			homeDir, err := os.UserHomeDir()
			cobra.CheckErr(err)
			bootstrapKubeconfig := homeDir + "/.kube/config"
			myCmd := exec.Command("nkp", "move", "capi-resources", "--from-kubeconfig", kubeconfig, "--to-kubeconfig", bootstrapKubeconfig)
			if err := runCommand(myCmd, true); err != nil {
				os.Exit(1)
			}

			// delete target cluster
			myCmd = exec.Command("nkp", "delete", "cluster", "--cluster-name", clusterName, "--kubeconfig", bootstrapKubeconfig)
			if err := runCommand(myCmd, true); err != nil {
				os.Exit(1)
			}

			// delete bootstrap cluster
			deleteBootstrap()

			// delete kubeconfig target cluster
			_ = os.Remove(kubeconfig) // Ignore errors

		} else {
			// delete target cluster
			cmdArgs := append([]string{"delete", "cluster"}, args...)
			myCmd := exec.Command("nkp", cmdArgs...)
			if err := runCommand(myCmd, true); err != nil {
				os.Exit(1)
			}
		}

	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)
}
