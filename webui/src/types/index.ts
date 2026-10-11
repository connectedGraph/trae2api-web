export interface AccountInfo {
  uid: string
  provider?: string
  nickname?: string
  enabled: boolean
  disabled: boolean
  cooling: boolean
  healthy: boolean
  reason?: string
  credits_remaining: number | null
  credits_total: number | null
  credits_used?: number | null
  credits_packs?: number
  expires_at?: number | null
  expired_soon?: boolean
  machine_id?: string
  device_id?: string
  failure_count?: number
}

export interface CreditAccount {
  uid: string
  nickname?: string
  remain: number
  limit: number
  used: number
  packs: number
  checked_in: boolean
  checkin_credits: number
  checkin_enable: boolean
  cooling: boolean
  disabled: boolean
  error?: string
}

export interface ModelInfo {
  id: string
  object?: string
  owned_by?: string
  name?: string
  context_length?: number
  max_output_tokens?: number
}

export interface ProviderInfo {
  name: string
  display_name: string
  ready: boolean
  default?: boolean
  capabilities: string[]
  status: Record<string, number | string | boolean>
}

export interface UsageSummary {
  total_requests: number
  total_tokens: number
  today_requests: number
  today_tokens: number
  by_model: { model: string; count: number; tokens: number }[]
  by_account: { account_uid: string; label?: string; count: number; tokens: number }[]
}

export interface UsageRecord {
  id: number
  ts: number
  model: string
  protocol: string
  account_uid: string
  input_tokens: number
  output_tokens: number
  total_tokens: number
  status: string
  latency_ms?: number
}

export interface Settings {
  listen: string
  callback_port: string
  auth_dir: string
  state_file: string
  default_model: string
  api_key_configured: boolean
  plan_credit: string
  soft_rate: string
  err_threshold: number
  err_cooldown: string
  upstream_timeout_seconds: number
}
