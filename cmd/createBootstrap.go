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
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// createBootstrapCmd represents the createBootstrap command
var createBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Create bootstrap cluster using Apple container",
	Run: func(cmd *cobra.Command, args []string) {
		err := createBootstrap()
		cobra.CheckErr(err)
	},
}

func init() {
	createCmd.AddCommand(createBootstrapCmd)
}

func createBootstrap() error {

	// Find nkp version
	spinner := DisplaySpinner("Checking nkp version...")
	versionCmd := exec.Command("nkp", "version")
	versionOutput, err := versionCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get nkp version: %v\n", err)
		os.Exit(1)
	}

	// Extract version from the line starting with "nkp:"
	nkpVersion := ""
	lines := strings.Split(string(versionOutput), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "nkp:") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				nkpVersion = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	close(spinner)
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("   ... %s ...\n", nkpVersion)

	// Test if bootstrap cluster already exists
	myCmd := exec.Command("container", "exec", nodeName, "hostname")
	if err := runCommand(myCmd, false); err == nil {
		cobra.CheckErr(fmt.Errorf("bootstrap cluster already exists"))
	}

	// Run container
	myCmd = exec.Command("container", "run",
		"-d",
		"--name", nodeName,
		"-m", "8G", "--disable-progress-updates",
		"-e", "KUBECONFIG=/etc/kubernetes/admin.conf",
		"-p", "127.0.0.1:6443:6443",
		fmt.Sprintf(nodeImage, nkpVersion),
	)

	spinner = DisplaySpinner("Creating a bootstrap cluster")
	if err := runCommand(myCmd, false); err != nil {
		close(spinner)
		return fmt.Errorf("failed to run container: %w", err)
	}

	// Set sysctl
	myCmd = exec.Command("container", "exec", nodeName, "sysctl", "-w", "net.ipv4.ip_forward=1")
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("failed to set sysctl: %w", err)
	}

	// Initialize cluster
	myCmd = exec.Command("container", "exec", nodeName, "kubeadm", "init", "--pod-network-cidr="+podCIDR)
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("failed to init cluster: %w", err)
	}

	// Remove taint
	myCmd = exec.Command("container", "exec", nodeName, "kubectl", "taint", "nodes", "--all", "node-role.kubernetes.io/control-plane-")
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("failed to remove taint: %w", err)
	}

	// Set up CNI
	cniCmd := fmt.Sprintf("sed -e 's@{{ .PodSubnet }}@%s@' /kind/manifests/default-cni.yaml | kubectl apply -f -", podCIDR)
	myCmd = exec.Command("container", "exec", nodeName, "sh", "-euc", cniCmd)
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("failed to set up CNI: %w", err)
	}

	// Set up StorageClass
	storageCmd := "cat /kind/manifests/default-storage.yaml | kubectl apply -f -"
	myCmd = exec.Command("container", "exec", nodeName, "sh", "-euc", storageCmd)
	if err := runCommand(myCmd, false); err != nil {
		return fmt.Errorf("failed to set up StorageClass: %w", err)
	}

	close(spinner)
	time.Sleep(500 * time.Millisecond)

	// Set up kubeconfig
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	kubeDir := filepath.Join(homeDir, ".kube")
	if err := os.MkdirAll(kubeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .kube directory: %w", err)
	}

	// Get kubeconfig from container
	myCmd = exec.Command("container", "exec", nodeName, "cat", "/etc/kubernetes/admin.conf")
	output, err := myCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	kubeconfigPath := filepath.Join(kubeDir, "config")
	if err := os.WriteFile(kubeconfigPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	// Install NKP Capi components
	myCmd = exec.Command("nkp", "create", "capi-components", "--kubeconfig", kubeconfigPath)
	if err := runCommand(myCmd, true); err != nil {
		return fmt.Errorf("failed to install NKP Capi components: %w", err)
	}

	return nil
}
