import { fastify } from "fastify";
import { fastifyConnectPlugin } from "@connectrpc/connect-fastify";
import { cors as connectCors } from "@connectrpc/connect";
import fastifyCors from "@fastify/cors";
import routes from "./grpc/connect";

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

  await server.register(fastifyConnectPlugin, { routes });

  return server;
}
