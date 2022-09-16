package internal

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(HandlerProvider),
)

func HandlerProvider(
	logger *zap.Logger,
) *Handler {
	return NewHandler(
		logger,
	)
}
