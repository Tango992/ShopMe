package config

import (
	"crypto/tls"
	"log"
	"os"
	"shopping/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func InitGrpc() (*grpc.ClientConn, pb.PaymentClient) {
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(os.Getenv("GRPC_URI"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	return conn, pb.NewPaymentClient(conn)
}