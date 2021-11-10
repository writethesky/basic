package service

import (
	mocks "basic/mock"
	v1 "basic/pb/token/v1"
	userV1 "basic/pb/user/v1"
	"basic/repository/dao"
	"basic/utility"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/agiledragon/gomonkey/v2"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
)

var tokenService = NewTokenService()

func TestTokenApplyWithoutUsernameAndPassword(t *testing.T) {
	res, err := tokenService.Apply(context.Background(), &v1.ApplyRequest{})
	assert.NotNil(t, err, res)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestTokenApplyExceptRight(t *testing.T) {
	userServiceClient := new(mocks.UserServiceClient)
	userServiceClient.On("CheckPassword", mock.Anything, mock.Anything).Return(&userV1.CheckPasswordResponse{
		IsOk:   true,
		UserId: 1,
	}, nil)
	patches := gomonkey.ApplyFunc(dao.GetToken, func(token string) (userID int, username string, expire time.Time, err error) {
		err = errors.New("")
		return
	}).ApplyFunc(dao.SaveToken, func(token string, expire time.Duration, userID int, username string) error {
		return nil
	})

	defer patches.Reset()

	tokenService := NewTokenService(userServiceClient)
	res, err := tokenService.Apply(context.Background(), &v1.ApplyRequest{
		Username: "admin",
		Password: "123456",
	})
	assert.Nil(t, err, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
	assert.Greater(t, time.Until(res.Expire.AsTime()), 10*time.Minute)
}

func TestGenerateTokenStringHappenNameRepeat(t *testing.T) {
	firstToken := utility.GenerateRandomString(tokenLength, tokenGenerateType)
	secondToken := utility.GenerateRandomString(tokenLength, tokenGenerateType)

	outputs := []gomonkey.OutputCell{
		{Values: gomonkey.Params{firstToken}},
		{Values: gomonkey.Params{secondToken}},
	}

	patches := gomonkey.ApplyFuncSeq(utility.GenerateRandomString, outputs).ApplyFunc(dao.GetToken, func(token string) (userID int, username string, expire time.Time, err error) {
		if firstToken == token {
			return
		}
		err = errors.New("")
		return
	}).ApplyFunc(dao.SaveToken, func(token string, expire time.Duration, userID int, username string) error {
		return nil
	})

	defer patches.Reset()

	tokenString := tokenService.generateTokenString(1, "admin", time.Minute)
	assert.Equal(t, secondToken, tokenString)
}

func TestTokenVerifyWithoutToken(t *testing.T) {
	res, err := tokenService.Parse(context.Background(), &v1.ParseRequest{})
	assert.NotNil(t, err, res)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestTokenVerifyWithNotExistOrExpiredToken(t *testing.T) {
	patches := gomonkey.ApplyFunc(dao.GetToken, func(token string) (userID int, username string, expire time.Time, err error) {
		return 0, "", time.Time{}, errors.New("")
	})
	defer patches.Reset()

	res, err := tokenService.Parse(context.Background(), &v1.ParseRequest{
		Token: "abc",
	})
	assert.NotNil(t, err, res)
	assert.Equal(t, codes.NotFound, status.Code(err), err)
}

func TestTokenVerifyWithRightToken(t *testing.T) {
	userID := 1
	username := "admin"
	patches := gomonkey.ApplyFunc(dao.SaveToken, func(token string, expire time.Duration, userID int, username string) error {
		return nil
	}).ApplyFunc(dao.GetToken, func(token string) (userID int, username string, expire time.Time, err error) {
		return 0, "", time.Time{}, errors.New("")
	})
	defer patches.Reset()

	tokenString := tokenService.generateTokenString(userID, username, time.Hour)
	patches.Reset()
	patches.ApplyFunc(dao.GetToken, func(token string) (_ int, _ string, _ time.Time, err error) {
		if tokenString == token {
			return userID, username, time.Time{}, nil
		}
		return 0, "", time.Time{}, errors.New("")
	})
	res, err := tokenService.Parse(context.Background(), &v1.ParseRequest{
		Token: tokenString,
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(userID), res.UserId)
	assert.Equal(t, username, res.Username)

}
