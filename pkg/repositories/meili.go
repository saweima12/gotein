package repositories

import (
	"fmt"
	"gotein/data"
	"gotein/data/models"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliRepo interface {
	GetUserList(docKey string) ([]string, error)
	SearchMedia(keyword string) ([]*models.MediaDoc, error)
}

func NewMeiliRepo(client meilisearch.ClientInterface) MeiliRepo {
	return &meiliRepo{
		cli: client,
	}
}

type meiliRepo struct {
	cli meilisearch.ClientInterface
}

func (repo *meiliRepo) SearchMedia(keyword string) ([]*models.MediaDoc, error) {
	index := repo.cli.Index(data.INDEX_CACHED)
	resp, err := index.Search(keyword, &meilisearch.SearchRequest{Limit: 20})
	fmt.Println(resp, err)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (repo *meiliRepo) GetUserList(docKey string) ([]string, error) {
	var resp *models.UserList

	index := repo.cli.Index(data.INDEX_CONFIG)
	err := index.GetDocument(docKey, nil, &resp)
	if err != nil {
		return []string{}, err
	}

	return resp.Users, nil
}
