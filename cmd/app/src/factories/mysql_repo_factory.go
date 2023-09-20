package factories

import (
	"database/sql"
	"exgamer.kz/example-service/cmd/app/src/modules/common/repository/mysql/postRepositories"
)

func NewMysqlRepositoryFactory(db *sql.DB) *MysqlRepositoryFactory {
	return &MysqlRepositoryFactory{
		PostRepository: postRepositories.NewPostRepository(db),
	}
}

type MysqlRepositoryFactory struct {
	PostRepository *postRepositories.PostRepository
}
