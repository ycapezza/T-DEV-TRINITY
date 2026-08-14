import api from "./axios";

export const authApi = {
  login: async (credentials) => {
    try {
      const csrfResponse = await api.get("/auth/csrf-token");
      const csrfToken = csrfResponse.headers["x-csrf-token"];

      if (!csrfToken) {
        throw new Error("No CSRF token in response headers");
      }

      const response = await api.post("/auth/login/admin", credentials, {
        headers: {
          "X-Csrf-Token": csrfToken,
        },
      });

      console.log("Received token:", response.data.token);

      return response;
    } catch (error) {
      console.error("Login error:", error);
      throw error;
    }
  },
};
