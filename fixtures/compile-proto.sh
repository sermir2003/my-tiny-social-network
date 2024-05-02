#/bin/bash
cd ./main_service/core
cp ../../proto/post/post.proto ./proto/post
./proto/generate.sh
cd ../..
echo "main service proto created"

cd ./ugc_service/core
cp ../../proto/post/post.proto ./proto/post
./proto/generate.sh
cd ../..
echo "ugc service proto created"
