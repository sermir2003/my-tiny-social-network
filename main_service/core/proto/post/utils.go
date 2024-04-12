package post

import (
	"fmt"
	"log"
	"main_service_core/utils"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Client PostServiceClient

func Connect() error {
	connection_line := fmt.Sprintf("ugc_core:%s", utils.GetenvSafe("UGC_SERVICE_PORT"))
	log.Printf("trying to connect to ugc_service at %s\n", connection_line)
	conn, err := grpc.Dial(connection_line, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	Client = NewPostServiceClient(conn)
	return nil
}
