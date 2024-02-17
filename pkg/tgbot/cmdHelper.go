package tgbot

type CommandMap map[string]func(m *MsgHelper, arg string)

func (c CommandMap) Notify(cmd string, arg string, m *MsgHelper) {
	if callback, ok := c[cmd]; ok {
		callback(m, arg)
	}
}
