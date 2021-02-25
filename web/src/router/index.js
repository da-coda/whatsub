import { createRouter, createWebHashHistory } from 'vue-router'
import StartScreen from '../views/StartScreen.vue'

const routes = [
  {
    path: '/',
    name: 'StartScreen',
    component: StartScreen
  },
  {
    path: '/newGame',
    name: 'NewGame',
    component: () => import(/* webpackChunkName: "about" */ '../views/NewGame.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
