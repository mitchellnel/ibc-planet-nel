package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/mitchellnel/ibc-planet-nel/testutil/keeper"
	"github.com/mitchellnel/ibc-planet-nel/x/blog/keeper"
	"github.com/mitchellnel/ibc-planet-nel/x/blog/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.BlogKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
