package postServices

import (
	"exgamer.kz/example-service/cmd/app/src/modules/common/repository/mysql/postRepositories/postRepoStructures"
	"exgamer.kz/example-service/cmd/app/src/modules/post/v1/postEntityManager"
	"github.com/exgamer/go-rest-sdk/pkg/exception"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
)

func NewPostService(postEntityManager *postEntityManager.PostEntityManager) *PostService {
	return &PostService{
		postEntityManager: postEntityManager,
	}
}

type PostService struct {
	postEntityManager *postEntityManager.PostEntityManager
}

func (service *PostService) All() (*[]postRepoStructures.Post, *exception.AppException) {
	posts, appException := service.postEntityManager.All()

	if appException != nil {
		logger.LogAppException(appException)

		return nil, appException
	}

	return posts, nil
}
