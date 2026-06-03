import { useState, useRef, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { listTracks, createTrack, deleteTrack, listEvents, trackEvent, createApiKey } from '../api/client';
import { getToken, getUser } from '../api/auth';
import Heatmap from '../components/Heatmap';
import './Dashboard.css';

const API_KEY_STORAGE = 'daytrack_api_key';

export default function DashboardPage() {
  const queryClient = useQueryClient();
  const token = getToken();
  const user = getUser();
  const [newTrackName, setNewTrackName] = useState('');
  const [showCreate, setShowCreate] = useState(false);
  const [expandedTrack, setExpandedTrack] = useState(null);
  const [apiKey, setApiKey] = useState(() => localStorage.getItem(API_KEY_STORAGE));
  const initRef = useRef(false);

  useEffect(() => {
    if (initRef.current || localStorage.getItem(API_KEY_STORAGE)) return;
    initRef.current = true;
    (async () => {
      try {
        const key = await createApiKey('web-client', token);
        localStorage.setItem(API_KEY_STORAGE, key.key);
        setApiKey(key.key);
      } catch {}
    })();
  }, [token]);

  const { data: tracksData } = useQuery({
    queryKey: ['tracks'],
    queryFn: () => listTracks(token),
  });

  const tracks = tracksData?.items || [];

  const createMutation = useMutation({
    mutationFn: (name) => createTrack(name, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tracks'] });
      setNewTrackName('');
      setShowCreate(false);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (name) => deleteTrack(name, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['tracks'] }),
  });

  return (
    <div className="dashboard">
      <div className="header-row">
        <h1>Dashboard</h1>
        <button onClick={() => setShowCreate(true)} className="btn btn-primary">
          + New Track
        </button>
      </div>

      {showCreate && (
        <form className="create-form" onSubmit={e => {
          e.preventDefault();
          if (newTrackName.trim()) createMutation.mutate(newTrackName.trim());
        }}>
          <input
            type="text"
            placeholder="Track name"
            value={newTrackName}
            onChange={e => setNewTrackName(e.target.value)}
            required
            minLength={3}
            pattern="[a-zA-Z0-9_-]+"
          />
          <button type="submit" className="btn btn-primary" disabled={createMutation.isPending}>
            Create
          </button>
          <button type="button" className="btn" onClick={() => setShowCreate(false)}>
            Cancel
          </button>
        </form>
      )}

      {tracks.length === 0 && !showCreate && (
        <p className="empty-state">No tracks yet. Create one to start tracking!</p>
      )}

      <div className="track-list">
        {tracks.map(track => (
          <TrackCard
            key={track.id}
            track={track}
            apiKey={apiKey}
            user={user}
            onDelete={() => deleteMutation.mutate(track.name)}
            isExpanded={expandedTrack?.id === track.id}
            onToggle={() => setExpandedTrack(expandedTrack?.id === track.id ? null : track)}
          />
        ))}
      </div>
    </div>
  );
}

function TrackCard({ track, apiKey, user, onDelete, isExpanded, onToggle }) {
  const [quantity, setQuantity] = useState(1);
  const [eventDate, setEventDate] = useState(() => new Date().toISOString().slice(0, 16));
  const [eventLoading, setEventLoading] = useState(false);
  const [viewMode, setViewMode] = useState('heatmap');
  const [deleteConfirm, setDeleteConfirm] = useState(false);
  const queryClient = useQueryClient();

  const { data: eventsData } = useQuery({
    queryKey: ['events', track.name],
    queryFn: () => listEvents(apiKey, user?.username, track.name, 'day'),
    enabled: !!apiKey,
  });

  const { data: rawEventsData } = useQuery({
    queryKey: ['events-raw', track.name],
    queryFn: () => listEvents(apiKey, user?.username, track.name, 'raw'),
    enabled: !!apiKey && isExpanded && viewMode === 'list',
  });

  const days = eventsData?.items || [];
  const rawEvents = rawEventsData?.items || [];

  function toLocalISO(d) {
    const pad = (n) => String(n).padStart(2, '0');
    return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
  }

  async function handleTrack(e) {
    e.preventDefault();
    if (!apiKey) return;
    setEventLoading(true);
    try {
      const dt = eventDate ? new Date(eventDate) : null;
      await trackEvent(apiKey, user?.username, track.name, quantity, dt);
      queryClient.invalidateQueries({ queryKey: ['events', track.name] });
      queryClient.invalidateQueries({ queryKey: ['events-raw', track.name] });
    } catch (err) {
      alert('Failed to register event: ' + err.message);
    } finally {
      setEventLoading(false);
    }
  }

  return (
    <div className={`track-card ${isExpanded ? 'expanded' : ''}`}>
      <div className="track-header" onClick={onToggle}>
        <div>
          <strong>{track.name}</strong>
          <span className={`badge ${track.visibility}`}>{track.visibility}</span>
        </div>
        <div className="track-actions">
          <button onClick={e => { e.stopPropagation(); setDeleteConfirm(true); }} className="btn btn-danger btn-sm">
            Delete
          </button>
        </div>
      </div>

      {deleteConfirm && (
        <div className="modal-overlay" onClick={() => setDeleteConfirm(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <h3>Delete Track</h3>
            <p className="modal-warning">
              Are you sure you want to delete <strong>{track.name}</strong>?<br />
              All events for this track will be lost. This cannot be undone.
            </p>
            <div className="modal-actions">
              <button className="btn" onClick={() => setDeleteConfirm(false)}>Cancel</button>
              <button className="btn btn-danger" onClick={() => { setDeleteConfirm(false); onDelete(); }}>Yes, delete it</button>
            </div>
          </div>
        </div>
      )}

      <div className="track-heatmap">
        <Heatmap days={days} />
      </div>

      {isExpanded && (
        <div className="track-detail">
          <form onSubmit={handleTrack} className="track-event-form">
            <input
              type="number"
              min="1"
              value={quantity}
              onChange={e => setQuantity(parseInt(e.target.value) || 1)}
              className="quantity-input"
              title="Quantity"
            />
            <input
              type="datetime-local"
              value={eventDate}
              onChange={e => setEventDate(e.target.value)}
              className="datetime-input"
            />
            <button type="submit" className="btn btn-primary btn-sm" disabled={eventLoading || !apiKey}>
              {eventLoading ? '...' : apiKey ? 'Track event' : 'Need API key'}
            </button>
          </form>

          <div className="view-toggle">
            <button
              className={`btn btn-sm ${viewMode === 'heatmap' ? 'btn-active' : ''}`}
              onClick={e => { e.stopPropagation(); setViewMode('heatmap'); }}
            >
              Heatmap
            </button>
            <button
              className={`btn btn-sm ${viewMode === 'list' ? 'btn-active' : ''}`}
              onClick={e => { e.stopPropagation(); setViewMode('list'); }}
            >
              Event List
            </button>
          </div>

          {viewMode === 'list' && (
            <div className="event-list">
              <h4>Recent Events</h4>
              {rawEvents.length === 0 && <p className="empty-state">No events yet.</p>}
              {rawEvents.map((ev, i) => (
                <div key={i} className="event-row">
                  <span className="event-date">{new Date(ev.date).toLocaleString()}</span>
                  <span className="event-qty">+{ev.quantity}</span>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
