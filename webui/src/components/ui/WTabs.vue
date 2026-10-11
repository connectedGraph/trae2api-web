<script setup lang="ts">
// 分段选项卡。v-model 绑定当前 key；tabs 为 {key,label} 数组。
defineProps<{ modelValue: string; tabs: { key: string; label: string }[] }>()
defineEmits<{ 'update:modelValue': [v: string] }>()
</script>

<template>
  <div class="inline-flex items-center gap-1 rounded-lg border border-line bg-bg/40 p-1">
    <button
      v-for="t in tabs"
      :key="t.key"
      type="button"
      class="rounded-md px-3 py-1 text-small font-medium transition-colors"
      :class="modelValue === t.key ? 'bg-brand/15 text-brand' : 'text-muted hover:text-ink'"
      @click="$emit('update:modelValue', t.key)"
    >
      {{ t.label }}
    </button>
  </div>
</template>
