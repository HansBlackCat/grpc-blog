package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	blogProto "github.com/HansBlackCat/grpc-blog/proto"
	zlog "github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct{}

func (_ *Server) DeleteBlog(ctx context.Context, req *blogProto.DeleteBlogRequest) (*blogProto.DeleteBlogResponse, error) {
	id := req.GetId()
	objectID, idErr := primitive.ObjectIDFromHex(id)
	if idErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Wrong given id: %v", idErr))
	}

	filter := primitive.M{
		"_id": objectID,
	}

	_, resErr := Collection.DeleteOne(ctx, filter)
	if resErr != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find with given id: %v", resErr))
	}

	return &blogProto.DeleteBlogResponse{}, nil
}

func (_ *Server) UpdateBlog(ctx context.Context, req *blogProto.UpdateBlogRequest) (*blogProto.UpdateBlogResponse, error) {
	blog := req.GetBlog()

	objectID, idErr := primitive.ObjectIDFromHex(blog.GetId())
	if idErr != nil {
		zlog.Info().Err(idErr).Msg("[UpdateBlog] Get wrong hexID")
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Wrong hexID: %v", idErr))
	}

	item := &BlogItem{}
	filter := primitive.M{
		"_id": objectID,
	}

	res := Collection.FindOne(context.Background(), filter)
	if err := res.Decode(item); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find BlogItem with given ID: %v", err))
	}

	item.AuthorID = blog.GetAuthorId()
	item.Content = blog.GetContent()
	item.Title = blog.GetTitle()

	updateResult, updateErr := Collection.ReplaceOne(context.Background(), filter, item)
	if updateErr != nil {
		return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Cannot update collection, DB currently unabavailable: %v", updateErr))
	}

	zlog.Info().Str("ModifiedCount", strconv.FormatInt(updateResult.ModifiedCount, 10))

	resp := &blogProto.UpdateBlogResponse{Blog: &blogProto.Blog{
		Id:       item.ID.Hex(),
		AuthorId: item.AuthorID,
		Title:    item.Title,
		Content:  item.Content,
	}}

	zlog.Info().Msgf("updateBlog was invoked successfully with: %v", resp)
	return resp, nil

}

func (_ *Server) ReadBlog(ctx context.Context, req *blogProto.ReadBlogRequest) (*blogProto.ReadBlogResponse, error) {
	hexID := req.GetId()

	objectID, idErr := primitive.ObjectIDFromHex(hexID)
	if idErr != nil {
		zlog.Info().Err(idErr).Msg("[ReadBlog] Get wrong hexID")
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Wrong hexID: %v", idErr))
	}

	item := &BlogItem{}
	filter := bson.M{
		"_id": objectID,
	}

	res := Collection.FindOne(context.Background(), filter)
	if err := res.Decode(item); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find BlogItem with given ID: %v", err))
	}

	resp := &blogProto.ReadBlogResponse{Blog: &blogProto.Blog{
		Id:       item.ID.Hex(),
		AuthorId: item.AuthorID,
		Content:  item.Content,
		Title:    item.Title,
	}}

	zlog.Info().Msgf("CreatingBlog was invoked successfully with: %v", resp)
	return resp, nil
}

func (_ *Server) CreateBlog(ctx context.Context, req *blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error) {
	// Check context
	if errors.Is(ctx.Err(), context.Canceled) {
		zlog.Err(ctx.Err()).Msg("Client cancels context")
		return nil, status.Errorf(codes.Canceled, fmt.Sprintf("Client cancels context: %v", ctx.Err()))
	}

	var collection *mongo.Collection = Collection
	b := req.GetBlog()

	item := BlogItem{
		// Omit ID
		AuthorID: b.GetAuthorId(),
		Content:  b.GetContent(),
		Title:    b.GetTitle(),
	}

	//if v := ctx.Value("collection"); v != nil {
	//	c, ok := v.(*mongo.Collection)
	//	if !ok {
	//		zlog.Fatal().Msg("Failed to pass *mongo.Collection with context")
	//	}
	//	collection = c
	//} else {
	//	zlog.Fatal().Msgf("Error while fetching value from context")
	//}

	res, err := collection.InsertOne(context.Background(), item)
	if err != nil {
		zlog.Err(err).Msgf("Fail to insert data: %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Fail to insert data: %v", err))
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		zlog.Err(errors.New("casting error"))
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Casting ID error"))
	}

	response := &blogProto.CreateBlogResponse{
		Blog: &blogProto.Blog{
			Id:       id.Hex(),
			AuthorId: b.GetAuthorId(),
			Title:    b.GetTitle(),
			Content:  b.GetContent(),
		},
	}

	zlog.Info().Msgf("CreatingBlog was invoked successfully with: %v", response)
	return response, nil
}
