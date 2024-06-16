#/bin/bash
cp ./proto/post/post.proto ./main_service/core/post
cp ./proto/post/post.proto ./ugc_service/core/post
cp ./proto/reaction/reaction.proto ./main_service/core/reaction
cp ./proto/reaction/reaction.proto ./stats_service/core/reaction
echo "proto files have been copied to the needed repositories"
