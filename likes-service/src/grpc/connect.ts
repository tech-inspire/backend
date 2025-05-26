import type { ConnectRouter } from "@connectrpc/connect";
import { LikesService } from "inspire-api-contracts/api/gen/ts/likes/v1/likes_pb";
import LikesServiceHandler from "./likesService";

export default (router: ConnectRouter) =>
  router.service(LikesService, new LikesServiceHandler());
