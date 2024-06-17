package stats

import (
	"fmt"
	"log"
	"main_service_core/utils"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Client StatsServiceClient

func Connect() error {
	connection_string := fmt.Sprintf("stats_core:%s", utils.GetenvSafe("STATS_SERVICE_GRPC_PORT"))
	log.Printf("trying to connect to stats_service at %s\n", connection_string)
	conn, err := grpc.Dial(connection_string, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	Client = NewStatsServiceClient(conn)
	return nil
}
