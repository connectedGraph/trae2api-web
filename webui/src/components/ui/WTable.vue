<script setup lang="ts">
// 轻量表格。传 columns（列定义）+ rows（数据），每列可用具名插槽
// #cell-<key> 自定义渲染，默认渲染 row[key]。表头吸顶，横向可滚动。
export interface Column {
  key: string
  label: string
  align?: 'left' | 'right' | 'center'
  mono?: boolean
  width?: string
  hint?: string
}
defineProps<{
  columns: Column[]
  rows: any[]
  rowKey?: string
  loading?: boolean
  minWidth?: string
}>()
const alignCls = { left: 'text-left', right: 'text-right', center: 'text-center' }
</script>

<template>
  <div class="overflow-x-auto">
    <table class="w-full border-collapse text-small" :style="minWidth ? { minWidth } : {}">
      <thead>
        <tr class="border-b border-line">
          <th
            v-for="c in columns"
            :key="c.key"
            class="whitespace-nowrap bg-elevated/60 px-3 py-2.5 text-micro font-semibold uppercase tracking-wide text-faint first:rounded-tl-lg last:rounded-tr-lg"
            :class="alignCls[c.align || 'left']"
            :style="c.width ? { width: c.width } : {}"
          >
            {{ c.label }}
            <span v-if="c.hint" class="ml-1 cursor-help text-faint/70" :title="c.hint">?</span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(row, i) in rows"
          :key="rowKey ? row[rowKey] : i"
          class="border-b border-line/60 transition-colors hover:bg-elevated/40"
        >
          <td
            v-for="c in columns"
            :key="c.key"
            class="px-3 py-2.5 align-middle"
            :class="[alignCls[c.align || 'left'], c.mono ? 'mono text-ink' : 'text-muted']"
          >
            <slot :name="`cell-${c.key}`" :row="row" :value="row[c.key]">{{ row[c.key] }}</slot>
          </td>
        </tr>
      </tbody>
    </table>
    <slot v-if="!rows.length && !loading" name="empty" />
  </div>
</template>
