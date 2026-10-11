<script setup lang="ts">
withDefaults(defineProps<{ label?: string; center?: boolean }>(), {})
</script>

<template>
  <div class="flex items-center justify-center gap-2 text-muted" :class="center ? 'py-14' : 'py-4'">
    <span class="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
    <span v-if="label" class="text-small">{{ label }}</span>
  </div>
</template>
