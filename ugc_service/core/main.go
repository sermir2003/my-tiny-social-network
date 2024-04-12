package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"ugc_service_core/db_utils"
	pb "ugc_service_core/proto/post"
	"ugc_service_core/utils"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPostServiceServer
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	return db_utils.Create(req)
}

func (s *server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	return db_utils.Update(req)
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return db_utils.Delete(req)
}

func (s *server) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetByIdResponse, error) {
	return db_utils.GetById(req)
}

func (s *server) GetPagination(ctx context.Context, req *pb.GetPaginationRequest) (*pb.GetPaginationResponse, error) {
	return db_utils.GetPagination(req)
}

func main() {
	err := db_utils.StartUpDB()
	for err != nil {
		log.Println(err, "retrying...")
		err = db_utils.StartUpDB()
	}

	listening_line := fmt.Sprintf(":%s", utils.GetenvSafe("UGC_SERVICE_PORT"))
	lis, err := net.Listen("tcp", listening_line)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, &server{})
	log.Printf("listening at %s\n", listening_line)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
