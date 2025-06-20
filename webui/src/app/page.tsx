import axiosApi from "~/lib/axiosApi";
import { auth } from "~/server/auth";

export default async function Home() {
  const session = await auth();
  const computers = await axiosApi.get("/api/v1/computers");
  return (
    <div>
      <h1>Home</h1>
      <p>{session?.user.refreshToken}</p>
    </div>
  );
}
