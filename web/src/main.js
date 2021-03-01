import { createApp } from 'vue'
import App from './App.vue'
import { ElButton, ElContainer, ElHeader, ElMain, ElFooter, ElRow, ElCol, ElMessageBox } from 'element-plus'
import router from './router'
import store from './store'

// Subscribe to store updates
store.subscribe((mutation, state) => {
  // Store the state object as a JSON string
  localStorage.setItem('store', JSON.stringify(state))
})
store.commit('initialiseStore')

const app = createApp(App)

app.use(store)
app.use(router)
app.use(whatsub)

app.use(ElButton)
  .use(ElContainer)
  .use(ElHeader)
  .use(ElMain)
  .use(ElFooter)
  .use(ElRow)
  .use(ElCol)
  .use(ElMessageBox)

app.mount('#app')
