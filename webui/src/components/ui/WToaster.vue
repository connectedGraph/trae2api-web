<script setup lang="ts">
// 全局 toast 渲染器。挂在 App 根部一次即可。
import { toasts, dismiss } from '@/lib/toast'
const tone: Record<string, string> = {
  info: 'border-route/40 text-route',
  success: 'border-live/40 text-live',
  error: 'border-fault/40 text-fault',
}
</script>

<template>
  <Teleport to="body">
    <div class="pointer-events-none fixed bottom-5 right-5 z-[60] flex flex-col gap-2">
      <TransitionGroup name="toast">
        <div
          v-for="t in toasts.items"
          :key="t.id"
          class="glass pointer-events-auto flex max-w-sm items-start gap-2.5 rounded-lg border-l-2 px-3.5 py-2.5 shadow-glass"
          :class="tone[t.tone]"
          @click="dismiss(t.id)"
        >
          <span class="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-current" />
          <span class="text-small text-ink">{{ t.text }}</span>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.25s ease; }
.toast-enter-from { opacity: 0; transform: translateX(16px); }
.toast-leave-to { opacity: 0; transform: translateX(16px); }
</style>
