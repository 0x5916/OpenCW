import { apiFetch } from './auth';

async function apiError(res: Response, fallback: string): Promise<never> {
  const body = await res.json().catch(() => ({}));
  throw new Error(body?.error ?? fallback);
}

export interface CWSettings {
  char_wpm: number;
  eff_wpm: number;
  freq: number;
  start_delay: number;
}

export interface UserInfo {
  username: string;
  email: string;
  created_at: string;
}

export async function getCWSettings(): Promise<CWSettings> {
  const res = await apiFetch('/cw/settings');
  if (!res.ok) return apiError(res, 'Failed to load settings');
  return res.json();
}

export async function saveCWSettings(settings: CWSettings): Promise<void> {
  const res = await apiFetch('/cw/settings', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(settings)
  });
  if (!res.ok) return apiError(res, 'Failed to save settings');
}

export async function getUserInfo(): Promise<UserInfo> {
  const res = await apiFetch('/user/me');
  if (!res.ok) return apiError(res, 'Failed to load user info');
  return res.json();
}

export async function updateEmail(email: string): Promise<void> {
  const res = await apiFetch('/user/email', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email })
  });
  if (!res.ok) return apiError(res, 'Failed to update email');
}

export interface ProgressRecord {
  lesson: string;
  char_wpm: number;
  eff_wpm: number;
  accuracy: number;
  created_at: string;
}

export async function getProgress(): Promise<ProgressRecord[]> {
  const res = await apiFetch('/cw/progress');
  if (!res.ok) return apiError(res, 'Failed to load progress');
  const data = await res.json();
  return data.data ?? [];
}

export async function submitProgress(
  lesson: string,
  charWpm: number,
  effWpm: number,
  accuracy: number
): Promise<void> {
  const res = await apiFetch('/cw/progress', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ lesson, char_wpm: charWpm, eff_wpm: effWpm, accuracy })
  });
  if (!res.ok) return apiError(res, 'Failed to submit progress');
}
