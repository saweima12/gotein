package services

import (
	"gotein/data"
	"gotein/logger"
	"gotein/pkg/repositories"
	"gotein/utils"
)

type UserServ interface {
	AddUser(doc string, userId string) error
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

func (us *userServ) AddUser(doc string, userId string) error {
	obj, err := us.meiliRepo.GetUserList(doc)
	if err != nil {
		return err
	}

	set := utils.SliceToSet[string](obj.Users)
	set[userId] = struct{}{}
	nlist := utils.SetToSlice[string](set)
	obj.Users = nlist

	err = us.meiliRepo.SetUserList(doc, obj)
	if err != nil {
		logger.Errorf("Error AddUser failed, err: %v", err)
		return err
	}
	return nil
}

func (us *userServ) IsAllowUser(uid string) bool {
	obj, err := us.meiliRepo.GetUserList(data.DOC_WHITELIST)
	if err != nil {
		logger.Errorf("Error CheckUserAllow failed err: %v", err)
		return false
	}
	usersSet := utils.SliceToSet[string](obj.Users)
	if _, ok := usersSet[uid]; !ok {
		return false
	}

	return true
}

func (us *userServ) IsListenUser(uid string) bool {
	obj, err := us.meiliRepo.GetUserList(data.DOC_LISTEN)
	if err != nil {
		logger.Errorf("Error CheckUserListen failed err: %v", err)
		return false
	}
	usersSet := utils.SliceToSet[string](obj.Users)
	if _, ok := usersSet[uid]; !ok {
		return false
	}

	return true
}
