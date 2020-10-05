package views

import (
	"fmt"

	"github.com/alienvspredator/simple-tgbot/internal/metrics/tgbot"
	"github.com/alienvspredator/simple-tgbot/internal/observability/common"
	"go.opencensus.io/stats/view"
)

func Register() error {
	if err := common.RegisterViews(); err != nil {
		return fmt.Errorf("register common metrics: %w", err)
	}
	if err := view.Register(tgbot.Views...); err != nil {
		return fmt.Errorf("register tgbot metrics: %w", err)
	}

	return nil
}
