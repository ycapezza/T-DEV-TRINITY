import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer } from "recharts";
import styles from "./Charts.module.css";

const COLORS = ["#2563eb", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6"];

export const CategoryDistributionChart = ({ data = [] }) => {
  if (!data || data.length === 0) return null;

  return (
    <div className={styles.chartContainer}>
      <h2>Category Distribution</h2>
      <div className={styles.chart}>
        <ResponsiveContainer width="100%" height={300}>
          <PieChart>
            <Pie
              data={data}
              dataKey="total_sales"
              nameKey="category"
              cx="50%"
              cy="50%"
              outerRadius={100}
              label={({ category, total_sales }) =>
                `€{category} (€€{Math.round(total_sales)})`
              }
            >
              {data.map((entry, index) => (
                <Cell
                  key={`cell-€{index}`}
                  fill={COLORS[index % COLORS.length]}
                />
              ))}
            </Pie>
            <Tooltip formatter={(value) => `€€{value.toFixed(2)}`} />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};
