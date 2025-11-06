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
	flagSend       = "send"
	flagMaxDeposit = "max-deposit"
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
		Use:   "add-package [pkgFolder] --max-deposit [coins] --send [coins] --from creator",
		Args:  cobra.ExactArgs(1),
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

			maxDepositStr, err := cmd.Flags().GetString(flagMaxDeposit)
			if err != nil {
				return err
			}
			maxDeposit, err := sdk.ParseCoinsNormalized(maxDepositStr)
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

			msg := types.NewMsgAddPackage(creator, send, maxDeposit, pkgJSON)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagSend, "", "Coins to send along with the package")
	cmd.Flags().String(flagMaxDeposit, "", "Maximum amount of coins to be spent for the storage fee (if empty the VM will use a default value)")

	return cmd
}

// NewCallCmd returns a CLI command handler for creating a MsgCall transaction.
func NewCallCmd(addressCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call [pkgPath] [function] [args] --max-deposit [coins] --send [coins] --from caller",
		Args:  cobra.MinimumNArgs(2),
		Short: "Call a package on the GnoVM",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			caller, err := addressCodec.BytesToString(clientCtx.GetFromAddress())
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

			maxDepositStr, err := cmd.Flags().GetString(flagMaxDeposit)
			if err != nil {
				return err
			}
			maxDeposit, err := sdk.ParseCoinsNormalized(maxDepositStr)
			if err != nil {
				return err
			}

			pkgPath := args[0]
			function := args[1]

			msg := types.NewMsgCall(caller, send, maxDeposit, pkgPath, function, args[2:])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagSend, "", "Coins to send along with the package")
	cmd.Flags().String(flagMaxDeposit, "", "Maximum amount of coins to be spent for the storage fee (if empty the VM will use a default value)")

	return cmd
}

// NewRunCmd returns a CLI command handler for creating a MsgRun transaction.
func NewRunCmd(addressCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [pkgFolder] --max-deposit [coins] --send [coins] --from caller",
		Args:  cobra.ExactArgs(1),
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

			sendStr, err := cmd.Flags().GetString(flagSend)
			if err != nil {
				return err
			}
			send, err := sdk.ParseCoinsNormalized(sendStr)
			if err != nil {
				return err
			}

			maxDepositStr, err := cmd.Flags().GetString(flagMaxDeposit)
			if err != nil {
				return err
			}
			maxDeposit, err := sdk.ParseCoinsNormalized(maxDepositStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgRun(caller, send, maxDeposit, pkgJSON)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagSend, "", "Coins to send along with the package")
	cmd.Flags().String(flagMaxDeposit, "", "Maximum amount of coins to be spent for the storage fee (if empty the VM will use a default value)")

	return cmd
}
