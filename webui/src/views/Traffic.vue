<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '@/api/client'
import type { UsageRecord, UsageSummary } from '@/types'
import { abbr, dt, int } from '@/lib/format'
import WPage from '@/components/ui/WPage.vue'
import WCard from '@/components/ui/WCard.vue'
import WStat from '@/components/ui/WStat.vue'
import WTable, { type Column } from '@/components/ui/WTable.vue'
import WTag from '@/components/ui/WTag.vue'
import WSpinner from '@/components/ui/WSpinner.vue'
import WEmpty from '@/components/ui/WEmpty.vue'
import WButton from '@/components/ui/WButton.vue'

const summary = ref<UsageSummary | null>(null)
const records = ref<UsageRecord[]>([])
const loading = ref(true)
const failed = ref(false)
const cols: Column[] = [
  { key: 'ts', label: '时间', mono: true },
  { key: 'model', label: '模型' },
  { key: 'protocol', label: '协议' },
  { key: 'account_uid', label: '账号', mono: true },
  { key: 'total_tokens', label: 'Tokens', align: 'right', mono: true },
  { key: 'status', label: '状态' },
]
const rows = computed(() => records.value as unknown as Record<string, any>[])
const topModels = computed(() => (summary.value?.by_model || []).slice(0, 8))
async function load() {
  loading.value = true
  try {
    const [s, r] = await Promise.all([api.usageSummary(), api.usageRecent()])
    summary.value = s
    records.value = r.records || []
    failed.value = false
  } catch { failed.value = true }
  finally { loading.value = false }
}
onMounted(load)
</script>

<template>
  <WPage title="流量" sub="调用量分析、调用日志与当前进程统计。">
    <template #actions><WButton :loading="loading" @click="load">刷新</WButton></template>
    <div v-if="failed" class="mb-4 rounded-lg border border-warn/40 bg-warn/5 px-3.5 py-2.5 text-small text-warn">流量接口暂不可用，可能是服务刚启动或尚未产生请求。</div>
    <WSpinner v-if="loading && !summary" center label="加载中" />
    <template v-else>
      <div class="mb-4 grid grid-cols-2 gap-4 lg:grid-cols-4"><WStat label="累计请求" :value="int(summary?.total_requests || 0)" tone="brand" /><WStat label="累计 Tokens" :value="abbr(summary?.total_tokens || 0)" tone="route" /><WStat label="今日请求" :value="int(summary?.today_requests || 0)" tone="live" /><WStat label="今日 Tokens" :value="abbr(summary?.today_tokens || 0)" /></div>
      <div class="mb-4 grid gap-4 lg:grid-cols-2">
        <WCard title="按模型" sub="请求数最多的模型。"><div v-if="!topModels.length" class="py-8 text-center text-small text-faint">暂无调用记录</div><div v-else class="space-y-3"><div v-for="item in topModels" :key="item.model" class="flex items-center gap-3"><span class="min-w-0 flex-1 truncate text-small text-muted">{{ item.model }}</span><span class="mono text-small text-ink">{{ int(item.count) }}</span><span class="mono w-16 text-right text-micro text-faint">{{ abbr(item.tokens) }}</span></div></div></WCard>
        <WCard title="采集范围" sub="统计存储在当前进程内，重启后重新开始。"><div class="space-y-3 text-small text-muted"><div class="flex justify-between gap-3"><span>协议</span><span class="mono text-ink">OpenAI Chat Completions</span></div><div class="flex justify-between gap-3"><span>统计粒度</span><span class="mono text-ink">请求与 Token</span></div><div class="flex justify-between gap-3"><span>保存位置</span><span class="mono text-ink">内存</span></div></div></WCard>
      </div>
      <WCard title="最近调用" sub="保留本次进程中最近的请求记录。" flush>
        <WTable :columns="cols" :rows="rows" row-key="id" min-width="780px"><template #cell-ts="{ value }"><span class="mono text-muted">{{ dt(value, 'MM-DD HH:mm:ss') }}</span></template><template #cell-total_tokens="{ value }"><span class="mono text-ink">{{ int(value) }}</span></template><template #cell-status="{ row }"><WTag :tone="row.status === 'ok' ? 'live' : row.status === 'error' ? 'fault' : 'muted'" dot>{{ row.status === 'ok' ? '成功' : row.status === 'error' ? '失败' : row.status }}</WTag></template><template #empty><WEmpty title="暂无调用记录" hint="通过 /v1/chat/completions 发起请求后，这里会显示最近调用。" /></template></WTable>
      </WCard>
    </template>
  </WPage>
</template>
