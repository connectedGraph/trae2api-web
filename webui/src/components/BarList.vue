<script setup lang="ts">
// 横向占比条列表。items 已排序；按最大值归一化条宽。
import { computed } from 'vue'
const props = defineProps<{
  items: { label: string; value: number; sub?: string }[]
  tone?: 'brand' | 'route' | 'live' | 'warn'
}>()
const max = computed(() => Math.max(1, ...props.items.map((i) => i.value)))
const barColor: Record<string, string> = { brand: 'bg-brand', route: 'bg-route', live: 'bg-live', warn: 'bg-warn' }
</script>

<template>
  <div class="space-y-2.5">
    <div v-for="(it, i) in items" :key="i">
      <div class="mb-1 flex items-center justify-between gap-2 text-small">
        <span class="truncate text-muted" :title="it.label">{{ it.label }}</span>
        <span class="mono shrink-0 text-ink">{{ it.value.toLocaleString('en-US') }}<span v-if="it.sub" class="ml-1 text-micro text-faint">{{ it.sub }}</span></span>
      </div>
      <div class="h-1.5 overflow-hidden rounded-full bg-elevated">
        <div class="h-full rounded-full" :class="barColor[tone || 'brand']" :style="{ width: (it.value / max) * 100 + '%' }" />
      </div>
    </div>
    <p v-if="!items.length" class="py-4 text-center text-small text-faint">暂无数据</p>
  </div>
</template>
