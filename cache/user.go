package cache

import (
	"fmt"
	"time"
)

func wrapUser(username string) string {
	return fmt.Sprintf("user:%s", username)
}

func (rep *MyRedis) SetToken(username, token string) (err error) {
	err = rep.Set(wrapUser(username), token, 2*time.Hour).Err()
	return
}
