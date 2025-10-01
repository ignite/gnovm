package types

import (
	"log/slog"

	"cosmossdk.io/log"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
)

func GnoContextFromSDKContext(ctx sdk.Context) gnosdk.Context {
	runMode := ctx.ExecMode()
	_ = ctx.MultiStore()

	return gnosdk.NewContext(convertExecMode(runMode), nil, nil /* todo */, slog.Default())
}

func SDKContextFromGnoContext(ctx gnosdk.Context) sdk.Context {
	var isCheckTx bool
	if ctx.Mode() == gnosdk.RunTxModeCheck {
		isCheckTx = true
	}

	return sdk.NewContext(nil, cmtproto.Header{}, isCheckTx, log.NewNopLogger())
}

func convertExecMode(execMode sdk.ExecMode) gnosdk.RunTxMode {
	switch execMode {
	case sdk.ExecModeCheck:
		return gnosdk.RunTxModeCheck
	case sdk.ExecModeReCheck:
		return gnosdk.RunTxModeCheck
	case sdk.ExecModeSimulate:
		return gnosdk.RunTxModeSimulate
	case sdk.ExecModePrepareProposal:
		return gnosdk.RunTxModeCheck
	case sdk.ExecModeProcessProposal:
		return gnosdk.RunTxModeDeliver
	case sdk.ExecModeVoteExtension:
		return gnosdk.RunTxModeCheck
	case sdk.ExecModeVerifyVoteExtension:
		return gnosdk.RunTxModeCheck
	case sdk.ExecModeFinalize:
		return gnosdk.RunTxModeDeliver
	}

	return gnosdk.RunTxModeCheck
}
