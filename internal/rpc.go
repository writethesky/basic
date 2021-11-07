package internal

import (
	userV1 "basic/pb/user/v1"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

var userServiceClient userV1.UserServiceClient

func NewUserServiceClient() userV1.UserServiceClient {
	if nil == userServiceClient {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", Config.UserServer.Host, Config.UserServer.Port), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		userServiceClient = userV1.NewUserServiceClient(conn)
	}
	return userServiceClient
}
