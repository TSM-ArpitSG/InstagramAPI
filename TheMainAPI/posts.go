package main

import (
"TheMainAPI/db"
"context"
"encoding/json"
"fmt"
"go.mongodb.org/mongo-driver/bson"
"io/ioutil"
"log"
"net/http"
"os"
"reflect"
"strconv"
"time"
)

//Struct for post collection as described in the Task
type postsDetails struct {

Id  int `json: "Id"`      // id and post id given unique integer id for furture uses
PostId  int `json: "PostId"`
Caption  string    `json: "Caption"`
ImageURL string   `json: "ImageURL"`
PostedTS time.Time `json: "PostedTS"` // Time selector used for time datatype
}
var PostArr []postsDetails  // all post details stored in this array
func createPosts(w http.ResponseWriter, r *http.Request){
if r.Method == "POST" {

  reqBody, _ := ioutil.ReadAll(r.Body)  // r used for request var.
  fmt.Fprintf(w, "%+v", string(reqBody))
  var posts postsDetails  // var to add and include postsDetails
  json.Unmarshal(reqBody, &posts)
  // update our global post array to include
  // our new post with the help of appendbackwards.
  PostArr = append(PostArr, posts)

  json.NewEncoder(w).Encode(posts)  // encoded the data in json format

  fmt.Println("Endpoint : posts")
  DatabaseName := "Instagram_DB" // move to env
  client, err := db.CreateDatabaseConnection(DatabaseName)
  if err != nil {   // here proccess is simimlar to main.go i.e. for the creation and connection to Server
    fmt.Println("Failed to connect to DB")
    panic(err)
  }
  defer client.Disconnect(context.TODO())

  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  col := client.Database(DatabaseName).Collection("users_posts")
  result, insertErr := col.InsertOne(ctx,posts)   // insertion of posts proccess
  if insertErr != nil {
    fmt.Println("InsertONE Error:", insertErr)
    os.Exit(1)
  } else {
    fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
    fmt.Println("InsertOne() api result type: ", result)

    newID := result.InsertedID      // seprate insertion Id assigned each Time
    fmt.Println("InsertedOne(), newID", newID)
    fmt.Println("InsertedOne(), newID type:", reflect.TypeOf(newID))

  }

}

}

func PostById(w http.ResponseWriter, r *http.Request) {   // function/logic to search the post with postid method

keys, ok := r.URL.Query()["postid"]

if !ok || len(keys[0]) < 1 {
  log.Println("Url Param 'key' is missing")
  return
}

// Query()["key"] will return an array of items,
// we only want the single item.
key := keys[0]

log.Println("Url Param 'key' is: " + string(key))
//var userDB UserDetails
DatabaseName := "Instagram_DB"
client, err := db.CreateDatabaseConnection(DatabaseName)
if err != nil {
  fmt.Println("Failed to connect to DB")
  panic(err)
}
defer client.Disconnect(context.TODO())
ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
col := client.Database(DatabaseName).Collection("users_posts")
id_key, err := strconv.Atoi(key)  // similar to main.go just instead of key here postid is used
filterCursor, err := col.Find(ctx, bson.M{"postid":id_key })
fmt.Println(filterCursor)
if err != nil {
  log.Fatal(err)
}
var postFiltered []bson.M
if err = filterCursor.All(ctx, &postFiltered); err != nil {
  log.Fatal(err)
}
fmt.Println(postFiltered) // mapped with postid and added to the collection in json format.
json.NewEncoder(w).Encode(postFiltered)
}


func returnAllPosts(w http.ResponseWriter, r *http.Request) {
keys, ok := r.URL.Query()["id"]

if !ok || len(keys[0]) < 1 {
  log.Println("Url Param 'key' is missing")
  return
}

// Query()["key"] will return an array of items,
// we only want the single item.
key := keys[0]

log.Println("Url Param 'key' is: " + string(key))
DatabaseName := "Instagram_DB" // move to env
client, err := db.CreateDatabaseConnection(DatabaseName)
if err != nil {
  fmt.Println("Failed to connect to DB")
  panic(err)
}
defer client.Disconnect(context.TODO())
ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
col := client.Database(DatabaseName).Collection("users_posts")
id_key, err := strconv.Atoi(key)
//var result UserDetails
cursor, err := col.Find(ctx, bson.M{"id":id_key })
if err != nil {
  log.Fatal(err)
}
var user_list []bson.M
if err = cursor.All(ctx, &user_list); err != nil {
  log.Fatal(err)
}
fmt.Println(user_list)
fmt.Println("Endpoint Hit: returnAllPosts")
json.NewEncoder(w).Encode(user_list)
}
//Arpit Singh
//19BCG10069
