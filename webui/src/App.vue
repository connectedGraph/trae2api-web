<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api, ApiError } from '@/api/client'
import { toast } from '@/lib/toast'
import WButton from '@/components/ui/WButton.vue'
import WIcon from '@/components/ui/WIcon.vue'
import WLed from '@/components/ui/WLed.vue'
import WToaster from '@/components/ui/WToaster.vue'
import WInput from '@/components/ui/WInput.vue'
import WConfirm from '@/components/ui/WConfirm.vue'

const route = useRoute()
const router = useRouter()
const nav = [
  { key: '/overview', label: '概览', icon: 'overview' },
  { key: '/accounts', label: '账号', icon: 'accounts' },
  { key: '/models', label: '模型', icon: 'models' },
  { key: '/keys', label: 'API 密钥', icon: 'keys' },
  { key: '/traffic', label: '流量', icon: 'traffic' },
  { key: '/settings', label: '设置', icon: 'settings' },
  { key: '/portal', label: '门户管理', icon: 'accounts' },
]
const activeKey = computed(() => '/' + (route.path.split('/')[1] || 'overview'))
const mobileOpen = ref(false)
function go(key: string) {
  if (route.path !== key) router.push(key)
  mobileOpen.value = false
}

const isLight = ref(document.documentElement.classList.contains('light'))
function toggleTheme() {
  isLight.value = !isLight.value
  document.documentElement.classList.toggle('light', isLight.value)
  localStorage.setItem('trae_theme', isLight.value ? 'light' : 'dark')
}

const checking = ref(true)
const authOk = ref(false)
const keyRequired = ref(false)
const keyInput = ref('')
const authError = ref('')
const loginBusy = ref(false)
async function checkAccess() {
  checking.value = true
  try {
    const info = await api.keyInfo()
    keyRequired.value = info.configured
    if (!info.configured) {
      authOk.value = true
      return
    }
    const saved = localStorage.getItem('trae_admin_key') || ''
    if (saved) {
      await api.verifyKey(saved)
      authOk.value = true
    }
  } catch {
    authOk.value = false
  } finally {
    checking.value = false
  }
}
async function login() {
  if (!keyInput.value.trim() || loginBusy.value) return
  loginBusy.value = true
  authError.value = ''
  try {
    await api.verifyKey(keyInput.value.trim())
    localStorage.setItem('trae_admin_key', keyInput.value.trim())
    keyInput.value = ''
    authOk.value = true
    toast.success('登录成功')
  } catch (error) {
    authError.value = error instanceof ApiError ? error.message : '登录失败，请检查 API Key'
  } finally {
    loginBusy.value = false
  }
}
function logout() {
  localStorage.removeItem('trae_admin_key')
  authOk.value = !keyRequired.value
  toast.success(keyRequired.value ? '已退出登录' : '本机免鉴权模式')
}

const healthy = ref<number | null>(null)
const total = ref<number | null>(null)
let timer: ReturnType<typeof setInterval> | null = null
async function refreshHealth() {
  if (!authOk.value) return
  try {
    const health = await api.health()
    healthy.value = health.available_accounts
    total.value = health.total_accounts
  } catch {
    healthy.value = null
    total.value = null
  }
}
const serving = computed(() => (healthy.value ?? 0) > 0)
const healthText = computed(() => healthy.value == null ? '状态未知' : serving.value ? '有可用账号' : '无可用账号')

onMounted(() => {
  void checkAccess().then(() => refreshHealth())
  timer = setInterval(refreshHealth, 30000)
})
onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<template>
  <div v-if="checking" class="flex min-h-screen items-center justify-center bg-bg">
    <span class="h-6 w-6 animate-spin rounded-full border-2 border-brand border-t-transparent" />
  </div>

  <div v-else-if="!authOk" class="flex min-h-screen items-center justify-center bg-bg p-4">
    <div class="glass w-full max-w-sm rounded-2xl p-7 shadow-glass">
      <div class="mb-6 flex items-center gap-2.5"><WLed tone="live" pulse /><span class="mono text-lg font-semibold tracking-tight text-ink">trae2api</span></div>
      <h1 class="mb-1 text-base font-semibold text-ink">管理员登录</h1>
      <p class="mb-5 text-micro text-faint">输入服务配置中的 TW2A_API_KEY，以访问账号管理和网关设置。</p>
      <WInput v-model="keyInput" type="password" placeholder="输入 TW2A_API_KEY" class="mb-3" @enter="login" />
      <p v-if="authError" class="mb-2 text-micro text-fault">{{ authError }}</p>
      <WButton variant="primary" block :loading="loginBusy" @click="login"><WIcon name="keys" :size="16" /> 登录</WButton>
    </div>
  </div>

  <div v-else class="flex h-screen overflow-hidden bg-bg text-ink">
    <div v-if="mobileOpen" class="fixed inset-0 z-30 bg-black/50 md:hidden" @click="mobileOpen = false" />
    <aside class="fixed inset-y-0 left-0 z-40 flex w-60 flex-col border-r border-line bg-surface/80 backdrop-blur-xl transition-transform md:static md:translate-x-0" :class="mobileOpen ? 'translate-x-0' : '-translate-x-full'">
      <div class="flex items-center gap-2.5 border-b border-line px-5 py-4">
        <WLed :tone="serving ? 'live' : healthy == null ? 'muted' : 'fault'" :pulse="serving" />
        <span class="mono text-[0.95rem] font-semibold tracking-tight text-ink">trae2api</span>
        <span class="ml-auto text-micro text-faint">网关</span>
      </div>
      <nav class="flex-1 space-y-1 overflow-y-auto p-3">
        <button v-for="item in nav" :key="item.key" class="group flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-small font-medium transition-colors" :class="activeKey === item.key ? 'bg-brand/12 text-brand' : 'text-muted hover:bg-elevated hover:text-ink'" @click="go(item.key)">
          <WIcon :name="item.icon" :size="18" /><span>{{ item.label }}</span><span v-if="activeKey === item.key" class="ml-auto h-1.5 w-1.5 rounded-full bg-brand" />
        </button>
      </nav>
      <div class="space-y-3 border-t border-line px-4 py-4">
        <div class="flex items-center gap-2 text-micro text-muted"><WLed :tone="serving ? 'live' : healthy == null ? 'muted' : 'fault'" /><span>{{ healthText }}</span><span v-if="healthy != null" class="mono ml-auto text-ink">{{ healthy }}<span class="text-faint">/{{ total }}</span></span></div>
        <div class="flex items-center gap-2">
          <button class="flex h-8 flex-1 items-center justify-center gap-1.5 rounded-lg border border-line text-micro text-muted transition-colors hover:border-brand hover:text-brand" @click="toggleTheme"><WIcon :name="isLight ? 'moon' : 'sun'" :size="15" /> {{ isLight ? '暗色' : '浅色' }}</button>
          <button v-if="keyRequired" class="flex h-8 flex-1 items-center justify-center gap-1.5 rounded-lg border border-line text-micro text-muted transition-colors hover:border-fault hover:text-fault" @click="logout"><WIcon name="keys" :size="15" /> 退出登录</button>
        </div>
      </div>
    </aside>
    <div class="flex min-w-0 flex-1 flex-col">
      <header class="flex items-center gap-3 border-b border-line px-4 py-3 md:hidden"><button class="text-muted" @click="mobileOpen = true"><WIcon name="menu" :size="22" /></button><span class="mono font-semibold">trae2api</span></header>
      <main class="flex-1 overflow-auto"><router-view /></main>
    </div>
  </div>
  <WToaster />
  <WConfirm />
</template>
