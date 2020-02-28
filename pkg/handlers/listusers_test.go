package handlers

import (
	"github.com/minio/minio/pkg/madmin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testUserLister struct {
}

func (ul testUserLister) ListUsers() (map[string]madmin.UserInfo, error) {
	returnMap := make(map[string]madmin.UserInfo)
	ui1 := madmin.UserInfo{
		PolicyName: "Policy1",
		Status:     madmin.AccountEnabled,
	}
	ui2 := madmin.UserInfo{
		PolicyName: "Policy2",
		Status:     madmin.AccountEnabled,
	}
	returnMap["user1"] = ui1
	returnMap["user2"] = ui2
	return returnMap, nil
}

func TestListUsers(t *testing.T) {
	t.Run("Should create new ListUsersHandler", func(t *testing.T) {
		listusersHandler, err := NewListUsersHandler(getTestAppConfig())
		assert.Nil(t, err)
		assert.NotNil(t, listusersHandler)
	})

	t.Run("Should return JSON of users", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost:8080/listusers", nil)
		response := httptest.NewRecorder()
		testUserLister := testUserLister{}
		listusersHandler := ListUsersHandler{
			userLister: testUserLister,
		}

		listusersHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), "user1")
		assert.Contains(t, response.Body.String(), "user2")
		assert.Contains(t, response.Body.String(), "Policy1")
		assert.Contains(t, response.Body.String(), "Policy2")
		assert.Contains(t, response.Body.String(), madmin.AccountEnabled)
	})

}
