package node

import (
	"fmt"
	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"strings"
)

type leaveOptions struct {
	nodes []string
}

func newLeaveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt leaveOptions

	cmd := &cobra.Command{
		Use:   "leave NODE [NODE...]",
		Short: "Make one or more nodes leave the cluster",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.nodes = args
			return runLeave(storageosCli, opt)
		},
	}

	return cmd
}

func runLeave(storageosCli *command.StorageOSCli, opt leaveOptions) error {
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
			Drain:       n.Drain,
			Health:      "left",
		})
		if err != nil {
			failed = append(failed, nodeID)
			continue
		}

		fmt.Fprintln(storageosCli.Out(), nodeID)
	}

	if len(failed) > 0 {
		return fmt.Errorf("Node failed to leave: %s", strings.Join(failed, ", "))
	}
	return nil
}
