package main

import (
	"TheMainAPI/DB"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)
//Struct defined as per the task here instead of struct users_collectionare used
type UserDetails struct {
	Id  int `json: "Id"`     //unique int type
	Name  string    `json: "Name"`   // rest all strings
	Email string   `json: "Email"`
	Password string `json: "Password"`
}
var UserArr []UserDetails             // here all user are stored with the help of array
//handle functions
func homePage(w http.ResponseWriter, r *http.Request) {       // Function to display few greeting lines
	fmt.Fprintf(w, "This is a Instagram API(Task 1)")
	fmt.Println("Endpoint: homePage")    //endpoint indication
}

func AddUserFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {                             // here Method post used to get the response from web users
		reqBody, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintf(w, "%+v", string(reqBody))
		var users UserDetails     // variable to strore the data
		json.Unmarshal(reqBody, &users)
		// update our global userdetail array to include
		// our new user
		UserArr = append(UserArr, users)    // adding entered details with appendbackwards in the golabal array

		json.NewEncoder(w).Encode(users)        // json used to encode the data frok users var

		fmt.Println("Endpoint: users")        //endpoint indication
		DatabaseName := "Instagram_DB"        // creating and moving the database to env
		client, err := db.CreateDatabaseConnection(DatabaseName)
		if err != nil {                              // checking connection to mongodb Server
			fmt.Println("Failed to connect to Database")
			panic(err)
		}
		defer client.Disconnect(context.TODO()) // this handels when user Disconnects
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		col := client.Database(DatabaseName).Collection("users_collection")
		HashPassword(users.Password)            //encrypting password for secruing user password
		fmt.Println(HashPassword(users.Password))   // showing user DEBUGed:Password
		result, insertErr := col.InsertOne(ctx,users)   // result and err used in collection col
		if insertErr != nil {
			fmt.Println("InsertONE Error:", insertErr) // checking id user entered the data correctly or not
			os.Exit(1)
		} else {
			fmt.Println("InsertOne() result type: ", reflect.TypeOf(result)) // inserting data in collection with reseult
			fmt.Println("InsertOne() api result type: ", result)

			newID := result.InsertedID   // assigning unique id each time at insertion
			fmt.Println("InsertedOne(), newID", newID)     // showing the ID
			fmt.Println("InsertedOne(), newID type:", reflect.TypeOf(newID))

		}

	}
}
func GetUserFunc(w http.ResponseWriter, r *http.Request) {    //function to perform the post and get tasks

		keys, ok := r.URL.Query()["id"]

		if !ok || len(keys[0]) < 1 {          // checking the parameters as passed
			log.Println("Url Param 'key' is missing")
			return
		}

		// Query()["key"] will return an array of items,
		// we only want the single item.
		key := keys[0]

		log.Println("Url Param 'key' is: " + string(key))
		//var userDB UserDetails
	DatabaseName := "Instagram_DB"
	client, err := db.CreateDatabaseConnection(DatabaseName) // create a database if all criterias are fullfiled
	if err != nil {
		fmt.Println("Failed to connect to DB")
		panic(err)  // else print Error
	}
	defer client.Disconnect(context.TODO())  // for disconnecting Client
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	col := client.Database(DatabaseName).Collection("users_collection")
	id_key, err := strconv.Atoi(key)
	filterCursor, err := col.Find(ctx, bson.M{"id":id_key }) // this maps the user id as stroed previously in col Collection
	fmt.Println(filterCursor)
	if err != nil {
		log.Fatal(err)  // if not found the key/id
	}
	var userFiltered []bson.M
	if err = filterCursor.All(ctx, &userFiltered); err != nil {
		log.Fatal(err)
	}
	fmt.Println(userFiltered)
	json.NewEncoder(w).Encode(userFiltered)  // encoded the data with json
}

func handlerequests(){    // the main handlerequests
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", AddUserFunc)
	http.HandleFunc("/users/",GetUserFunc)
	http.HandleFunc("/posts",createPosts)
	http.HandleFunc("/posts/",PostById)
	http.HandleFunc("/posts/users/",returnAllPosts)

	log.Fatal(http.ListenAndServe(":12345", nil))



}
//password encryption genrating methods
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
// checking for Password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func main() {

handlerequests()    // function call in main Method

}
//Arpit Singh
//19BCG10069
