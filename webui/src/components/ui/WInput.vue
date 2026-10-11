<script setup lang="ts">
// 输入框。支持 v-model、前置图标插槽、password 切换由外部控制 type。
withDefaults(
  defineProps<{ modelValue?: string | number; placeholder?: string; type?: string; disabled?: boolean; autocomplete?: string; ariaLabel?: string }>(),
  { type: 'text' },
)
defineEmits<{ 'update:modelValue': [v: string]; enter: [] }>()
</script>

<template>
  <label
    class="flex h-9 items-center gap-2 rounded-lg border border-line bg-bg/40 px-3 text-small transition-colors focus-within:border-brand"
  >
    <slot name="prefix" />
    <input
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :autocomplete="autocomplete"
      :aria-label="ariaLabel"
      class="min-w-0 flex-1 bg-transparent text-ink placeholder:text-faint focus:outline-none disabled:opacity-50"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      @keyup.enter="$emit('enter')"
    />
    <slot name="suffix" />
  </label>
</template>
