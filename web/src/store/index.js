import { createStore } from 'vuex'

export default createStore({
  state: {
    gameId: null,
    isGameHead: false
  },
  getters: {
    isGameHead: state => {
      return state.isGameHead
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
        this.replaceState(
          Object.assign(state, JSON.parse(localStorage.getItem('store')))
        )
      }
    },
    /**
     * Sets the uuid of the current game
     * @param state
     * @param {uuid} gameId
     */
    setGameId (state, gameId) {
      state.gameId = gameId
    },
    /**
     * Set if current user started the game
     * @param state
     * @param {boolean} isHead
     */
    setGameHead (state, isHead) {
      state.isGameHead = isHead
    }
  },
  actions: {
  },
  modules: {
  }
})
