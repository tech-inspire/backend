// auth.ts
import { createRemoteJWKSet, jwtVerify } from "jose";
import {
  Code,
  ConnectError,
  createContextKey,
  Interceptor,
} from "@connectrpc/connect";

const JWKS_URL = new URL(
  "http://auth-service-1:5080/auth/.well-known/jwks.json",
);

type User = { userID: string; isAnonymous: boolean };

export const userContextKey = createContextKey<User>(
  { userID: "null", isAnonymous: true }, // Default value
);

export function authInterceptor(
  jwksUrl: string,
  allowedMethods: string[] = [],
): Interceptor {
  const jwks = createRemoteJWKSet(JWKS_URL);

  return (next) => async (request) => {
    if (allowedMethods.includes(request.method.name))
      return await next(request);

    const authHeader = request.header.get("Authorization");

    if (!authHeader?.startsWith("Bearer ")) {
      throw new ConnectError(
        "Missing or invalid Authorization header",
        Code.Unauthenticated,
      );
    }

    const token = authHeader.split(" ")[1];

    try {
      const { payload } = await jwtVerify(token, jwks, {
        algorithms: ["EdDSA"],
      });
      if (!payload.sub) {
        throw new ConnectError(
          "Missing or invalid sub field",
          Code.Unauthenticated,
        );
      }
      const user: User = {
        userID: payload.sub,
        isAnonymous: false,
      };
      request.contextValues.set(userContextKey, user);
    } catch (error) {
      let message = "Authentication failed";
      if (error instanceof Error) {
        message = error.message;
      }

      throw new ConnectError(message, Code.Unauthenticated);
    }

    return await next(request);
  };
}
