import styles from "./Charts.module.css";

export const KPICards = ({ salesReport, stockAlertsCount }) => {
  return (
    <div className={styles.kpiGrid}>
      <div className={styles.kpiCard}>
        <h3>Revenue</h3>
        <p className={styles.kpiValue}>
          €{salesReport?.total_revenue?.toLocaleString() ?? "0"}
        </p>
      </div>
      <div className={styles.kpiCard}>
        <h3>Orders</h3>
        <p className={styles.kpiValue}>
          {salesReport?.total_orders?.toLocaleString() ?? "0"}
        </p>
      </div>
      <div className={styles.kpiCard}>
        <h3>Average Order Value</h3>
        <p className={styles.kpiValue}>
          €{salesReport?.average_order_size?.toFixed(2) ?? "0"}
        </p>
      </div>
      <div className={styles.kpiCard}>
        <h3>Low Stock Items</h3>
        <p className={styles.kpiValue}>{stockAlertsCount}</p>
      </div>
    </div>
  );
};
