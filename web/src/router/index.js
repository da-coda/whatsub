import { createRouter, createWebHashHistory } from 'vue-router'
import StartScreen from '../views/StartScreen.vue'

const routes = [
  {
    path: '/',
    name: 'StartScreen',
    component: StartScreen
  },
  {
    path: '/game/join/:code',
    name: 'NewGame',
    props: true,
    component: () => import(/* webpackChunkName: "GameLobby" */ '../views/GameLobby.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
