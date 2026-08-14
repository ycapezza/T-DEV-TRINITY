import { useMutation, useQueryClient } from "@tanstack/react-query";
import { invoicesApi } from "../../api";
import InvoiceStatusBadge from "../InvoiceStatusBadge/InvoiceStatusBadge";
import styles from "./InvoiceDetailModal.module.css";

const InvoiceDetailModal = ({ invoice, onClose }) => {
  const queryClient = useQueryClient();

  const updateStatusMutation = useMutation({
    mutationFn: (status) => invoicesApi.updateInvoice(invoice.id, { status }),
    onSuccess: () => {
      queryClient.invalidateQueries(["invoices"]);
    },
  });

  const handleStatusChange = (e) => {
    updateStatusMutation.mutate(e.target.value);
  };

  return (
    <div className={styles.modalOverlay}>
      <div className={styles.modal}>
        <div className={styles.modalHeader}>
          <h2>Invoice #{invoice.id}</h2>
          <button className={styles.closeButton} onClick={onClose}>
            ×
          </button>
        </div>

        <div className={styles.modalContent}>
          <div className={styles.section}>
            <h3>Customer Information</h3>
            <div className={styles.info}>
              <p>
                <strong>Name:</strong> {invoice.user.first_name}{" "}
                {invoice.user.last_name}
              </p>
              <p>
                <strong>Email:</strong> {invoice.user.email}
              </p>
              <p>
                <strong>Phone:</strong> {invoice.user.phone_number || "-"}
              </p>
              {invoice.user.address && (
                <p>
                  <strong>Address:</strong> {invoice.user.address},{" "}
                  {invoice.user.city}, {invoice.user.country}
                </p>
              )}
            </div>
          </div>

          <div className={styles.section}>
            <div className={styles.statusHeader}>
              <h3>Status</h3>
              <select
                value={invoice.status}
                onChange={handleStatusChange}
                className={styles.statusSelect}
                disabled={updateStatusMutation.isPending}
              >
                <option value="pending">Pending</option>
                <option value="processing">Processing</option>
                <option value="completed">Completed</option>
                <option value="cancelled">Cancelled</option>
              </select>
            </div>
            <div className={styles.info}>
              <p>
                <strong>Date:</strong>{" "}
                {new Date(invoice.created_at).toLocaleDateString()}
              </p>
              <p>
                <strong>Current Status:</strong>{" "}
                <InvoiceStatusBadge status={invoice.status} />
              </p>
            </div>
          </div>

          <div className={styles.section}>
            <h3>Items</h3>
            <table className={styles.itemsTable}>
              <thead>
                <tr>
                  <th>Product</th>
                  <th>Quantity</th>
                  <th>Price</th>
                  <th>Total</th>
                </tr>
              </thead>
              <tbody>
                {invoice.items.map((item) => (
                  <tr key={item.product_id}>
                    <td>{item.product.name}</td>
                    <td>{item.quantity}</td>
                    <td>${item.price.toFixed(2)}</td>
                    <td>${(item.quantity * item.price).toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
              <tfoot>
                <tr>
                  <td colSpan="3">
                    <strong>Total</strong>
                  </td>
                  <td>
                    <strong>${invoice.total.toFixed(2)}</strong>
                  </td>
                </tr>
              </tfoot>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};

export default InvoiceDetailModal;
