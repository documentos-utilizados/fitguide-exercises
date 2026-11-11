import { createRouter, createWebHistory } from 'vue-router'
import TreinosView from '../views/TreinosView.vue'
import NotificacoesView from '../views/NotificacoesView.vue'

const routes = [
  {
    path: '/',
    redirect: '/treinos'
  },
  {
    path: '/treinos',
    name: 'Treinos',
    component: TreinosView
  },
  {
    path: '/notificacoes',
    name: 'Notificacoes',
    component: NotificacoesView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
