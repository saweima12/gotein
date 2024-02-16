package tgbot

import (
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
)

const (
	M_ANIMATION = "animation"
	M_AUDIO     = "audio"
	M_STICKER   = "sticker"
	M_DICE      = "dice"
	M_DOCUMENT  = "document"
	M_TEXT      = "text"
	M_PHOTO     = "photo"
	M_VIDEO     = "video"
	M_UNKNOWN   = "unknown"
)

func GetMsgHelper(m *telego.Message) *MsgHelper {
	return &MsgHelper{
		Message: m,
	}
}

type MsgHelper struct {
	*telego.Message
	BotID int64
}

func (in *MsgHelper) FromIdStr() string {
	return strconv.Itoa(int(in.From.ID))
}

func (in *MsgHelper) FullName() string {
	builder := strings.Builder{}
	builder.WriteString(in.From.FirstName)
	builder.WriteString(" ")
	builder.WriteString(in.From.LastName)
	return builder.String()
}

func (in *MsgHelper) IsMedia() bool {
	switch in.ContentType() {
	case M_STICKER:
		return true
	case M_ANIMATION:
		return true
	}
	return false
}

func (in *MsgHelper) ContentType() string {

	if in.Sticker != nil {
		return M_STICKER
	}

	if in.Dice != nil {
		return M_DICE
	}

	if in.Video != nil {
		return M_VIDEO
	}

	if in.Animation != nil {
		return M_ANIMATION
	}

	if in.Photo != nil {
		return M_PHOTO
	}

	if in.Audio != nil {
		return M_AUDIO
	}

	if in.Document != nil {
		return M_DOCUMENT
	}

	if in.Text != "" {
		return M_TEXT
	}

	return M_UNKNOWN
}
