package services

import (
	"gotein/data"
	"gotein/data/models"
	"gotein/pkg/repositories"
	"gotein/utils"
)

type MeiliServ interface {
	IsAllowUser(uid string) bool
	IsListenUser(uid string) bool
	InsertAllow(uid string)
	InsertListen(uid string)
	AddMedia(data *models.MediaDoc)
	AddKeyword(data *models.MediaDoc)
	SearchMedia(keyword string) []*models.MediaDoc
}

func NewMeiliServ(meili repositories.MeiliRepo) MeiliServ {
	return &meiliServ{
		meiliRepo: meili,
	}

}

type meiliServ struct {
	meiliRepo repositories.MeiliRepo
}

func (me *meiliServ) IsAllowUser(uid string) bool {
	usersList, err := me.meiliRepo.GetUserList(data.DOC_WHITELIST)
	if err != nil {
		return false
	}
	usersSet := utils.SliceToSet[string](usersList)
	if _, ok := usersSet[uid]; !ok {
		return false
	}

	return true
}

func (me *meiliServ) IsListenUser(uid string) bool {
	usersList, err := me.meiliRepo.GetUserList(data.DOC_LISTEN)
	if err != nil {
		return false
	}
	usersSet := utils.SliceToSet[string](usersList)
	if _, ok := usersSet[uid]; !ok {
		return false
	}

	return true
}

func (me *meiliServ) InsertAllow(uid string) {
	panic("not implemented") // TODO: Implement
}

func (me *meiliServ) InsertListen(uid string) {
	panic("not implemented") // TODO: Implement
}

func (me *meiliServ) AddMedia(data *models.MediaDoc) {
	panic("not implemented") // TODO: Implement
}

func (me *meiliServ) AddKeyword(data *models.MediaDoc) {
	panic("not implemented") // TODO: Implement
}

func (me *meiliServ) SearchMedia(keyword string) []*models.MediaDoc {
	me.meiliRepo.SearchMedia(keyword)
	return nil
}
