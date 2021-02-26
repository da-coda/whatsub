import { createApp } from 'vue'
import App from './App.vue'
import { ElButton, ElContainer, ElHeader, ElMain, ElFooter, ElRow, ElCol, ElMessageBox } from 'element-plus'
import router from './router'
import store from './store'
import axios from 'axios'
import VueAxios from 'vue-axios'

const app = createApp(App)
app.use(store) // TODO replace with react if possible
app.use(router)
app.use(VueAxios, axios)

app.use(ElButton)
  .use(ElContainer)
  .use(ElHeader)
  .use(ElMain)
  .use(ElFooter)
  .use(ElRow)
  .use(ElCol)
  .use(ElMessageBox)

app.mount('#app')
