import { createApp } from 'vue'
import { ElButton, ElContainer, ElHeader, ElMain, ElFooter, ElRow, ElCol, ElMessageBox, ElDivider, ElLoading, ElDialog, ElForm, ElInput, ElMessage } from 'element-plus'
import '../element-variables.scss'
import App from './App.vue'

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
