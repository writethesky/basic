package service

import (
	"basic/internal"
	tokenV1 "basic/pb/token/v1"
	userV1 "basic/pb/user/v1"
	"basic/repository/dao"
	"basic/utility"
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenService struct {
	tokenV1.UnimplementedTokenServiceServer
	userServiceClient userV1.UserServiceClient
}

func NewTokenService(doubles ...interface{}) *TokenService {
	tokenService := new(TokenService)
	tokenService.userServiceClient = internal.NewUserServiceClient()
	for _, double := range doubles {
		if userServiceClientDouble, isOk := double.(userV1.UserServiceClient); isOk {
			tokenService.userServiceClient = userServiceClientDouble
		}
	}
	return tokenService
}

const (
	tokenLength       = 24
	tokenExpire       = 7 * 12 * time.Hour
	tokenGenerateType = utility.GenerateTypeNumber | utility.GenerateTypeUpperLetter | utility.GenerateTypeLowerLetter
)

func (receiver TokenService) Apply(ctx context.Context, in *tokenV1.ApplyRequest) (*tokenV1.ApplyResponse, error) {

	if in.Username == "" || in.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username and password is required")
	}
	// check username and password is right
	res, err := receiver.userServiceClient.CheckPassword(ctx, &userV1.CheckPasswordRequest{
		Username: in.Username,
		Password: in.Password,
	})
	if nil != err {
		return nil, err
	}
	if !res.IsOk {
		return nil, status.Errorf(codes.InvalidArgument, "password not right")
	}

	expire := time.Now().Add(tokenExpire)
	tokenString := receiver.generateTokenString(int(res.UserId), in.Username, tokenExpire)
	return &tokenV1.ApplyResponse{
		Token:  tokenString,
		Expire: timestamppb.New(expire),
	}, nil
}

func (receiver TokenService) Parse(_ context.Context, in *tokenV1.ParseRequest) (*tokenV1.ParseResponse, error) {
	if "" == in.Token {
		return nil, status.Errorf(codes.InvalidArgument, "token is required")
	}
	userID, username, expire, err := dao.GetToken(in.Token)
	if nil != err {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}
	return &tokenV1.ParseResponse{
		UserId:   int64(userID),
		Username: username,
		Expire:   timestamppb.New(expire),
	}, nil

}

func (receiver TokenService) generateTokenString(userID int, username string, expire time.Duration) string {
	tokenString := utility.GenerateRandomString(tokenLength, tokenGenerateType)
	_, _, _, err := dao.GetToken(tokenString)

	if nil == err {
		return receiver.generateTokenString(userID, username, expire)
	}
	_ = dao.SaveToken(tokenString, expire, userID, username)
	return tokenString
}
