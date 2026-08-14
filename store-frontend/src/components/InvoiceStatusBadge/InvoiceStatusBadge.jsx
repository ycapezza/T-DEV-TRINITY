import styles from "./InvoiceStatusBadge.module.css";

const statusColors = {
  pending: styles.pending,
  processing: styles.processing,
  completed: styles.completed,
  cancelled: styles.cancelled,
};

const InvoiceStatusBadge = ({ status }) => {
  return (
    <span className={`${styles.badge} ${statusColors[status.toLowerCase()]}`}>
      {status}
    </span>
  );
};

export default InvoiceStatusBadge;
