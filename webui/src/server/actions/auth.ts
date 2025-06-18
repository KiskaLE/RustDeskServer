"use server";
import { cookies } from "next/headers";

export async function storeTokens(tokens: {
  token: string;
  refreshToken: string;
}) {
  const cookieStore = cookies();
  (await cookieStore).set("apiToken", tokens.token, {
    httpOnly: true,
    path: "/",
  });
  (await cookieStore).set("refreshToken", tokens.refreshToken, {
    httpOnly: true,
    path: "/",
  });
  return true;
}
