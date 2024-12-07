package services

import (
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type UserService struct {
	store *db.UserStore
}

func ReadAllUsers() []Data.User {
}

func ReadUserByUserID(userid int) Data.User {
}

func CreateUser(user Data.User) {
}

func DeleteUserByUserID(userid int) {
}

func UpdateUser(userid int, user Data.User) {

}
