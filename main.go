package main

import (
	"basic/internal"
	tokenV1 "basic/pb/token/v1"
	userV1 "basic/pb/user/v1"
	"basic/repository/entity"
	"basic/service"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

var wg sync.WaitGroup

type RegisterType uint

func (registerType RegisterType) String() string {
	switch registerType {
	case registerTokenServiceServer:
		return "token service"
	case registerUserServiceServer:
		return "user service"
	}
	return ""
}

const (
	registerTokenServiceServer RegisterType = iota
	registerUserServiceServer
)

func main() {
	wg.Add(1)
	go runTokenServer()
	go runUserServer()
	migrate()
	wg.Wait()
}

func runTokenServer() {
	runRPCServer(internal.Config.TokenServer.Port, registerTokenServiceServer)
}

func runUserServer() {
	runRPCServer(internal.Config.UserServer.Port, registerUserServiceServer)
}

func runRPCServer(port int, registerType RegisterType) {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening and serving RPC(%s) on :%d", registerType, port)
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	switch registerType {
	case registerTokenServiceServer:
		tokenV1.RegisterTokenServiceServer(grpcServer, service.NewTokenService())
	case registerUserServiceServer:
		userV1.RegisterUserServiceServer(grpcServer, service.NewUserService())
	default:
		panic("unknown server")
	}

	err = grpcServer.Serve(lis)
	if nil != err {
		log.Fatalf("failed to listen: %v", err)
	}
	wg.Done()

}

func migrate() {
	err := internal.DB.AutoMigrate(new(entity.User))
	if nil != err {
		panic(err)
	}
}
