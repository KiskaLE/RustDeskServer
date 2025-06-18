import axios, { AxiosError, type InternalAxiosRequestConfig } from "axios";
import { refreshSessionTokens } from "~/client/actions/session";
import { env } from "~/env";
import { auth, signOut } from "~/server/auth";

let refreshPromise: Promise<{ token: string; refreshToken: string }> | null =
  null;

const axiosApi = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: false,
});

axiosApi.interceptors.request.use(async (config) => {
  // get API JWT from session
  const session = await auth();
  if (session?.user?.apiToken) {
    config.headers.Authorization = `Bearer ${session.user.apiToken}`;
  }
  return config;
});

axiosApi.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    const session = await auth();

    if (
      error.response?.status === 401 &&
      !originalRequest._retry &&
      session?.user?.refreshToken
    ) {
      originalRequest._retry = true;
      try {
        refreshPromise =
          refreshPromise ??
          axios
            .post(`${env.NEXT_PUBLIC_API_URL}/api/v1/account/refresh-token`, {
              refresh_token: session.user.refreshToken,
            })
            .then((r) => {
              return {
                token: r.data.token,
                refreshToken: r.data.refresh_token,
              };
            })
            .finally(() => {
              refreshPromise = null;
            });

        const newTokens = await refreshPromise;

        refreshSessionTokens(newTokens);

        refreshPromise = null;

        // update header on both the cached request and the API instance
        axiosApi.defaults.headers.common.Authorization = `Bearer ${newTokens.token}`;
        if (typeof originalRequest.headers.set === "function") {
          // Axios 1.x path
          originalRequest.headers.set(
            "Authorization",
            `Bearer ${newTokens.token}`,
          );
        } else {
          // Fallback for older Axios versions
          originalRequest.headers["Authorization"] =
            `Bearer ${newTokens.token}`;
        }
        return axiosApi(originalRequest);
      } catch (error) {
        refreshPromise = null;
        signOut();
        return Promise.reject(error);
      }
    }

    return Promise.reject(error);
  },
);

export default axiosApi;
