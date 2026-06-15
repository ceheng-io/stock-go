import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'
import Dashboard from '@/pages/Dashboard/Dashboard.vue'

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: AppLayout,
      children: [
        { path: '', name: 'dashboard', component: Dashboard },
        { path: 'heatmap', name: 'heatmap', component: () => import('@/pages/Heatmap/Heatmap.vue') },
        { path: 'rankings', name: 'rankings', component: () => import('@/pages/Rankings/Rankings.vue') },
        { path: 'boards', name: 'boards', component: () => import('@/pages/Boards/Boards.vue') },
        {
          path: 'boards/:type/:code',
          name: 'board-detail',
          component: () => import('@/pages/Boards/BoardDetail.vue'),
        },
        { path: 'watchlist', name: 'watchlist', component: () => import('@/pages/Watchlist/Watchlist.vue') },
        { path: 'scanner', name: 'scanner', component: () => import('@/pages/Scanner/Scanner.vue') },
        {
          path: 'eod-picker',
          name: 'eod-picker',
          component: () => import('@/pages/EndOfDayPicker/EndOfDayPicker.vue'),
        },
        { path: 'settings', name: 'settings', component: () => import('@/pages/Settings/Settings.vue') },
        { path: 's/:code', name: 'stock-detail', component: () => import('@/pages/StockDetail/StockDetail.vue') },
      ],
    },
  ],
})
