import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { usersApi } from "../../api";
import UserModal from "../../components/UserModal/UserModal";
import styles from "./Users.module.css";

const Users = () => {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [userToEdit, setUserToEdit] = useState(null);
  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useQuery({
    queryKey: ["users", page],
    queryFn: () => usersApi.getUsers(page, 10),
  });

  const deleteMutation = useMutation({
    mutationFn: (userId) => usersApi.deleteUser(userId),
    onSuccess: () => {
      queryClient.invalidateQueries(["users"]);
    },
  });

  const handleEdit = (user) => {
    setUserToEdit(user);
    setIsModalOpen(true);
  };

  const handleDelete = async (userId) => {
    if (window.confirm("Are you sure you want to delete this user?")) {
      await deleteMutation.mutate(userId);
    }
  };

  const filteredUsers = data?.data.filter((user) => {
    if (!searchQuery) return true;
    const searchLower = searchQuery.toLowerCase();
    return (
      user.first_name.toLowerCase().includes(searchLower) ||
      user.last_name.toLowerCase().includes(searchLower) ||
      user.email.toLowerCase().includes(searchLower)
    );
  });

  return (
    <div className={styles.users}>
      <div className={styles.header}>
        <h1>Users</h1>
        <button
          className={styles.addButton}
          onClick={() => setIsModalOpen(true)}
        >
          Add User
        </button>
      </div>

      <div className={styles.toolbar}>
        <input
          type="text"
          placeholder="Search users..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className={styles.searchInput}
        />
      </div>

      <div className={styles.tableContainer}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Phone</th>
              <th>Location</th>
              <th>Role</th>
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
            ) : filteredUsers?.length === 0 ? (
              <tr>
                <td colSpan="6" className={styles.noData}>
                  No users found
                </td>
              </tr>
            ) : (
              filteredUsers?.map((user) => (
                <tr key={user.id}>
                  <td>
                    {user.first_name} {user.last_name}
                  </td>
                  <td>{user.email}</td>
                  <td>{user.phone_number || "-"}</td>
                  <td>{user.city ? `${user.city}, ${user.country}` : "-"}</td>
                  <td>
                    <span
                      className={
                        user.is_admin ? styles.adminBadge : styles.userBadge
                      }
                    >
                      {user.is_admin ? "Admin" : "User"}
                    </span>
                  </td>
                  <td>
                    <div className={styles.actions}>
                      <button
                        className={styles.editButton}
                        onClick={() => handleEdit(user)}
                      >
                        Edit
                      </button>
                      <button
                        className={styles.deleteButton}
                        onClick={() => handleDelete(user.id)}
                        disabled={deleteMutation.isPending || user.is_admin}
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

      {isModalOpen && (
        <UserModal
          user={userToEdit}
          onClose={() => {
            setIsModalOpen(false);
            setUserToEdit(null);
          }}
        />
      )}
    </div>
  );
};

export default Users;
