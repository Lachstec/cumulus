package services

import (
	Data "data"
)

func ReadNumOfUsers() int {
	return len(Data.Users)
}

func ReadAllUsers() []Data.User {
	return Data.Users
}

func ReadUserByUserID(userid int) Data.User {
	return Data.Users[userid]
}

func CreateUser(user Data.User) {
	Data.Users = append(Data.Users, user)
}

func DeleteUserByUserID(userid int) {
	Data.Users = append(Data.Users[:userid], Data.Users[userid + 1:]...)
}

func UpdateUser(userid int, user Data.User) {
	switch {
		case user.Name != "":
			Data.Users[userid].Name = user.Name
		case user.Role != "":
			Data.Users[userid].Role = user.Role
	}
}