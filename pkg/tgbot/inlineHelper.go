package tgbot

import (
	"strings"

	"github.com/mymmrac/telego"
)

func GetInlineHelper(q *telego.InlineQuery) *InlineHelper {
	return &InlineHelper{
		q,
	}
}

type InlineHelper struct {
	*telego.InlineQuery
}

func (in *InlineHelper) FullName() string {
	builder := strings.Builder{}
	builder.WriteString(in.From.FirstName)
	builder.WriteString(" ")
	builder.WriteString(in.From.LastName)
	return builder.String()
}
