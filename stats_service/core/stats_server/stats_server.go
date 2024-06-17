package stats_server

import (
	context "context"
	"fmt"
	"log"
	"net"
	"stats_service_core/db_utils"
	stats_pb "stats_service_core/stats"
	"stats_service_core/utils"

	"google.golang.org/grpc"
)

type server struct {
	stats_pb.UnimplementedStatsServiceServer
}

func (s *server) GetPostStats(ctx context.Context, req *stats_pb.GetPostStatsRequest) (*stats_pb.GetPostStatsResponse, error) {
	return db_utils.GetPostStats(req.PostId)
}

func (s *server) GetTopPosts(ctx context.Context, req *stats_pb.GetTopPostsRequest) (*stats_pb.GetTopPostsResponse, error) {
	return db_utils.GetTopPosts(req)
}

func (s *server) GetTopUsers(ctx context.Context, req *stats_pb.GetTopUsersRequest) (*stats_pb.GetTopUsersResponse, error) {
	return db_utils.GetTopUsers(req)
}

func RunServer() {
	listening_line := fmt.Sprintf(":%s", utils.GetenvSafe("STATS_SERVICE_GRPC_PORT"))
	lis, err := net.Listen("tcp", listening_line)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	stats_pb.RegisterStatsServiceServer(s, &server{})
	log.Printf("listening grpc at %s\n", listening_line)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
