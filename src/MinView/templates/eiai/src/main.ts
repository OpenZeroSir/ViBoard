import { createApp } from 'vue'
// import './style.css'
import App from './App.vue'

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import { createRouter, createWebHashHistory } from "vue-router";
import { md } from 'vuetify/iconsets/md'
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg'
import 'material-design-icons-iconfont/dist/material-design-icons.css'
// import { aliases, md } from 'vuetify/iconsets/md'

import CompCenter from "./components/CompCenter.vue"
import DataInfoCenter from "./components/DataInfoCenter.vue"

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', component: DataInfoCenter },
    { path: '/comp', component: CompCenter },
  ],
});

const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      md,mdi,
    },
  },
})

createApp(App).use(vuetify).use(router).mount('#app')
