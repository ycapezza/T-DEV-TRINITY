import { Outlet, NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../../utils/AuthContext";
import styles from "./Layout.module.css";

const Layout = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <div className={styles.layout}>
      <aside className={styles.sidebar}>
        <div className={styles.sidebarHeader}>
          <h1>Store Admin</h1>
        </div>

        <nav className={styles.nav}>
          <NavLink
            to="/"
            end
            className={({ isActive }) =>
              isActive ? styles.activeLink : styles.link
            }
          >
            Dashboard
          </NavLink>

          <NavLink
            to="/products"
            className={({ isActive }) =>
              isActive ? styles.activeLink : styles.link
            }
          >
            Products
          </NavLink>

          <NavLink
            to="/users"
            className={({ isActive }) =>
              isActive ? styles.activeLink : styles.link
            }
          >
            Users
          </NavLink>

          <NavLink
            to="/invoices"
            className={({ isActive }) =>
              isActive ? styles.activeLink : styles.link
            }
          >
            Invoices
          </NavLink>

          <NavLink
            to="/reports"
            className={({ isActive }) =>
              isActive ? styles.activeLink : styles.link
            }
          >
            Reports
          </NavLink>
        </nav>
      </aside>

      <main className={styles.main}>
        <header className={styles.header}>
          <div className={styles.userInfo}>
            <span>
              {user?.first_name} {user?.last_name}
            </span>
            <button onClick={handleLogout} className={styles.logoutButton}>
              Logout
            </button>
          </div>
        </header>

        <div className={styles.content}>
          <Outlet />
        </div>
      </main>
    </div>
  );
};

export default Layout;
