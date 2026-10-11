<script setup lang="ts">
defineProps<{ modelValue?: string; placeholder?: string; rows?: number }>()
defineEmits<{ 'update:modelValue': [v: string] }>()
</script>

<template>
  <textarea
    :value="modelValue"
    :placeholder="placeholder"
    :rows="rows || 4"
    class="w-full resize-y rounded-lg border border-line bg-bg/40 px-3 py-2 text-small leading-relaxed text-ink placeholder:text-faint transition-colors focus:border-brand focus:outline-none"
    @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
  />
</template>
