package cmd

import "github.com/spf13/cobra"

// TxCommand regesters a new tx command.
func txCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "query things about a chain",
	}

	cmd.AddCommand(bankSendCmd())
	cmd.AddCommand(stakingDelegateCmd())
	cmd.AddCommand(stakingRedelegateCmd())
	cmd.AddCommand(distributionWithdrawRewardsCmd())

	return cmd
}

func bankSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [from] [to] [amount]",
		Short: "send coins from one address to another",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := config.GetDefaultClient()
			var (
				fromAddress sdk.AccAddress
				err         error
			)
			if cl.KeyExists(args[0]) {
				fromAddress, err = cl.GetKeyAddress()
			} else {
				fromAddress, err = cl.DecodeBech32AccAddr(args[0])
			}
			if err != nil {
				return err
			}
			toAddr, err := cl.DecodeBech32AccAddr(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			from, err := cl.EncodeBech32AccAddr(fromAddress)
			if err != nil {
				return err
			}

			to, err := cl.EncodeBech32AccAddr(toAddr)
			if err != nil {
				return err
			}

			msg := &types.MsgSend{
				FromAddress: from,
				ToAddress:   to,
				Amount:      coins,
			}

			res, ok, err := cl.SendMsg(cmd.Context(), msg)
			if err != nil || !ok {
				if res != nil {
					return fmt.Errorf("failed to send coins: code(%d) msg(%s)", res.Code, res.Logs)
				}
				return fmt.Errorf("failed to send coins: err(%w)", err)
			}

			bz, err := cl.Codec.Marshaler.MarshalJSON(res)
			if err != nil {
				return err
			}

			var out = bytes.NewBuffer([]byte{})
			if err := json.Indent(out, bz, "", "  "); err != nil {
				return err
			}
			fmt.Println(out.String())
			return nil

		},
	}
	return cmd
}