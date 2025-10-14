package gnovm

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/ignite/gnovm/x/gnovm/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "AddPackage",
					Use:            "add-package [deposit]",
					Short:          "Send a AddPackage tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "deposit"}},
				},
				{
					RpcMethod: "Call",
					Use:       "call [send] [pkg-path] [function] [args]",
					Short:     "Send a Call tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "send"},
						{ProtoField: "pkg_path"},
						{ProtoField: "function"},
						{ProtoField: "args"},
					},
				},
				{
					RpcMethod: "Run",
					Use:       "run [send] [pkg]",
					Short:     "Send a Run tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "send"},
						{ProtoField: "pkg"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
