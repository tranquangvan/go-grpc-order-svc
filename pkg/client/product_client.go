package client

import (
	"context"
	"fmt"

	"github.com/tranquangvan/go-grpc-order-svc/pkg/pb"
	"google.golang.org/grpc"
)

type ProductServiceClient struct {
	Client pb.ProductServiceClient
}

func InitProductServiceClient(url string) ProductServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
	}
	c := ProductServiceClient{
		Client: pb.NewProductServiceClient(cc),
	}
	return c
}

func (c *ProductServiceClient) FindOne(productId int64) (*pb.FindOneResponse, error) {

	res, err := c.Client.FindOne(context.Background(), &pb.FindOneRequest{
		Id: productId,
	})
	return res, err

}

func (c *ProductServiceClient) DecreaseStock(productId int64, orderId int64, quantity int64) (*pb.DecreaseStockResponse, error) {
	res, err := c.Client.DecreaseStock(context.Background(), &pb.DecreaseStockRequest{
		Id:       productId,
		OrderId:  orderId,
		Quantity: quantity,
	})
	return res, err
}
