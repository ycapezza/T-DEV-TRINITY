import api from "./axios";

export const invoicesApi = {
  getInvoices: (page = 1, pageSize = 10) =>
    api.get(`/api/admin/invoices?page=${page}&page_size=${pageSize}`),
  getInvoice: (id) => api.get(`/api/admin/invoices/${id}`),
  updateInvoice: (id, invoiceData) =>
    api.put(`/api/admin/invoices/${id}`, invoiceData),
  deleteInvoice: (id) => api.delete(`/api/admin/invoices/${id}`),
};
