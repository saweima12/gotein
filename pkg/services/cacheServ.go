package services

import (
	"gotein/cfg"
	"gotein/logger"
	"os"
	"path"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type CacheServ interface {
	CacheMedia(fileId string, destination string)
}

func NewCacheServ(api *telego.Bot, cfg *cfg.Configuration) CacheServ {
	return &cacheServ{
		api:       api,
		cacheDir:  cfg.CacheDir,
		domainURL: cfg.DomainURL,
	}
}

type cacheServ struct {
	api       *telego.Bot
	cacheDir  string
	domainURL string
}

func (ca *cacheServ) CacheMedia(fileId string, destination string) {
	resp, err := ca.api.GetFile(&telego.GetFileParams{
		FileID: fileId,
	})

	if err != nil {
		logger.Errorf("Error: Get file path failed, err: %v", err)
		return
	}

	downloadURL := ca.api.FileDownloadURL(resp.FilePath)
	bytes, err := tu.DownloadFile(downloadURL)
	if err != nil {
		logger.Errorf("Error: download file failed, err: %v", err)
		return
	}

	distPath := ca.cacheDir + destination
	baseDir := path.Dir(distPath)
	err = os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		logger.Error(err)
	}

	err = os.WriteFile(distPath, bytes, os.ModePerm)
	if err != nil {
		logger.Errorf("Error: WriteFile err: %v, path: %s", err, distPath)
	}
}
