import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080",
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(
  async (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    console.log("Request headers:", config.headers);

    if (config.method !== "get" && !config.headers["X-Csrf-Token"]) {
      try {
        const csrfResponse = await axios.get(
          "http://localhost:8080/auth/csrf-token",
          {
            withCredentials: true,
          }
        );
        const csrfToken = csrfResponse.headers["x-csrf-token"];
        if (csrfToken) {
          config.headers["X-Csrf-Token"] = csrfToken;
        }
      } catch (error) {
        console.error("Failed to get CSRF token:", error);
      }
    }

    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 403) {
      console.error("Auth error:", error.response);
    }
    return Promise.reject(error);
  }
);

export default api;
