package postResponses

import "exgamer.kz/example-service/cmd/app/src/modules/common/repository/mysql/postRepositories/postRepoStructures"

type PostResponse struct {
	Data    []postRepoStructures.Post `json:"data"`
	Success bool                      `json:"success"`
}
