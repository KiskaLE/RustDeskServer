"use client";

import { useSession, signOut } from "next-auth/react";
import { useEffect } from "react";

export default function Session() {
  const { data: session } = useSession();

  useEffect(() => {
    if (session?.error === "RefreshAccessTokenError") {
      // Clear the invalid session and redirect to sign in
      signOut();
    }
  }, [session]);
  return <></>;
}
