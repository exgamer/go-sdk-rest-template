package factories

import "exgamer.kz/example-service/cmd/app/src/modules/post/v1/postServices"

func NewServiceFactory(entityManagersFactory *EntityManagersFactory) *ServicesFactory {
	return &ServicesFactory{
		PostService: postServices.NewPostService(entityManagersFactory.PostEntityManager),
	}
}

type ServicesFactory struct {
	PostService *postServices.PostService
}
