package gnovm

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	gnovmsimulation "github.com/ignite/gnovm/x/gnovm/simulation"
	"github.com/ignite/gnovm/x/gnovm/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	gnovmGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&gnovmGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgAddPackage          = "op_weight_msg_gnovm"
		defaultWeightMsgAddPackage int = 100
	)

	var weightMsgAddPackage int
	simState.AppParams.GetOrGenerate(opWeightMsgAddPackage, &weightMsgAddPackage, nil,
		func(_ *rand.Rand) {
			weightMsgAddPackage = defaultWeightMsgAddPackage
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddPackage,
		gnovmsimulation.SimulateMsgAddPackage(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCall          = "op_weight_msg_gnovm"
		defaultWeightMsgCall int = 100
	)

	var weightMsgCall int
	simState.AppParams.GetOrGenerate(opWeightMsgCall, &weightMsgCall, nil,
		func(_ *rand.Rand) {
			weightMsgCall = defaultWeightMsgCall
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCall,
		gnovmsimulation.SimulateMsgCall(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRun          = "op_weight_msg_gnovm"
		defaultWeightMsgRun int = 100
	)

	var weightMsgRun int
	simState.AppParams.GetOrGenerate(opWeightMsgRun, &weightMsgRun, nil,
		func(_ *rand.Rand) {
			weightMsgRun = defaultWeightMsgRun
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRun,
		gnovmsimulation.SimulateMsgRun(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
