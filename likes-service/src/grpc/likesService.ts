import {
  getLikesCount,
  hasUserLikedPost,
  likePost,
  unlikePost,
  getUserLikedPosts,
} from "../db/likesRepository";

import {
  GetLikesCountRequest,
  HasUserLikedPostRequest,
  LikePostRequest,
  UnlikePostRequest,
  GetUserLikedPostsRequest,
} from "inspire-api-contracts/api/gen/ts/likes/v1/likes_pb"; // adjust path accordingly

class LikesServiceHandler {
  async getLikesCount(req: GetLikesCountRequest) {
    const postId = req.postId;
    const count = await getLikesCount(postId);
    return {
      likesCount: BigInt(count),
    };
  }

  async hasUserLikedPost(req: HasUserLikedPostRequest) {
    const userId = req.userId;
    const postId = req.postId;
    const liked = await hasUserLikedPost(userId, postId);
    return { liked };
  }

  async likePost(req: LikePostRequest) {
    const userId = req.userId;
    const postId = req.postId;

    await likePost(userId, postId);
    return {};
  }

  async unlikePost(req: UnlikePostRequest) {
    const userId = req.userId;
    const postId = req.postId;

    await unlikePost(userId, postId);
    return {};
  }

  async getUserLikedPosts(req: GetUserLikedPostsRequest) {
    const userId = req.userId;
    const limit = req.limit;
    const offset = req.offset;
    const postIds = await getUserLikedPosts(userId, limit, offset);
    return { postIds };
  }
}

export default LikesServiceHandler;
