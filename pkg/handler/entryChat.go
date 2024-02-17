package handler

import (
	"fmt"
	"gotein/cfg"
	"gotein/data"
	"gotein/data/models"
	"gotein/libs/cjson"
	"gotein/logger"
	"gotein/pkg/tgbot"
	"gotein/utils"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (bo *BotHandler) handleChatMessage(m *tgbot.MsgHelper) {
	// Check command is avaliable.
	if m.Text != "" {
		bo.handleChatCommand(m)
		return
	}

	if !m.IsMedia() {
		return
	}

	// Check permission.
	if !bo.userServ.IsAllowUser(m.FromIdStr()) {
		return
	}

	// Get file id and type.
	fUid, fId, contentType := handleFileIDs(m)
	if fUid == "" || fId == "" || contentType == tgbot.M_UNKNOWN {
		return
	}

	// Try to queryMedia from meilisearch
	if params := bo.queryMedia(m, fUid); params != nil {
		_, err := bo.API().SendMessage(params)
		if err != nil {
			logger.Errorf("Error: SendMessage failed, err: %v", err)
		}
		return
	}

	// If the item hasn't been added yet, store it in the unchecked list.
	doc := models.MediaDoc{
		FileID:    fId,
		Uid:       fUid,
		MediaType: contentType,
		Keywords:  []string{},
	}

	if params := bo.cacheData(m, &doc); params != nil {
		_, err := bo.API().SendMessage(params)
		if err != nil {
			logger.Errorf("Error: cacheData failed, err: %v", err)
		}
	}

	logger.Infof("[%s] AddImage %s", m.FullName(), doc.Uid)
	bo.mediaServ.InsertMedia(&doc)
}

func (bo *BotHandler) queryMedia(m *tgbot.MsgHelper, fUid string) *telego.SendMessageParams {
	// if the item has been checked, return it.
	if doc := bo.mediaServ.GetMedia(data.INDEX_CACHED, fUid); doc != nil {

		text := cfg.GetText(cfg.MEDIA_UNCHECKED)
		if len(doc.Keywords) > 0 {
			keywords, _ := cjson.MarshalToString(doc.Keywords)
			text = fmt.Sprintf(cfg.GetText(cfg.MEDIA_CHECKED), keywords)
		}

		return tu.Message(m.Chat.ChatID(), text).
			WithReplyParameters(&telego.ReplyParameters{
				MessageID: m.MessageID,
			})
	}

	return nil
}

func (bo *BotHandler) cacheData(m *tgbot.MsgHelper, doc *models.MediaDoc) *telego.SendMessageParams {

	var distPath, downloadFid string
	if m.ContentType() == tgbot.M_ANIMATION {
		if m.Animation.Thumbnail != nil {
			distPath = fmt.Sprintf("/gif/%s.webp", doc.Uid)
			downloadFid = m.Animation.Thumbnail.FileID
		} else {
			distPath = fmt.Sprintf("/gif/%s.mp4", doc.Uid)
			downloadFid = m.Animation.FileID
		}
	}

	if m.ContentType() == tgbot.M_STICKER {
		setName := m.Sticker.SetName
		distPath = fmt.Sprintf("/%s/%s.webp", setName, doc.Uid)
		if m.Sticker.Thumbnail != nil {
			downloadFid = m.Sticker.Thumbnail.FileID
		} else {
			downloadFid = m.Sticker.FileID
		}
	}

	if distPath == "" || downloadFid == "" {
		return nil
	}

	doc.FilePath = distPath
	doc.CacheURL = bo.cfg.DomainURL + "/asset" + distPath
	go bo.cacheServ.CacheMedia(downloadFid, distPath)

	return tu.Messagef(m.Chat.ChatID(), cfg.GetText(cfg.MEDIA_NEW)).
		WithReplyParameters(&telego.ReplyParameters{
			MessageID: m.MessageID,
		})
}

func (bo *BotHandler) handleChatCommand(m *tgbot.MsgHelper) {
	cmd, argStr, ok := utils.GetCommand("/", m.Text)
	if !ok {
		return
	}
	bo.chatCmdMap.Notify(cmd, argStr, m)
}

func handleFileIDs(m *tgbot.MsgHelper) (fUid, fId, contentType string) {
	if m.Animation != nil {
		return m.Animation.FileUniqueID, m.Animation.FileID, m.ContentType()
	}

	if m.Sticker != nil {
		return m.Sticker.FileUniqueID, m.Sticker.FileID, m.ContentType()
	}

	return fUid, fId, contentType
}
