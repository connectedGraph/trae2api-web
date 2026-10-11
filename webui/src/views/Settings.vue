<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '@/api/client'
import type { Settings as SettingsData } from '@/types'
import { toast } from '@/lib/toast'
import WPage from '@/components/ui/WPage.vue'
import WCard from '@/components/ui/WCard.vue'
import WButton from '@/components/ui/WButton.vue'
import WSpinner from '@/components/ui/WSpinner.vue'
import WTag from '@/components/ui/WTag.vue'

const settings = ref<SettingsData | null>(null)
const loading = ref(true)
async function load() {
  loading.value = true
  try { settings.value = await api.settings() }
  catch (error) { toast.error(error instanceof Error ? error.message : '配置读取失败') }
  finally { loading.value = false }
}
onMounted(load)
</script>

<template>
  <WPage title="设置" sub="服务监听、账号池、刷新策略与 API 兼容配置。">
    <template #actions><WButton :loading="loading" @click="load">刷新配置</WButton></template>
    <WSpinner v-if="loading" center label="读取服务配置" />
    <div v-else-if="settings" class="grid gap-4 lg:grid-cols-2">
      <WCard title="服务" sub="当前进程启动时读取的配置。"><div class="space-y-3 text-small"><div class="flex justify-between gap-4"><span class="text-muted">监听地址</span><code class="mono text-ink">{{ settings.listen }}</code></div><div class="flex justify-between gap-4"><span class="text-muted">登录回调端口</span><code class="mono text-ink">{{ settings.callback_port }}</code></div><div class="flex justify-between gap-4"><span class="text-muted">默认模型</span><code class="mono text-ink">{{ settings.default_model }}</code></div><div class="flex justify-between gap-4"><span class="text-muted">API Key</span><WTag :tone="settings.api_key_configured ? 'live' : 'warn'" dot>{{ settings.api_key_configured ? '已配置' : '未配置' }}</WTag></div></div></WCard>
      <WCard title="账号池策略" sub="错误与权益限制会触发账号冷却和轮换。"><div class="space-y-3 text-small"><div class="flex justify-between gap-4"><span class="text-muted">权益不足冷却</span><code class="mono text-ink">{{ settings.plan_credit }}</code></div><div class="flex justify-between gap-4"><span class="text-muted">429 短冷却</span><code class="mono text-ink">{{ settings.soft_rate }}</code></div><div class="flex justify-between gap-4"><span class="text-muted">错误阈值</span><code class="mono text-ink">{{ settings.err_threshold }}</code></div><div class="flex justify-between gap-4"><span class="text-muted">错误冷却</span><code class="mono text-ink">{{ settings.err_cooldown }}</code></div></div></WCard>
      <WCard title="路径" sub="凭证和池状态文件由启动配置决定。"><div class="space-y-3 text-small"><div><div class="mb-1 text-muted">凭证目录</div><code class="mono break-all text-ink">{{ settings.auth_dir }}</code></div><div><div class="mb-1 text-muted">状态文件</div><code class="mono break-all text-ink">{{ settings.state_file }}</code></div></div></WCard>
      <WCard title="上游请求" sub="网络请求超时和协议适配由服务端统一处理。"><div class="space-y-3 text-small"><div class="flex justify-between gap-4"><span class="text-muted">短请求超时</span><code class="mono text-ink">{{ settings.upstream_timeout_seconds }} 秒</code></div><div class="flex justify-between gap-4"><span class="text-muted">支持流式输出</span><WTag tone="live" dot>已启用</WTag></div><div class="flex justify-between gap-4"><span class="text-muted">模型目录</span><WTag tone="route">动态 + 内置后备</WTag></div></div></WCard>
    </div>
  </WPage>
</template>
