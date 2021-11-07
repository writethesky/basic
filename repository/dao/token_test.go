package dao

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserToRedisTokenString(t *testing.T) {

	//gomonkey.ApplyGlobalVar(internal.RDB, nil)
	expire := time.Now().Add(time.Hour).Unix()
	redisTokenString := userToRedisTokenString(1, "admin", time.Hour)
	assert.Equal(t, fmt.Sprintf("1 admin %d", expire), redisTokenString)
}

func TestParseTokenReturnNotExpire(t *testing.T) {
	//gomonkey.ApplyFunc()
	tokenString := userToRedisTokenString(1, "admin", time.Hour)
	userID, username, expire, err := parseToken(tokenString)
	assert.Nil(t, err)
	assert.NotZero(t, userID)
	assert.NotEmpty(t, username)
	assert.Greater(t, time.Until(expire).Seconds(), 0.0)

}

func TestParseTokenReturnExpired(t *testing.T) {
	tokenString := userToRedisTokenString(1, "admin", 0)
	_, _, _, err := parseToken(tokenString)
	assert.NotNil(t, err)

}
