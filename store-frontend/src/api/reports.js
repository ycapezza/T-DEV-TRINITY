import api from "./axios";

export const reportsApi = {
  getSalesReport: async (startDate, endDate) => {
    const response = await api.get("/api/admin/reports/sales", {
      params: { start_date: startDate, end_date: endDate },
    });
    return response;
  },
  getTopProducts: async (limit = 5) => {
    const response = await api.get("/api/admin/reports/top-products", {
      params: { limit },
    });
    return response;
  },
  getCategoryPerformance: async () => {
    const response = await api.get("/api/admin/reports/categories");
    const categories = new Map();

    response.data?.forEach((product) => {
      product.categories?.forEach((category) => {
        if (!categories.has(category)) {
          categories.set(category, {
            category: category,
            total_sales: product.total_sales || 0,
            order_count: product.order_count || 0,
          });
        } else {
          const existing = categories.get(category);
          categories.set(category, {
            category: category,
            total_sales:
              (existing.total_sales || 0) + (product.total_sales || 0),
            order_count:
              (existing.order_count || 0) + (product.order_count || 0),
          });
        }
      });
    });

    return { data: Array.from(categories.values()) };
  },
  getStockAlerts: async (minimumStock = 10) => {
    const response = await api.get("/api/admin/reports/stock-alerts", {
      params: { minimum_stock: minimumStock },
    });
    return response;
  },
  getSalesEvolution: async (period = "daily") => {
    const response = await api.get("/api/admin/reports/sales-evolution", {
      params: { period },
    });
    return response;
  },
};
