#/bin/bash
cd ./main_service/core
cp ../../proto/post/post.proto ./post
./post/generate.sh
cd ../..
echo "main service post proto created"

cd ./ugc_service/core
cp ../../proto/post/post.proto ./post
./post/generate.sh
cd ../..
echo "ugc service post proto created"

cd ./main_service/core
cp ../../proto/reaction/reaction.proto ./reaction
./reaction/generate.sh
cd ../..
echo "main service reaction proto created"

cd ./stats_service/core
cp ../../proto/reaction/reaction.proto ./reaction
./reaction/generate.sh
cd ../..
echo "stats service reaction proto created"
