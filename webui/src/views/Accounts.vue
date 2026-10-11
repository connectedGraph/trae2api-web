<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { api } from '@/api/client'
import type { AccountInfo, CreditAccount } from '@/types'
import { confirm } from '@/lib/confirm'
import { toast } from '@/lib/toast'
import { credits as fmtCredits, dt } from '@/lib/format'
import WPage from '@/components/ui/WPage.vue'
import WCard from '@/components/ui/WCard.vue'
import WButton from '@/components/ui/WButton.vue'
import WInput from '@/components/ui/WInput.vue'
import WTextarea from '@/components/ui/WTextarea.vue'
import WToggle from '@/components/ui/WToggle.vue'
import WTable, { type Column } from '@/components/ui/WTable.vue'
import WTag from '@/components/ui/WTag.vue'
import WModal from '@/components/ui/WModal.vue'
import WSpinner from '@/components/ui/WSpinner.vue'
import WIcon from '@/components/ui/WIcon.vue'

const accounts = ref<AccountInfo[]>([])
const credits = ref<Record<string, CreditAccount>>({})
const loading = ref(true)
const busy = ref(false)
const lastFetched = ref('')
const importText = ref('')
const importBusy = ref(false)
const loginBusy = ref(false)
const loginStatus = ref('')
const loginError = ref(false)
const pendingId = ref('')
let pollTimer: ReturnType<typeof setTimeout> | null = null

const jsonOpen = ref(false)
const jsonTitle = ref('凭证预览')
const jsonBody = ref('')
const renameOpen = ref(false)
const renameUID = ref('')
const renameValue = ref('')

const rows = computed(() => accounts.value)
const cols: Column[] = [
  { key: 'nickname', label: '账号' },
  { key: 'uid', label: 'UID', mono: true },
  { key: 'healthy', label: '状态' },
  { key: 'credits_remaining', label: '剩余额度', align: 'right', mono: true },
  { key: 'expires_at', label: 'Token 有效期', mono: true },
  { key: 'actions', label: '', align: 'right' },
]

function credit(uid: string): CreditAccount | undefined { return credits.value[uid] }
function statusTone(a: AccountInfo): 'live' | 'warn' | 'fault' | 'muted' {
  if (a.healthy) return 'live'
  if (a.disabled) return 'fault'
  if (a.cooling) return 'warn'
  return 'muted'
}
function statusLabel(a: AccountInfo): string {
  if (a.healthy) return '正常'
  if (a.disabled) return '会话失效'
  if (a.cooling) return '冷却中'
  if (!a.enabled) return '已停用'
  return '不可用'
}

async function load() {
  loading.value = true
  try {
    const [accountData, creditData] = await Promise.all([api.accounts(), api.credits()])
    accounts.value = accountData.accounts || []
    credits.value = Object.fromEntries((creditData.accounts || []).map((item) => [item.uid, item]))
    lastFetched.value = creditData.fetched_at || new Date().toLocaleTimeString()
  } catch (error) {
    toast.error(error instanceof Error ? error.message : '账号加载失败')
  } finally {
    loading.value = false
  }
}
async function refreshCredits() {
  busy.value = true
  try {
    const data = await api.refreshCredits()
    credits.value = Object.fromEntries((data.accounts || []).map((item) => [item.uid, item]))
    await load()
    toast.success('额度已刷新')
  } finally { busy.value = false }
}
async function checkin() {
  busy.value = true
  try { await api.checkin(); await load(); toast.success('签到请求已完成') }
  finally { busy.value = false }
}
async function toggleEnabled(a: AccountInfo, enabled: boolean) {
  try { await api.setEnabled(a.uid, enabled); await load(); toast.success(enabled ? '账号已启用' : '账号已停用') }
  catch (error) { toast.error(error instanceof Error ? error.message : '更新失败') }
}
function openRename(a: AccountInfo) {
  renameUID.value = a.uid
  renameValue.value = a.nickname || ''
  renameOpen.value = true
}
async function saveRename() {
  if (!renameUID.value) return
  busy.value = true
  try { await api.setAlias(renameUID.value, renameValue.value.trim()); renameOpen.value = false; await load(); toast.success('账号名称已保存') }
  finally { busy.value = false }
}
async function refreshAccount(a: AccountInfo) {
  busy.value = true
  try { await api.refreshAccount(a.uid); await load(); toast.success('Token 已刷新') }
  catch (error) { toast.error(error instanceof Error ? error.message : 'Token 刷新失败') }
  finally { busy.value = false }
}
async function removeAccount(a: AccountInfo) {
  const ok = await confirm({ title: `删除账号「${a.nickname || a.uid}」？`, body: '删除后凭证文件也会移除，且无法从本地面板恢复。', tone: 'danger', okText: '删除' })
  if (!ok) return
  busy.value = true
  try { await api.deleteAccount(a.uid); await load(); toast.success('账号已删除') }
  finally { busy.value = false }
}
async function showJSON(a: AccountInfo) {
  try {
    const data = await api.accountJSON(a.uid)
    jsonTitle.value = `凭证预览（脱敏） · ${a.nickname || a.uid}`
    jsonBody.value = JSON.stringify(data, null, 2)
    jsonOpen.value = true
  } catch (error) { toast.error(error instanceof Error ? error.message : '凭证读取失败') }
}
async function importAccount() {
  const value = importText.value.trim()
  if (!value || importBusy.value) return
  importBusy.value = true
  try { const data = await api.importAccount(value); importText.value = ''; await load(); toast.success(`导入成功：${data.nickname || data.uid}`) }
  catch (error) { toast.error(error instanceof Error ? error.message : '导入失败') }
  finally { importBusy.value = false }
}
async function startLogin() {
  loginBusy.value = true
  loginError.value = false
  loginStatus.value = '正在生成 TRAE 登录链接…'
  try {
    const data = await api.startLogin()
    pendingId.value = data.pending_id
    loginStatus.value = '已打开登录窗口，完成登录后此页会自动刷新。'
    window.open(data.login_url, '_blank', 'noopener,noreferrer')
    pollLogin()
  } catch (error) {
    loginError.value = true
    loginStatus.value = error instanceof Error ? error.message : '登录链接生成失败'
  } finally { loginBusy.value = false }
}
function stopPoll() {
  if (pollTimer) clearTimeout(pollTimer)
  pollTimer = null
}
function pollLogin() {
  stopPoll()
  const run = async () => {
    if (!pendingId.value) return
    try {
      const data = await api.pollLogin(pendingId.value)
      if (data.state === 'success') {
        loginError.value = false
        loginStatus.value = `登录成功：${data.nickname || ''}（${data.uid || ''}）`
        pendingId.value = ''
        await load()
        return
      }
      if (data.state === 'failed' || data.state === 'canceled') {
        loginError.value = true
        loginStatus.value = data.error || '登录未完成'
        pendingId.value = ''
        return
      }
      pollTimer = setTimeout(run, 2000)
    } catch (error) {
      loginError.value = true
      loginStatus.value = error instanceof Error ? error.message : '登录状态查询失败'
      pendingId.value = ''
    }
  }
  void run()
}
async function cancelLogin() {
  if (!pendingId.value) return
  await api.cancelLogin(pendingId.value).catch(() => {})
  pendingId.value = ''
  stopPoll()
  loginStatus.value = '已取消登录'
}
onMounted(load)
onUnmounted(stopPoll)
</script>

<template>
  <WPage title="账号" sub="TRAE 登录、凭证、额度和账号池状态。">
    <template #actions>
      <WButton :loading="busy" @click="checkin"><WIcon name="check" :size="15" />签到</WButton>
      <WButton :loading="busy" @click="refreshCredits"><WIcon name="refresh" :size="15" />刷新额度</WButton>
      <WButton :loading="loading" @click="load">刷新</WButton>
    </template>

    <div class="mb-4 grid gap-4 lg:grid-cols-2">
      <WCard title="TRAE 登录" sub="打开官方登录页，回调会自动写入当前账号池。">
        <div class="flex flex-wrap items-center gap-3"><WButton variant="primary" :loading="loginBusy" :disabled="!!pendingId" @click="startLogin">添加账号</WButton><WButton v-if="pendingId" variant="subtle" @click="cancelLogin">取消登录</WButton><span v-if="loginStatus" class="text-small" :class="loginError ? 'text-fault' : 'text-muted'">{{ loginStatus }}</span></div>
      </WCard>
      <WCard title="手动导入凭证" sub="支持回调链接或 JSON；令牌只在服务端解析和保存。">
        <WTextarea v-model="importText" :rows="3" placeholder="粘贴 http://127.0.0.1:18080/authorize?... 或 JSON 凭证" />
        <div class="mt-3 flex items-center justify-between gap-3"><span class="text-micro text-faint">写操作使用当前管理员 API Key</span><WButton variant="primary" :loading="importBusy" :disabled="!importText.trim()" @click="importAccount">导入凭证</WButton></div>
      </WCard>
    </div>

    <WCard title="账号池" :sub="lastFetched ? `额度更新于 ${lastFetched}` : '正在读取账号状态'" flush>
      <WSpinner v-if="loading && !accounts.length" center label="加载中" />
      <WTable v-else :columns="cols" :rows="rows" row-key="uid" min-width="980px">
        <template #cell-nickname="{ row }"><div class="font-medium text-ink">{{ row.nickname || '未命名账号' }}</div><div v-if="row.reason" class="mt-0.5 max-w-52 truncate text-micro text-faint" :title="row.reason">{{ row.reason }}</div></template>
        <template #cell-healthy="{ row }"><WTag :tone="statusTone(row)" dot>{{ statusLabel(row) }}</WTag></template>
        <template #cell-credits_remaining="{ row }"><div class="mono text-ink">{{ fmtCredits(credit(row.uid)?.remain ?? row.credits_remaining) }}</div><div v-if="credit(row.uid)?.limit" class="text-micro text-faint">/{{ fmtCredits(credit(row.uid)?.limit) }}</div></template>
        <template #cell-expires_at="{ row }"><span class="mono text-muted">{{ row.expires_at ? dt(row.expires_at, 'YYYY-MM-DD HH:mm') : '—' }}</span></template>
        <template #cell-actions="{ row }"><div class="flex flex-wrap justify-end gap-1.5"><WToggle :model-value="Boolean(row.enabled && !row.disabled)" :disabled="row.disabled || busy" @update:model-value="toggleEnabled(row, $event)" /><WButton size="sm" :disabled="busy" @click="openRename(row)">改名</WButton><WButton size="sm" :disabled="busy" @click="refreshAccount(row)">刷新 Token</WButton><WButton size="sm" @click="showJSON(row)">JSON</WButton><WButton size="sm" variant="danger" :disabled="busy" @click="removeAccount(row)">删除</WButton></div></template>
        <template #empty><div class="px-4 py-12 text-center text-small text-faint">暂无账号，请使用上方 TRAE 登录或手动导入。</div></template>
      </WTable>
    </WCard>

    <WModal v-model:open="jsonOpen" :title="jsonTitle" size="md"><pre class="max-h-[55vh] overflow-auto rounded-lg border border-line bg-bg p-3 text-micro leading-relaxed text-muted">{{ jsonBody }}</pre></WModal>
    <WModal v-model:open="renameOpen" title="设置账号名称" size="sm">
      <label class="mb-1.5 block text-small text-muted">显示名称</label><WInput v-model="renameValue" placeholder="例如：主账号" />
      <template #footer><WButton variant="subtle" @click="renameOpen = false">取消</WButton><WButton variant="primary" :loading="busy" @click="saveRename">保存</WButton></template>
    </WModal>
  </WPage>
</template>
