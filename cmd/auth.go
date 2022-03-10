package cmd

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func authAccountCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account [address]",
		Aliases: []string{"acc"},
		Short:   "query an account for its number and sequence or pass no arguement to query default account",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			keyNameOrAddress := ""
			if len(args) == 0 {
				keyNameOrAddress = cl.Config.Key
			} else {
				keyNameOrAddress = args[0]
			}
			address, err := cl.AccountFromKeyOrAddress(keyNameOrAddress)
			if err != nil {
				return err
			}
			req := &authtypes.QueryAccountRequest{Address: cl.MustEncodeAccAddr(address)}
			res, err := authtypes.NewQueryClient(cl).Account(cmd.Context(), req)
			if err != nil {
				return err
			}
			return cl.PrintObject(res)
		},
	}
	return cmd
}

func authAccountsCmd(v *viper.Viper, lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "accounts",
		Aliases: []string{"accs"},
		Short:   "query all accounts on a given chain w/ pagination",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			pr, err := ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &authtypes.QueryAccountsRequest{Pagination: pr}
			res, err := authtypes.NewQueryClient(cl).Accounts(cmd.Context(), req)
			if err != nil {
				return err
			}
			return cl.PrintObject(res)
		},
	}
	return paginationFlags(cmd, v)
}

func authParamsCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"param", "params", "p"},
		Short:   "query the current auth parameters",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			res, err := authtypes.NewQueryClient(cl).Params(cmd.Context(), &authtypes.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return cl.PrintObject(res)
		},
	}
	return cmd
}
