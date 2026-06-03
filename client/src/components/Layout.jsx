import { Link, Outlet, useNavigate } from 'react-router-dom';
import { clearToken, isAuthenticated } from '../api/auth';
import './Layout.css';

export default function Layout() {
  const navigate = useNavigate();

  function handleLogout() {
    clearToken();
    navigate('/login');
  }

  if (!isAuthenticated()) {
    return null;
  }

  return (
    <div className="layout">
      <aside className="sidebar">
        <div className="sidebar-brand">
          <img src="/logo.svg" alt="Daytrack" className="sidebar-logo" />
          <h2>Daytrack</h2>
        </div>
        <nav>
          <Link to="/">Dashboard</Link>
          <Link to="/api-keys">API Keys</Link>
        </nav>
        <button className="logout-btn" onClick={handleLogout}>Log out</button>
      </aside>
      <main className="content">
        <Outlet />
      </main>
    </div>
  );
}
