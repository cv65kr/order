package internal

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
	"time"
)

type (
	WorkflowHandler struct {
		log            *zap.Logger
		temporalClient client.Client
	}

	WorkflowSignal struct {
		Message string
	}
)

func NewWorkflowHandler(log *zap.Logger, temporalClient client.Client) *WorkflowHandler {
	return &WorkflowHandler{
		log:            log,
		temporalClient: temporalClient,
	}
}

func Workflow(ctx workflow.Context) (string, error) {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    3,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, CreatePayment).Get(ctx, nil)
	if err != nil {
		return "", err
	}

	// Send signal to order service
	signal := WorkflowSignal{
		Message: "Payment approved",
	}

	err = workflow.SignalExternalWorkflow(ctx, "create-order-workflow", "", "workflow-signal", signal).Get(ctx, nil)
	if err != nil {
		return "", err
	}

	return "Finished", nil
}
