#/bin/bash
sudo rm -rf main_service/users_db/data
sudo rm -rf ugc_service/ugc_db/data
echo "databases storages deleted"

cd ./main_service/core/proto/post
rm *.proto
rm *.pb.go
cd ../../../..
echo "main service proto deleted"

cd ./ugc_service/core/proto/post
rm *.proto
rm *.pb.go
cd ../../../..
echo "ugc service proto deleted"
