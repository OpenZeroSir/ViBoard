<script lang="ts" setup>
// import HelloWorld from './components/HelloWorld.vue'
import Home from './components/home/Home.vue';
import { getCurrentInstance, computed } from 'vue';

// let rnd = computed(() => {
//   return Math.random()
// })

</script>

<template>
  <main>
    <v-app>
      <!-- img id="logo" alt="Wails logo" src="./assets/images/logo-universal.png"/ -->
      <!-- <Home  /> -->
      <!-- <keep-alive > -->
        <router-view ></router-view>
      <!-- </keep-alive> -->
    </v-app>
  </main>
</template>

<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
