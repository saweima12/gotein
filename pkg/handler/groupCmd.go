package handler

import (
	"gotein/cfg"
	"gotein/data"
	"gotein/data/models"
	"gotein/libs/cjson"
	"gotein/logger"
	"gotein/pkg/tgbot"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (bo *BotHandler) listenCommand(m *tgbot.MsgHelper, arg string) {

	reply := m.ReplyToMessage
	if reply == nil {
		return
	}

	if m.From.ID != bo.cfg.OwnerID {
		return
	}

	targetId := strconv.Itoa(int(reply.From.ID))
	err := bo.userServ.AddUser(data.DOC_LISTEN, targetId)
	if err != nil {
		logger.Errorf("Error: ListenUser failed, err: %v", err)
		msg := tu.Message(m.Chat.ChatID(), cfg.GetText(cfg.ERR_INTERNAL))
		bo.API().SendMessage(msg)
		return
	}

	helper := tgbot.GetMsgHelper(reply)
	logger.Infof("[%s] Listen success -- ", helper.FullName())

	msg := tu.Messagef(m.Chat.ChatID(), cfg.GetText(cfg.LISTEN_SUCCESS), m.FullName())
	bo.API().SendMessage(msg)

}

func (bo *BotHandler) addkeywordCommand(m *tgbot.MsgHelper, arg string) {

	reply := m.ReplyToMessage
	if reply == nil {
		return
	}

	// Check permission.
	if !bo.userServ.IsListenUser(m.FromIdStr()) {
		return
	}

	helper := tgbot.GetMsgHelper(reply)
	fUid, fId, contentType := handleFileIDs(helper)

	doc := &models.MediaDoc{
		Uid:       fUid,
		FileID:    fId,
		MediaType: contentType,
	}
	nKeywords := strings.Split(arg, ",")
	doc.Keywords = nKeywords

	if ok := bo.mediaServ.HasMedia(data.INDEX_CACHED, fUid); !ok {
		bo.cacheData(helper, doc)
		if err := bo.mediaServ.InsertMedia(doc); err != nil {
			return
		}
	} else {
		if rtn := bo.mediaServ.AddKeyword(doc); rtn == nil {
			return
		}
	}
	kStr, _ := cjson.MarshalToString(nKeywords)

	logger.Infof("[%s] AddKeyword for %s input: %s", m.FullName(), doc.Uid, kStr)
	msg := tu.Messagef(m.Chat.ChatID(), cfg.GetText(cfg.GROUP_ADD), kStr).
		WithReplyParameters(&telego.ReplyParameters{
			MessageID: reply.MessageID,
		})
	bo.API().SendMessage(msg)

}
