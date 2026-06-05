import { Link, Outlet, useNavigate } from 'react-router-dom';
import { clearToken, isAuthenticated } from '../api/auth';
import './Layout.css';

const API_KEY_STORAGE = 'daytrack_api_key';

export default function Layout() {
  const navigate = useNavigate();

  function handleLogout() {
    clearToken();
    localStorage.removeItem(API_KEY_STORAGE);
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
        <nav className="sidebar-nav">
          <Link to="/">Dashboard</Link>
          <Link to="/api-keys">API Keys</Link>
        </nav>
        <button className="logout-btn" onClick={handleLogout}>Log out</button>
      </aside>
      <div className="layout-body">
        <main className="content">
          <Outlet />
        </main>
        <nav className="mobile-nav">
          <Link to="/" className={({ isActive }) => isActive ? 'active' : ''}>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
            <span>Dashboard</span>
          </Link>
          <Link to="/api-keys" className={({ isActive }) => isActive ? 'active' : ''}>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>
            <span>API Keys</span>
          </Link>
        </nav>
      </div>
    </div>
  );
}
