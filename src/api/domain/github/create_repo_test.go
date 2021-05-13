package github

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepoRequestAsJson(t *testing.T) {
	request := CreateRepoRequest{
		Name:        "Aditya Test Repo",
		Description: "This is a test git repository",
		Homepage:    "https://aditya.com",
		Private:     true,
		HasIssues:   true,
		HasProjects: true,
		HasWiki:     true,
	}

	// Marshal takes an input interface and attempts to create a valid JSON string
	bytes, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	// t.Log(string(bytes))
	// assert.EqualValues(t, `{"name":"Aditya Test Repo","description":"This is a test git repository","homepage":"https://aditya.com","private":true,"has_issues":true,"has_projects":true,"has_wiki":true}`, string(bytes))

	var target CreateRepoRequest
	// Unmarshal takes input byte array (JSON) and a pointer to struct we want to populate with JSON fields
	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
}
