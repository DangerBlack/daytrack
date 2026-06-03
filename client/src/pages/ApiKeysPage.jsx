import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { listApiKeys, createApiKey, deleteApiKey } from '../api/client';
import { getToken, getUser } from '../api/auth';
import './Dashboard.css';

export default function ApiKeysPage() {
  const queryClient = useQueryClient();
  const token = getToken();
  const user = getUser();
  const [name, setName] = useState('');
  const [showCreate, setShowCreate] = useState(false);
  const [newKeyData, setNewKeyData] = useState(null);
  const [copied, setCopied] = useState(null);
  const [tab, setTab] = useState('keys');
  const [deleteConfirm, setDeleteConfirm] = useState(null);

  const { data } = useQuery({
    queryKey: ['apiKeys'],
    queryFn: () => listApiKeys(token),
  });

  const keys = data?.items || [];

  const createMutation = useMutation({
    mutationFn: (n) => createApiKey(n, token),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['apiKeys'] });
      setNewKeyData(data);
      setName('');
      setShowCreate(false);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id) => deleteApiKey(id, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apiKeys'] });
      setDeleteConfirm(null);
    },
  });

  async function copyKey(key) {
    try {
      await navigator.clipboard.writeText(key);
      setCopied(key);
      setTimeout(() => setCopied(null), 2000);
    } catch {}
  }

  const baseUrl = `${window.location.protocol}//${window.location.host}`;
  const trackExample = `${baseUrl}/v1/events/${user?.username}/YOUR_TRACK_NAME?key=YOUR_API_KEY&quantity=1`;
  const trackWithDate = `${baseUrl}/v1/events/${user?.username}/YOUR_TRACK_NAME?key=YOUR_API_KEY&quantity=1&created_at=2026-06-03T10:00:00Z`;
  const listDay = `${baseUrl}/v1/events/${user?.username}/YOUR_TRACK_NAME?key=YOUR_API_KEY&list_by=day`;
  const listRaw = `${baseUrl}/v1/events/${user?.username}/YOUR_TRACK_NAME?key=YOUR_API_KEY&list_by=raw`;

  return (
    <div className="dashboard">
      <div className="header-row">
        <h1>API Keys</h1>
        <div className="header-actions">
          <button onClick={() => setTab('keys')} className={`btn btn-sm ${tab === 'keys' ? 'btn-active' : ''}`}>
            Keys
          </button>
          <button onClick={() => setTab('usage')} className={`btn btn-sm ${tab === 'usage' ? 'btn-active' : ''}`}>
            Usage
          </button>
        </div>
      </div>

      {tab === 'keys' && (
        <>
          <button onClick={() => setShowCreate(true)} className="btn btn-primary" style={{ marginBottom: '1rem' }}>
            + New Key
          </button>

          {showCreate && (
            <form className="create-form" onSubmit={e => {
              e.preventDefault();
              if (name.trim()) createMutation.mutate(name.trim());
            }}>
              <input
                type="text"
                placeholder="Key name"
                value={name}
                onChange={e => setName(e.target.value)}
                required
              />
              <button type="submit" className="btn btn-primary" disabled={createMutation.isPending}>
                Create
              </button>
              <button type="button" className="btn" onClick={() => setShowCreate(false)}>
                Cancel
              </button>
            </form>
          )}

          {newKeyData && (
            <div className="modal-overlay" onClick={() => setNewKeyData(null)}>
              <div className="modal" onClick={e => e.stopPropagation()}>
                <h3>API Key Created</h3>
                <p className="modal-warning">
                  This is the <strong>only time</strong> you will see this key. Save it somewhere safe.
                </p>
                <div className="modal-key-row">
                  <code className="modal-key">{newKeyData.key}</code>
                  <button className="btn btn-sm btn-primary" onClick={() => copyKey(newKeyData.key)}>
                    {copied === newKeyData.key ? 'Copied!' : 'Copy'}
                  </button>
                </div>
                <button className="btn" onClick={() => setNewKeyData(null)} style={{ marginTop: '1rem', width: '100%' }}>
                  I've saved the key
                </button>
              </div>
            </div>
          )}

          <div className="track-list">
            {keys.map(key => (
              <div key={key.id} className="track-card">
                <div className="track-header" style={{ cursor: 'default' }}>
                  <div>
                    <strong>{key.name}</strong>
                    <code className="key-value">{key.key}</code>
                  </div>
                  <button
                    onClick={() => setDeleteConfirm(key)}
                    className="btn btn-danger btn-sm"
                  >
                    Delete
                  </button>
                </div>
                <div className="track-detail">
                  <small>Created: {new Date(key.created_at).toLocaleString()}</small>
                </div>
              </div>
            ))}
            {keys.length === 0 && <p className="empty">No API keys yet.</p>}
          </div>

          {deleteConfirm && (
            <div className="modal-overlay" onClick={() => setDeleteConfirm(null)}>
              <div className="modal" onClick={e => e.stopPropagation()}>
                <h3>Delete API Key</h3>
                <p className="modal-warning">
                  Are you sure you want to delete the key <strong>{deleteConfirm.name}</strong>?<br />
                  This action cannot be undone. Any services using this key will stop working.
                </p>
                <div className="modal-actions">
                  <button className="btn" onClick={() => setDeleteConfirm(null)}>
                    Cancel
                  </button>
                  <button className="btn btn-danger" onClick={() => deleteMutation.mutate(deleteConfirm.id)}>
                    {deleteMutation.isPending ? 'Deleting...' : 'Yes, delete it'}
                  </button>
                </div>
              </div>
            </div>
          )}
        </>
      )}

      {tab === 'usage' && (
        <div className="usage-section">
          <h3>Using the API Key</h3>
          <p className="usage-subtitle">
            The API key is passed as a query parameter <code>key</code>. All event endpoints are unauthenticated (API key only).
          </p>

          <h4>Track an event (curl)</h4>
          <pre className="code-block">{`curl -X POST "${trackExample}"`}</pre>

          <h4>Track an event with a custom timestamp</h4>
          <pre className="code-block">{`curl -X POST "${trackWithDate}"`}</pre>

          <h4>List events (grouped by day)</h4>
          <pre className="code-block">{`curl "${listDay}"`}</pre>

          <h4>List raw events (each event individually)</h4>
          <pre className="code-block">{`curl "${listRaw}"`}</pre>
        </div>
      )}
    </div>
  );
}
