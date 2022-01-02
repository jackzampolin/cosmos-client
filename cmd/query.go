package cmd

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// queryCmd represents the keys command
func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "query things about a chain",
	}

	cmd.AddCommand(
		queryBalanceCmd(),
		queryAccountCmd(),
		getAuthQueryCmd(),
		getDistributionQueryCmd(),
	)

	return cmd
}

func queryBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balance [key-or-address]",
		Aliases: []string{"bal"},
		Short:   "query the account balance for a key or address, if none is passed will query the balance of the default account",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := config.GetDefaultClient()
			var (
				keyNameOrAddress = ""
				address          sdk.AccAddress
				err              error
			)
			if len(args) == 0 {
				keyNameOrAddress = cl.Config.Key
			} else {
				keyNameOrAddress = args[0]
			}
			if cl.KeyExists(keyNameOrAddress) {
				cl.Config.Key = keyNameOrAddress

				address, err = cl.GetKeyAddress()
			} else {
				address, err = cl.DecodeBech32AccAddr(keyNameOrAddress)
			}
			if err != nil {
				return err
			}
			balance, err := cl.QueryBalance(address, false)
			if err != nil {
				return err
			}
			bz, err := json.MarshalIndent(balance, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
		},
	}
	return cmd
}

func queryAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account [key-or-address]",
		Aliases: []string{"acc"},
		Short:   "query the account details for a key or address, if none is passed will query the balance of the default account",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := config.GetDefaultClient()
			var (
				keyNameOrAddress = ""
				address          sdk.AccAddress
				err              error
			)
			if len(args) == 0 {
				keyNameOrAddress = cl.Config.Key
			} else {
				keyNameOrAddress = args[0]
			}
			if cl.KeyExists(keyNameOrAddress) {
				cl.Config.Key = keyNameOrAddress
				address, err = cl.GetKeyAddress()
			} else {
				address, err = cl.DecodeBech32AccAddr(keyNameOrAddress)
			}
			if err != nil {
				return err
			}
			account, err := cl.QueryAccount(address)
			if err != nil {
				return err
			}
			bz, err := cl.Codec.Marshaler.MarshalInterfaceJSON(account)
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
		},
	}
	return cmd
}
