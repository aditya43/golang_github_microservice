package app

import (
	"github.com/aditya43/golang_github_microservice/src/api/controllers/polo"
	"github.com/aditya43/golang_github_microservice/src/api/controllers/repositories"
)

func mapUrls() {
	router.GET("/marco", polo.Marco)
	router.POST("/repository", repositories.CreateRepo)
	router.POST("/repositories", repositories.CreateRepos)
}
