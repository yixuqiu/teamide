package service

import (
	"base"
)

func LoginByAccount(account string, password string) (user *base.UserEntity, err error) {
	user, err = UserGetByAccount(account)

	if err != nil {
		return
	}
	return
}
