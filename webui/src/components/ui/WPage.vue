<script setup lang="ts">
// 页面容器：统一标题区（标题 + 说明 + 右侧操作插槽）与内边距、入场动画。
defineProps<{ title: string; sub?: string }>()
</script>

<template>
  <div class="mx-auto max-w-[1320px] animate-fade-up px-4 py-6 sm:px-6 lg:px-8">
    <header class="mb-5 flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-xl font-semibold tracking-tight text-ink">{{ title }}</h1>
        <p v-if="sub" class="mt-1 text-small text-muted">{{ sub }}</p>
      </div>
      <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2"><slot name="actions" /></div>
    </header>
    <slot />
  </div>
</template>
