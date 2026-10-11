<script setup lang="ts">
// 空状态：一句方向性说明 + 可选操作插槽（而非单纯"暂无数据"）。
defineProps<{ title?: string; hint?: string }>()
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-2 px-6 py-14 text-center">
    <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-elevated text-faint">
      <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9.75h16.5M3.75 15h16.5M6 4.5h12A2.25 2.25 0 0120.25 6.75v10.5A2.25 2.25 0 0118 19.5H6a2.25 2.25 0 01-2.25-2.25V6.75A2.25 2.25 0 016 4.5z" />
      </svg>
    </div>
    <p class="text-small font-medium text-muted">{{ title || '这里还没有内容' }}</p>
    <p v-if="hint" class="max-w-xs text-micro text-faint">{{ hint }}</p>
    <div v-if="$slots.default" class="mt-1"><slot /></div>
  </div>
</template>
