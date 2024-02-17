package services

import (
	"gotein/cfg"
	"gotein/data"
	"gotein/data/models"
	"gotein/logger"
	"gotein/pkg/repositories"
	"gotein/utils"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type MediaServ interface {
	SearchMedia(keyword string) []telego.InlineQueryResult
	HasMedia(indexUid, mId string) bool
	GetMedia(indexUid, mId string) *models.MediaDoc
	DeleteMedia(uid string) error
	InsertMedia(item *models.MediaDoc) error
	AddKeyword(item *models.MediaDoc) *models.MediaDoc
	SetKeyword(item *models.MediaDoc) *models.MediaDoc
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
		logger.Errorf("Error: GetMedia failed, err: %v", err.Error())
		return nil
	}
	return m
}

func (me *mediaServ) DeleteMedia(uid string) error {
	err := me.meiliRepo.DeleteMedia(uid)
	if err != nil {
		logger.Errorf("Error: DeleteMedia failed, err: %v", err)
		return err
	}
	return nil
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

func (me *mediaServ) InsertMedia(item *models.MediaDoc) error {
	err := me.meiliRepo.InsertMedia(item)
	if err != nil {
		logger.Errorf("Error: InsertMedia failed, err: %v", err)
		return err
	}
	return nil
}

func (me *mediaServ) AddKeyword(item *models.MediaDoc) *models.MediaDoc {
	record, err := me.meiliRepo.GetMedia(data.INDEX_CACHED, item.Uid)
	if err != nil {
		logger.Errorf("Error Addkeyword GetMedia failed, err: %v", err)
		return nil
	}

	// Add keyword to set.
	keywords := utils.SliceToSet[string](record.Keywords)
	for _, val := range item.Keywords {
		keywords[val] = struct{}{}
	}
	record.Keywords = utils.SetToSlice[string](keywords)

	// Replace document.
	err = me.meiliRepo.PutMedia(record)
	if err != nil {
		logger.Errorf("Error Addkeyword put media failed, err: %v", err)
		return nil
	}
	return record
}

func (me *mediaServ) SetKeyword(item *models.MediaDoc) *models.MediaDoc {
	record, err := me.meiliRepo.GetMedia(data.INDEX_CACHED, item.Uid)
	if err != nil {
		logger.Errorf("Error SetKeyword GetMedia failed, err: %v", err)
		return nil
	}

	// Add keyword to set.
	record.Keywords = item.Keywords

	// Replace document.
	err = me.meiliRepo.PutMedia(record)
	if err != nil {
		logger.Errorf("Error SetKeyword put media failed, err: %v", err)
		return nil
	}
	return record
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
