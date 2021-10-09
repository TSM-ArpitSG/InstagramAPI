package main //executable

import (
	"context"
//	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"     // for processing http requests
	//"go.mongodb.org/mongo-driver/bson" // for bson formating used for MangoDB interpretation
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"   // for MangoDB connection
)
//Data models
type User struct {                      // defining as asked in task
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
  Password  string             `json:"Password,omitempty" bson:"Password,omitempty"`
}
type  Post struct{                      // defining as asked in task
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption    string             `json:"Caption,omitempty" bson:"Caption,omitempty"`
	Image_url  string             `json:"Image_url,omitempty" bson:"Image_url,omitempty"`
  Posted_time string            `json:"Posted_time,omitempty" bson:"Posted_time,omitempty"`
}
var client *mongo.Client      // a client would be setup for each time of implementation
func main() {
	fmt.Println("Starting the application...")   //just to show start of application
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)   // for establishing connection (context)
	client, _ = mongo.Connect(ctx,("mongodb://localhost:27017"))  // mongo Connect
  router := mux.NewRouter()     // setup router http properties with the help of Mux
  http.ListenAndServe(":12345", router) // http properties set here along with port and router that we created
}
