package bosh

import (
	"bytes"
	"encoding/json"
)

// Instance represents a vm deployed by BOSH
type Instance struct {
	Name  string
	IP    string
	State string
}

// Instances returns the list of Concourse VMs
func (client *Client) Instances() ([]Instance, error) {
	output := bytes.NewBuffer(nil)

	if err := client.director.RunAuthenticatedCommand(
		output,
		client.stderr,
		"--deployment",
		concourseDeploymentName,
		"instances",
		"--json",
	); err != nil {
		return nil, err
	}

	jsonOutput := struct {
		Tables []struct {
			Rows []struct {
				Instance     string `json:"instance"`
				IPs          string `json:"ips"`
				ProcessState string `json:"process_state"`
			} `json:"Rows"`
		} `json:"Tables"`
	}{}

	if err := json.NewDecoder(output).Decode(&jsonOutput); err != nil {
		return nil, err
	}

	instances := []Instance{}

	for _, table := range jsonOutput.Tables {
		for _, row := range table.Rows {
			instances = append(instances, Instance{
				Name:  row.Instance,
				IP:    row.IPs,
				State: row.ProcessState,
			})
		}
	}

	return instances, nil
}
