import { createApp } from 'vue'
import App from './App.vue'
import { initFaro } from './faro'

initFaro()

const app = createApp(App)
app.mount('#app')