package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/HansBlackCat/grpc-blog/server/dbctrl"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	blogProto "github.com/HansBlackCat/grpc-blog/proto"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// BlogItem is elements mongodb hold
type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

// Global variables for test
var (
	hostURL        string = "mongodb://localhost:27017"
	dbName         string = "grpcblog"
	collectionName string = "blog"
	Collection     *mongo.Collection
)

func main() {
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
	zlog.Info().Msg("Starting blog server")

	zlog.Info().Msg("Connecting to mongodb")
	dbClient, dbErr := dbctrl.CreateNewClient(hostURL)
	if dbErr != nil {
		zlog.Fatal().Msgf("Fail to create mongodb client: %v", dbErr)
	}

	zlog.Info().Msg("Fetching collection from db")
	Collection = dbctrl.FetchCollection(dbClient, dbName, collectionName)

	listen, listenErr := net.Listen("tcp", "0.0.0.0:50051")
	if listenErr != nil {
		zlog.Fatal().Msgf("net Listen failed: %v", listenErr)
	}

	// check if TLS is available
	var options []grpc.ServerOption
	_, tslExistErr := os.Stat("/tls")
	if os.IsExist(tslExistErr) {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			zlog.Fatal().Msgf("Failed loading certificates: %v", sslErr)
		}
		options = append(options, grpc.Creds(creds))
	} else if os.IsNotExist(tslExistErr) {
		// tls is currently unavailable
		zlog.Warn().Msg("Unable to use TLS for grpc")
		zlog.Warn().Msg("Its highly recommend to use TLS on real service")
	} else {
		zlog.Fatal().Msgf("Fail to access to file system, check permission: %v", tslExistErr)
	}

	s := grpc.NewServer(options...)
	blogProto.RegisterBlogServiceServer(s, &Server{})

	// Channel for catching SIGINT
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		zlog.Info().Msgf("Starting Server")
		if err := s.Serve(listen); err != nil {
			zlog.Fatal().Msgf("Failed to serve: %v", err)
		}
	}()

	<-ch
	zlog.Info().Msg("Stop server gracefully...")
	s.Stop()
	zlog.Info().Msg("Closing listener")
	_ = listen.Close()
	zlog.Info().Msg("Closing db")
	err := dbClient.Disconnect(context.Background())
	if err != nil {
		zlog.Fatal().Msgf("Fail to close mongodb gracefully: %v", err)
	}
}
