import { createRouter, createWebHashHistory } from 'vue-router'
import StartScreen from '../views/StartScreen.vue'

const routes = [
  {
    path: '/',
    name: 'StartScreen',
    component: StartScreen
  },
  {
    path: '/game/:code/join',
    name: 'JoinScreen',
    props: true,
    component: () => import(/* webpackChunkName: "JoinScreen" */ '../views/JoinScreen.vue')
  },
  {
    path: '/game/:code/lobby',
    name: 'GameLobby',
    props: true,
    component: () => import(/* webpackChunkName: "GameLobby" */ '../views/GameLobby.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
