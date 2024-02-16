package services

import (
	"gotein/data"
	"gotein/pkg/repositories"
	"gotein/utils"
)

type UserServ interface {
	AddUser(uid string)
	AddListener(uid string)
	IsAllowUser(uid string) bool
	IsListenUser(uid string) bool
}

func NewUserServ(meiliRepo repositories.MeiliRepo) UserServ {
	return &userServ{
		meiliRepo: meiliRepo,
	}
}

type userServ struct {
	meiliRepo repositories.MeiliRepo
}

func (us *userServ) AddUser(uid string) {
	panic("not implemented") // TODO: Implement
}

func (us *userServ) AddListener(uid string) {
	panic("not implemented") // TODO: Implement
}

func (us *userServ) IsAllowUser(uid string) bool {
	usersList, err := us.meiliRepo.GetUserList(data.DOC_WHITELIST)
	if err != nil {
		return false
	}
	usersSet := utils.SliceToSet[string](usersList)
	if _, ok := usersSet[uid]; !ok {
		return false
	}

	return true
}

func (us *userServ) IsListenUser(uid string) bool {
	usersList, err := us.meiliRepo.GetUserList(data.DOC_LISTEN)
	if err != nil {
		return false
	}
	usersSet := utils.SliceToSet[string](usersList)
	if _, ok := usersSet[uid]; !ok {
		return false
	}

	return true
}
