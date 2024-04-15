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
	connection_string := fmt.Sprintf("ugc_core:%s", utils.GetenvSafe("UGC_SERVICE_PORT"))
	log.Printf("trying to connect to ugc_service at %s\n", connection_string)
	conn, err := grpc.Dial(connection_string, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	Client = NewPostServiceClient(conn)
	return nil
}
