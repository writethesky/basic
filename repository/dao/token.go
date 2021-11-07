package dao

import (
	"basic/internal"
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	tokenSetName = "token"
)

func SaveToken(token string, expire time.Duration, userID int, username string) error {
	tokenValue := userToRedisTokenString(userID, username, expire)
	return internal.RDB.HSet(context.Background(), tokenSetName, token, tokenValue).Err()
}

func GetToken(token string) (userID int, username string, expire time.Time, err error) {
	cleanExpiredToken()

	tokenInfoString, err := internal.RDB.HGet(context.Background(), tokenSetName, token).Result()
	if nil != err {
		return
	}
	return parseToken(tokenInfoString)
}

func cleanExpiredToken() {
	// Clearing expired tokens is triggered with a 1% probability
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) != 0 {
		return
	}

	res, err := internal.RDB.HGetAll(context.Background(), tokenSetName).Result()
	if nil != err {
		log.Println(err)
		return
	}
	for token, tokenInfoString := range res {
		fmt.Println(token, tokenInfoString)
	}
}

func userToRedisTokenString(userID int, username string, expire time.Duration) string {
	return fmt.Sprintf("%d %s %d", userID, username, time.Now().Add(expire).Unix())
}

func parseToken(tokenString string) (userID int, username string, expire time.Time, err error) {
	tokenInfoArray := strings.Split(tokenString, " ")
	if 3 != len(tokenInfoArray) {
		err = errors.New("token not right")
		return
	}
	userID, err = strconv.Atoi(tokenInfoArray[0])
	if nil != err {
		return
	}
	username = tokenInfoArray[1]
	expireInt, err := strconv.Atoi(tokenInfoArray[2])
	if nil != err {
		return
	}
	expire = time.Unix(int64(expireInt), 0)

	if time.Until(expire).Seconds() <= 0 {
		err = errors.New("token has expired")
	}
	return
}
