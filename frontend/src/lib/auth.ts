import { PUBLIC_API_BASE } from '$env/static/public';
import { writable } from 'svelte/store';

export interface AuthUser {
  username: string;
}

export const user = writable<AuthUser | null>(null);

const API_BASE = PUBLIC_API_BASE;

/** Rehydrate user from stored tokens on app start */
export function initAuth() {
  const token = localStorage.getItem('access_token');
  const username = localStorage.getItem('username');
  if (token && username) {
    user.set({ username });
  }
}

export async function register(username: string, email: string, password: string): Promise<void> {
  const response = await fetch(`${API_BASE}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, email, password })
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.error ?? 'Registration failed');
  }

  const data = await response.json();
  const resolvedUsername = await resolveUsername(data.access_token, username);
  persistTokens(data.access_token, data.refresh_token, resolvedUsername);
  user.set({ username: resolvedUsername });
}

export async function login(username: string, password: string): Promise<void> {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identifier: username, password })
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.error ?? 'Login failed');
  }

  const data = await response.json();
  const resolvedUsername = await resolveUsername(data.access_token, username);
  persistTokens(data.access_token, data.refresh_token, resolvedUsername);
  user.set({ username: resolvedUsername });
}

export async function refreshTokens(): Promise<boolean> {
  const refreshToken = localStorage.getItem('refresh_token');
  if (!refreshToken) return false;

  const response = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken })
  });

  if (!response.ok) {
    logout();
    return false;
  }

  const data = await response.json();
  const username = localStorage.getItem('username') ?? '';
  persistTokens(data.access_token, data.refresh_token, username);
  return true;
}

/** Decode JWT payload and return true if the token is expired (or invalid). */
function isTokenExpired(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    // exp is in seconds; subtract 10 s buffer for clock skew
    return Date.now() / 1000 >= payload.exp - 10;
  } catch {
    return true;
  }
}

/**
 * Fetch wrapper that proactively refreshes an expired access token before
 * sending, and retries once on a 401 as a safety net.
 */
export async function apiFetch(input: string, init: RequestInit = {}): Promise<Response> {
  let token = localStorage.getItem('access_token');
  const headers = new Headers(init.headers);

  // Proactively refresh if the token is already expired
  if (token && isTokenExpired(token)) {
    const refreshed = await refreshTokens();
    token = refreshed ? localStorage.getItem('access_token') : null;
  }

  if (token) headers.set('Authorization', `Bearer ${token}`);

  let response = await fetch(`${API_BASE}${input}`, { ...init, headers });

  // Safety net: server says 401 even though we thought the token was valid
  if (response.status === 401) {
    const refreshed = await refreshTokens();
    if (refreshed) {
      headers.set('Authorization', `Bearer ${localStorage.getItem('access_token')}`);
      response = await fetch(`${API_BASE}${input}`, { ...init, headers });
    }
  }

  return response;
}

export function logout() {
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  localStorage.removeItem('username');
  user.set(null);
}

function persistTokens(accessToken: string, refreshToken: string, username: string) {
  localStorage.setItem('access_token', accessToken);
  localStorage.setItem('refresh_token', refreshToken);
  localStorage.setItem('username', username);
}

async function resolveUsername(accessToken: string, fallback: string): Promise<string> {
  try {
    const payload = JSON.parse(atob(accessToken.split('.')[1]));
    if (typeof payload?.username === 'string' && payload.username.trim() !== '') {
      return payload.username;
    }
  } catch {
    // Ignore token parse failure and fallback to API lookup.
  }

  try {
    const meRes = await fetch(`${API_BASE}/user/me`, {
      headers: { Authorization: `Bearer ${accessToken}` }
    });
    if (!meRes.ok) return fallback;
    const me = await meRes.json();
    if (typeof me?.username === 'string' && me.username.trim() !== '') {
      return me.username;
    }
  } catch {
    // Ignore lookup failure and fallback.
  }

  return fallback;
}

