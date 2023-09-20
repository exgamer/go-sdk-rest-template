package factories

import "exgamer.kz/example-service/cmd/app/src/modules/post/v1/postHttp"

func NewControllersFactory(servicesFactory *ServicesFactory) *ControllersFactory {
	return &ControllersFactory{
		PostController: postHttp.NewPostController(servicesFactory.PostService),
	}
}

type ControllersFactory struct {
	PostController *postHttp.PostController
}
