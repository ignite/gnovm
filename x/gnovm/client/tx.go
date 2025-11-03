package client

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/pkg/gnomod"

	"github.com/ignite/gnovm/x/gnovm/types"
)

const (
	flagSend = "send"
)

// NewTxCmd returns a root CLI command handler for gnovm transaction commands with a better UX than with AutoCLI.
func NewTxCmd(addressCodec address.Codec) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "GnoVM transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	rootCmd.AddCommand(
		NewAddPackageCmd(addressCodec),
		NewCallCmd(addressCodec),
		NewRunCmd(addressCodec),
	)

	return rootCmd
}

// NewAddPackageCmd returns a CLI command handler for creating a MsgAddPackage transaction.
func NewAddPackageCmd(addressCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-package [pkgFolder] [deposit] --send [coins] --from creator",
		Args:  cobra.ExactArgs(2),
		Short: "Add a new package to the GnoVM",
		Long:  "Add a new package to the GnoVM. Currently only one package can be added at a time.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator, err := addressCodec.BytesToString(clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			folderPath, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			gnoMod, err := gnomod.ParseDir(folderPath)
			if err != nil {
				return err
			}

			memPkg, err := gnolang.ReadMemPackage(folderPath, gnoMod.Module, gnolang.MPAnyAll)
			if err != nil {
				return fmt.Errorf("failed to read package")
			}

			pkgJSON, err := json.Marshal(memPkg)
			if err != nil {
				return fmt.Errorf("failed to marshal package: %v", err)
			}

			deposit, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			sendStr, err := cmd.Flags().GetString(flagSend)
			if err != nil {
				return err
			}

			send, err := sdk.ParseCoinsNormalized(sendStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddPackage(creator, send, deposit, pkgJSON)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagSend, "", "Coins to send along with the package")

	return cmd
}

// NewCallCmd returns a CLI command handler for creating a MsgCall transaction.
func NewCallCmd(addressCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call [send] [pkgPath] [function] [args] --deposit [coins] --from caller",
		Args:  cobra.MinimumNArgs(3),
		Short: "Call a package on the GnoVM",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			send, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			caller, err := addressCodec.BytesToString(clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString("deposit")
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			pkgPath := args[1]
			function := args[2]

			msg := types.NewMsgCall(caller, send, deposit, pkgPath, function, args[3:])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("deposit", "", "Coins to deposit with the call")

	return cmd
}

// NewRunCmd returns a CLI command handler for creating a MsgRun transaction.
func NewRunCmd(addressCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [pkgFolder] [deposit] --send [coins] --from caller",
		Args:  cobra.ExactArgs(2),
		Short: "Run a tx on the GnoVM",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			caller, err := addressCodec.BytesToString(clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			folderPath, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			gnoMod, err := gnomod.ParseDir(folderPath)
			if err != nil {
				return err
			}

			memPkg, err := gnolang.ReadMemPackage(folderPath, gnoMod.Module, gnolang.MPAnyAll)
			if err != nil {
				return fmt.Errorf("failed to read package")
			}

			pkgJSON, err := json.Marshal(memPkg)
			if err != nil {
				return fmt.Errorf("failed to marshal package: %v", err)
			}

			deposit, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			sendStr, err := cmd.Flags().GetString(flagSend)
			if err != nil {
				return err
			}

			send, err := sdk.ParseCoinsNormalized(sendStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgRun(caller, send, deposit, pkgJSON)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagSend, "", "Coins to send with the run")

	return cmd
}
