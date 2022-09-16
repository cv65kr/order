package internal

import (
	"github.com/cv65kr/order/internal/zapdapter"
	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(ZapLoggerProvider),
	fx.Provide(TemporalClientProvider),
	fx.Provide(WorkflowHandlerProvider),
)

func ZapLoggerProvider() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func TemporalClientProvider(log *zap.Logger) client.Client {
	c, err := client.Dial(client.Options{
		Logger: zapadapter.NewZapAdapter(log),
	})
	if err != nil {
		log.Fatal("unable to create Temporal client", zap.Error(err))
	}
	return c
}

func WorkflowHandlerProvider(log *zap.Logger, temporalClient client.Client) *WorkflowHandler {
	return NewWorkflowHandler(
		log,
		temporalClient,
	)
}
