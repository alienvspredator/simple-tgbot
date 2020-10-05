package tgbot

import (
	"github.com/alienvspredator/simple-tgbot/internal/metrics"
	"go.opencensus.io/stats"
)

var (
	tgbotMetricsPrefix = metrics.MetricRoot + "tgbot/"

	FailedReply = stats.Int64(
		tgbotMetricsPrefix+"reply_failed",
		"Replies with error messages", stats.UnitDimensionless,
	)

	SendFailed = stats.Int64(
		tgbotMetricsPrefix+"send_failed",
		"Errors when sending replies", stats.UnitDimensionless,
	)

	MessageSaveFailed = stats.Int64(
		tgbotMetricsPrefix+"message_save_failed",
		"Errors when saving message", stats.UnitDimensionless,
	)

	IncomingMessage = stats.Int64(
		tgbotMetricsPrefix+"incoming_message",
		"Incoming messages", stats.UnitDimensionless,
	)

	OutgoingMessage = stats.Int64(
		tgbotMetricsPrefix+"outgoing_message",
		"Outgoing messages", stats.UnitDimensionless,
	)
)
