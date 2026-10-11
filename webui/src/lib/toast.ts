// 轻量全局提示（替代 antd message）。框架无关：一个响应式队列，
// 由 WToaster 渲染；client.ts 等非组件模块也能直接调用。
import { reactive } from 'vue'

export type ToastTone = 'info' | 'success' | 'error'
export interface Toast {
  id: number
  tone: ToastTone
  text: string
}

const state = reactive<{ items: Toast[] }>({ items: [] })
let seq = 0

// 去重：同文案在窗口期内只弹一次（概览页轮询失败时避免刷屏）
const DEDUPE_MS = 4000
let lastText = ''
let lastAt = 0

function push(tone: ToastTone, text: string, ttl = 3600) {
  const now = Date.now()
  if (text === lastText && now - lastAt < DEDUPE_MS) return
  lastText = text
  lastAt = now
  const id = ++seq
  state.items.push({ id, tone, text })
  setTimeout(() => dismiss(id), ttl)
}

export function dismiss(id: number) {
  const i = state.items.findIndex((t) => t.id === id)
  if (i >= 0) state.items.splice(i, 1)
}

export const toasts = state
export const toast = {
  info: (t: string) => push('info', t),
  success: (t: string) => push('success', t, 2600),
  error: (t: string) => push('error', t, 4600),
}
