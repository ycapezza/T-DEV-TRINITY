import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { invoicesApi } from "../../api";
import InvoiceStatusBadge from "../../components/InvoiceStatusBadge/InvoiceStatusBadge";
import InvoiceDetailModal from "../../components/InvoiceDetailModal/InvoiceDetailModal";
import styles from "./Invoices.module.css";

const Invoices = () => {
  const [page, setPage] = useState(1);
  const [selectedInvoice, setSelectedInvoice] = useState(null);
  const [filters, setFilters] = useState({
    status: "",
    dateRange: "all",
  });

  const { data, isLoading } = useQuery({
    queryKey: ["invoices", page],
    queryFn: () => invoicesApi.getInvoices(page, 10),
  });

  const handleViewInvoice = (invoice) => {
    setSelectedInvoice(invoice);
  };

  const filterInvoices = (invoices) => {
    if (!invoices) return [];

    return invoices.filter((invoice) => {
      let passesFilter = true;

      if (filters.status && invoice.status !== filters.status) {
        passesFilter = false;
      }

      if (filters.dateRange !== "all") {
        const invoiceDate = new Date(invoice.created_at);
        const today = new Date();
        const diffTime = Math.abs(today - invoiceDate);
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        switch (filters.dateRange) {
          case "today":
            passesFilter = diffDays <= 1;
            break;
          case "week":
            passesFilter = diffDays <= 7;
            break;
          case "month":
            passesFilter = diffDays <= 30;
            break;
          default:
            break;
        }
      }

      return passesFilter;
    });
  };

  const filteredInvoices = filterInvoices(data?.data);

  return (
    <div className={styles.invoices}>
      <div className={styles.header}>
        <h1>Invoices</h1>
        <div className={styles.filters}>
          <select
            value={filters.status}
            onChange={(e) =>
              setFilters((prev) => ({ ...prev, status: e.target.value }))
            }
            className={styles.filterSelect}
          >
            <option value="">All Statuses</option>
            <option value="pending">Pending</option>
            <option value="processing">Processing</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>

          <select
            value={filters.dateRange}
            onChange={(e) =>
              setFilters((prev) => ({ ...prev, dateRange: e.target.value }))
            }
            className={styles.filterSelect}
          >
            <option value="all">All Time</option>
            <option value="today">Today</option>
            <option value="week">Last 7 Days</option>
            <option value="month">Last 30 Days</option>
          </select>
        </div>
      </div>

      <div className={styles.tableContainer}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Invoice #</th>
              <th>Customer</th>
              <th>Date</th>
              <th>Status</th>
              <th>Total</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <tr>
                <td colSpan="6" className={styles.loading}>
                  Loading...
                </td>
              </tr>
            ) : filteredInvoices.length === 0 ? (
              <tr>
                <td colSpan="6" className={styles.noData}>
                  No invoices found
                </td>
              </tr>
            ) : (
              filteredInvoices.map((invoice) => (
                <tr key={invoice.id}>
                  <td>#{invoice.id}</td>
                  <td>
                    <div className={styles.customer}>
                      {invoice.user.first_name} {invoice.user.last_name}
                      <span className={styles.email}>{invoice.user.email}</span>
                    </div>
                  </td>
                  <td>{new Date(invoice.created_at).toLocaleDateString()}</td>
                  <td>
                    <InvoiceStatusBadge status={invoice.status} />
                  </td>
                  <td>€{invoice.total.toFixed(2)}</td>
                  <td>
                    <button
                      className={styles.viewButton}
                      onClick={() => handleViewInvoice(invoice)}
                    >
                      View Details
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {data?.pagination && (
        <div className={styles.pagination}>
          <button
            disabled={page === 1}
            onClick={() => setPage(page - 1)}
            className={styles.pageButton}
          >
            Previous
          </button>
          <span className={styles.pageInfo}>
            Page {page} of{" "}
            {Math.ceil(data.pagination.total / data.pagination.page_size)}
          </span>
          <button
            disabled={page * data.pagination.page_size >= data.pagination.total}
            onClick={() => setPage(page + 1)}
            className={styles.pageButton}
          >
            Next
          </button>
        </div>
      )}

      {selectedInvoice && (
        <InvoiceDetailModal
          invoice={selectedInvoice}
          onClose={() => setSelectedInvoice(null)}
        />
      )}
    </div>
  );
};

export default Invoices;
