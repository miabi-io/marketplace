import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: Home },
    {
      path: '/templates/:name',
      name: 'template',
      component: () => import('@/views/TemplateDetail.vue'),
      props: true,
    },
    { path: '/about', name: 'about', component: () => import('@/views/About.vue') },
    { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('@/views/NotFound.vue') },
  ],
  scrollBehavior(to, _from, saved) {
    if (saved) return saved
    if (to.name === 'home' && to.query.page) return {}
    return { top: 0 }
  },
})

export default router
