package postEntityManager

import (
	"exgamer.kz/example-service/cmd/app/src/modules/common/repository/mysql/postRepositories"
	"exgamer.kz/example-service/cmd/app/src/modules/common/repository/mysql/postRepositories/postRepoStructures"
	"github.com/exgamer/go-rest-sdk/pkg/exception"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
)

func NewPostEntityManager(
	postRepository *postRepositories.PostRepository,
) *PostEntityManager {
	return &PostEntityManager{
		postRepository: postRepository,
	}
}

type PostEntityManager struct {
	postRepository *postRepositories.PostRepository
}

func (manager *PostEntityManager) All() (*[]postRepoStructures.Post, *exception.AppException) {
	result, err := manager.postRepository.All()

	if err != nil {
		logger.LogAppException(err)

		return nil, err
	}

	return result, nil
}
