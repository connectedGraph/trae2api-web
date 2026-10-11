/** work2api WebUI — Tailwind 配置
 * 设计系统：青色系（teal）品牌 + 深蓝灰中性色，暗色优先、class 切换。
 * 语义色用 CSS 变量驱动（style.css 里按明/暗两套定义），这样组件只写
 * `text-brand` / `bg-surface` 就能同时适配两种模式，不用到处写 dark: 前缀。
 */
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        // 中性层（背景/面板/边框/文字）——全部走 CSS 变量
        bg: 'rgb(var(--c-bg) / <alpha-value>)',
        surface: 'rgb(var(--c-surface) / <alpha-value>)',
        elevated: 'rgb(var(--c-elevated) / <alpha-value>)',
        line: 'rgb(var(--c-line) / <alpha-value>)',
        ink: 'rgb(var(--c-ink) / <alpha-value>)', // 主文字
        muted: 'rgb(var(--c-muted) / <alpha-value>)', // 次要文字
        faint: 'rgb(var(--c-faint) / <alpha-value>)', // 最弱文字/占位
        // 品牌 + 语义信号
        brand: 'rgb(var(--c-brand) / <alpha-value>)',
        'brand-soft': 'rgb(var(--c-brand-soft) / <alpha-value>)',
        live: 'rgb(var(--c-live) / <alpha-value>)', // 健康/成功
        warn: 'rgb(var(--c-warn) / <alpha-value>)', // 额度预警
        fault: 'rgb(var(--c-fault) / <alpha-value>)', // 错误/冷却
        route: 'rgb(var(--c-route) / <alpha-value>)', // 链接/信息/选中
      },
      fontFamily: {
        sans: ['Inter Variable', 'Inter', 'system-ui', 'PingFang SC', 'Microsoft YaHei', 'sans-serif'],
        mono: ['JetBrains Mono', 'ui-monospace', 'SFMono-Regular', 'Menlo', 'monospace'],
      },
      fontSize: {
        micro: ['0.75rem', { lineHeight: '1rem' }],
        small: ['0.8125rem', { lineHeight: '1.25rem' }],
      },
      borderRadius: { xl: '0.875rem', '2xl': '1.125rem' },
      boxShadow: {
        glass: '0 1px 0 0 rgb(var(--c-hairline) / 0.6) inset, 0 12px 40px -24px rgb(0 0 0 / 0.6)',
        glow: '0 0 0 1px rgb(var(--c-brand) / 0.35), 0 12px 48px -20px rgb(var(--c-brand) / 0.45)',
      },
      keyframes: {
        'fade-up': { '0%': { opacity: '0', transform: 'translateY(6px)' }, '100%': { opacity: '1', transform: 'none' } },
        pulse: { '0%,100%': { opacity: '1' }, '50%': { opacity: '0.35' } },
      },
      animation: { 'fade-up': 'fade-up .28s ease both', pulse: 'pulse 2s ease-in-out infinite' },
    },
  },
  plugins: [],
}
