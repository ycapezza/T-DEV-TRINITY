import { useState, useEffect, useMemo } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { productsApi } from "../../api";
import ProductModal from "../../components/ProductModal/ProductModal";
import styles from "./Products.module.css";

const Products = () => {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [filters, setFilters] = useState({
    name: "",
    category: "",
  });
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [productToEdit, setProductToEdit] = useState(null);

  const { data: { data: products = [], pagination } = {}, isLoading } =
    useQuery({
      queryKey: ["products", page, filters],
      queryFn: () =>
        productsApi.getProducts({
          page,
          page_size: 10,
          ...filters,
        }),
    });

  const availableCategories = useMemo(() => {
    const categorySet = new Set();
    products.forEach((product) => {
      product.categories.forEach((category) => categorySet.add(category));
    });
    return Array.from(categorySet).sort();
  }, [products]);

  const deleteMutation = useMutation({
    mutationFn: (id) => productsApi.deleteProduct(id),
    onSuccess: () => {
      queryClient.invalidateQueries(["products"]);
    },
  });

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters((prev) => ({
      ...prev,
      [name]: value,
    }));
    setPage(1);
  };

  const handleDelete = async (id) => {
    if (window.confirm("Are you sure you want to delete this product?")) {
      await deleteMutation.mutate(id);
    }
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setProductToEdit(null);
  };

  return (
    <div className={styles.products}>
      <div className={styles.header}>
        <h1>Products</h1>
        <button
          className={styles.addButton}
          onClick={() => setIsModalOpen(true)}
        >
          Add Product
        </button>
      </div>

      <div className={styles.filters}>
        <div className={styles.filterGroup}>
          <input
            type="text"
            name="name"
            placeholder="Search by name..."
            value={filters.name}
            onChange={handleFilterChange}
            className={styles.filterInput}
          />
        </div>

        <div className={styles.filterGroup}>
          <select
            name="category"
            value={filters.category}
            onChange={handleFilterChange}
            className={styles.filterInput}
          >
            <option value="">All Categories</option>
            {availableCategories.map((category) => (
              <option key={category} value={category}>
                {category}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className={styles.tableContainer}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Name</th>
              <th>Categories</th>
              <th>Price</th>
              <th>Stock</th>
              <th>Brand</th>
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
            ) : products.length === 0 ? (
              <tr>
                <td colSpan="6" className={styles.noData}>
                  No products found
                </td>
              </tr>
            ) : (
              products.map((product) => (
                <tr key={product.id}>
                  <td>{product.name}</td>
                  <td>
                    <div className={styles.categories}>
                      {product.categories.map((category, index) => (
                        <span key={index} className={styles.category}>
                          {category}
                        </span>
                      ))}
                    </div>
                  </td>
                  <td>€{product.price.toFixed(2)}</td>
                  <td>
                    <span
                      className={
                        product.stock_quantity <= 10 ? styles.lowStock : ""
                      }
                    >
                      {product.stock_quantity}
                    </span>
                  </td>
                  <td>{product.brand}</td>
                  <td>
                    <div className={styles.actions}>
                      <button
                        className={styles.editButton}
                        onClick={() => {
                          setProductToEdit(product);
                          setIsModalOpen(true);
                        }}
                      >
                        Edit
                      </button>
                      <button
                        className={styles.deleteButton}
                        onClick={() => handleDelete(product.id)}
                        disabled={deleteMutation.isPending}
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {pagination && (
        <div className={styles.pagination}>
          <button
            disabled={page === 1}
            onClick={() => setPage(page - 1)}
            className={styles.pageButton}
          >
            Previous
          </button>
          <span className={styles.pageInfo}>
            Page {page} of {Math.ceil(pagination.total / pagination.page_size)}
          </span>
          <button
            disabled={page * pagination.page_size >= pagination.total}
            onClick={() => setPage(page + 1)}
            className={styles.pageButton}
          >
            Next
          </button>
        </div>
      )}

      {isModalOpen && (
        <ProductModal product={productToEdit} onClose={handleCloseModal} />
      )}
    </div>
  );
};

export default Products;
