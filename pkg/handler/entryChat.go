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
	// Check content type.
	if m.Text != "" {
		// Check command is avaliable.
		cmd, argStr, ok := utils.GetCommand("/", m.Text)
		if ok {
			bo.handleChatCommand(cmd, argStr)
		}
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

	// if the item has been checked, return it.
	if doc := bo.mediaServ.GetMedia(data.INDEX_CACHED, fUid); doc != nil {
		text := cfg.GetText(cfg.MEDIA_CHECKED)
		keywords, _ := cjson.MarshalToString(doc.Keywords)

		msg := tu.Messagef(m.Chat.ChatID(), text, keywords).
			WithReplyParameters(&telego.ReplyParameters{
				MessageID: m.MessageID,
			})
		_, err := bo.API().SendMessage(msg)
		if err != nil {
			logger.Error(err)
		}
		return
	}

	// if the item has been  added, but not checked, return it.
	if doc := bo.mediaServ.GetMedia(data.INDEX_UNCHECK, fUid); doc != nil {
		text := cfg.GetText(cfg.MEDIA_UNCHECKED)
		msg := tu.Message(m.Chat.ChatID(), text).
			WithReplyParameters(&telego.ReplyParameters{
				MessageID: m.MessageID,
			})
		bo.API().SendMessage(msg)
		return
	}

	// Generate download request.

	// If the item hasn't been added yet, store it in the unchecked list.
	doc := &models.MediaDoc{
		FileID:    fId,
		Uid:       fUid,
		MediaType: contentType,
	}
	fmt.Println(doc)

}

func (bo *BotHandler) handleChatCommand(cmd string, argStr string) {

}

func (bo *BotHandler) cacheData(m *tgbot.MsgHelper, doc *models.MediaDoc) *models.MediaDoc {

	return doc
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
