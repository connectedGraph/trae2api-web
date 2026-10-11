<script setup lang="ts">
// 状态标签。tone 决定配色；默认淡底描边风，dot 可选前置圆点。
import { computed } from 'vue'
const props = withDefaults(
  defineProps<{ tone?: 'brand' | 'live' | 'warn' | 'fault' | 'route' | 'muted'; dot?: boolean }>(),
  { tone: 'muted' },
)
const map: Record<string, string> = {
  brand: 'text-brand bg-brand/10 ring-brand/25',
  live: 'text-live bg-live/10 ring-live/25',
  warn: 'text-warn bg-warn/10 ring-warn/25',
  fault: 'text-fault bg-fault/10 ring-fault/25',
  route: 'text-route bg-route/10 ring-route/25',
  muted: 'text-muted bg-muted/10 ring-muted/20',
}
const cls = computed(() => map[props.tone])
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 rounded-md px-2 py-0.5 text-micro font-medium ring-1 ring-inset"
    :class="cls"
  >
    <span v-if="dot" class="h-1.5 w-1.5 rounded-full bg-current" />
    <slot />
  </span>
</template>
