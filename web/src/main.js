import { createApp } from 'vue'
import { ElButton, ElContainer, ElHeader, ElMain, ElFooter, ElRow, ElCol, ElMessageBox, ElDivider, ElLoading, ElDialog, ElForm, ElInput, ElMessage } from 'element-plus'
import '../element-variables.scss'
import App from './App.vue'

import router from './router'
import store from './store'
import { joinGame } from '@/lib/whatsub'

// Subscribe to store updates
store.subscribe((mutation, state) => {
  // Store the state object as a JSON string
  localStorage.setItem('store', JSON.stringify(state))
})
store.commit('initialiseStore')

console.log('Game was running: ' + store.getters.isExistingGameFound)

if (store.getters.isExistingGameFound) {
  const playerName = store.state.playerName
  const playerUUID = store.state.playerUUID
  const gameShortId = store.state.gameShortId
  const webSocket = joinGame(playerName, gameShortId, playerUUID)

  webSocket.onopen = ev => {
    store.commit('setWebsocketConnection', webSocket)
    store.commit('setGameData', {
      gameShortId: gameShortId,
      playerUUID: playerUUID,
      playerName: playerName,
      isGameHead: store.state.isGameHead
    })
  }

  webSocket.onerror = ev => {
    console.log('Existing game found, but websocket connection could not be established. Going back to start screen.')
    store.commit('clearGameData')
    router.push('/')
  }
}
const app = createApp(App)

app.use(store)
app.use(router)

app.use(ElButton)
  .use(ElContainer)
  .use(ElHeader)
  .use(ElMain)
  .use(ElFooter)
  .use(ElRow)
  .use(ElCol)
  .use(ElMessageBox)
  .use(ElDivider)
  .use(ElLoading)
  .use(ElDialog)
  .use(ElForm)
  .use(ElInput)
  .use(ElMessage)

app.mount('#app')
