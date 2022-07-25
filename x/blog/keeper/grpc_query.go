package keeper

import (
	"github.com/mitchellnel/ibc-planet-nel/x/blog/types"
)

var _ types.QueryServer = Keeper{}
