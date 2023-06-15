package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sinde530/go-mancer/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var GroupCollection *mongo.Collection

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
	UserCollection = Client.Database("chattings").Collection("users")
	GroupCollection = Client.Database("chattings").Collection("groups")
}

func CheckUser(email string) error {
	var result model.RegisterRequest
	err := UserCollection.FindOne(context.Background(),
		bson.M{"email": email}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return fmt.Errorf("email already exists")
}

func SaveUser(request *model.User) error {
	err := CheckUser(request.Email)
	if err != nil {
		return err
	}

	_, err = UserCollection.InsertOne(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// AuthenticateUser ...
func AuthenticateUser(email, password string) (*model.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		if err.Error() == "user not found" {

			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect email or password")
	}

	return user, nil
}

func SaveGroup(group *model.Group) error {
	_, err := GroupCollection.InsertOne(context.Background(), group)
	if err != nil {
		return err
	}
	return nil
}

func SendGroups() ([]*model.Group, error) {
	groups := []*model.Group{}

	// Define options to sort the groups by a field (e.g., createdAt)
	options := options.Find().SetSort(bson.D{{"createdAt", -1}})

	// Find all groups in the database
	cursor, err := GroupCollection.Find(context.Background(), bson.D{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and decode each group
	for cursor.Next(context.Background()) {
		var group model.Group
		if err := cursor.Decode(&group); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	// Check if any error occurred during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func UpdateUser(user *model.User) error {
	update := bson.M{
		"$set": user,
	}

	_, err := UserCollection.UpdateOne(context.Background(), bson.M{"uid": user.UID}, update)
	if err != nil {
		return err
	}
	return nil
}
