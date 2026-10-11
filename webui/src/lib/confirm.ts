// 全局确认对话框（替代 antd popconfirm）。用法：
//   if (await confirm({ title: '删除账号？', tone: 'danger', okText: '删除' })) { ... }
import { reactive } from 'vue'

export interface ConfirmOptions {
  title: string
  body?: string
  okText?: string
  cancelText?: string
  tone?: 'brand' | 'danger'
}

interface ConfirmState extends ConfirmOptions {
  open: boolean
  _resolve?: (ok: boolean) => void
}

export const confirmState = reactive<ConfirmState>({ open: false, title: '' })

export function confirm(opts: ConfirmOptions): Promise<boolean> {
  return new Promise((resolve) => {
    Object.assign(confirmState, opts, { open: true, _resolve: resolve })
  })
}

export function resolveConfirm(ok: boolean) {
  confirmState.open = false
  confirmState._resolve?.(ok)
  confirmState._resolve = undefined
}
