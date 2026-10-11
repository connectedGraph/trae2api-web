<script setup lang="ts">
// 状态指示灯（LED 圆点）。live 态带呼吸动画表示"在线/服务中"。
import { computed } from 'vue'
const props = withDefaults(
  defineProps<{ tone?: 'live' | 'warn' | 'fault' | 'route' | 'muted'; pulse?: boolean }>(),
  { tone: 'muted' },
)
const color = computed(
  () =>
    ({ live: 'bg-live', warn: 'bg-warn', fault: 'bg-fault', route: 'bg-route', muted: 'bg-faint' })[
      props.tone
    ],
)
</script>

<template>
  <span class="relative inline-flex h-2.5 w-2.5">
    <span
      v-if="pulse"
      class="absolute inline-flex h-full w-full animate-ping rounded-full opacity-60"
      :class="color"
    />
    <span class="relative inline-flex h-2.5 w-2.5 rounded-full" :class="color" />
  </span>
</template>
