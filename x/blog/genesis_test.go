package blog_test

import (
	"testing"

	keepertest "github.com/mitchellnel/ibc-planet-nel/testutil/keeper"
	"github.com/mitchellnel/ibc-planet-nel/testutil/nullify"
	"github.com/mitchellnel/ibc-planet-nel/x/blog"
	"github.com/mitchellnel/ibc-planet-nel/x/blog/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BlogKeeper(t)
	blog.InitGenesis(ctx, *k, genesisState)
	got := blog.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
