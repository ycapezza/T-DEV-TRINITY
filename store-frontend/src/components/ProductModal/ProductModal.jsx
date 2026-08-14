import { useState, useEffect } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { productsApi } from "../../api";
import styles from "./ProductModal.module.css";

const AVAILABLE_CATEGORIES = [
  "dairy",
  "beverages",
  "snacks",
  "produce",
  "bakery",
  "organic",
  "fruits",
  "fresh",
  "protein",
  "grains",
  "super foods",
  "yogurt",
  "nuts",
  "juice",
  "whole grain",
  "sweet",
  "bread",
];

const ProductModal = ({ product, onClose }) => {
  const queryClient = useQueryClient();
  const [mode, setMode] = useState("manual");
  const [formData, setFormData] = useState(
    product || {
      name: "",
      price: "",
      brand: "",
      categories: [],
      nutritional_info: "",
      stock_quantity: "",
    }
  );
  const [barcode, setBarcode] = useState("");
  const [selectedCategories, setSelectedCategories] = useState(
    new Set(product?.categories || [])
  );

  const mutation = useMutation({
    mutationFn: async (data) => {
      if (mode === "barcode" && !product) {
        return productsApi.createProductByBarcode({
          barcode: barcode,
          price: parseFloat(data.price),
          stock_quantity: parseInt(data.stock_quantity),
        });
      }

      const productData = {
        ...data,
        price: parseFloat(data.price),
        stock_quantity: parseInt(data.stock_quantity),
        categories: Array.from(selectedCategories),
      };

      if (product) {
        return productsApi.updateProduct(product.id, productData);
      }
      return productsApi.createProduct(productData);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["products"]);
      onClose();
    },
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    const submitData = {
      ...formData,
      categories: Array.from(selectedCategories),
    };
    mutation.mutate(submitData);
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const toggleCategory = (category) => {
    setSelectedCategories((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(category)) {
        newSet.delete(category);
      } else {
        newSet.add(category);
      }
      return newSet;
    });
  };

  return (
    <div className={styles.modalOverlay}>
      <div className={styles.modal}>
        <div className={styles.modalHeader}>
          <h2>{product ? "Edit Product" : "Add Product"}</h2>
          <button className={styles.closeButton} onClick={onClose}>
            ×
          </button>
        </div>

        {!product && (
          <div className={styles.modeToggle}>
            <button
              className={`${styles.modeButton} ${
                mode === "manual" ? styles.active : ""
              }`}
              onClick={() => setMode("manual")}
              type="button"
            >
              Manual Entry
            </button>
            <button
              className={`${styles.modeButton} ${
                mode === "barcode" ? styles.active : ""
              }`}
              onClick={() => setMode("barcode")}
              type="button"
            >
              Barcode
            </button>
          </div>
        )}

        <form onSubmit={handleSubmit} className={styles.form}>
          {mode === "barcode" && !product ? (
            <>
              <div className={styles.formGroup}>
                <label htmlFor="barcode">Barcode</label>
                <input
                  type="text"
                  id="barcode"
                  value={barcode}
                  onChange={(e) => setBarcode(e.target.value)}
                  required
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="price">Price</label>
                <input
                  type="number"
                  id="price"
                  name="price"
                  value={formData.price}
                  onChange={handleChange}
                  step="0.01"
                  required
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="stock_quantity">Stock Quantity</label>
                <input
                  type="number"
                  id="stock_quantity"
                  name="stock_quantity"
                  value={formData.stock_quantity}
                  onChange={handleChange}
                  required
                />
              </div>
            </>
          ) : (
            <>
              <div className={styles.formGroup}>
                <label htmlFor="name">Name</label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="price">Price</label>
                <input
                  type="number"
                  id="price"
                  name="price"
                  value={formData.price}
                  onChange={handleChange}
                  step="0.01"
                  required
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="brand">Brand</label>
                <input
                  type="text"
                  id="brand"
                  name="brand"
                  value={formData.brand}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className={styles.formGroup}>
                <label>Categories</label>
                <div className={styles.categoriesGrid}>
                  {AVAILABLE_CATEGORIES.map((category) => (
                    <label key={category} className={styles.categoryCheckbox}>
                      <input
                        type="checkbox"
                        checked={selectedCategories.has(category)}
                        onChange={() => toggleCategory(category)}
                      />
                      {category}
                    </label>
                  ))}
                </div>
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="stock_quantity">Stock Quantity</label>
                <input
                  type="number"
                  id="stock_quantity"
                  name="stock_quantity"
                  value={formData.stock_quantity}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="nutritional_info">Nutritional Info</label>
                <textarea
                  id="nutritional_info"
                  name="nutritional_info"
                  value={formData.nutritional_info}
                  onChange={handleChange}
                  rows="3"
                />
              </div>
            </>
          )}

          <div className={styles.formActions}>
            <button
              type="button"
              onClick={onClose}
              className={styles.cancelButton}
            >
              Cancel
            </button>
            <button
              type="submit"
              className={styles.submitButton}
              disabled={mutation.isPending}
            >
              {mutation.isPending ? "Saving..." : product ? "Update" : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ProductModal;
