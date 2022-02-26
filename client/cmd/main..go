package main

import (
	"os"
	"time"

	blogProto "github.com/HansBlackCat/grpc-blog/proto"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
	zlog.Info().Msg("Starting blog client...")

	// Check TLS
	var dialOption grpc.DialOption
	_, tlsErr := os.Stat("/tls")
	if os.IsExist(tlsErr) {
		certFile := "ssl/ca.crt"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			zlog.Fatal().Err(sslErr).Msg("Cannot load CA cert from file")
		}
		dialOption = grpc.WithTransportCredentials(creds)
	} else if os.IsNotExist(tlsErr) {
		zlog.Warn().Msg("Unable to use TLS for grpc")
		zlog.Warn().Msg("Its highly recommend to use TLS on real service")
		dialOption = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		zlog.Fatal().Msgf("Fail to access to file system, check permission: %v", tlsErr)
	}

	conn, dialErr := grpc.Dial("127.0.0.1:50051", dialOption)
	if dialErr != nil {
		zlog.Fatal().Msg("Fail to establish dial")
	}
	defer conn.Close()

	c := blogProto.NewBlogServiceClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()

	//// TODO make client form
	//blog := &blogProto.Blog{
	//	AuthorId: "AAAA",
	//	Title:    "This is for test",
	//	Content:  "This is test content",
	//}
	//var temp string
	//if res, err := requestCreateBlog(ctx, c, &blogProto.CreateBlogRequest{Blog: blog}); err != nil {
	//	zlog.Err(err).Msg("request failed")
	//} else {
	//	temp = res.GetBlog().GetId()
	//	zlog.Info().Msgf("Get response: %v", temp)
	//}

	//ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel2()
	//if _, err := requestReadBlog(ctx2, c, &blogProto.ReadBlogRequest{Id: "6214e6b9a524da25c9d74117"}); err != nil {
	//	zlog.Err(err).Msg("request failed")
	//}

	//blogUpdate := &blogProto.Blog{
	//	Id:       "6214e6296a8f5971f0703e1b",
	//	AuthorId: "ABBB",
	//	Title:    "I'm B",
	//	Content:  "Github",
	//}
	//ctx3, cancel3 := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel3()
	//if _, err := requestUpdateBlog(ctx3, c, &blogProto.UpdateBlogRequest{Blog: blogUpdate}); err != nil {
	//	zlog.Err(err).Msg("request failed")
	//}

	if _, err := requestDeleteBlog(context.Background(), c, &blogProto.DeleteBlogRequest{Id: "6214e6b6a524da25c9d74116"}); err != nil {
		zlog.Err(err).Msg("request failed")
	}
}

func requestCreateBlog(ctx context.Context, c blogProto.BlogServiceClient, req *blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error) {
	createBlogResponse, resErr := c.CreateBlog(ctx, req)
	if resErr != nil {
		zlog.Err(resErr).Msg("Fail to createBlog")
		return nil, resErr
	}
	zlog.Info().Msgf("Successfully create blog: %v", createBlogResponse)
	return createBlogResponse, nil
}

func requestReadBlog(ctx context.Context, c blogProto.BlogServiceClient, req *blogProto.ReadBlogRequest) (*blogProto.ReadBlogResponse, error) {
	readBlogResponse, resErr := c.ReadBlog(ctx, req)
	if resErr != nil {
		zlog.Err(resErr).Msg("Fail to readBlog")
		return nil, resErr
	}
	zlog.Info().Msgf("Successfully read blog: %v", readBlogResponse)
	return readBlogResponse, nil
}

func requestUpdateBlog(ctx context.Context, c blogProto.BlogServiceClient, req *blogProto.UpdateBlogRequest) (*blogProto.UpdateBlogResponse, error) {
	updateBlogRes, resErr := c.UpdateBlog(ctx, req)
	if resErr != nil {
		zlog.Err(resErr).Msg("Fail to updateBlog")
		return nil, resErr
	}
	zlog.Info().Msgf("Successfully update blog: %v", updateBlogRes)
	return updateBlogRes, nil
}

func requestDeleteBlog(ctx context.Context, c blogProto.BlogServiceClient, req *blogProto.DeleteBlogRequest) (*blogProto.DeleteBlogResponse, error) {
	deleteRes, resErr := c.DeleteBlog(ctx, req)
	if resErr != nil {
		zlog.Err(resErr).Msg("Fail to deleteBlog")
		return nil, resErr
	}
	zlog.Info().Msgf("Successfully delete blog: %v", deleteRes)
	return deleteRes, nil
}
