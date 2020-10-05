package tgbot

import (
	"github.com/alienvspredator/simple-tgbot/internal/metrics"
	"go.opencensus.io/stats/view"
)

var (
	Views = []*view.View{
		{
			Name:        metrics.MetricRoot + "failed_reply_count",
			Description: "Total number of replies with errors",
			Measure:     FailedReply,
			Aggregation: view.Sum(),
		},
		{
			Name:        metrics.MetricRoot + "send_failed_count",
			Description: "Total number of message sending errors",
			Measure:     SendFailed,
			Aggregation: view.Sum(),
		},
		{
			Name:        metrics.MetricRoot + "message_save_failed_count",
			Description: "Total number of message saving errors",
			Measure:     MessageSaveFailed,
			Aggregation: view.Sum(),
		},
		{
			Name:        metrics.MetricRoot + "incoming_messages_count",
			Description: "Total number of incoming messages",
			Measure:     IncomingMessage,
			Aggregation: view.Sum(),
		},
		{
			Name:        metrics.MetricRoot + "outgoing_message",
			Description: "Total number of outgoing messages",
			Measure:     OutgoingMessage,
			Aggregation: view.Sum(),
		},
	}
)
