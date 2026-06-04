const BASE_URL = import.meta.env.VITE_API_URL || '';

function isAuthPath(path) {
  return path.startsWith('/v1/users/challenge') || path.startsWith('/v1/users/signin') || path.startsWith('/v1/users/signup');
}

async function request(path, options = {}) {
  const url = `${BASE_URL}${path}`;
  const config = {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  };

  if (options.token) {
    config.headers['Authorization'] = `Bearer ${options.token}`;
    delete config.token;
  }

  const res = await fetch(url, config);

  if (res.status === 401 && !isAuthPath(path)) {
    localStorage.removeItem('daytrack_token');
    localStorage.removeItem('daytrack_user');
    window.location.href = '/login';
    return;
  }

  const text = await res.text();

  let data;
  try {
    data = JSON.parse(text);
  } catch {
    data = text;
  }

  if (!res.ok) {
    throw new Error(data?.message || `Request failed: ${res.status}`);
  }

  return data;
}

export function generateChallenge(email) {
  return request(`/v1/users/challenge?email=${encodeURIComponent(email)}`);
}

export function signup(username, email, publicKey, salt) {
  return request('/v1/users/signup', {
    method: 'POST',
    body: JSON.stringify({ username, email, public_key: publicKey, salt }),
  });
}

export function signin(email, challenge, signedChallenge) {
  return request('/v1/users/signin', {
    method: 'POST',
    body: JSON.stringify({ email, challenge, signed_challenge: signedChallenge }),
  });
}

export function createTrack(name, token, visibility) {
  const body = { name };
  if (visibility) {
    body.visibility = visibility;
  }
  return request('/v1/tracks', {
    method: 'POST',
    body: JSON.stringify(body),
    token,
  });
}

export function listTracks(token) {
  return request('/v1/tracks', { token });
}

export function updateTrack(name, data, token) {
  return request(`/v1/tracks/${encodeURIComponent(name)}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
    token,
  });
}

export function deleteTrack(name, token) {
  return request(`/v1/tracks/${encodeURIComponent(name)}`, {
    method: 'DELETE',
    token,
  });
}

export function createApiKey(name, token) {
  return request(`/v1/api_keys?name=${encodeURIComponent(name)}`, {
    method: 'POST',
    token,
  });
}

export function listApiKeys(token) {
  return request('/v1/api_keys', { token });
}

export function deleteApiKey(id, token) {
  return request(`/v1/api_keys/${id}`, {
    method: 'DELETE',
    token,
  });
}

export function trackEvent(apiKey, username, trackName, quantity, createdAt) {
  let path = `/v1/events/${encodeURIComponent(username)}/${encodeURIComponent(trackName)}?key=${apiKey}&quantity=${quantity}`;
  if (createdAt) {
    path += `&created_at=${encodeURIComponent(createdAt.toISOString())}`;
  }
  return request(path, { method: 'POST' });
}

export function listEvents(apiKey, username, trackName, listBy) {
  let path = `/v1/events/${encodeURIComponent(username)}/${encodeURIComponent(trackName)}?key=${apiKey}`;
  if (listBy) {
    path += `&list_by=${listBy}`;
  }
  return request(path);
}


