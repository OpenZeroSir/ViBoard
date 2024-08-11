<template>
    <v-dialog v-model="store.state.viShowHtml" width="auto">
        <v-card>

            <v-card-text>
                <div id="codeblockinshowhtmlinfo" v-html="store.state.viHtmlContext"></div>
            </v-card-text>
            <div style="background-color: #01579B;padding:10px">
                <SmallInfo />
            </div>
        </v-card>
    </v-dialog>
</template>
  
<script lang="ts" setup>
import { store } from '../../../store/index';
import hljs from 'highlight.js';
import 'highlight.js/styles/default.css';
import { nextTick, onBeforeMount, watch } from 'vue';
import SmallInfo from './SmallInfo.vue';

watch(() => store.state.viShowHtml, async (newValue) => {
    if (newValue) {
        await nextTick(); // Wait for the next DOM update
        const container = document.getElementById("codeblockinshowhtmlinfo");
        if (container) {
            const codeBlocks = container.querySelectorAll('code');
            codeBlocks.forEach((codeBlock) => {
                hljs.highlightBlock(codeBlock);
            });
        }
    }
})
</script>



  