package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/gnovm/x/gnovm/keeper"
	"github.com/ignite/gnovm/x/gnovm/types"
)

func TestParamsQuery(t *testing.T) {
	f := initFixture(t)

	q := keeper.NewQueryServerImpl(&f.keeper)
	params := types.DefaultParams()
	require.NoError(t, f.keeper.Params.Set(f.ctx, params))

	response, err := q.Params(f.ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
