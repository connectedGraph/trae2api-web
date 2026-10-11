// 通用格式化工具：数值缩写、额度、时间、相对时间。
import dayjs from 'dayjs'

// 大数缩写：1234 → 1.23K，1200000 → 1.2M
export function abbr(n: number | undefined | null, digits = 1): string {
  if (n == null || isNaN(n)) return '—'
  const abs = Math.abs(n)
  if (abs < 1000) return String(Math.round(n))
  const units = [
    { v: 1e9, s: 'B' },
    { v: 1e6, s: 'M' },
    { v: 1e3, s: 'K' },
  ]
  for (const u of units) {
    if (abs >= u.v) return (n / u.v).toFixed(digits).replace(/\.0+$/, '') + u.s
  }
  return String(n)
}

// 千分位整数
export function int(n: number | undefined | null): string {
  if (n == null || isNaN(n)) return '—'
  return Math.round(n).toLocaleString('en-US')
}

// 额度：保留两位，去掉多余 0
export function credits(n: number | undefined | null): string {
  if (n == null || isNaN(n)) return '—'
  return n.toLocaleString('en-US', { maximumFractionDigits: 2 })
}

// unix 秒 → 日期时间
export function dt(sec: number | undefined | null, fmt = 'MM-DD HH:mm'): string {
  if (!sec) return '—'
  return dayjs(sec * 1000).format(fmt)
}

// unix 秒 → 相对现在（如 "3 分钟前" / "2 天后"）
export function rel(sec: number | undefined | null): string {
  if (!sec) return '—'
  const diff = sec * 1000 - Date.now()
  const past = diff < 0
  const a = Math.abs(diff)
  const m = Math.round(a / 60000)
  if (m < 1) return past ? '刚刚' : '即将'
  if (m < 60) return past ? `${m} 分钟前` : `${m} 分钟后`
  const h = Math.round(m / 60)
  if (h < 24) return past ? `${h} 小时前` : `${h} 小时后`
  const d = Math.round(h / 24)
  return past ? `${d} 天前` : `${d} 天后`
}

// 冷却剩余：unix 秒 → "还剩 42s" / 空
export function cooldown(until: number | undefined | null): string {
  if (!until) return ''
  const s = Math.round(until - Date.now() / 1000)
  if (s <= 0) return ''
  if (s < 60) return `${s}s`
  return `${Math.round(s / 60)}m`
}

// 延迟毫秒 → "1.2s" / "340ms"
export function latency(ms: number | undefined | null): string {
  if (ms == null) return '—'
  return ms >= 1000 ? (ms / 1000).toFixed(2) + 's' : ms + 'ms'
}

// 缓存命中率 0~1 → "87%"；null/undefined → null（调用方渲染占位）
export function pct(n: number | undefined | null): string | null {
  if (n == null || !Number.isFinite(n) || n < 0 || n > 1) return null
  return Math.round(n * 100) + '%'
}
