package common

import (
	"fmt"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

func RegisterViews() error {
	httpViews := []*view.View{
		ochttp.ClientCompletedCount,
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
	}
	if err := view.Register(httpViews...); err != nil {
		return fmt.Errorf("register http views: %w", err)
	}

	return nil
}
