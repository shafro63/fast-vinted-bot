package database

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	_        = godotenv.Load()
	mongoURI = os.Getenv("MONGODB_URI")
	Client   *mongo.Client
)

func InitMongoClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		slog.Error("unable to connect to mongoDB", "error", err)
		os.Exit(1)
	}
	// Verify connection sucess (ping)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		slog.Error("can't send ping to mongoDB", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to mongoDB !")
	Client = client
}
