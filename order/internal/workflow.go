package internal

import (
	"context"
	"github.com/cv65kr/order/api/common/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
	"log"
	"time"
)

type (
	WorkflowHandler struct {
		log            *zap.Logger
		temporalClient client.Client
	}
)

func NewWorkflowHandler(log *zap.Logger, temporalClient client.Client) *WorkflowHandler {
	return &WorkflowHandler{
		log:            log,
		temporalClient: temporalClient,
	}
}

func (h *WorkflowHandler) TriggerWorkflow(customer *common.Customer) string {
	workflowOptions := client.StartWorkflowOptions{
		ID:        "create-order-workflow",
		TaskQueue: "order-service",
	}

	we, err := h.temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, Workflow, customer)
	if err != nil {
		log.Fatal("Unable to execute workflow", zap.Error(err))
	}

	return we.GetID()
}

func Workflow(ctx workflow.Context, customer *common.Customer) error {
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

	err := workflow.ExecuteActivity(ctx, SaveOrder, customer).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Call payment service
	options = workflow.ActivityOptions{
		TaskQueue:           "payment-service",
		StartToCloseTimeout: 5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	err = workflow.ExecuteActivity(ctx, "CreatePayment").Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
