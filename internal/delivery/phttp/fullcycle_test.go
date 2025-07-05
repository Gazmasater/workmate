package phttp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/gaz358/myprog/workmate/repository/memory"
	"github.com/gaz358/myprog/workmate/usecase"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *httptest.Server {
	repo := memory.NewInMemoryRepo()
	uc := usecase.NewTaskUseCase(repo, 200*time.Millisecond)
	handler := NewHandler(uc)
	return httptest.NewServer(handler.Routes())
}

func TestTaskHandler_FullCycleWithCancel(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create task
	resp, err := http.Post(server.URL+"/", "application/json", nil)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var created domain.Task
	err = json.Unmarshal(body, &created)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)

	// Cancel task
	cancelReq, err := http.NewRequest(http.MethodPut, server.URL+"/"+created.ID+"/cancel", nil)
	assert.NoError(t, err)
	cancelResp, err := http.DefaultClient.Do(cancelReq)
	assert.NoError(t, err)
	defer cancelResp.Body.Close()
	assert.Equal(t, http.StatusOK, cancelResp.StatusCode)

	body, _ = io.ReadAll(cancelResp.Body)
	var cancelResult map[string]string
	err = json.Unmarshal(body, &cancelResult)
	assert.NoError(t, err)
	assert.Equal(t, "canceled", cancelResult["status"])

	// Get task by ID
	getResp, err := http.Get(server.URL + "/" + created.ID)
	assert.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	body, _ = io.ReadAll(getResp.Body)
	var fetched domain.Task
	err = json.Unmarshal(body, &fetched)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, domain.StatusCancelled, fetched.Status)

	// Delete task
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/"+created.ID, nil)
	assert.NoError(t, err)
	delResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer delResp.Body.Close()
	assert.Equal(t, http.StatusNoContent, delResp.StatusCode)

	// Get deleted task
	getDeletedResp, err := http.Get(server.URL + "/" + created.ID)
	assert.NoError(t, err)
	defer getDeletedResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, getDeletedResp.StatusCode)
}
