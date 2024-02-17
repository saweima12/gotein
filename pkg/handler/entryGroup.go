package handler

import (
	"gotein/pkg/tgbot"
	"gotein/utils"
)

func (bo *BotHandler) handleGroupMessage(m *tgbot.MsgHelper) {
	if m.Text == "" {
		return
	}
	bo.handleGroupCommand(m)
}

func (bo *BotHandler) handleGroupCommand(m *tgbot.MsgHelper) {
	cmd, argStr, ok := utils.GetCommand("!", m.Text)

	if !ok {
		return
	}

	bo.groupCmdMap.Notify(cmd, argStr, m)
}
