package services

import (
	"net/http"
	"sync"

	"github.com/aditya43/golang_github_microservice/src/api/config"
	"github.com/aditya43/golang_github_microservice/src/api/domain/github"
	"github.com/aditya43/golang_github_microservice/src/api/domain/repositories"
	"github.com/aditya43/golang_github_microservice/src/api/log/zap_logger"
	"github.com/aditya43/golang_github_microservice/src/api/providers/github_provider"
	"github.com/aditya43/golang_github_microservice/src/api/utils/errors"
)

type reposService struct{}

type reposServiceInterface interface {
	CreateRepo(clientId string, request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

var (
	RepositoryService reposServiceInterface
)

func init() {
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(clientId string, input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}
	zap_logger.Info("about to send request to external api",
		zap_logger.Field("client_id", clientId),
		zap_logger.Field("status", "pending"),
		zap_logger.Field("authenticated", clientId != ""))

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		zap_logger.Error("response obtained from external api", err,
			zap_logger.Field("client_id", clientId),
			zap_logger.Field("status", "error"),
			zap_logger.Field("authenticated", clientId != ""))
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	zap_logger.Info("response obtained from external api",
		zap_logger.Field("client_id", clientId),
		zap_logger.Field("status", "success"),
		zap_logger.Field("authenticated", clientId != ""))

	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}
	return &result, nil
}

func (s *reposService) CreateRepos(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {
	input := make(chan repositories.CreateRespositoriesResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	var wg sync.WaitGroup
	go s.handleRepoResults(&wg, input, output)

	for _, current := range requests {
		wg.Add(1)
		go s.createRepoConcurrent(current, input)
	}

	wg.Wait()
	close(input)

	result := <-output

	successCreations := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successCreations++
		}
	}
	if successCreations == 0 {
		result.StatusCode = result.Results[0].Error.Status()
	} else if successCreations == len(requests) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}
	return result, nil
}

func (s *reposService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateRespositoriesResult, output chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse
	for incomingEvent := range input {
		repoResult := repositories.CreateRespositoriesResult{
			Response: incomingEvent.Response,
			Error:    incomingEvent.Error,
		}
		results.Results = append(results.Results, repoResult)
		wg.Done()
	}
	output <- results
}

func (s *reposService) createRepoConcurrent(input repositories.CreateRepoRequest, output chan repositories.CreateRespositoriesResult) {
	if err := input.Validate(); err != nil {
		output <- repositories.CreateRespositoriesResult{Error: err}
		return
	}
	result, err := s.CreateRepo("TODO_client_id", input)
	if err != nil {
		output <- repositories.CreateRespositoriesResult{Error: err}
		return
	}
	output <- repositories.CreateRespositoriesResult{Response: result}
}
