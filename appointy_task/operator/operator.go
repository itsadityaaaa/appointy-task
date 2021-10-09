package operator

import (
	"Rest_api/Structure"
	"Rest_api/encrypt"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	get_user_request      = regexp.MustCompile(`^\/users\/(\w+)$`)
	get_post_request      = regexp.MustCompile(`^\/posts\/(\w+)$`)
	get_user_post_request = regexp.MustCompile(`\/posts/users\/(\w+)$`)
)

func ConnectDB() (*mongo.Collection, *mongo.Collection) {

	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection_1 := client.Database("aditya").Collection("users")
	collection_2 := client.Database("aditya").Collection("posts")

	return collection_1, collection_2
}

var collection_1, collection_2 = ConnectDB()

func Createuserendpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user Structure.User

	_ = json.NewDecoder(r.Body).Decode(&user)

	key := "123456789012345678901234"
	hashed_password := encrypt.Encrypt(key, user.Password)
	user.Password = hashed_password

	result, err := collection_1.InsertOne(context.TODO(), &user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(result)

}

func Userbyidendpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user Structure.User
	Path := get_user_request.FindStringSubmatch(r.URL.Path)
	id := Path[1]

	filter := bson.M{"_id": id}
	err := collection_1.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func Createpostendpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var post Structure.Post
	post.TimeStamp = time.Now()

	_ = json.NewDecoder(r.Body).Decode(&post)
	result, err := collection_2.InsertOne(context.TODO(), &post)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func Getpostbyidendpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var post Structure.Post

	Path := get_post_request.FindStringSubmatch(r.URL.Path)
	id := Path[1]

	filter := bson.M{"_id": id}
	err := collection_2.FindOne(context.TODO(), filter).Decode(&post)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return

	}

	json.NewEncoder(w).Encode(post)
}

func Getuserspostbyidendpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var posts []Structure.Post

	Path := get_user_post_request.FindStringSubmatch(r.URL.Path)

	id := Path[1]

	cur, err := collection_2.Find(context.TODO(), bson.M{})

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for cur.Next(context.TODO()) {
		var single_post Structure.Post

		err := cur.Decode(&single_post)
		if err != nil {
			log.Fatal(err)
		}
		if (single_post.UserID) == id {
			posts = append(posts, single_post)
		}
	}

	if err := cur.Err(); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(posts)
}