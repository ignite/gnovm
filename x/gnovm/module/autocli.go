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
					Short:     "Shows the parameters of the module.",
				},
				{
					RpcMethod:      "Info",
					Use:            "info [pkg-path]",
					Short:          "Query raw info from a package.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pkg_path"}},
				},

				{
					RpcMethod:      "RealmStorage",
					Use:            "realm-storage [pkg-path]",
					Short:          "Query storage from a realm.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pkg_path"}},
				},

				{
					RpcMethod:      "Eval",
					Use:            "eval [pkg-path] [expr]",
					Short:          "Evaluates any expression in readonly mode and returns the results.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pkg_path"}, {ProtoField: "expr"}},
				},

				{
					RpcMethod: "Render",
					Use:       "render [pkg-path] <args>",
					Short:     "Queries Render method on the contract given the package path and arguments.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "pkg_path"},
						{ProtoField: "args", Varargs: true},
					},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:   "UpdateParams",
					GovProposal: true,
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
