package postRepositories

import (
	"database/sql"
	"encoding/json"
	"exgamer.kz/example-service/cmd/app/src/modules/common/repository/mysql/postRepositories/postRepoStructures"
	"github.com/exgamer/go-rest-sdk/pkg/db/mysql"
	"github.com/exgamer/go-rest-sdk/pkg/exception"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
	"net/http"
)

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

type PostRepository struct {
	db *sql.DB
}

func (repo *PostRepository) All() (*[]postRepoStructures.Post, *exception.AppException) {
	queryBuilder := mysql.NewQueryBuilder()

	res, appException := queryBuilder.
		SetDb(repo.db).
		Table("posts").
		TableAlias("p").
		Select(
			"p.id",
			"p.title",
			"p.description",
		).
		All()

	if appException != nil {
		logger.LogAppException(appException)

		return nil, appException
	}

	entities, appException := repo.mapToEntityArray(res)

	if appException != nil {
		logger.LogAppException(appException)

		return nil, appException
	}

	return entities, nil
}

func (repo *PostRepository) mapToEntityArray(res *[]map[string]interface{}) (*[]postRepoStructures.Post, *exception.AppException) {
	jsonRes, err := json.Marshal(res)

	if err != nil {
		logger.LogError(err)

		return nil, exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	var entities []postRepoStructures.Post
	json.Unmarshal(jsonRes, &entities)

	return &entities, nil
}
