import { createStore } from 'vuex'

export default createStore({
  state: {
    gameId: null
  },
  mutations: {
    setGameId (state, gameId) {
      state.gameId = gameId
    }
  },
  actions: {
  },
  modules: {
  }
})
