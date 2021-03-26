import { createStore } from 'vuex'

export default createStore({
  state: {
    gameShortId: null,
    playerUUID: null,
    playerName: null,
    isGameHead: false,
    websocketConnection: null,
    scoreBoard: {}
  },
  getters: {
    isExistingGameFound: state => {
      return [state.gameId, state.playerUUID, state.playerName].every(gameAttr => gameAttr !== null)
    }
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
    clearGameData (state) {
      state.gameShortId = null
      state.playerUUID = null
      state.isGameHead = false
      state.websocketConnection = null
      state.scoreBoard = { }
    },
    /**
     * Save player game data
     * @param state
     * @param {object} payload
     */
    setGameData (state, payload) {
      state.gameShortId = payload.gameShortId
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
    },
    updateScoreBoard (state, scores) {
      state.scoreBoard = scores
    }
  },
  actions: {
  },
  modules: {
  }
})
