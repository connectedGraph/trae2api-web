<script setup lang="ts">
// 按钮：一套克制的变体。primary=青色实心，ghost=描边，subtle=淡底，
// danger=错误色描边。所有变体统一高度与圆角，图标用默认插槽前置。
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'ghost' | 'subtle' | 'danger'
    size?: 'sm' | 'md'
    disabled?: boolean
    loading?: boolean
    block?: boolean
    type?: 'button' | 'submit'
  }>(),
  { variant: 'ghost', size: 'md', type: 'button' },
)

const cls = computed(() => {
  const base =
    'inline-flex items-center justify-center gap-1.5 rounded-lg font-medium transition-colors select-none whitespace-nowrap disabled:opacity-45 disabled:cursor-not-allowed'
  const size = props.size === 'sm' ? 'h-8 px-3 text-small' : 'h-9 px-3.5 text-small'
  const v: Record<string, string> = {
    primary: 'bg-brand text-white hover:bg-brand-soft shadow-sm',
    ghost: 'border border-line text-ink hover:border-brand hover:text-brand bg-transparent',
    subtle: 'bg-elevated text-muted hover:text-ink hover:bg-line/40',
    danger: 'border border-fault/50 text-fault hover:bg-fault/10',
  }
  return [base, size, v[props.variant], props.block ? 'w-full' : ''].join(' ')
})
</script>

<template>
  <button :type="type" :class="cls" :disabled="disabled || loading">
    <span v-if="loading" class="h-3.5 w-3.5 animate-spin rounded-full border-2 border-current border-t-transparent" />
    <slot />
  </button>
</template>
