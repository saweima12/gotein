package models

type MediaDoc struct {
	MediaType string   `json:"media_type"`
	Keywords  []string `json:"keywords"`
	Uid       string   `json:"uid"`
	FileID    string   `json:"file_id"`
	FilePath  string   `json:"file_path"`
	CacheURL  string   `json:"cache_url"`
}
