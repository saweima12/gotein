package tgbot

import "github.com/mymmrac/telego"

func GetMsgHelper(m *telego.Message) *MsgHelper {
	return &MsgHelper{
		m,
	}
}

type MsgHelper struct {
	*telego.Message
}
