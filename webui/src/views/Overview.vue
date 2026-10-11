<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { api } from '@/api/client'
import type { AccountInfo, ModelInfo, UsageSummary } from '@/types'
import { abbr, credits as fmtCredits, int } from '@/lib/format'
import { toast } from '@/lib/toast'
import WPage from '@/components/ui/WPage.vue'
import WCard from '@/components/ui/WCard.vue'
import WStat from '@/components/ui/WStat.vue'
import WLed from '@/components/ui/WLed.vue'
import WTable, { type Column } from '@/components/ui/WTable.vue'
import WButton from '@/components/ui/WButton.vue'
import WSpinner from '@/components/ui/WSpinner.vue'
import WTag from '@/components/ui/WTag.vue'

const accounts = ref<AccountInfo[]>([])
const models = ref<ModelInfo[]>([])
const usage = ref<UsageSummary | null>(null)
const credits = ref({ remaining: 0, total: 0, used: 0 })
const loading = ref(true)
const failed = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

const healthy = computed(() => accounts.value.filter((a) => a.healthy).length)
const serving = computed(() => healthy.value > 0)
const remaining = computed(() => credits.value.remaining)
const recentAccounts = computed(() => accounts.value.slice(0, 8) as unknown as Record<string, any>[])
const cols: Column[] = [
  { key: 'nickname', label: '账号' },
  { key: 'uid', label: 'UID', mono: true },
  { key: 'credits_remaining', label: '剩余额度', align: 'right', mono: true },
  { key: 'healthy', label: '状态' },
]

async function load() {
  try {
    const [overview, modelData] = await Promise.all([api.overview(), api.models()])
    accounts.value = overview.accounts || []
    models.value = modelData.models || overview.models || []
    usage.value = overview.usage || null
    credits.value = overview.credits || { remaining: 0, total: 0, used: 0 }
    failed.value = false
  } catch (error) {
    failed.value = true
    if (!accounts.value.length) toast.error(error instanceof Error ? error.message : '概览加载失败')
  } finally {
    loading.value = false
  }
}
async function refresh() {
  loading.value = true
  await load()
}
onMounted(() => { void load(); timer = setInterval(() => void load(), 30000) })
onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<template>
  <WPage title="概览" sub="TRAE SOLO 网关运行状态与额度续航一览。">
    <template #actions><WButton :loading="loading" @click="refresh"><span>刷新</span></WButton></template>
    <div v-if="failed" class="mb-4 rounded-lg border border-warn/40 bg-warn/5 px-3.5 py-2.5 text-small text-warn">部分数据刷新失败，当前显示上次获取的结果。</div>
    <WSpinner v-if="loading && !accounts.length" center label="加载中" />
    <template v-else>
      <div class="mb-4 grid gap-4 lg:grid-cols-[1.4fr_1fr]">
        <WCard hover>
          <div class="flex items-center gap-2.5"><WLed :tone="serving ? 'live' : 'fault'" :pulse="serving" /><span class="text-small font-medium" :class="serving ? 'text-live' : 'text-fault'">{{ serving ? '网关正在服务' : '暂无健康账号，网关等待恢复' }}</span></div>
          <div class="mt-4"><div class="text-micro uppercase tracking-wide text-faint">当前可用剩余额度</div><div class="mono mt-1 text-4xl font-semibold leading-none text-brand">{{ fmtCredits(remaining) }}</div><div class="mt-2 flex flex-wrap gap-x-6 gap-y-1 text-small text-muted"><span>总额度 <span class="mono text-ink">{{ fmtCredits(credits.total) }}</span></span><span>已消耗 <span class="mono text-ink">{{ fmtCredits(credits.used) }}</span></span></div></div>
        </WCard>
        <div class="grid grid-cols-2 gap-4">
          <WStat label="健康账号" :value="`${healthy}/${accounts.length}`" :tone="serving ? 'live' : 'fault'" />
          <WStat label="模型数" :value="int(models.length)" tone="route" sub="TRAE SOLO 目录" />
          <WStat label="今日请求" :value="int(usage?.today_requests || 0)" tone="brand" />
          <WStat label="今日 Tokens" :value="abbr(usage?.today_tokens || 0)" />
        </div>
      </div>

      <div class="mb-4 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <div class="glass flex items-center gap-3 rounded-xl p-4"><div class="flex h-9 w-9 items-center justify-center rounded-lg bg-brand/10 text-small font-semibold text-brand">TR</div><div class="min-w-0 flex-1"><div class="flex items-center gap-2"><span class="font-medium text-ink">TRAE SOLO</span><WLed :tone="serving ? 'live' : 'fault'" :pulse="serving" /></div><div class="truncate text-micro text-faint">{{ accounts.length }} 个账号 · {{ models.length }} 个模型</div></div><WTag :tone="serving ? 'live' : 'muted'" dot>{{ serving ? '就绪' : '等待账号' }}</WTag></div>
        <div class="glass rounded-xl p-4"><div class="text-micro uppercase tracking-wide text-faint">OpenAI 兼容</div><div class="mt-2 flex flex-wrap gap-2"><WTag tone="route">Chat Completions</WTag><WTag tone="route">Models</WTag><WTag tone="muted">Streaming</WTag></div></div>
        <div class="glass rounded-xl p-4"><div class="text-micro uppercase tracking-wide text-faint">默认模型</div><div class="mono mt-2 truncate text-small text-ink" title="{{ models[0]?.id }}">{{ models[0]?.id || 'glm-5.2' }}</div><div class="mt-1 text-micro text-faint">可在设置中查看服务配置</div></div>
      </div>

      <WCard title="账号状态" sub="按当前账号池状态排序展示。" flush>
        <WTable :columns="cols" :rows="recentAccounts" row-key="uid" min-width="620px">
          <template #cell-nickname="{ row }"><span class="font-medium text-ink">{{ row.nickname || '未命名账号' }}</span></template>
          <template #cell-credits_remaining="{ value }"><span class="mono text-ink">{{ fmtCredits(value) }}</span></template>
          <template #cell-healthy="{ row }"><WTag :tone="row.healthy ? 'live' : row.disabled ? 'fault' : 'warn'" dot>{{ row.healthy ? '可用' : row.disabled ? '已禁用' : row.cooling ? '冷却中' : '不可用' }}</WTag></template>
          <template #empty><div class="px-4 py-10 text-center text-small text-faint">暂无账号，请先到「账号」页登录或导入 TRAE 凭证。</div></template>
        </WTable>
      </WCard>
    </template>
  </WPage>
</template>
