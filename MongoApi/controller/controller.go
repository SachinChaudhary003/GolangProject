package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sachinchaudhary003/MongoApi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoconn = "mongodb+srv://Godb:password@cluster0.0r5zajl.mongodb.net/"
const dbName = "netflix"
const colName = "watchlist"

var collection *mongo.Collection

//connect with mongodb

func init() {
	//client option

	clientOption := options.Client().ApplyURI(mongoconn)

	//connect to mongo
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mogo connected success")

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection intance is ready")

}

//mongo helpers

// inser1record
func insertOneMovie(movie model.Netflix) {

	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 movie ", inserted.InsertedID)
}

// update one movie
func updateOneMovie(movieid string) {
	id, err := primitive.ObjectIDFromHex(movieid)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modifies", result.ModifiedCount)
}

//delete one movie

func deleteOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	deletecnt, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted movie count", deletecnt)
}

func deleteAllMovie() int64 {

	filter := bson.D{{}}
	deletecnt, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("delted all movie count", deletecnt.DeletedCount)
	return deletecnt.DeletedCount
}

//get all movies

func getAllMovies() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M
	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	defer cur.Close(context.Background())
	return movies
}

//

func GetAllMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencodr")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencodr")
	w.Header().Set("Allow-Control_Allow_Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)

}

func MarkedAsWatched(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencodr")
	w.Header().Set("Allow-Control_Allow_Methods", "PUT")

	params := mux.Vars(r)

	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencodr")
	w.Header().Set("Allow-Control_Allow_Methods", "DELETE")
	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencodr")
	w.Header().Set("Allow-Control_Allow_Methods", "POST")
	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}
