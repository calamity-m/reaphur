package cmd

import (
	"github.com/calamity-m/reaphur/central"
	"github.com/calamity-m/reaphur/gw"
	"github.com/calamity-m/reaphur/pkg/bindings"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "reaphur",
	Short: "reaphur short",
	Long:  "reaphur long",
}

func Execute() error {
	RootCommand.PersistentFlags().BoolVarP(&bindings.Debug, "debug", "d", false, "Force debug")

	RootCommand.AddCommand(central.CentralCommand)
	RootCommand.AddCommand(gw.GRPCGatewayCommand)

	return RootCommand.Execute()
}
