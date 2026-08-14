import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { reportsApi } from "../../api";
import { SalesEvolutionChart } from "../../components/Charts/SalesEvolutionChart";
import { TopProductsChart } from "../../components/Charts/TopProductsChart";
import { CategoryDistributionChart } from "../../components/Charts/CategoryDistributionChart";
import { KPICards } from "../../components/Charts/KPICards";
import styles from "./Reports.module.css";

const Reports = () => {
  const [dateRange, setDateRange] = useState({
    startDate: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)
      .toISOString()
      .split("T")[0],
    endDate: new Date().toISOString().split("T")[0],
  });

  const { data: salesReport } = useQuery({
    queryKey: ["salesReport", dateRange],
    queryFn: () =>
      reportsApi.getSalesReport(dateRange.startDate, dateRange.endDate),
    select: (response) => response?.data,
  });

  const { data: topProducts = [] } = useQuery({
    queryKey: ["topProducts"],
    queryFn: () => reportsApi.getTopProducts(5),
    select: (response) => response?.data || [],
  });

  const { data: categories = [] } = useQuery({
    queryKey: ["categories"],
    queryFn: () => reportsApi.getCategoryPerformance(),
    select: (response) => response?.data || [],
  });

  const { data: stockAlerts = [] } = useQuery({
    queryKey: ["stockAlerts"],
    queryFn: () => reportsApi.getStockAlerts(),
    select: (response) => response?.data || [],
  });

  const { data: salesEvolution = [] } = useQuery({
    queryKey: ["salesEvolution"],
    queryFn: () => reportsApi.getSalesEvolution("daily"),
    select: (response) => response?.data || [],
  });

  return (
    <div className={styles.reports}>
      <div className={styles.header}>
        <h1>Reports & Analytics</h1>
        <div className={styles.dateFilters}>
          <input
            type="date"
            name="startDate"
            value={dateRange.startDate}
            onChange={(e) =>
              setDateRange((prev) => ({ ...prev, startDate: e.target.value }))
            }
            className={styles.dateInput}
          />
          <span>to</span>
          <input
            type="date"
            name="endDate"
            value={dateRange.endDate}
            onChange={(e) =>
              setDateRange((prev) => ({ ...prev, endDate: e.target.value }))
            }
            className={styles.dateInput}
          />
        </div>
      </div>

      <KPICards
        salesReport={salesReport}
        stockAlertsCount={stockAlerts.length}
      />

      <SalesEvolutionChart data={salesEvolution} />

      <div className={styles.chartGrid}>
        <TopProductsChart data={topProducts} />
        <CategoryDistributionChart data={categories} />
      </div>

      {stockAlerts.length > 0 && (
        <div className={styles.stockAlerts}>
          <h2>Low Stock Alerts</h2>
          <div className={styles.tableContainer}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Product</th>
                  <th>Current Stock</th>
                  <th>Minimum Stock</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {stockAlerts.map((alert) => (
                  <tr key={alert.product_id}>
                    <td>{alert.product_name}</td>
                    <td>{alert.current_stock}</td>
                    <td>{alert.minimum_stock}</td>
                    <td>
                      <span
                        className={
                          alert.current_stock === 0
                            ? styles.outOfStock
                            : styles.lowStock
                        }
                      >
                        {alert.current_stock === 0
                          ? "Out of Stock"
                          : "Low Stock"}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
};

export default Reports;
