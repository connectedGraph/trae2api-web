<script setup lang="ts">
import { ref } from 'vue'
import WButton from './ui/WButton.vue'
import WInput from './ui/WInput.vue'
const props = defineProps<{ change: (oldPassword: string, newPassword: string) => Promise<unknown>; forced?: boolean }>()
const emit = defineEmits<{ done: [] }>()
const oldPassword = ref(''), newPassword = ref(''), confirm = ref(''), error = ref(''), busy = ref(false)
async function submit() {
  if (busy.value) return
  error.value = ''
  if (!oldPassword.value || !newPassword.value) { error.value = '请填写当前密码和新密码'; return }
  if (newPassword.value !== confirm.value) { error.value = '两次新密码不一致'; return }
  if (oldPassword.value === newPassword.value) { error.value = '新密码必须与当前密码不同'; return }
  busy.value = true
  try { await props.change(oldPassword.value, newPassword.value); oldPassword.value = newPassword.value = confirm.value = ''; emit('done') }
  catch (e: any) { error.value = e?.response?.data?.error?.message || e?.message || '修改失败' }
  finally { busy.value = false }
}
</script>
<template>
  <form class="space-y-3" @submit.prevent="submit">
    <p class="text-small text-muted">{{ forced ? '首次登录或管理员重置密码后，请先设置自己的密码。' : '修改密码后，其他设备需要重新登录。管理入口也需要重新登录。' }}</p>
    <div class="text-small">当前密码<WInput v-model="oldPassword" aria-label="当前密码" type="password" autocomplete="current-password" class="mt-1" /></div>
    <div class="text-small">新密码<WInput v-model="newPassword" aria-label="新密码" type="password" autocomplete="new-password" placeholder="至少 8 个字符，不超过 72 字节" class="mt-1" /></div>
    <div class="text-small">确认新密码<WInput v-model="confirm" aria-label="确认新密码" type="password" autocomplete="new-password" class="mt-1" @enter="submit" /></div>
    <p v-if="error" role="alert" class="text-small text-fault">{{ error }}</p>
    <WButton variant="primary" :loading="busy" @click="submit">保存新密码</WButton>
  </form>
</template>
