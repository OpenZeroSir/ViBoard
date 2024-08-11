import { createApp, defineAsyncComponent } from 'vue'
import App from './App.vue'
import { store } from './store/index'

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import 'vuetify/dist/vuetify.css'

// import { fa } from 'vuetify/iconsets/fa'
import { md } from 'vuetify/iconsets/md'
import 'material-design-icons-iconfont/dist/material-design-icons.css'
// import { aliases, mdi } from 'vuetify/iconsets/mdi'
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg'
// import '@mdi/font/css/materialdesignicons.css'
import VuetifyUseDialog from 'vuetify-use-dialog'

import { ComponentsName } from "./components/custom/index"

import { createRouter, createWebHashHistory } from "vue-router";
import Home from "./components/home/Home.vue"
import Display from "./components/display/Display.vue"


import 'video.js/dist/video-js.css'
import 'echarts/extension/bmap/bmap'
import * as echarts from 'echarts';

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', component: Home },
    { path: '/display', component: Display },
  ],
});

const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      md, mdi
    },
  },

})

declare global{
  interface Window {
      BMAP_PROTOCOL: string;
      BMap_loadScriptTime: number;
  }
}

function drawPolygon(points:any[],x:number,y:number,width:number,height:number){
  return echarts.graphic.clipPointsByRect(points, {
      x: x,
      y: y,
      width: width,
      height: height
  })
}

const app = createApp(App)

ComponentsName.forEach(name => {
  app.component(name,
    defineAsyncComponent(() => import(`./components/custom/${name}/index.vue`)))
})
app.provide('draw', drawPolygon);
app.use(router).use(store).use(vuetify).use(VuetifyUseDialog).mount('#app')
