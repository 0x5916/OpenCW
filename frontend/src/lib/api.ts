import { apiFetch } from './auth';
import { extractErrorCodeFromBody } from './errorCode';

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
  throw new ApiError(extractErrorCodeFromBody(body) ?? fallback, res.status, body);
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

async function apiSendJson<T = void>(
  path: string,
  method: 'POST' | 'PUT',
  body: unknown,
  fallback: string
): Promise<T> {
  const res = await apiFetch(path, {
    method,
    headers: JSON_HEADERS,
    body: JSON.stringify(body)
  });

  if (!res.ok) return throwApiError(res, fallback);

  try {
    return await parseJsonBody<T>(res);
  } catch {
    throw new ApiError(fallback, res.status, null);
  }
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

export interface ForumCategory {
  id: string;
  name: string;
  description?: string | null;
  created_at?: string;
  updated_at?: string;
}

export interface ForumThread {
  id: string;
  category_id: string;
  author_id?: string;
  author?: string;
  username?: string;
  title: string;
  is_pinned: boolean;
  is_locked: boolean;
  created_at: string;
  updated_at?: string;
  latest_activity?: string;
  post_count?: number;
}

export interface ForumPost {
  id: string;
  thread_id: string;
  author_id?: string;
  author?: string;
  username?: string;
  body: string;
  parent_id?: string | null;
  created_at: string;
  updated_at?: string;
}

export interface ForumPage<T> {
  data: T[];
  page: number;
  limit: number;
  total?: number;
  total_pages?: number;
}

export interface ForumThreadCreated {
  thread: ForumThread;
  first_post: ForumPost;
}

function normalizeForumPage<T>(
  value: T[] | { data?: T[]; page?: number; limit?: number; total?: number; total_pages?: number },
  page: number,
  limit: number
): ForumPage<T> {
  if (Array.isArray(value)) return { data: value, page, limit };

  return {
    data: value.data ?? [],
    page: value.page ?? page,
    limit: value.limit ?? limit,
    total: value.total,
    total_pages: value.total_pages
  };
}

function forumQuery(page: number, limit: number): string {
  const safePage = Math.max(1, Math.floor(page));
  const safeLimit = Math.min(100, Math.max(1, Math.floor(limit)));
  return `?page=${safePage}&limit=${safeLimit}`;
}

export async function getForumCategories(): Promise<ForumCategory[]> {
  const response = await apiGetJson<ForumCategory[] | { data?: ForumCategory[] }>(
    '/forum/categories',
    'FORUM_CATEGORIES_FETCH_FAILED'
  );
  return Array.isArray(response) ? response : (response.data ?? []);
}

export async function getForumCategoryThreads(
  categoryId: string,
  page = 1,
  limit = 20
): Promise<ForumPage<ForumThread>> {
  const response = await apiGetJson<
    | ForumThread[]
    | { data?: ForumThread[]; page?: number; limit?: number; total?: number; total_pages?: number }
  >(
    `/forum/categories/${encodeURIComponent(categoryId)}/threads${forumQuery(page, limit)}`,
    'FORUM_THREADS_FETCH_FAILED'
  );
  return normalizeForumPage(response, page, limit);
}

export async function getForumThread(threadId: string): Promise<ForumThread> {
  const response = await apiGetJson<{ data?: ForumThread } | ForumThread>(
    `/forum/threads/${encodeURIComponent(threadId)}`,
    'FORUM_THREAD_NOT_FOUND'
  );
  if ('data' in response && response.data) return response.data;
  return response as ForumThread;
}

export async function getForumThreadPosts(
  threadId: string,
  page = 1,
  limit = 20
): Promise<ForumPage<ForumPost>> {
  const response = await apiGetJson<
    | ForumPost[]
    | { data?: ForumPost[]; page?: number; limit?: number; total?: number; total_pages?: number }
  >(
    `/forum/threads/${encodeURIComponent(threadId)}/posts${forumQuery(page, limit)}`,
    'FORUM_POSTS_FETCH_FAILED'
  );
  return normalizeForumPage(response, page, limit);
}

export async function createForumThread(
  categoryId: string,
  title: string,
  body: string
): Promise<ForumThreadCreated> {
  const response = await apiSendJson<{
    data?: ForumThreadCreated;
    thread?: ForumThread;
    first_post?: ForumPost;
  }>(
    '/forum/threads',
    'POST',
    { category_id: categoryId, title, body },
    'FORUM_THREAD_CREATE_FAILED'
  );
  if (response.data) return response.data;
  return { thread: response.thread as ForumThread, first_post: response.first_post as ForumPost };
}

export async function createForumPost(
  threadId: string,
  body: string,
  parentId?: string
): Promise<ForumPost> {
  const response = await apiSendJson<{ data?: ForumPost } | ForumPost>(
    `/forum/threads/${encodeURIComponent(threadId)}/posts`,
    'POST',
    { body, ...(parentId ? { parent_id: parentId } : {}) },
    'FORUM_POST_CREATE_FAILED'
  );
  if ('data' in response && response.data) return response.data;
  return response as ForumPost;
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

  await apiSendJson('/cw/progress', 'PUT', payload, 'PROGRESS_CREATE_FAILED');
}

export async function updatePassword(oldPassword: string, newPassword: string): Promise<void> {
  await apiSendJson(
    '/user/password',
    'PUT',
    { old_password: oldPassword, new_password: newPassword },
    'INTERNAL_SERVER_ERROR'
  );
}
