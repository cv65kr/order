package internal

import (
	"context"
	"github.com/cv65kr/order/api/common/v1"
	pb_order "github.com/cv65kr/order/api/orderservice/v1"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	Handler struct {
		log *zap.Logger
	}
)

var (
	_ pb_order.OrderServiceServer = (*Handler)(nil)
)

func NewHandler(log *zap.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func (h Handler) CreateOrder(ctx context.Context, request *pb_order.CreateOrderRequest) (*pb_order.CreateOrderResponse, error) {
	return &pb_order.CreateOrderResponse{
		Id: uuid.New().String(),
		Customer: &common.Customer{
			FirstName: request.Customer.FirstName,
			LastName:  request.Customer.LastName,
		},
	}, nil
}
