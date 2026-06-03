import { getPublicKey, getPublicKeyAsync, signAsync, utils } from '@noble/ed25519';

export const AUTH_TOKEN_KEY = 'daytrack_token';
export const AUTH_USER_KEY = 'daytrack_user';

export function getToken() {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

export function setToken(token) {
  localStorage.setItem(AUTH_TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem(AUTH_USER_KEY);
}

export function getUser() {
  const raw = localStorage.getItem(AUTH_USER_KEY);
  return raw ? JSON.parse(raw) : null;
}

export function setUser(user) {
  localStorage.setItem(AUTH_USER_KEY, JSON.stringify(user));
}

export function isAuthenticated() {
  return !!getToken();
}

async function sha256(message) {
  const encoder = new TextEncoder();
  const data = encoder.encode(message);
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  return new Uint8Array(hashBuffer);
}

function bytesToBase64(bytes) {
  return btoa(String.fromCharCode(...new Uint8Array(bytes)));
}

export async function generateKeypair(password, salt) {
  const hash = await sha256(password + salt);
  const seed = hash.slice(0, 32);

  const privateKey = seed;
  const publicKey = await getPublicKeyAsync(privateKey);

  return {
    publicKey: bytesToBase64(publicKey),
    privateKey,
  };
}

export async function signChallenge(privateKey, challenge) {
  const encoder = new TextEncoder();
  const signature = await signAsync(encoder.encode(challenge), privateKey);
  return bytesToBase64(signature);
}
