package factories

import "exgamer.kz/example-service/cmd/app/src/modules/post/v1/postEntityManager"

func NewEntityManagersFactory(mysqlRepositoryFactory *MysqlRepositoryFactory) *EntityManagersFactory {
	return &EntityManagersFactory{
		PostEntityManager: postEntityManager.NewPostEntityManager(
			mysqlRepositoryFactory.PostRepository,
		),
	}
}

type EntityManagersFactory struct {
	PostEntityManager *postEntityManager.PostEntityManager
}
