<script setup lang="ts">
import WPage from '@/components/ui/WPage.vue'
import WCard from '@/components/ui/WCard.vue'
import WTag from '@/components/ui/WTag.vue'
</script>

<template>
  <WPage title="门户管理" sub="网关访问范围与兼容能力说明。">
    <div class="grid gap-4 lg:grid-cols-2">
      <WCard title="当前部署模式" sub="trae2api 保持轻量的本机管理台模型。"><div class="flex items-center gap-3"><WTag tone="live" dot>本机管理台</WTag><span class="text-small text-muted">管理台与 OpenAI API 使用同一进程</span></div><p class="mt-4 text-small leading-relaxed text-muted">账号凭证保存在服务端配置的 auths 目录中，管理台写操作由 TW2A_API_KEY 保护。需要对外提供服务时，请在反向代理层继续限制管理台路径。</p></WCard>
      <WCard title="API 访问范围" sub="客户端通过统一的 /v1 路径调用。"><div class="space-y-3 text-small text-muted"><div class="flex justify-between gap-3"><span>模型目录</span><code class="mono text-ink">GET /v1/models</code></div><div class="flex justify-between gap-3"><span>对话补全</span><code class="mono text-ink">POST /v1/chat/completions</code></div><div class="flex justify-between gap-3"><span>健康检查</span><code class="mono text-ink">GET /health</code></div></div></WCard>
    </div>
  </WPage>
</template>
