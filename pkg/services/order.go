package services

import (
	"context"
	"net/http"

	"github.com/tranquangvan/go-grpc-order-svc/pkg/client"
	"github.com/tranquangvan/go-grpc-order-svc/pkg/db"
	"github.com/tranquangvan/go-grpc-order-svc/pkg/models"
	"github.com/tranquangvan/go-grpc-order-svc/pkg/pb"
)

type Server struct {
	H          db.Handler
	ProductSvc client.ProductServiceClient
	pb.UnimplementedOrderServiceServer
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	product, err := s.ProductSvc.FindOne(req.ProductId)
	if err != nil || product.Error != "" {
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	} else if product.Data.Stock < req.Quantity {
		return &pb.CreateOrderResponse{
			Status: http.StatusConflict,
			Error:  "Stock too less",
		}, nil
	}
	order := &models.Order{
		Price:     product.Data.Price,
		ProductId: req.ProductId,
		UserId:    req.UserId,
	}
	s.H.DB.Create(&order)

	res, err := s.ProductSvc.DecreaseStock(req.ProductId, order.Id, req.Quantity)
	if err != nil || res.Error != "" {
		s.H.DB.Delete(&models.Order{}, order.Id)
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}
	return &pb.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     order.Id,
	}, nil
}
