package metricsware

import (
	"context"

	"github.com/alienvspredator/simple-tgbot/internal/metrics/tgbot"
	"go.opencensus.io/stats"
)

func (m Middleware) RecordFailedReply(ctx context.Context) {
	stats.Record(ctx, tgbot.FailedReply.M(1))
}

func (m Middleware) RecordSendFailure(ctx context.Context) {
	stats.Record(ctx, tgbot.SendFailed.M(1))
}

func (m Middleware) RecordMessageSaveFailure(ctx context.Context) {
	stats.Record(ctx, tgbot.MessageSaveFailed.M(1))
}

func (m Middleware) RecordIncomingMessage(ctx context.Context) {
	stats.Record(ctx, tgbot.IncomingMessage.M(1))
}

func (m Middleware) RecordOutgoingMessage(ctx context.Context) {
	stats.Record(ctx, tgbot.OutgoingMessage.M(1))
}
