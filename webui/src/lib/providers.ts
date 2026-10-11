// 供应商元数据：统一的展示名、短码与配色，让"供应商"在全站作为一致的
// 筛选/标签维度出现（而非各自为政的独立页面）。
export interface ProviderMeta {
  name: string
  label: string
  short: string
  tone: 'brand' | 'route' | 'warn' | 'live'
  chip: string // 图标底色+文字色（静态 class，避免 Tailwind 动态拼接被 purge）
}

const REGISTRY: Record<string, ProviderMeta> = {
  workbuddy: { name: 'workbuddy', label: 'WorkBuddy', short: 'WB', tone: 'brand', chip: 'bg-brand/12 text-brand' },
  qoder: { name: 'qoder', label: 'Qoder', short: 'QD', tone: 'route', chip: 'bg-route/12 text-route' },
  opencode: { name: 'opencode', label: 'OpenCode', short: 'OC', tone: 'warn', chip: 'bg-warn/12 text-warn' },
}

export function providerMeta(name: string): ProviderMeta {
  return (
    REGISTRY[name] || {
      name,
      label: name ? name[0].toUpperCase() + name.slice(1) : '未知',
      short: (name || '?').slice(0, 2).toUpperCase(),
      tone: 'live',
      chip: 'bg-live/12 text-live',
    }
  )
}

// 站点/区域中文名
export function siteLabel(site: string | undefined): string {
  const m: Record<string, string> = {
    domestic: '国内',
    international: '国际',
    'qoder-cn': 'Qoder 国内',
    'qoder-global': 'Qoder 国际',
    cn: '国内',
    global: '国际',
  }
  return site ? m[site] || site : ''
}
