package handler

import (
	"gotein/cfg"
	"gotein/data"
	"gotein/data/models"
	"gotein/libs/cjson"
	"gotein/logger"
	"gotein/pkg/tgbot"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (bo *BotHandler) akCommand(m *tgbot.MsgHelper, arg string) {
	// Check permission.
	if !bo.userServ.IsAllowUser(m.FromIdStr()) {
		return
	}

	reply := m.ReplyToMessage
	if reply == nil {
		msg := tu.Message(m.Chat.ChatID(), cfg.GetText(cfg.ERR_NEED_REPLY))
		bo.API().SendMessage(msg)
		return
	}

	nKeywords := strings.Split(arg, ",")
	helper := tgbot.GetMsgHelper(reply)
	fUid, _, _ := handleFileIDs(helper)

	req := &models.MediaDoc{
		Uid:      fUid,
		Keywords: nKeywords,
	}
	cur := bo.mediaServ.AddKeyword(req)

	kStr, _ := cjson.MarshalToString(cur.Keywords)
	msg := tu.Messagef(m.Chat.ChatID(), cfg.GetText(cfg.AK_SUCCESS), kStr)
	_, err := bo.API().SendMessage(msg)

	logger.Infof("[%s] AddKeyword for %s input: %s", m.FullName(), req.Uid, kStr)
	if err != nil {
		logger.Errorf("Error Addkeyword response send failed, err: %v", err)
	}
}

func (bo *BotHandler) skCommand(m *tgbot.MsgHelper, arg string) {
	// Check permission.
	if !bo.userServ.IsAllowUser(m.FromIdStr()) {
		return
	}

	reply := m.ReplyToMessage
	if reply == nil {
		msg := tu.Message(m.Chat.ChatID(), cfg.GetText(cfg.ERR_NEED_REPLY))
		bo.API().SendMessage(msg)
		return
	}

	// media parameter
	helper := tgbot.GetMsgHelper(reply)
	fUid, _, _ := handleFileIDs(helper)

	// get all keywords.
	nKeywords := strings.Split(arg, ",")
	req := &models.MediaDoc{
		Uid:      fUid,
		Keywords: nKeywords,
	}
	cur := bo.mediaServ.SetKeyword(req)
	kStr, _ := cjson.MarshalToString(cur.Keywords)

	logger.Infof("[%s] SetKeyword for %s input: %s", m.FullName(), req.Uid, kStr)
	msg := tu.Messagef(m.Chat.ChatID(), cfg.GetText(cfg.SK_SUCCESS), kStr).
		WithReplyParameters(&telego.ReplyParameters{
			MessageID: reply.MessageID,
		})
	_, err := bo.API().SendMessage(msg)

	if err != nil {
		logger.Errorf("Error Addkeyword response send failed, err: %v", err)
	}
}

func (bo *BotHandler) rmCommand(m *tgbot.MsgHelper, arg string) {
	// Check permission.
	if !bo.userServ.IsAllowUser(m.FromIdStr()) {
		return
	}

	// Get reply data.
	reply := m.ReplyToMessage
	if reply == nil {
		msg := tu.Message(m.Chat.ChatID(), cfg.GetText(cfg.ERR_NEED_REPLY))
		bo.API().SendMessage(msg)
		return
	}

	helper := tgbot.GetMsgHelper(reply)
	fUid, _, _ := handleFileIDs(helper)

	// Try to delete media.
	err := bo.mediaServ.DeleteMedia(fUid)
	if err != nil {
		msg := tu.Message(m.Chat.ChatID(), cfg.GetText(cfg.ERR_NOT_FOUND))
		bo.API().SendMessage(msg)
		return
	}
	// Send message.
	logger.Infof("[%s] Remove media for %s", m.FullName(), fUid)

	msg := tu.Message(m.Chat.ChatID(), cfg.GetText(cfg.RM_SUCCESS)).
		WithReplyParameters(&telego.ReplyParameters{
			MessageID: reply.MessageID,
		})
	_, err = bo.API().SendMessage(msg)
	if err != nil {
		logger.Errorf("Error Send DeleteMessage failed, err: %v", err)
	}
}

func (bo *BotHandler) registerCommand(m *tgbot.MsgHelper, arg string) {
	if arg != bo.cfg.RegisterKey {
		return
	}

	targetId := m.FromIdStr()
	if err := bo.userServ.AddUser(data.DOC_WHITELIST, targetId); err != nil {
		msg := tu.Message(m.Chat.ChatID(), cfg.GetText(cfg.ERR_INTERNAL))
		bo.API().SendMessage(msg)
		return
	}

	logger.Infof("[%s] Register success -- ", m.FullName())
	msg := tu.Message(m.Chat.ChatID(), cfg.GetText(cfg.REGISTER_SUCCESS))
	bo.API().SendMessage(msg)
}
