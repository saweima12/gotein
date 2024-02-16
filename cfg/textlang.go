package cfg

type textLangMap map[string]string

const (
	MEDIA_CHECKED    = "media_checked"
	MEDIA_UNCHECKED  = "media_unchecked"
	MEDIA_NEW        = "media_new"
	HELP_MSG         = "help_msg"
	REGISTER_SUCCESS = "register_success"

	ERR_INTERNAL = "err_internal"
)

func (te textLangMap) GetText(key string) string {
	if val, ok := te[key]; ok {
		return val
	}
	return "Keyword not found."
}
