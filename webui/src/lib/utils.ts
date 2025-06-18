import axios from "axios";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { env } from "~/env";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

type APITokens = {
  token: string;
  refreshToken: string;
};

export async function refreshAPIToken(
  refreshToken: string,
): Promise<APITokens | null> {
  const apiUrl = `${env.NEXT_PUBLIC_API_URL}/api/v1/account/refresh-token`;
  const res = await axios.post(apiUrl, {
    refresh_token: refreshToken,
  });

  if (res.status !== 200) {
    return null;
  }

  return {
    token: res.data.token,
    refreshToken: res.data.refresh_token,
  };
}
