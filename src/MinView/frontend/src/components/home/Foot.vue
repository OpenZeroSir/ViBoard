<template>
    <v-footer class=" text-center d-flex flex-column">
        <span style="color:#01579B;">沉浸式数字末端编辑与显示服务</span>
        <div>
            <v-btn label link ripple size="24px" class="ma-2 " icon elevation="0" @click="showcopyright">
                <v-icon size="20px" color="light-blue-darken-4">{{ mdiCopyright }}</v-icon>
            </v-btn>
            <v-btn label link ripple size="24px" class="ma-2 " icon elevation="0" @click="showHelp">
                <v-icon size="20px" color="light-blue-darken-4">{{ mdiHelpCircleOutline }}</v-icon>
            </v-btn>

            <v-btn label link ripple size="24px" class="ma-2 " icon elevation="0" @click="showlicense">
                <v-icon size="20px" color="light-blue-darken-4">{{ mdiLicense }}</v-icon>
            </v-btn>
            <v-btn label link ripple size="24px" class="ma-2 " icon elevation="0" @click="showupdate">
                <v-icon size="20px" color="light-blue-darken-4">{{ mdiArrowUpBoldCircleOutline }}</v-icon>
            </v-btn>
        </div>
        <!-- <v-label :style="{
        }">愿图表离您更近一点点</v-label> -->
        <!-- <div class="pt-0">
            <span :color="light-blue-darken-4" style="font-size: small;">
                愿图表离您更近一点点</span>
        </div> -->
        <!-- <div>{{ store.state.editor_scale }}</div> -->
    </v-footer>
</template>

<script lang="ts" setup>
import { store } from "../../store";
import Editor from "./editor/index.vue"
import { mdiAt, mdiHelpCircleOutline, mdiWeb, mdiArrowUpBoldCircleOutline, mdiLicense, mdiCopyright } from '@mdi/js'
// import { store } from "../../store/store";

function showcopyright() {
    store.commit("set_viShowCopyRight", true)
}

function showHelp() {
    store.commit("set_viShowHelp", true)
}


function showupdate() {
    store.commit("set_viShowUpdate", true)
}


function showlicense() {
    store.commit("set_viShowLicense", true)
}
</script>