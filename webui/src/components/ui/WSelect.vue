<script setup lang="ts">
// 原生 select 的样式化封装。options 为 {value,label} 数组；空值项由 placeholder 提供。
defineProps<{
  modelValue?: string | number
  options: { value: string | number; label: string }[]
  placeholder?: string
  disabled?: boolean
}>()
defineEmits<{ 'update:modelValue': [v: string] }>()
</script>

<template>
  <div class="relative">
    <select
      :value="modelValue"
      :disabled="disabled"
      class="h-9 w-full appearance-none rounded-lg border border-line bg-bg/40 pl-3 pr-8 text-small text-ink transition-colors focus:border-brand focus:outline-none disabled:opacity-50"
      @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
    >
      <option v-if="placeholder" value="">{{ placeholder }}</option>
      <option v-for="o in options" :key="String(o.value)" :value="o.value">{{ o.label }}</option>
    </select>
    <svg class="pointer-events-none absolute right-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-faint" viewBox="0 0 20 20" fill="currentColor">
      <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.17l3.71-3.94a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
    </svg>
  </div>
</template>
