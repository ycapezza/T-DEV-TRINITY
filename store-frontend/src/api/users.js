import api from "./axios";

export const usersApi = {
  getUsers: (page = 1, pageSize = 10) =>
    api.get(`/api/admin/users?page=${page}&page_size=${pageSize}`),
  getUser: (id) => api.get(`/api/admin/users/${id}`),
  createUser: (userData) => api.post("/api/admin/users", userData),
  updateUser: (id, userData) => api.put(`/api/admin/users/${id}`, userData),
  deleteUser: (id) => api.delete(`/api/admin/users/${id}`),
};
