import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import '@fontsource-variable/inter'
import '@fontsource/jetbrains-mono/400.css'
import '@fontsource/jetbrains-mono/500.css'
import '@fontsource/jetbrains-mono/700.css'
import './styles/app.css'
import App from './App.vue'

const savedTheme = localStorage.getItem('trae_theme')
const wantLight = savedTheme
  ? savedTheme === 'light'
  : !!window.matchMedia?.('(prefers-color-scheme: light)').matches
document.documentElement.classList.toggle('light', wantLight)

const router = createRouter({
  history: createWebHashHistory('/admin-ui/'),
  routes: [
    { path: '/', redirect: '/overview' },
    { path: '/overview', component: () => import('./views/Overview.vue'), meta: { title: '概览' } },
    { path: '/accounts', component: () => import('./views/Accounts.vue'), meta: { title: '账号' } },
    { path: '/models', component: () => import('./views/Models.vue'), meta: { title: '模型' } },
    { path: '/keys', component: () => import('./views/Keys.vue'), meta: { title: 'API 密钥' } },
    { path: '/traffic', component: () => import('./views/Traffic.vue'), meta: { title: '流量' } },
    { path: '/settings', component: () => import('./views/Settings.vue'), meta: { title: '设置' } },
    { path: '/portal', component: () => import('./views/Portal.vue'), meta: { title: '门户管理' } },
    { path: '/:pathMatch(.*)*', redirect: '/overview' },
  ],
})

createApp(App).use(router).mount('#app')
