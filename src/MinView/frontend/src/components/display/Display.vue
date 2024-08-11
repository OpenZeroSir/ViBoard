<template>
    <!-- <div style="height: 100%;"> -->
        <MainMenu id="displayMenu" />
    <DisplayEditor />
    <Foot id="displayFoot" />
    <CopyRightInfo />
    <LoaderDialog />
    <AskPasswd />
    <ChangePasswdDialog />
    <InfoDialog />
    <!-- </div> -->
</template>

<script lang="ts" setup>
import { getCurrentInstance,onMounted } from 'vue';
import MainMenu from './MainMenu.vue'
import InfoDialog from './utils/InfoDialog.vue';
import AskPasswd from './utils/AskPasswd.vue';
import LoaderDialog from './utils/LoaderDialog.vue';
// import SessionInfo from './SessionInfo.vue'
import Foot from './Foot.vue';
import DisplayEditor from "./editor/index.vue"
import CopyRightInfo from './utils/CopyRightInfo.vue';
import ChangePasswdDialog from './utils/ChangePasswdDialog.vue';
import { EventsOn } from '../../../wailsjs/runtime/runtime';

const instance = getCurrentInstance();
// function compo_callback(data :any){
//     // console.log("compo update",data)
// }
onMounted(()=>{
    // EventsOn("compo_3SPMaDm5kz8ID4Aftz16X",compo_callback)
})


</script>