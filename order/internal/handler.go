package internal

import (
	"context"
	"github.com/cv65kr/order/api/common/v1"
	pb_order "github.com/cv65kr/order/api/orderservice/v1"
)

var (
	_ pb_order.OrderServiceServer = (*WorkflowHandler)(nil)
)

func (h WorkflowHandler) CreateOrder(ctx context.Context, request *pb_order.CreateOrderRequest) (*pb_order.CreateOrderResponse, error) {

	workflowId := h.TriggerWorkflow(request.Customer)

	return &pb_order.CreateOrderResponse{
		Id: workflowId,
		Customer: &common.Customer{
			FirstName: request.Customer.FirstName,
			LastName:  request.Customer.LastName,
		},
	}, nil
}
