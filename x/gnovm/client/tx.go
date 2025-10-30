package client

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/tm2/pkg/std"

	"github.com/ignite/gnovm/x/gnovm/types"
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
		Use:   "add-package [pkgFolder] [deposit] --from creator",
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

			gnoMod, err := parseGnoMod(filepath.Join(folderPath, gnoModName))
			if err != nil {
				return err
			}

			memPkg, err := gnolang.ReadMemPackage(folderPath, gnoMod.Module, gnolang.MPAnyAll)
			if err != nil {
				return fmt.Errorf("failed to read package")
			}

			deposit, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgAddPackage(creator, sdk.Coins{deposit}, deposit, toPkg(memPkg))

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewCallCmd returns a CLI command handler for creating a MsgCall transaction.
func NewCallCmd(addressCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call [send] [pkgPath] [function] [args] --from caller",
		Args:  cobra.MinimumNArgs(3),
		Short: "Call a package on the GnoVM",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			caller, err := addressCodec.BytesToString(clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			pkgPath := args[1]
			function := args[2]

			msg := types.NewMsgCall(caller, sdk.Coins{amount}, amount, pkgPath, function, args[3:])
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func toPkg(mp *std.MemPackage) *types.Package {
	p := &types.Package{
		Name: mp.Name,
		Path: mp.Path,
	}
	for _, f := range mp.Files {
		p.Files = append(p.Files, &types.File{
			Name: f.Name,
			Body: f.Body,
		})
	}
	return p
}

// NewRunCmd returns a CLI command handler for creating a MsgRun transaction.
func NewRunCmd(addressCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [pkgFolder] [deposit] --from caller",
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

			gnoMod, err := parseGnoMod(filepath.Join(folderPath, gnoModName))
			if err != nil {
				return err
			}

			memPkg, err := gnolang.ReadMemPackage(folderPath, gnoMod.Module, gnolang.MPAnyAll)
			if err != nil {
				return fmt.Errorf("failed to read package")
			}

			deposit, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgRun(caller, sdk.Coins{deposit}, deposit, toPkg(memPkg))

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
