<script setup lang="ts">
// 模态框。v-model:open 控制显隐；title + 默认插槽 + footer 插槽。
// 点遮罩/ESC 关闭。宽度经 size 控制。
import { onMounted, onUnmounted } from 'vue'
const props = withDefaults(defineProps<{ open: boolean; title?: string; size?: 'sm' | 'md' | 'lg'; dismissible?: boolean }>(), {
  size: 'md',
  dismissible: true,
})
const emit = defineEmits<{ 'update:open': [v: boolean] }>()
const width = { sm: 'max-w-md', md: 'max-w-xl', lg: 'max-w-3xl' }
function close() {
  if (!props.dismissible) return
  emit('update:open', false)
}
function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.open) close()
}
onMounted(() => window.addEventListener('keydown', onKey))
onUnmounted(() => window.removeEventListener('keydown', onKey))
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-black/55 p-4 backdrop-blur-sm sm:items-center"
        @click.self="close"
      >
        <div class="glass w-full rounded-2xl shadow-glass" :class="width[size]">
          <header v-if="title || $slots.title" class="flex items-center justify-between border-b border-line px-5 py-3.5">
            <slot name="title"><h3 class="text-[0.95rem] font-semibold text-ink">{{ title }}</h3></slot>
            <button class="text-faint hover:text-ink disabled:opacity-45" :disabled="!dismissible" @click="close" aria-label="关闭">
              <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor"><path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" /></svg>
            </button>
          </header>
          <div class="max-h-[70vh] overflow-y-auto px-5 py-4"><slot /></div>
          <footer v-if="$slots.footer" class="flex justify-end gap-2 border-t border-line px-5 py-3.5">
            <slot name="footer" />
          </footer>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
