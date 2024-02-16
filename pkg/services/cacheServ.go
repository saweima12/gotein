package services

import (
	"github.com/mymmrac/telego"
)

type CacheServ interface{}

func NewCacheServ(api *telego.Bot) CacheServ {
	return &cacheServ{
		api: api,
	}
}

type cacheServ struct {
	api *telego.Bot
}

func (ca *cacheServ) DownloadFile() {

}
