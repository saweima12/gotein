package services

import (
	"gotein/cfg"
	"gotein/data/models"
	"gotein/logger"
	"gotein/pkg/repositories"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type MediaServ interface {
	SearchMedia(keyword string) []telego.InlineQueryResult
	HasMedia(indexUid, mId string) bool
	GetMedia(indexUid, mId string) *models.MediaDoc
	InsertUnchecked(item *models.MediaDoc) error
}

func NewMediaServ(
	meili repositories.MeiliRepo,
	cacheServ CacheServ,
	cfg *cfg.Configuration,
) MediaServ {
	return &mediaServ{
		meiliRepo: meili,
		cacheServ: cacheServ,
		domainURL: cfg.DomainURL,
	}
}

type mediaServ struct {
	meiliRepo repositories.MeiliRepo
	cacheServ CacheServ
	domainURL string
}

func (me *mediaServ) HasMedia(indexUid, mId string) bool {
	doc := me.GetMedia(indexUid, mId)
	return doc != nil
}

func (me *mediaServ) GetMedia(indexUid, mId string) *models.MediaDoc {
	m, err := me.meiliRepo.GetMedia(indexUid, mId)
	if err != nil {
		logger.Errorf("Error: GetMedia failed, err: %v", err)
		return nil
	}
	return m
}

func (me *mediaServ) SearchMedia(keyword string) []telego.InlineQueryResult {
	items, err := me.meiliRepo.SearchMedia(keyword)
	if err != nil {
		logger.Errorf("Error: search media failed, err: %v", err)
		return []telego.InlineQueryResult{}
	}

	result := make([]telego.InlineQueryResult, 0, len(items))
	for i := range items {
		media := getInlineMedia(items[i].Uid, items[i].MediaType, items[i].FileID)
		if media == nil {
			continue
		}
		result = append(result, media)
	}
	return result
}

func (me *mediaServ) InsertUnchecked(item *models.MediaDoc) error {

	return nil
}

func getInlineMedia(uid, mediaType, fid string) telego.InlineQueryResult {
	switch mediaType {
	case "sticker":
		return tu.ResultCachedSticker(uid, fid)
	case "animation":
		return tu.ResultCachedMpeg4Gif(uid, fid)
	default:
		return nil
	}
}
