"use client";

import { useSession } from "next-auth/react";

type SessionTokens = {
  token: string;
  refreshToken: string;
};

export function refreshSessionTokens({ token, refreshToken }: SessionTokens) {
  const { data: session, status, update } = useSession();
  update({
    apiToken: token,
    refreshToken: refreshToken,
  });
}
