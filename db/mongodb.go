package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sinde530/go-mancer/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection

func init() {
	// Load env from .env file
	if gin.Mode() != gin.ReleaseMode {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}

		log.Printf("Successd env load")
	}

	uri := os.Getenv("MONGOURI")

	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// checked connection
	err = Client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	// get Collection
	Collection = Client.Database("chattings").Collection("users")
}

func CheckUser(email string) error {
	var result model.RegisterRequest
	err := Collection.FindOne(context.Background(),
		bson.M{"email": email}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return fmt.Errorf("Email already exists")
}

func SaveUser(request *model.RegisterRequest) error {
	err := CheckUser(request.Email)
	if err != nil {
		return err
	}

	_, err = Collection.InsertOne(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*model.RegisterRequest, error) {
	var result model.RegisterRequest
	err := Collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("User not found")
		}
		return nil, err
	}
	return &result, nil
}
