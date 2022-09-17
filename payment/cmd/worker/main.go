package main

import (
	"context"
	"github.com/cv65kr/payment/internal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		internal.Module,
		fx.Invoke(TemporalWorkerHooks),
	).Run()
}

func TemporalWorkerHooks(lifecycle fx.Lifecycle, log *zap.Logger, temporalClient client.Client) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				w := worker.New(temporalClient, "payment-service", worker.Options{})

				// register workflow
				w.RegisterWorkflow(internal.Workflow)
				w.RegisterActivity(internal.CreatePayment)

				err := w.Run(worker.InterruptCh())
				if err != nil {
					log.Fatal("Unable to start worker", zap.Error(err))
				}
				return nil
			},
			OnStop: func(context.Context) error {
				temporalClient.Close()
				return log.Sync()
			},
		},
	)
}
