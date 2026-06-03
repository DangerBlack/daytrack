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
  const [trackError, setTrackError] = useState('');
  const [expandedTrack, setExpandedTrack] = useState(null);
  const [apiKey, setApiKey] = useState(null);
  const initRef = useRef(false);

  useEffect(() => {
    if (initRef.current) return;
    initRef.current = true;
    (async () => {
      const cached = localStorage.getItem(API_KEY_STORAGE);
      if (cached) {
        setApiKey(cached);
        return;
      }
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
      setTrackError('');
    },
    onError: (err) => {
      setTrackError(err.message);
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
            onChange={e => { setNewTrackName(e.target.value); setTrackError(''); }}
            required
            minLength={3}
            maxLength={255}
            pattern="[a-zA-Z0-9_-]+"
            title="Letters, numbers, underscores and hyphens only (3-255 characters)"
          />
          {trackError && <p className="form-error">{trackError}</p>}
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
  const [quickLoading, setQuickLoading] = useState(false);
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
    enabled: !!apiKey && isExpanded,
  });

  const days = eventsData?.items || [];
  const rawEvents = rawEventsData?.items || [];

  let daysSince = null;
  if (days.length > 0) {
    const lastDate = new Date(days[0].date);
    const diff = Math.floor((Date.now() - lastDate.getTime()) / 86400000);
    daysSince = diff;
  }

  function invalidateEvents() {
    queryClient.invalidateQueries({ queryKey: ['events', track.name] });
    queryClient.invalidateQueries({ queryKey: ['events-raw', track.name] });
  }

  async function handleTrack(e) {
    e.preventDefault();
    if (!apiKey) return;
    setEventLoading(true);
    try {
      const dt = eventDate ? new Date(eventDate) : null;
      await trackEvent(apiKey, user?.username, track.name, quantity, dt);
      invalidateEvents();
    } catch (err) {
      alert('Failed to register event: ' + err.message);
    } finally {
      setEventLoading(false);
    }
  }

  async function handleQuickTrack() {
    if (!apiKey) return;
    setQuickLoading(true);
    try {
      await trackEvent(apiKey, user?.username, track.name, 1, new Date());
      invalidateEvents();
    } catch (err) {
      alert('Failed to register event: ' + err.message);
    } finally {
      setQuickLoading(false);
    }
  }

  return (
    <div className={`track-card ${isExpanded ? 'expanded' : ''}`}>
      <div className="track-header">
        <div className="track-info">
          <strong>{track.name}</strong>
          <span className={`badge ${track.visibility}`}>{track.visibility}</span>
          {daysSince !== null && (
            <span className="days-since">{daysSince === 0 ? 'today' : `${daysSince}d ago`}</span>
          )}
        </div>
        <div className="track-actions">
          {!isExpanded && (
            <button onClick={handleQuickTrack} className="btn btn-sm btn-quick" disabled={quickLoading || !apiKey}>
              {quickLoading ? '...' : '+1'}
            </button>
          )}
          <button onClick={onToggle} className="btn btn-sm">
            {isExpanded ? 'Show less' : 'Show more'}
          </button>
          <button onClick={() => setDeleteConfirm(true)} className="btn btn-danger btn-sm">
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

          <div className="event-list">
            <h4>Recent Events</h4>
            {rawEvents.length === 0 && <p className="empty-state">No events yet.</p>}
            {rawEvents.slice(0, 20).map((ev, i) => (
              <div key={i} className="event-row">
                <span className="event-date">{new Date(ev.date).toLocaleString()}</span>
                <span className="event-qty">+{ev.quantity}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
