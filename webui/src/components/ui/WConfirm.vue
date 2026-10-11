<script setup lang="ts">
// 全局确认对话框宿主。挂在 App 根部一次；由 confirm() 驱动。
import WModal from './WModal.vue'
import WButton from './WButton.vue'
import { confirmState, resolveConfirm } from '@/lib/confirm'
</script>

<template>
  <WModal :open="confirmState.open" size="sm" :title="confirmState.title" @update:open="(v) => !v && resolveConfirm(false)">
    <p v-if="confirmState.body" class="text-small text-muted">{{ confirmState.body }}</p>
    <template #footer>
      <WButton variant="subtle" @click="resolveConfirm(false)">{{ confirmState.cancelText || '取消' }}</WButton>
      <WButton :variant="confirmState.tone === 'danger' ? 'danger' : 'primary'" @click="resolveConfirm(true)">
        {{ confirmState.okText || '确定' }}
      </WButton>
    </template>
  </WModal>
</template>
