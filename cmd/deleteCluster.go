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
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

// deleteClusterCmd represents the deleteCluster command
var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Delete a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		clusterName, _ := cmd.Flags().GetString("cluster-name")

		spinner := displaySpinner("Deleting cluster " + clusterName)
		myCmd := exec.Command("nkp", "delete", "cluster", "-c", clusterName)
		if err := runCommand(myCmd, false); err != nil {
			close(spinner)
			time.Sleep(500 * time.Millisecond)
			cobra.CheckErr(err)
		}
		close(spinner)
	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteClusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deleteClusterCmd.Flags().StringP("cluster-name", "c", "", "Name used to prefix the cluster and all the created resources.")
	deleteClusterCmd.MarkFlagRequired("cluster-name")
}
