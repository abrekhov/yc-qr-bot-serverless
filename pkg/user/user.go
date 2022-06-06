package user

import "time"

type User struct {
	UserID   int64     `dynamo:",hash"`
	Created  time.Time `dynamo:",range"`
	Username string    `dynamo:""`
}
