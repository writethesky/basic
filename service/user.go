package service

import (
	userV1 "basic/pb/user/v1"
	"basic/repository/dao"
	"basic/repository/entity"
	"basic/utility"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	userV1.UnimplementedUserServiceServer
}

const (
	usernameMinLength = 6
	usernameMaxLength = 20
	passwordMinLength = 10
	passwordMaxLength = 20
)

func NewUserService() *UserService {
	return new(UserService)
}

func (UserService) Register(_ context.Context, in *userV1.RegisterRequest) (_ *emptypb.Empty, err error) {
	if "" == in.Username || "" == in.Password {
		return nil, status.Errorf(codes.InvalidArgument, "username and password are required")
	}

	if isOk := utility.ValidateStringRange(in.Username, usernameMinLength, usernameMaxLength); !isOk {
		return nil, status.Errorf(codes.InvalidArgument, "minimum length of the username is %d and maximum length of the username is %d", usernameMinLength, usernameMaxLength)
	}

	if isOk := utility.ValidateStringRange(in.Password, passwordMinLength, passwordMaxLength); !isOk {
		return nil, status.Errorf(codes.InvalidArgument, "minimum length of the password is %d and maximum length of the password is %d", passwordMinLength, passwordMaxLength)
	}

	if !utility.ValidateStringHasTypes(in.Password, utility.GenerateTypeNumber|utility.GenerateTypeLowerLetter|utility.GenerateTypeUpperLetter|utility.GenerateTypeSpecialCharacter) {
		return nil, status.Errorf(codes.InvalidArgument, "The password must contain digits, uppercase letters, lowercase letters, and special characters")
	}

	// username already existing
	_, err = dao.GetUserIDByUsername(in.Username)
	if nil == err {
		return nil, status.Errorf(codes.AlreadyExists, "the username already exists")
	}

	// save user info
	salt := utility.GenerateRandomString(6, utility.GenerateTypeNumber|utility.GenerateTypeUpperLetter|utility.GenerateTypeLowerLetter)
	password := encryptPassword(in.Password, salt)
	err = dao.CreateUser(entity.User{
		Username: in.Username,
		Password: password,
		Salt:     salt,
	})
	if nil != err {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	return new(emptypb.Empty), nil
}

func encryptPassword(password, salt string) string {
	return utility.EncryptMd5(utility.EncryptMd5(password) + salt)
}

func (UserService) CheckPassword(_ context.Context, in *userV1.CheckPasswordRequest) (*userV1.CheckPasswordResponse, error) {

	userInfo, err := dao.GetUserInfo(in.Username)
	if nil != err {
		return nil, status.Errorf(codes.NotFound, "username not found")
	}

	if encryptPassword(in.Password, userInfo.Salt) == userInfo.Password {
		return &userV1.CheckPasswordResponse{
			IsOk:   true,
			UserId: int64(userInfo.ID),
		}, nil
	}
	return &userV1.CheckPasswordResponse{
		IsOk: false,
	}, nil
}

func (UserService) SetPassword(_ context.Context, in *userV1.SetPasswordRequest) (*emptypb.Empty, error) {
	if 0 == in.UserId {
		return new(emptypb.Empty), status.Errorf(codes.InvalidArgument, "The user ID cannot be 0")
	}

	if isOk := utility.ValidateStringRange(in.Password, passwordMinLength, passwordMaxLength); !isOk {
		return new(emptypb.Empty), status.Errorf(codes.InvalidArgument, "minimum length of the password is %d and maximum length of the password is %d", passwordMinLength, passwordMaxLength)
	}

	if !utility.ValidateStringHasTypes(in.Password, utility.GenerateTypeNumber|utility.GenerateTypeLowerLetter|utility.GenerateTypeUpperLetter|utility.GenerateTypeSpecialCharacter) {
		return new(emptypb.Empty), status.Errorf(codes.InvalidArgument, "The password must contain digits, uppercase letters, lowercase letters, and special characters")
	}

	salt := utility.GenerateRandomString(6, utility.GenerateTypeNumber|utility.GenerateTypeLowerLetter|utility.GenerateTypeUpperLetter)
	in.Password = encryptPassword(in.Password, salt)
	err := dao.SetPassword(in.UserId, in.Password, salt)
	if nil != err {
		return new(emptypb.Empty), status.Errorf(codes.Unavailable, err.Error())
	}
	return new(emptypb.Empty), nil
}
