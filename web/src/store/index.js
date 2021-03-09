import { createStore } from 'vuex'

export default createStore({
  state: {
    gameId: null,
    playerUUID: null,
    playerName: null,
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
        const loadedState = JSON.parse(localStorage.getItem('store'))
        loadedState.websocketConnection = null
        this.replaceState(
          Object.assign(state, loadedState)
        )
      }
    },
    /**
     * Save player game data
     * @param state
     * @param {object} payload
     */
    setGameData (state, payload) {
      state.gameId = payload.gameUUID
      state.playerUUID = payload.playerUUID
      state.playerName = payload.playerName
      state.isGameHead = payload.isGameHead
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
