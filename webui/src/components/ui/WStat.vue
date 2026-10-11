<script setup lang="ts">
// 读数瓦片：小标签 + 大号等宽数值 + 可选副文/趋势。用于概览 KPI 行。
withDefaults(defineProps<{ label: string; value: string | number; sub?: string; tone?: 'ink' | 'brand' | 'live' | 'warn' | 'fault' | 'route' }>(), {
  tone: 'ink',
})
const toneCls: Record<string, string> = {
  ink: 'text-ink',
  brand: 'text-brand',
  live: 'text-live',
  warn: 'text-warn',
  fault: 'text-fault',
  route: 'text-route',
}
</script>

<template>
  <div class="glass rounded-xl px-4 py-3.5">
    <div class="text-micro font-medium uppercase tracking-wide text-faint">{{ label }}</div>
    <div class="mono mt-1.5 text-2xl font-semibold leading-none" :class="toneCls[tone]">{{ value }}</div>
    <div v-if="sub" class="mt-1.5 text-micro text-muted">{{ sub }}</div>
  </div>
</template>
