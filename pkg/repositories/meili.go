package repositories

import (
	"gotein/data"
	"gotein/data/models"
	"gotein/libs/cjson"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliRepo interface {
	SearchMedia(keyword string) ([]*models.MediaDoc, error)
	GetUserList(docKey string) (*models.UserList, error)
	SetUserList(docKey string, val *models.UserList) error
	GetMedia(indexUid, uid string) (*models.MediaDoc, error)
	InsertMedia(doc *models.MediaDoc) error
	PutMedia(doc *models.MediaDoc) error
	DeleteMedia(uid string) error
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

	rtn := make([]*models.MediaDoc, 0, len(resp.Hits))

	for i := range resp.Hits {
		var media *models.MediaDoc
		bytes, _ := cjson.Marshal(resp.Hits[i])
		err := cjson.Unmarshal(bytes, &media)
		if err != nil {
			continue
		}
		rtn = append(rtn, media)
	}

	if err != nil {
		return nil, err
	}

	return rtn, nil
}

func (repo *meiliRepo) GetUserList(docKey string) (*models.UserList, error) {
	var resp *models.UserList

	index := repo.cli.Index(data.INDEX_CONFIG)
	err := index.GetDocument(docKey, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (repo *meiliRepo) SetUserList(docKey string, val *models.UserList) error {
	index := repo.cli.Index(data.INDEX_CONFIG)

	_, err := index.UpdateDocuments([]any{val})
	if err != nil {
		return err
	}

	return nil
}

func (repo *meiliRepo) GetMedia(indexUid string, uid string) (*models.MediaDoc, error) {
	index := repo.cli.Index(indexUid)

	var resp models.MediaDoc
	err := index.GetDocument(uid, nil, &resp)

	if err != nil {
		castErr := err.(*meilisearch.Error)
		if castErr.StatusCode == 404 {
			return nil, nil
		}
		return nil, err
	}

	return &resp, nil
}

func (repo *meiliRepo) InsertMedia(item *models.MediaDoc) error {
	index := repo.cli.Index(data.INDEX_CACHED)
	_, err := index.AddDocuments([]any{item})
	if err != nil {
		return err
	}
	return nil
}

func (repo *meiliRepo) PutMedia(doc *models.MediaDoc) error {
	index := repo.cli.Index(data.INDEX_CACHED)
	_, err := index.UpdateDocuments([]any{doc})
	if err != nil {
		return err
	}
	return nil
}

func (repo *meiliRepo) DeleteMedia(uid string) error {
	index := repo.cli.Index(data.INDEX_CACHED)
	_, err := index.DeleteDocument(uid)
	if err != nil {
		return err
	}
	return nil
}
