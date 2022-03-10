/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// tendermintCmd represents the tendermint command
func tendermintCmd(v *viper.Viper, lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tendermint",
		Aliases: []string{"tm"},
		Short:   "all tendermint query commands",
	}
	cmd.AddCommand(
		abciInfoCmd(lc),
		abciQueryCmd(lc),
		blockCmd(lc),
		blockByHashCmd(lc),
		blockResultsCmd(lc),
		blockSearchCmd(),
		consensusParamsCmd(lc),
		consensusStateCmd(lc),
		dumpConsensusStateCmd(lc),
		healthCmd(lc),
		netInfoCmd(v, lc),
		numUnconfirmedTxs(v, lc),
		statusCmd(lc),
		queryTxCmd(v, lc),
	)
	return cmd
}

func abciInfoCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "abci-info",
		Aliases: []string{"abcii"},
		Short:   "queries for block height, app name and app hash",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			info, err := cl.RPCClient.ABCIInfo(cmd.Context())
			if err != nil {
				return err
			}

			if err := writeJSON(cmd.OutOrStdout(), info); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func abciQueryCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "abci-query [path] [data] [height]",
		Aliases: []string{"qabci"},
		Short:   "query the abci interface for tendermint directly",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			path := args[0]
			data := []byte(args[1])
			height, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			// TODO: wire up prove
			opts := rpcclient.ABCIQueryOptions{
				Height: height,
				Prove:  false,
			}

			info, err := cl.RPCClient.ABCIQueryWithOptions(cmd.Context(), path, data, opts)
			if err != nil {
				return err
			}
			if err := writeJSON(cmd.OutOrStdout(), info); err != nil {
				return err
			}
			return nil
		},
	}
	// TODO: add prove flag
	return cmd
}

func blockCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		// TODO: make this use a height flag and make height arg optional
		Use:     "block [height]",
		Aliases: []string{"bl"},
		Short:   "query tendermint data for a block at given height",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			block, err := cl.RPCClient.Block(cmd.Context(), &height)
			if err != nil {
				return err
			}
			if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func blockByHashCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "block-by-hash [hash]",
		Aliases: []string{"blhash", "blh"},
		Short:   "query tendermint for a given block by hash",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			h, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}
			block, err := cl.RPCClient.BlockByHash(cmd.Context(), h)
			if err != nil {
				return err
			}
			if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func blockResultsCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "block-results [height]",
		Aliases: []string{"blres"},
		Short:   "query tendermint tx results for a given block by height",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			block, err := cl.RPCClient.BlockResults(cmd.Context(), &height)
			if err != nil {
				return err
			}
			// TODO: figure out how to fix the base64 output here
			if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func blockSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "block-search [query] [page] [per-page]",
		Aliases: []string{"bls", "bs", "blsearch"},
		Short:   "search blocks with given query",
		// TODO: long explaination and example should include example queries
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "TODO")
			return nil
		},
	}
	// TODO: order by flag
	return cmd
}

func consensusParamsCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		// TODO: make this use a height flag and make height arg optional
		Use:     "consensus-params [height]",
		Aliases: []string{"csparams", "cs-params"},
		Short:   "query tendermint consensus params at a given height",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			block, err := cl.RPCClient.ConsensusParams(cmd.Context(), &height)
			if err != nil {
				return err
			}
			// TODO: figure out how to fix the base64 output here
			if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func consensusStateCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		// TODO: add special flag to this for network startup
		// that runs query on timer and shows a progress bar
		// _{*extraCredit*}_
		Use:     "consensus-state",
		Aliases: []string{"csstate", "cs-state"},
		Short:   "query current tendermint consensus state",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			block, err := cl.RPCClient.ConsensusState(cmd.Context())
			if err != nil {
				return err
			}
			if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func dumpConsensusStateCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dump-consensus-state",
		Aliases: []string{"dump-cs", "csdump", "cs-dump", "dumpcs"},
		Short:   "query detailed version of current tendermint consensus state",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			block, err := cl.RPCClient.DumpConsensusState(cmd.Context())
			if err != nil {
				return err
			}
			if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func healthCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "health",
		Aliases: []string{"h", "ok"},
		Short:   "query to see if node server is online",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			block, err := cl.RPCClient.Health(cmd.Context())
			if err != nil {
				return err
			}
			if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func netInfoCmd(v *viper.Viper, lc *lensConfig) *cobra.Command {
	// TODO: add flag for pulling out comma seperated list of peers
	// and also filter out private IPs and other ill formed peers
	// _{*extraCredit*}_
	cmd := &cobra.Command{
		Use:     "net-info",
		Aliases: []string{"ni", "net", "netinfo", "peers"},
		Short:   "query for p2p network connection information",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			peers, err := cmd.Flags().GetBool("peers")
			if err != nil {
				return err
			}
			block, err := cl.RPCClient.NetInfo(cmd.Context())
			if err != nil {
				return err
			}
			if !peers {
				if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
					return err
				}
				return nil
			}
			peersList := make([]string, 0, len(block.Peers))
			for _, peer := range block.Peers {
				url, err := url.Parse(peer.NodeInfo.ListenAddr)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "error parsing addr %q: %v\n", peer.NodeInfo.ListenAddr, err)
					continue
				}
				peersList = append(peersList, fmt.Sprintf("%s@%s:%s", peer.NodeInfo.ID(), peer.RemoteIP, url.Port()))
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(peersList, ","))
			return nil
		},
	}
	return peersFlag(cmd, v)
}

func numUnconfirmedTxs(v *viper.Viper, lc *lensConfig) *cobra.Command {
	// TODO: add example for parsing these txs
	// _{*extraCredit*}_
	cmd := &cobra.Command{
		Use:     "mempool",
		Aliases: []string{"unconfirmed", "mem"},
		Short:   "query for number of unconfirmed txs",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				return err
			}
			block, err := cl.RPCClient.UnconfirmedTxs(cmd.Context(), &limit)
			if err != nil {
				return err
			}
			// for _, txbz := range block.Txs {
			// 	fmt.Printf("%X\n", tmtypes.Tx(txbz).Hash())
			// }
			if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
				return err
			}
			return nil
		},
	}
	return limitFlag(cmd, v)
}

func statusCmd(lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"stat", "s"},
		Short:   "query status of the node",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			block, err := cl.RPCClient.Status(cmd.Context())
			if err != nil {
				return err
			}

			if err := writeJSON(cmd.OutOrStdout(), block); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func queryTxCmd(v *viper.Viper, lc *lensConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx [hash]",
		Short: "query for a transaction by hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cl := lc.config.GetDefaultClient()
			prove, err := cmd.Flags().GetBool("prove")
			if err != nil {
				return err
			}
			h, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}
			block, err := cl.RPCClient.Tx(cmd.Context(), h, prove)
			if err != nil {
				return err
			}
			return cl.PrintObject(block)
		},
	}
	return proveFlag(cmd, v)
}
