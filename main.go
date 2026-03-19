package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Todo struct {
	ID        bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool          `json:"completed"`
	Body      string        `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("server run")

	// 1. Load env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// 2. Get URI
	MONGODB_URI := os.Getenv("MONGODB_URI")
	if MONGODB_URI == "" {
		log.Fatal("MONGODB_URI not found in .env")
	}
	// 3. Connect to MongoDB (v2 Style)
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	// In v2, mongo.Connect returns ONLY the client, not an error.
	// It connects lazily (on demand).
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 4. Verification (Ping)
	// In v2, Ping requires a Context and a ReadPreference (nil is okay for default)
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}
	// 5. Initialize the collection variable
	collection = client.Database("golang_db").Collection("todos")
	fmt.Println("Successfully connected to mongoDB!")
	// Don't forget to close the connection when the app stops
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
	//--------------- connect to spcfic db and collection
	collection = client.Database("golang_db").Collection("todo")
	// start create fiber app
	app := fiber.New()
	app.Get("/api/v1/todos", getTodos)
	app.Post("/api/v1/todos", createTodo)
	app.Patch("/api/v1/todos/:id", updateTodo)
	app.Delete("/api/v1/todos/:id", deleteTodo)
	// port
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}
	log.Fatal(app.Listen(":" + PORT))
}

func getTodos(c fiber.Ctx) error {
	var todos []Todo
	// 1. Fetch data from MongoDB
	// bson.M{} means "no filter" (give me everything)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// 2. IMPORTANT: Close the cursor when the function is done
	defer cursor.Close(context.Background())

	// 3. Loop through the database results
	// one by one and decode it
	for cursor.Next(context.Background()) {
		var todo Todo
		// Decode turns the BSON (binary) from DB into our Go Struct
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		// Add to our list
		todos = append(todos, todo)
	}
	//

	// 4. Send the list back to the user as JSON
	return c.Status(200).JSON(todos)
}

func createTodo(c fiber.Ctx) error {
	todo := new(Todo)
	c.Bind().Body(&todo)

	if err := c.Bind().Body(&todo); err != nil {
		return err
	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "body is required"})
	}
	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}
	todo.ID = insertResult.InsertedID.(bson.ObjectID)

	return c.Status(201).JSON(todo)
}
func updateTodo(c fiber.Ctx) error {
    id := c.Params("id")
    objectID, err := bson.ObjectIDFromHex(id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"msg": "invalid todo ID"})
    }
    filter := bson.M{"_id": objectID}
    update := bson.M{"$set": bson.M{"completed": true}}
    _, err = collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }
    return c.Status(200).JSON(fiber.Map{"msg": "Todo updated successfully"})
}
func deleteTodo (c fiber.Ctx)error{
    id := c.Params("id")
    objectID, err := bson.ObjectIDFromHex(id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"msg": "invalid todo ID"})
    }
    filter := bson.M{"_id": objectID}
    _, err = collection.DeleteOne(context.Background(), filter)
    if err != nil {
        return err
    }
    return c.Status(200).JSON(fiber.Map{"msg": "Todo deleted successfully"})
}
