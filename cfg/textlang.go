package cfg

type textLangMap map[string]string

func (te textLangMap) GetText(key string) string {
	if val, ok := te[key]; ok {
		return val
	}
	return "Keyword not found."
}
