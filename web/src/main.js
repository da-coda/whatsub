import { createApp } from 'vue'
import App from './App.vue'
import { ElButton, ElContainer, ElHeader, ElMain, ElFooter, ElRow, ElCol, ElMessageBox } from 'element-plus'
import router from './router'
import store from './store'

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

app.mount('#app')
