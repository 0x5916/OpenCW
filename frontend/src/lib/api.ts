import { apiFetch } from './auth';

class ApiError extends Error {
  status: number;

  body: unknown;

  code: string;

  constructor(code: string, status: number, body: unknown) {
    super(code);
    this.name = 'ApiError';
    this.status = status;
    this.body = body;
    this.code = code;
  }
}

const JSON_HEADERS = {
  'Content-Type': 'application/json',
  Accept: 'application/json'
} as const;
const inFlightGetRequests = new Map<string, Promise<unknown>>();

function isErrorCode(value: string): boolean {
  return /^[A-Z0-9_]+$/.test(value);
}

function pickErrorCode(body: unknown, fallback: string): string {
  if (!body || typeof body !== 'object') return fallback;

  const code = (body as Record<string, unknown>).code;
  if (typeof code === 'string' && code.trim() !== '') return code;

  const error = (body as Record<string, unknown>).error;
  if (typeof error === 'string' && error.trim() !== '' && isErrorCode(error.trim())) {
    return error;
  }

  const message = (body as Record<string, unknown>).message;
  if (typeof message === 'string' && message.trim() !== '' && isErrorCode(message.trim())) {
    return message;
  }

  const detail = (body as Record<string, unknown>).detail;
  if (typeof detail === 'string' && detail.trim() !== '' && isErrorCode(detail.trim())) {
    return detail;
  }

  return fallback;
}

async function parseJsonBody<T>(res: Response): Promise<T> {
  if (res.status === 204) {
    return undefined as T;
  }

  const contentLength = res.headers.get('content-length');
  if (contentLength === '0') {
    return undefined as T;
  }

  const text = await res.text();
  if (text.trim() === '') {
    return undefined as T;
  }

  return JSON.parse(text) as T;
}

async function throwApiError(res: Response, fallback: string): Promise<never> {
  const body = await parseJsonBody<unknown>(res).catch(() => null);
  throw new ApiError(pickErrorCode(body, fallback), res.status, body);
}

async function apiGetJson<T>(path: string, fallback: string): Promise<T> {
  const requestKey = `GET:${path}`;
  const existing = inFlightGetRequests.get(requestKey) as Promise<T> | undefined;
  if (existing) return existing;

  const request = (async () => {
    const res = await apiFetch(path);
    if (!res.ok) return throwApiError(res, fallback);

    try {
      return await parseJsonBody<T>(res);
    } catch {
      throw new ApiError(fallback, res.status, null);
    }
  })();

  inFlightGetRequests.set(requestKey, request);

  try {
    return await request;
  } finally {
    inFlightGetRequests.delete(requestKey);
  }
}

async function apiSendJson(path: string, method: 'POST' | 'PUT', body: unknown, fallback: string) {
  const res = await apiFetch(path, {
    method,
    headers: JSON_HEADERS,
    body: JSON.stringify(body)
  });

  if (!res.ok) return throwApiError(res, fallback);
}

async function apiPost(path: string, fallback: string): Promise<void> {
  const res = await apiFetch(path, {
    method: 'POST',
    headers: { Accept: 'application/json' }
  });

  if (!res.ok) return throwApiError(res, fallback);
}

export interface CWSettings {
  char_wpm: number;
  eff_wpm: number;
  freq: number;
  start_delay: number;
  updated_at?: string;
}

export interface PageSettings {
  theme: 'auto' | 'dark' | 'light';
  language: string;
  cur_lesson: number;
  updated_at?: string;
}

export interface UserInfo {
  call_sign: string | null;
  username: string;
  email: string;
  email_verified: boolean;
  created_at: string;
}

export interface CombinedSettings {
  cw_settings: CWSettings;
  page_settings: PageSettings;
}

export async function getCWSettings(): Promise<CWSettings> {
  return apiGetJson<CWSettings>('/settings/cw', 'SETTINGS_FETCH_FAILED');
}

export async function saveCWSettings(settings: CWSettings): Promise<void> {
  await apiSendJson('/settings/cw', 'POST', settings, 'SETTINGS_UPDATE_FAILED');
}

export async function getSettings(): Promise<CombinedSettings> {
  return apiGetJson<CombinedSettings>('/settings/all', 'SETTINGS_FETCH_FAILED');
}

export async function getUserInfo(): Promise<UserInfo> {
  return apiGetJson<UserInfo>('/user/me', 'INTERNAL_SERVER_ERROR');
}

export async function savePageSettings(settings: PageSettings): Promise<void> {
  await apiSendJson('/settings/page', 'POST', settings, 'SETTINGS_UPDATE_FAILED');
}

export async function getPageSettings(): Promise<PageSettings> {
  return apiGetJson<PageSettings>('/settings/page', 'SETTINGS_FETCH_FAILED');
}

export async function updateEmail(email: string): Promise<void> {
  await apiSendJson('/user/email', 'PUT', { email }, 'INTERNAL_SERVER_ERROR');
}

export async function sendVerificationEmail(): Promise<void> {
  await apiPost('/auth/send-verification-email', 'VERIFICATION_SEND_FAILED');
}

export async function verifyEmail(code: string): Promise<void> {
  await apiSendJson('/auth/verify-email', 'POST', { code }, 'VERIFICATION_CODE_INVALID');
}

export async function updateCallSign(callSign: string): Promise<void> {
  await apiSendJson('/user/callsign', 'PUT', { call_sign: callSign }, 'INTERNAL_SERVER_ERROR');
}

export interface ProgressRecord {
  lesson: string;
  char_wpm: number;
  eff_wpm: number;
  accuracy: number;
  created_at: string;
  client_created_at?: string;
}

export async function getProgress(): Promise<ProgressRecord[]> {
  const data = await apiGetJson<{
    data?: Array<Omit<ProgressRecord, 'lesson'> & { lesson: string | number }>;
  }>('/cw/progress', 'PROGRESS_QUERY_FAILED');

  return (data.data ?? []).map((record) => ({
    ...record,
    lesson: String(record.lesson)
  }));
}

export async function submitProgress(
  lesson: number,
  charWpm: number,
  effWpm: number,
  accuracy: number,
  clientCreatedAt?: string
): Promise<void> {
  const payload: {
    lesson: number;
    char_wpm: number;
    eff_wpm: number;
    accuracy: number;
    client_created_at?: string;
  } = {
    lesson,
    char_wpm: charWpm,
    eff_wpm: effWpm,
    accuracy
  };

  if (clientCreatedAt) {
    payload.client_created_at = clientCreatedAt;
  }

  await apiSendJson(
    '/cw/progress',
    'PUT',
    payload,
    'PROGRESS_CREATE_FAILED'
  );
}

export async function updatePassword(oldPassword: string, newPassword: string): Promise<void> {
  await apiSendJson(
    '/user/password',
    'PUT',
    { old_password: oldPassword, new_password: newPassword },
    'INTERNAL_SERVER_ERROR'
  );
}
