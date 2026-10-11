<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '@/api/client'
import type { ModelInfo } from '@/types'
import { toast } from '@/lib/toast'
import { int } from '@/lib/format'
import WPage from '@/components/ui/WPage.vue'
import WCard from '@/components/ui/WCard.vue'
import WButton from '@/components/ui/WButton.vue'
import WInput from '@/components/ui/WInput.vue'
import WTag from '@/components/ui/WTag.vue'
import WSpinner from '@/components/ui/WSpinner.vue'
import WEmpty from '@/components/ui/WEmpty.vue'
import WIcon from '@/components/ui/WIcon.vue'

const models = ref<ModelInfo[]>([])
const source = ref('')
const search = ref('')
const loading = ref(true)
const refreshing = ref(false)
const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  return q ? models.value.filter((m) => `${m.id} ${m.name || ''}`.toLowerCase().includes(q)) : models.value
})
async function load() {
  loading.value = true
  try { const data = await api.models(); models.value = data.models || []; source.value = data.source || '' }
  catch (error) { toast.error(error instanceof Error ? error.message : '模型加载失败') }
  finally { loading.value = false }
}
async function refresh() {
  refreshing.value = true
  try { const data = await api.refreshModels(); models.value = data.models || []; source.value = data.source || ''; toast.success('模型目录已刷新') }
  catch (error) { toast.error(error instanceof Error ? error.message : '模型刷新失败') }
  finally { refreshing.value = false }
}
onMounted(load)
</script>

<template>
  <WPage title="模型" sub="TRAE SOLO 可用模型目录与上下文能力。">
    <template #actions><WInput v-model="search" placeholder="搜索模型" class="w-56" /><WButton variant="primary" :loading="refreshing" @click="refresh"><WIcon name="refresh" :size="15" />刷新目录</WButton></template>
    <div class="mb-4 grid grid-cols-2 gap-4 sm:grid-cols-3"><WCard><div class="text-micro uppercase tracking-wide text-faint">模型总数</div><div class="mono mt-1.5 text-2xl font-semibold text-brand">{{ int(models.length) }}</div></WCard><WCard><div class="text-micro uppercase tracking-wide text-faint">目录来源</div><div class="mt-2"><WTag :tone="source === 'dynamic' ? 'live' : 'muted'" dot>{{ source === 'dynamic' ? '上游实时' : '内置后备' }}</WTag></div></WCard><WCard><div class="text-micro uppercase tracking-wide text-faint">默认协议</div><div class="mono mt-2 text-small text-ink">OpenAI Chat</div></WCard></div>
    <WSpinner v-if="loading" center label="加载模型目录" />
    <WEmpty v-else-if="!filtered.length" title="没有匹配模型" hint="调整搜索词，或刷新上游模型目录。" />
    <div v-else class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      <WCard v-for="model in filtered" :key="model.id" hover>
        <div class="flex items-start justify-between gap-3"><div class="min-w-0"><h2 class="break-all font-medium text-ink">{{ model.name || model.id }}</h2><p class="mono mt-1 break-all text-micro text-faint">{{ model.id }}</p></div><WTag tone="route">TRAE</WTag></div>
        <div class="mt-4 flex flex-wrap gap-2"><WTag tone="live" dot>可调用</WTag><WTag tone="muted">{{ model.context_length ? `${int(model.context_length)} context` : '上下文信息未声明' }}</WTag></div>
      </WCard>
    </div>
  </WPage>
</template>
