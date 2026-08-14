import api from "./axios";

export const productsApi = {
  getProducts: async (params) => {
    const response = await api.get("/api/admin/products", { params });
    return {
      data: response.data.data || [],
      pagination: response.data.pagination,
    };
  },
  getProduct: (id) => api.get(`/api/admin/products/${id}`),
  createProduct: (productData) => api.post("/api/admin/products", productData),
  createProductByBarcode: (barcodeData) =>
    api.post("/api/admin/products/barcode", barcodeData),
  updateProduct: (id, productData) =>
    api.put(`/api/admin/products/${id}`, productData),
  deleteProduct: (id) => api.delete(`/api/admin/products/${id}`),
  searchProducts: (name) =>
    api.get("/api/admin/products/search", { params: { name } }),
};
