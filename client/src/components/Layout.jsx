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
        <h2>Daytrack</h2>
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
