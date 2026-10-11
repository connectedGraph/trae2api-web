<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '@/api/client'
import { toast } from '@/lib/toast'
import WPage from '@/components/ui/WPage.vue'
import WCard from '@/components/ui/WCard.vue'
import WButton from '@/components/ui/WButton.vue'
import WTag from '@/components/ui/WTag.vue'
import WIcon from '@/components/ui/WIcon.vue'

const configured = ref(false)
const loading = ref(true)
const baseURL = window.location.origin
async function load() {
  loading.value = true
  try { configured.value = (await api.keyInfo()).configured }
  catch { configured.value = false }
  finally { loading.value = false }
}
async function copy(value: string) {
  await navigator.clipboard.writeText(value)
  toast.success('已复制')
}
const curl = `curl ${baseURL}/v1/chat/completions \\
  -H "Authorization: Bearer $TW2A_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"glm-5.2","messages":[{"role":"user","content":"你好"}]}'`
onMounted(load)
</script>

<template>
  <WPage title="API 密钥" sub="统一管理 TRAE 网关的 OpenAI 兼容访问凭证。">
    <div class="grid gap-4 lg:grid-cols-[1.1fr_1fr]">
      <WCard title="服务密钥" sub="当前版本使用单个 TW2A_API_KEY，写入环境变量后重启服务生效。">
        <div class="flex items-center gap-3"><WTag :tone="configured ? 'live' : 'warn'" dot>{{ configured ? '已配置' : '未配置（本机免鉴权）' }}</WTag><span class="text-small text-muted">{{ configured ? 'API 和管理写操作均受 Bearer Key 保护' : '仅建议在本机回环环境使用' }}</span></div>
        <div class="mt-5 rounded-lg border border-line bg-bg/40 p-3 text-small text-muted"><p><code class="mono text-brand">TW2A_API_KEY</code> 用于：</p><ul class="mt-2 list-disc space-y-1 pl-5"><li>调用 <code class="mono text-ink">/v1/models</code> 和 <code class="mono text-ink">/v1/chat/completions</code></li><li>账号导入、删除、刷新和管理台写操作</li><li>管理台登录会话的访问验证</li></ul></div>
      </WCard>
      <WCard title="调用地址" sub="与 OpenAI SDK、Chatbox、Cline 等客户端兼容。">
        <div class="flex items-center gap-2 rounded-lg border border-line bg-bg/40 p-3"><code class="mono min-w-0 flex-1 break-all text-small text-brand">{{ baseURL }}/v1</code><WButton size="sm" @click="copy(`${baseURL}/v1`)"><WIcon name="copy" :size="15" />复制</WButton></div>
        <p class="mt-3 text-micro text-faint">模型列表：{{ baseURL }}/v1/models · 对话：{{ baseURL }}/v1/chat/completions</p>
      </WCard>
    </div>
    <WCard class="mt-4" title="调用示例" sub="将环境变量替换为实际的 TW2A_API_KEY。">
      <div class="relative"><pre class="overflow-x-auto rounded-lg border border-line bg-bg p-4 text-small leading-relaxed text-muted">{{ curl }}</pre><WButton class="absolute right-3 top-3" size="sm" @click="copy(curl)"><WIcon name="copy" :size="15" />复制</WButton></div>
    </WCard>
    <WCard v-if="loading" class="mt-4"><span class="text-small text-faint">正在读取密钥状态…</span></WCard>
  </WPage>
</template>
