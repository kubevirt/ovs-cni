/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Red Hat, Inc.
 *
 */

package node

import (
	"fmt"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// RunOnNode performs bash shell command on node
// TODO: Use job with a node affinity instead
func RunOnNode(node string, command string) (string, error) {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("docker ps | grep %s | awk '{ print $1}'", node)).CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("failed to run docker ps error output: %s", string(out)))
	}

	out, err = exec.Command("docker", "exec", string(out[:12]), "ssh.sh", command).CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("failed to run docker exec command error output: %s", string(out)))
	}
	outString := string(out)
	outLines := strings.Split(outString, "\n")
	// first two lines of output indicate that connection was successful
	outStripped := outLines[2:]
	outStrippedString := strings.Join(outStripped, "\n")

	return outStrippedString, err
}

// RemoveOvsBridgeOnNode removes ovs bridge on the node
func RemoveOvsBridgeOnNode(bridgeName string) {
	By("Removing ovs-bridge on the node")
	out, err := RunOnNode("node01", "sudo ovs-vsctl --if-exists del-br "+bridgeName)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to run command on node. stdout: %s", out))
}

// AddOvsBridgeOnNode add ovs bridge on the node
func AddOvsBridgeOnNode(bridgeName string) {
	By("Adding ovs-bridge on the node")
	out, err := RunOnNode("node01", "sudo ovs-vsctl add-br "+bridgeName)
	Expect(err).ToNot(HaveOccurred(), fmt.Sprintf("Failed to run command on node. stdout: %s", out))
}
