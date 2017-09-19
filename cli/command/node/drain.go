package node

import (
	"fmt"
	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"strings"
)

type drainOptions struct {
	nodes []string
}

func newDrainCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt drainOptions

	cmd := &cobra.Command{
		Use:   "drain NODE [NODE...]",
		Short: "Drain the volumes from one or more nodes",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.nodes = args
			return runDrain(storageosCli, opt)
		},
	}

	return cmd
}

func runDrain(storageosCli *command.StorageOSCli, opt drainOptions) error {
	client := storageosCli.Client()
	failed := make([]string, 0, len(opt.nodes))

	for _, nodeID := range opt.nodes {
		n, err := client.Controller(nodeID)
		if err != nil {
			failed = append(failed, nodeID)
			continue
		}

		_, err = client.ControllerUpdate(types.ControllerUpdateOptions{
			ID:          n.ID,
			Name:        n.Name,
			Description: n.Description,
			Labels:      n.Labels,
			Cordon:      n.Cordon,
			Drain:       true,
		})
		if err != nil {
			failed = append(failed, nodeID)
			continue
		}

		fmt.Fprintln(storageosCli.Out(), nodeID)
	}

	if len(failed) > 0 {
		return fmt.Errorf("Failed to drain: %s", strings.Join(failed, ", "))
	}
	return nil
}
