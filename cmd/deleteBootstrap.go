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
	"time"

	"github.com/spf13/cobra"
)

// deleteBootstrapCmd represents the deleteBootstrap command
var deleteBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Delete bootstrap cluster using Apple container",
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteBootstrap()
		cobra.CheckErr(err)
	},
}

func init() {
	deleteCmd.AddCommand(deleteBootstrapCmd)
}

func deleteBootstrap() error {

	spinner := displaySpinner("Deleting bootstrap cluster")
	// Remove container
	myCmd := exec.Command("container", "rm", "-f", nodeName)
	_ = myCmd.Run() // Ignore errors

	// Remove kubeconfig
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	kubeconfigPath := filepath.Join(homeDir, ".kube", "config")
	_ = os.Remove(kubeconfigPath) // Ignore errors

	close(spinner)
	time.Sleep(500 * time.Millisecond)
	return nil
}
