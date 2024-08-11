
<template>
    <v-dialog v-model="show_dialog" width="auto">
        <v-card>
            <v-card-text>
                <div>
                    <div style="margin-left: 15px;" id="charthelpinfomation">
                        <!-- <pre>
                        <code class="markdown" v-html="store.state.viMessage"></code>
                        </pre> -->

                        <!-- {{ store.state.viMessage }} -->
                    </div>
                    <br/>
                    <div style="background-color: #01579B;">
                        <SmallInfo/>
                    </div>
                </div>

            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { ref } from 'vue'
import { store } from '../../../store/index';

import hljs from 'highlight.js';
import 'highlight.js/styles/default.css';
import { watch } from 'vue';
import { nextTick } from 'vue';
import * as marked from 'marked';
import SmallInfo from './SmallInfo.vue';
let show_dialog = computed({
    get: () => store.state.viShowChartHelp,
    set: (val) => {
        store.commit("set_viShowChartHelp", val)
    }
});
watch(() => store.state.viShowChartHelp, async (newValue) => {
    if (newValue) {
        await nextTick(); // Wait for the next DOM update
        const container = document.getElementById("charthelpinfomation");
        if (container) {
            // const codeBlocks = container.querySelectorAll('code');
            // codeBlocks.forEach((codeBlock) => {
            //     hljs.highlightBlock(codeBlock);
            // });
            // hljs.highlightBlock(container);
            // let a=  hljs.highlightAuto(store.state.viMessage,["markdown"])
            // marked(store.state.viMessage)
            container.innerHTML = marked.parse(store.state.viMessage) as any
            const codeBlocks = container.querySelectorAll('code');
            codeBlocks.forEach((codeBlock) => {
                hljs.highlightBlock(codeBlock);
            });
            // hljs.highlightBlock(container);

        }
    }
})
</script>