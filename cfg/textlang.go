package cfg

type textLangMap map[string]string

const (
	MEDIA_CHECKED    = "media_checked"
	MEDIA_UNCHECKED  = "media_unchecked"
	MEDIA_NEW        = "media_new"
	HELP_MSG         = "help_msg"
	LISTEN_SUCCESS   = "listen_success"
	REGISTER_SUCCESS = "register_success"
	GROUP_ADD        = "group_add"

	ERR_INTERNAL    = "err_internal"
	ERR_NEED_REPLY  = "err_need_reply"
	ERR_PARAMS_NEED = "err_params_need"
	ERR_NOT_FOUND   = "err_not_found"

	SK_SUCCESS = "sk_success"
	AK_SUCCESS = "ak_success"
	RM_SUCCESS = "rm_success"
)

func (te textLangMap) GetText(key string) string {
	if val, ok := te[key]; ok {
		return val
	}
	return "Keyword not found."
}
