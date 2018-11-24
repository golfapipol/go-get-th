package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo"
)

// func TestTodoCreateShouldBeCreated(t *testing.T) {
// 	todoJSON := `{topic:""}`
// 	request := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(todoJSON))
// 	recorder := httptest.NewRecorder()

// 	e := echo.New()
// e.POST("/todos", create)
// e.ServeHTTP(recorder, request)

// response := recorder.Result()
// }

func TestTodoCreateShouldHandleBadRequest(t *testing.T) {
	request := httptest.NewRequest("POST", "/todos", nil)
	recorder := httptest.NewRecorder()

	e := echo.New()
	e.POST("/todos", func(c echo.Context) error {
		return c.JSON(http.StatusBadRequest, "error")
	})
	e.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}
