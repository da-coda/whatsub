import { createRouter, createWebHashHistory } from 'vue-router'
import StartScreen from '../views/StartScreen.vue'

const routes = [
  {
    path: '/',
    name: 'StartScreen',
    component: StartScreen
  },
  {
    path: '/game/join',
    name: 'JoinScreenByCode',
    component: () => import(/* webpackChunkName: "JoinScreen" */ '../views/JoinScreen.vue')
  },
  {
    path: '/game/:code/join',
    name: 'JoinScreenByLink',
    props: true,
    component: () => import(/* webpackChunkName: "JoinScreen" */ '../views/JoinScreen.vue')
  },
  {
    path: '/game/:code/create',
    name: 'JoinScreenCreated',
    props (route) {
      return {
        code: route.params.code,
        isGameHead: true
      }
    },
    component: () => import(/* webpackChunkName: "JoinScreen" */ '../views/JoinScreen.vue')
  },
  {
    path: '/game/:code/lobby',
    name: 'GameLobby',
    props: true,
    component: () => import(/* webpackChunkName: "GameLobby" */ '../views/GameLobby.vue')
  },
  {
    path: '/game/:code',
    name: 'Game',
    props: true,
    component: () => import(/* webpackChunkName: "Game" */ '../views/Game.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
