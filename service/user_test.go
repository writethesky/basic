package service

import (
	v1 "basic/pb/user/v1"
	"basic/repository/dao"
	"basic/repository/entity"
	"basic/utility"
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var userService = NewUserService()
var rightPassword, rightUsername string

func init() {
	rightUsername = utility.GenerateRandomString(usernameMaxLength-1, utility.GenerateTypeLowerLetter)
	rightPassword = utility.GenerateRandomString(passwordMaxLength-1, utility.GenerateTypeNumber|utility.GenerateTypeLowerLetter|utility.GenerateTypeUpperLetter|utility.GenerateTypeSpecialCharacter)
}
func TestUserRegisterWithoutUsername(t *testing.T) {
	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Password: rightPassword,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestUserRegisterWithoutPassword(t *testing.T) {
	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: rightUsername,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestUserRegisterWithTooShortUsername(t *testing.T) {
	tooShortUsername := utility.GenerateRandomString(usernameMinLength-1, utility.GenerateTypeLowerLetter)
	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: tooShortUsername,
		Password: rightPassword,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestUserRegisterWithTooLongUsername(t *testing.T) {
	tooLongUsername := utility.GenerateRandomString(usernameMaxLength+1, utility.GenerateTypeLowerLetter)
	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: tooLongUsername,
		Password: rightPassword,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestUserRegisterWithExistingUsername(t *testing.T) {
	patches := gomonkey.ApplyFunc(dao.GetUserIDByUsername, func(username string) (userID int64, err error) {
		return 1, nil
	})
	defer patches.Reset()
	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: rightUsername,
		Password: rightPassword,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err), err)
}

func TestUserRegisterWithTooShortPassword(t *testing.T) {
	tooShortPassword := utility.GenerateRandomString(passwordMinLength-1, utility.GenerateTypeLowerLetter)

	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: rightUsername,
		Password: tooShortPassword,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestUserRegisterWithTooLongPassword(t *testing.T) {
	tooLongPassword := utility.GenerateRandomString(passwordMaxLength+1, utility.GenerateTypeLowerLetter)

	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: rightUsername,
		Password: tooLongPassword,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestUserRegisterWithPasswordDoesNotContainLowerLetter(t *testing.T) {
	tooLongPassword := utility.GenerateRandomString(passwordMaxLength-1, utility.GenerateTypeNumber|utility.GenerateTypeUpperLetter|utility.GenerateTypeSpecialCharacter)

	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: rightUsername,
		Password: tooLongPassword,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestUserRegisterWithPasswordDoesNotContainUpperLetter(t *testing.T) {
	tooLongPassword := utility.GenerateRandomString(passwordMaxLength-1, utility.GenerateTypeNumber|utility.GenerateTypeLowerLetter|utility.GenerateTypeSpecialCharacter)

	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: rightUsername,
		Password: tooLongPassword,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestUserRegisterWithPasswordDoesNotContainSpecialSymbol(t *testing.T) {
	tooLongPassword := utility.GenerateRandomString(passwordMaxLength-1, utility.GenerateTypeNumber|utility.GenerateTypeLowerLetter|utility.GenerateTypeUpperLetter)

	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: rightUsername,
		Password: tooLongPassword,
	})
	assert.NotNil(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err), err)
}

func TestUserRegisterWithRightParams(t *testing.T) {
	patches := gomonkey.ApplyFunc(dao.GetUserIDByUsername, func(username string) (userID int64, err error) {
		return 0, errors.New("")
	}).ApplyFunc(dao.CreateUser, func(user entity.User) (err error) {
		return nil
	})
	defer patches.Reset()
	_, err := userService.Register(context.Background(), &v1.RegisterRequest{
		Username: rightUsername,
		Password: rightPassword,
	})
	assert.Nil(t, err, err)
}

func Test_encryptPassword(t *testing.T) {
	salt := utility.GenerateRandomString(6, utility.GenerateTypeNumber|utility.GenerateTypeUpperLetter|utility.GenerateTypeLowerLetter)
	p1 := encryptPassword(rightPassword, salt)
	p2 := encryptPassword(rightPassword, salt)
	assert.Equal(t, p1, p2)
}

func TestUserService_CheckPassword_exceptOk(t *testing.T) {
	salt := utility.GenerateRandomString(6, utility.GenerateTypeNumber|utility.GenerateTypeUpperLetter|utility.GenerateTypeLowerLetter)
	p := encryptPassword(rightPassword, salt)

	patches := gomonkey.ApplyFunc(dao.GetUserInfo, func(username string) (userInfo entity.User, err error) {
		return entity.User{
			Username: rightUsername,
			Password: p,
			Salt:     salt,
		}, nil
	})
	defer patches.Reset()
	res, err := userService.CheckPassword(context.Background(), &v1.CheckPasswordRequest{
		Username: rightUsername,
		Password: rightPassword,
	})
	assert.Nil(t, err, err)
	assert.True(t, res.IsOk)
}

func TestUserService_CheckPassword_withNotExistUsername(t *testing.T) {

	patches := gomonkey.ApplyFunc(dao.GetUserInfo, func(username string) (userInfo entity.User, err error) {
		return entity.User{}, errors.New("")
	})
	defer patches.Reset()
	res, err := userService.CheckPassword(context.Background(), &v1.CheckPasswordRequest{
		Username: rightUsername,
		Password: rightPassword,
	})
	assert.NotNil(t, err, res)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestUserService_CheckPassword_withFailedPassword(t *testing.T) {
	salt := utility.GenerateRandomString(6, utility.GenerateTypeNumber|utility.GenerateTypeUpperLetter|utility.GenerateTypeLowerLetter)
	p := encryptPassword(rightPassword, salt)

	patches := gomonkey.ApplyFunc(dao.GetUserInfo, func(username string) (userInfo entity.User, err error) {
		return entity.User{
			Username: rightUsername,
			Password: p,
			Salt:     salt,
		}, nil
	})
	defer patches.Reset()
	res, err := userService.CheckPassword(context.Background(), &v1.CheckPasswordRequest{
		Username: rightUsername,
		Password: rightPassword + "1",
	})
	assert.Nil(t, err, err)
	assert.False(t, res.IsOk)
}
