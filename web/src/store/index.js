import { createStore } from 'vuex'

export default createStore({
  state: {
    gameId: null,
    playerUUID: null,
    isGameHead: false,
    websocketConnection: null
  },
  mutations: {
    /**
     * Run before vue app is created and get the state from local storage if
     * available
     * @param state
     */
    initialiseStore (state) {
      if (localStorage.getItem('store')) {
        this.replaceState(
          Object.assign(state, JSON.parse(localStorage.getItem('store')))
        )
      }
    },
    /**
     * Save player game data
     * @param state
     * @param {uuid} gameId
     * @param {uuid} playerUUID
     * @param {boolean} isGameHead
     */
    setGameData (state, gameId, playerUUID, isGameHead) {
      state.gameId = gameId
      state.playerUUID = playerUUID
      state.isGameHead = isGameHead
    },
    /**
     * Save WebSocketConnection
     * @param state
     * @param websocketConnection
     */
    setWebsocketConnection (state, websocketConnection) {
      state.websocketConnection = websocketConnection
    }
  },
  actions: {
  },
  modules: {
  }
})
