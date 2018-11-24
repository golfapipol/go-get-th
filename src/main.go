package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())

	// Routes
	e.POST("/todos", create)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

type todo struct {
	ID    string `json:"id"`
	Topic string `json:"topic"`
	Done  bool   `json:"done"`
}

// Handler
func create(c echo.Context) error {
	var newTodo todo
	if err := c.Bind(&newTodo); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, newTodo)
}
