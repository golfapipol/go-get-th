package main

import (
	"github.com/globalsign/mgo/bson"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/globalsign/mgo"
)

const (
	url            = "localhost:27018"
	database       = "gogetth"
	todoCollection = "todos"
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
	ID    bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
}

// Handler
func create(c echo.Context) error {
	var newTodo todo
	if err := c.Bind(&newTodo); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	newTodo.ID = bson.NewObjectId()
	collection := session.DB(database).C(todoCollection)

	if err = collection.Insert(&newTodo); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newTodo)
}
