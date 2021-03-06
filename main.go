package main

import (
	"strings"
	"log"
	"github.com/globalsign/mgo/bson"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

const (
	database      = "gogetth"
	todoCollection = "todos"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	url := viper.GetString("mongo.url")
	port := ":" + viper.GetString("port")

	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	// Echo instance
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	// Routes
	handler := &Handler{session: session}
	e.GET("/todos", handler.list)
	e.POST("/todos", handler.create)
	e.GET("/todos/:id", handler.view)
	e.PUT("/todos/:id", handler.done)
	e.DELETE("/todos/:id", handler.delete)

	// Start server
	e.Logger.Fatal(e.Start(port))
}

type todo struct {
	ID    bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
}

type Handler struct {
	session *mgo.Session
}

// Handler
func (handler Handler)list(c echo.Context) error {
	session := handler.session.Copy()
	defer session.Close()
	
	var list []todo
	collection := session.DB(database).C(todoCollection)
	
	if err := collection.Find(nil).All(&list); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, list)
}
func (handler Handler)create(c echo.Context) error {
	session := handler.session.Copy()
	defer session.Close()
	
	var newTodo todo
	if err := c.Bind(&newTodo); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	
	newTodo.ID = bson.NewObjectId()
	collection := session.DB(database).C(todoCollection)
	
	if err := collection.Insert(&newTodo); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newTodo)
}

func (handler Handler)view(c echo.Context) error {
	session := handler.session.Copy()
	defer session.Close()
	
	id := bson.ObjectIdHex(c.Param("id"))

	var todo todo
	collection := session.DB(database).C(todoCollection)
	
	if err := collection.FindId(id).One(&todo); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, todo)
}

func  (handler Handler)done(c echo.Context) error {
	session := handler.session.Copy()
	defer session.Close()
	
	id := bson.ObjectIdHex(c.Param("id"))

	var todo todo
	collection := session.DB(database).C(todoCollection)
	
	if err := collection.FindId(id).One(&todo); err != nil {
		return err
	}

	todo.Done = true

	if err := collection.UpdateId(id, todo); err != nil {
		return err
	}

	
	return c.JSON(http.StatusOK, todo)
}

func (handler Handler)delete(c echo.Context) error {
	session := handler.session.Copy()
	defer session.Close()
	
	id := bson.ObjectIdHex(c.Param("id"))

	collection := session.DB(database).C(todoCollection)
	
	if err := collection.RemoveId(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}