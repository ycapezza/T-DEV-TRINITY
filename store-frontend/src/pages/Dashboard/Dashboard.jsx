import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { reportsApi } from "../../api";
import { SalesEvolutionChart } from "../../components/Charts/SalesEvolutionChart";
import { TopProductsChart } from "../../components/Charts/TopProductsChart";
import { CategoryDistributionChart } from "../../components/Charts/CategoryDistributionChart";
import { KPICards } from "../../components/Charts/KPICards";
import styles from "./Dashboard.module.css";

const Dashboard = () => {
  const [dateRange, setDateRange] = useState({
    startDate: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
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
    <div className={styles.dashboard}>
      <h1>Dashboard</h1>

      <KPICards
        salesReport={salesReport}
        stockAlertsCount={stockAlerts.length}
      />

      <div className={styles.chartGrid}>
        <SalesEvolutionChart data={salesEvolution} />
        <TopProductsChart data={topProducts} />
        <CategoryDistributionChart data={categories} />
      </div>
    </div>
  );
};

export default Dashboard;
