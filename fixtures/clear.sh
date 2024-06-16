#/bin/bash
sudo rm -rf main_service/users_db/data
sudo rm -rf ugc_service/ugc_db/data
sudo rm -rf stats_service/stats_db/data
echo "databases storages deleted"

rm ./main_service/core/post/*.proto
rm ./main_service/core/post/*.pb.go
echo "main service post proto deleted"

rm ./ugc_service/core/post/*.proto
rm ./ugc_service/core/post/*.pb.go
echo "ugc service post proto deleted"

rm ./main_service/core/reaction/*.proto
rm ./main_service/core/reaction/*.pb.go
echo "main service reaction proto deleted"


rm ./stats_service/core/reaction/*.proto
rm ./stats_service/core/reaction/*.pb.go
echo "stats service reaction proto deleted"
