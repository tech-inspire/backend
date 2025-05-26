import redis from "./redisClient";

function userLikedPostsKey(userId: string): string {
  return `user:${userId}:liked_posts`;
}

function postLikedUsersKey(postId: string): string {
  return `post:${postId}:liked_users`;
}

function postLikesCountKey(postId: string): string {
  return `post:${postId}:likes_count`;
}

export async function getLikesCount(postId: string): Promise<number> {
  const countStr = await redis.get(postLikesCountKey(postId));
  return countStr ? parseInt(countStr, 10) : 0;
}

export async function hasUserLikedPost(
  userId: string,
  postId: string,
): Promise<boolean> {
  const score = await redis.zscore(userLikedPostsKey(userId), postId);
  return score !== null;
}

export async function likePost(userId: string, postId: string) {
  const userLikesKey = userLikedPostsKey(userId);
  const postLikesKey = postLikedUsersKey(postId);
  const countKey = postLikesCountKey(postId);

  const alreadyLiked = await redis.zscore(userLikesKey, postId);
  if (alreadyLiked !== null) {
    return;
  }

  const pipeline = redis.pipeline();
  pipeline.zadd(userLikesKey, Date.now(), postId);
  pipeline.sadd(postLikesKey, userId.toString());
  pipeline.incr(countKey);
  await pipeline.exec();
}

export async function unlikePost(userId: string, postId: string) {
  const userLikesKey = userLikedPostsKey(userId);
  const postLikesKey = postLikedUsersKey(postId);
  const countKey = postLikesCountKey(postId);

  const liked = await redis.zscore(userLikesKey, postId);
  if (liked === null) {
    return;
  }

  const pipeline = redis.pipeline();
  pipeline.zrem(userLikesKey, postId);
  pipeline.srem(postLikesKey, userId.toString());
  pipeline.decr(countKey);
  await pipeline.exec();
}

export async function getUserLikedPosts(
  userId: string,
  limit: number,
  offset: number,
): Promise<string[]> {
  return redis.zrevrange(userLikedPostsKey(userId), offset, offset + limit - 1);
}
