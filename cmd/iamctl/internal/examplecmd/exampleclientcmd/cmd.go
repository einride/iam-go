package exampleclientcmd

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "client",
	Short: "IAM example client",
}

func init() {
	// shippers
	Command.AddCommand(createShipperCommand)
	Command.AddCommand(deleteShipperCommand)
	Command.AddCommand(getShipperCommand)
	Command.AddCommand(listShippersCommand)
	Command.AddCommand(updateShipperCommand)
	// sites
	Command.AddCommand(batchGetSitesCommand)
	Command.AddCommand(createSiteCommand)
	Command.AddCommand(deleteSiteCommand)
	Command.AddCommand(getSiteCommand)
	Command.AddCommand(listSitesCommand)
	Command.AddCommand(searchSitesCommand)
	Command.AddCommand(updateSiteCommand)
	// shipments
	Command.AddCommand(batchGetShipmentsCommand)
	Command.AddCommand(createShipmentCommand)
	Command.AddCommand(deleteShipmentCommand)
	Command.AddCommand(getShipmentCommand)
	Command.AddCommand(listShipmentsCommand)
	Command.AddCommand(updateShipmentCommand)
}
