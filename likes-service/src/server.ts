import { fastify } from "fastify";
import { fastifyConnectPlugin } from "@connectrpc/connect-fastify";
import { cors as connectCors } from "@connectrpc/connect";
import fastifyCors from "@fastify/cors";
import routes from "./api/likes";
import { authInterceptor } from "./api/auth";

if (process.argv[1] === new URL(import.meta.url).pathname) {
  const PORT = parseInt(process.env.PORT || "40051");
  const HOST = process.env.HOST || "localhost";
  const server = await build();
  await server.listen({ host: HOST, port: PORT });
}

export async function build() {
  const server = fastify();

  server.get("/health", async () => {
    return { status: "ok" };
  });

  await server.register(fastifyCors, {
    origin: true,
    methods: [...connectCors.allowedMethods],
    allowedHeaders: ["Authorization", ...connectCors.allowedHeaders],
    exposedHeaders: [...connectCors.exposedHeaders],
  });

  const JWKS_URL = process.env.JWKS_URL;
  if (!JWKS_URL) {
    throw new Error("Invalid JWKS_URL");
  }

  const allowedProcedures: string[] = ["GetLikesCount", "GetUserLikedPosts"];

  await server.register(fastifyConnectPlugin, {
    routes,
    interceptors: [authInterceptor(JWKS_URL, allowedProcedures)],
  });

  return server;
}
