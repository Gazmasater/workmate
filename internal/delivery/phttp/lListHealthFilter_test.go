package phttp

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gaz358/myprog/workmate/domain"
	"github.com/stretchr/testify/assert"
)

func TestTaskHandler_ListAndFilterAndHealth(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create two задачи
	resp, err := http.Post(server.URL+"/", "application/json", nil)
	assert.NoError(t, err)
	defer resp.Body.Close()
	var task1 domain.Task
	_ = json.NewDecoder(resp.Body).Decode(&task1)

	resp2, err := http.Post(server.URL+"/", "application/json", nil)
	assert.NoError(t, err)
	defer resp2.Body.Close()
	var task2 domain.Task
	_ = json.NewDecoder(resp2.Body).Decode(&task2)

	// Список всех задач
	listResp, err := http.Get(server.URL + "/")
	assert.NoError(t, err)
	defer listResp.Body.Close()
	assert.Equal(t, http.StatusOK, listResp.StatusCode)

	var list []map[string]interface{}
	err = json.NewDecoder(listResp.Body).Decode(&list)
	assert.NoError(t, err)
	assert.True(t, len(list) >= 2)

	// Фильтр по id
	filterResp, err := http.Get(server.URL + "/filter?id=" + task1.ID)
	assert.NoError(t, err)
	defer filterResp.Body.Close()
	assert.Equal(t, http.StatusOK, filterResp.StatusCode)

	var filtered []map[string]interface{}
	err = json.NewDecoder(filterResp.Body).Decode(&filtered)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, task1.ID, filtered[0]["id"])

	// Health
	healthResp, err := http.Get(server.URL + "/health")
	assert.NoError(t, err)
	defer healthResp.Body.Close()
	body, _ := io.ReadAll(healthResp.Body)
	assert.Equal(t, "ok", string(body))
}
