package tgbot

import "github.com/mymmrac/telego"

func GetInlineHelper(q *telego.InlineQuery) *InlineHelper {
	return &InlineHelper{
		q,
	}
}

type InlineHelper struct {
	*telego.InlineQuery
}
