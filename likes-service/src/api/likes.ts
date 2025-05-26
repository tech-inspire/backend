import type { ConnectRouter, HandlerContext } from "@connectrpc/connect";
import {
  GetLikesCountRequest,
  GetUserLikedPostsRequest,
  HasUserLikedPostRequest,
  LikePostRequest,
  LikesService,
  UnlikePostRequest,
} from "inspire-api-contracts/api/gen/ts/likes/v1/likes_pb";

import {
  getLikesCount,
  hasUserLikedPost,
  likePost,
  unlikePost,
  getUserLikedPosts,
} from "../db/likesRepository";

import { userContextKey } from "./auth";

export default (router: ConnectRouter) =>
  router.service(LikesService, {
    async getLikesCount(req: GetLikesCountRequest) {
      const count = await getLikesCount(req.postId);
      return {
        likesCount: BigInt(count),
      };
    },

    async hasUserLikedPost(req: HasUserLikedPostRequest) {
      const liked = await hasUserLikedPost(req.userId, req.postId);
      return { liked };
    },

    async likePost(req: LikePostRequest, context: HandlerContext) {
      const user = context.values.get(userContextKey);
      await likePost(user.userID, req.postId);
      return {};
    },

    async unlikePost(req: UnlikePostRequest, context: HandlerContext) {
      const user = context.values.get(userContextKey);
      await unlikePost(user.userID, req.postId);
      return {};
    },

    async getUserLikedPosts(req: GetUserLikedPostsRequest) {
      const postIds = await getUserLikedPosts(
        req.userId,
        req.limit,
        req.offset,
      );
      return { postIds };
    },
  });
