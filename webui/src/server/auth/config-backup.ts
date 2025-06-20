import { type DefaultSession, type NextAuthConfig } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";
import { env } from "~/env";

/**
 * Module augmentation for `next-auth` types. Allows us to add custom properties to the `session`
 * object and keep type safety.
 *
 * @see https://next-auth.js.org/getting-started/typescript#module-augmentation
 */
declare module "next-auth" {
  interface Session extends DefaultSession {
    user: {
      id: string;
      apiToken: string;
      refreshToken: string;
      apiTokenExp: Date;
      refreshTokenExp: Date;
    } & DefaultSession["user"];
    error?: string;
  }

  interface User {
    apiToken: string;
    refreshToken: string;
    apiTokenExp: Date;
    refreshTokenExp: Date;
  }
}

async function refreshAccessToken({ token }: { token: any }) {
  console.log("refresh access token", token.refreshToken);
  try {
    const response = await fetch(
      `${env.NEXT_PUBLIC_API_URL}/api/v1/account/refresh-token`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          refresh_token: token.refreshToken,
        }),
      },
    );

    const refreshedTokens = await response.json();

    if (!response.ok) {
      throw refreshedTokens;
    }

    return {
      ...token,
      apiToken: refreshedTokens.token,
      refreshToken: refreshedTokens.refresh_token,
      apiTokenExp: refreshedTokens.token_exp,
      refreshTokenExp: refreshedTokens.refresh_token_exp,
      error: undefined,
    };
  } catch (error) {
    console.log("Token refresh failed:", error);
    return { ...token, error: "RefreshAccessTokenError" };
  }
}

/**
 * Options for NextAuth.js used to configure adapters, providers, callbacks, etc.
 *
 * @see https://next-auth.js.org/configuration/options
 */
export const authConfig = {
  providers: [
    CredentialsProvider({
      name: "Credentials",
      credentials: {
        email: { label: "Email", type: "text", placeholder: "your@email.com" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        const apiUrl = `${env.NEXT_PUBLIC_API_URL}/api/v1/account/login`;
        try {
          const response = await fetch(apiUrl, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              email: credentials.email,
              password: credentials.password,
            }),
          });

          const data = await response.json();

          if (response.ok) {
            return {
              id: credentials.email as string,
              apiToken: data.token,
              refreshToken: data.refresh_token,
              apiTokenExp: data.token_exp,
              refreshTokenExp: data.refresh_token_exp,
            };
          }
          return null;
        } catch (error) {
          console.log("Authorization failed:", error);
          return null;
        }
      },
    }),
  ],
  session: {
    strategy: "jwt",
    maxAge: 60 * 60, // 1 hour
  },
  callbacks: {
    async jwt({ token, user, account }) {
      // Initial sign in
      if (account && user) {
        return {
          ...token,
          sub: user.id,
          apiToken: user.apiToken,
          refreshToken: user.refreshToken,
          apiTokenExp: user.apiTokenExp,
          refreshTokenExp: user.refreshTokenExp,
        };
      }

      // Return previous token if the access token has not expired yet
      if (Date.now() < new Date(token.apiTokenExp as string).getTime()) {
        return token;
      }

      // Access token has expired, try to refresh it
      return await refreshAccessToken({ token });
    },
    session: ({ session, token }) => {
      if (token.error) {
        return {
          ...session,
          error: token.error as string,
          user: {
            ...session.user,
            id: token.sub,
          },
        };
      }

      return {
        ...session,
        user: {
          ...session.user,
          id: token.sub,
          apiToken: token.apiToken as string,
          refreshToken: token.refreshToken as string,
          apiTokenExp: token.apiTokenExp as Date,
          refreshTokenExp: token.refreshTokenExp as Date,
        },
      };
    },
  },
  pages: {
    signIn: "/auth/signin", // Optional: custom sign-in page
    error: "/auth/error", // Optional: custom error page
  },
} satisfies NextAuthConfig;
