import type { AccountInfo, CreditAccount, ModelInfo, ProviderInfo, Settings, UsageRecord, UsageSummary } from '@/types'

export class ApiError extends Error {
  status: number
  data: any
  constructor(status: number, message: string, data: any = {}) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }
}

type RequestOptions = Omit<RequestInit, 'body'> & { body?: BodyInit | Record<string, unknown> }

function storedKey(): string {
  return localStorage.getItem('trae_admin_key') || ''
}

async function request<T>(path: string, init: RequestOptions = {}, explicitKey?: string): Promise<T> {
  const headers = new Headers(init.headers)
  const key = explicitKey ?? storedKey()
  if (key) headers.set('Authorization', `Bearer ${key}`)
  let body = init.body
  if (body && typeof body !== 'string' && !(body instanceof FormData) && !(body instanceof Blob) && !(body instanceof ArrayBuffer)) {
    headers.set('Content-Type', 'application/json')
    body = JSON.stringify(body)
  }
  const response = await fetch(path, { ...init, body: body as BodyInit | undefined, headers, credentials: 'same-origin' })
  const text = await response.text()
  let data: any = {}
  try { data = text ? JSON.parse(text) : {} } catch { data = { message: text } }
  if (!response.ok) {
    const message = data?.error?.message || data?.message || `HTTP ${response.status}`
    throw new ApiError(response.status, message, data)
  }
  return data as T
}

export const api = {
  keyInfo: () => request<{ configured: boolean }>('/admin/api/key'),
  verifyKey: (key: string) => request<{ ok: boolean }>('/admin/api/key/check', { method: 'POST' }, key),
  health: () => request<{ available_accounts: number; total_accounts: number; available_providers: number; total_providers: number }>('/admin/health'),
  overview: () => request<{ accounts: AccountInfo[]; models: ModelInfo[]; model_count: number; credits: { remaining: number; total: number; used: number }; usage: UsageSummary }>('/admin/overview'),
  accounts: () => request<{ accounts: AccountInfo[] }>('/admin/accounts'),
  credits: () => request<{ fetched_at: string; accounts: CreditAccount[] }>('/admin/api/credits'),
  setEnabled: (uid: string, enabled: boolean) => request<{ ok: boolean }>(`/admin/accounts/${encodeURIComponent(uid)}/${enabled ? 'enable' : 'disable'}`, { method: 'POST' }),
  setAlias: (uid: string, alias: string) => request<{ ok: boolean }>(`/admin/accounts/${encodeURIComponent(uid)}/alias`, { method: 'POST', body: { alias } }),
  deleteAccount: (uid: string) => request<{ ok: boolean }>(`/admin/accounts/${encodeURIComponent(uid)}`, { method: 'DELETE' }),
  refreshAccount: (uid: string) => request<{ ok: boolean }>(`/admin/api/accounts/${encodeURIComponent(uid)}/refresh`, { method: 'POST' }),
  accountJSON: (uid: string) => request<Record<string, unknown>>(`/admin/api/accounts/${encodeURIComponent(uid)}/json`),
  refreshCredits: () => request<{ ok: boolean; accounts: CreditAccount[] }>('/admin/credits/refresh', { method: 'POST' }),
  checkin: () => request<{ ok: boolean; accounts: CreditAccount[] }>('/admin/checkin', { method: 'POST' }),
  startLogin: () => request<{ login_url: string; pending_id: string }>('/admin/api/login', { method: 'POST', body: {} }),
  pollLogin: (pendingId: string) => request<{ state: string; uid?: string; nickname?: string; error?: string }>(`/admin/api/login/result?pending_id=${encodeURIComponent(pendingId)}`),
  cancelLogin: (pendingId: string) => request<{ canceled: boolean }>('/admin/api/login/cancel', { method: 'POST', body: { pending_id: pendingId } }),
  importAccount: (value: string) => request<{ uid: string; nickname?: string; action: string }>('/admin/api/accounts/import', {
    method: 'POST',
    body: value.trim().startsWith('{') ? { json: value.trim() } : { callback_url: value.trim() },
  }),
  models: () => request<{ models: ModelInfo[]; source: string }>('/admin/models'),
  modelCatalog: () => request<{ models: ModelInfo[] }>('/admin/models/catalog'),
  refreshModels: () => request<{ ok: boolean; models: ModelInfo[]; source: string }>('/admin/models/refresh', { method: 'POST' }),
  providers: () => request<{ providers: ProviderInfo[] }>('/admin/providers'),
  provider: (name: string) => request<{ accounts: AccountInfo[]; models: ModelInfo[]; provider: ProviderInfo }>(`/admin/providers/${encodeURIComponent(name)}`),
  usageSummary: () => request<UsageSummary>('/admin/usage/summary'),
  usageTimeseries: () => request<{ data: { bucket: string; bucket_ts: number; count: number; tokens: number }[] }>('/admin/usage/timeseries'),
  usageRecent: () => request<{ records: UsageRecord[]; total: number }>('/admin/usage/recent'),
  settings: () => request<Settings>('/admin/settings'),
  saveSettings: (data: Partial<Settings>) => request<Settings & { ok: boolean }>('/admin/settings', { method: 'POST', body: data }),
}
