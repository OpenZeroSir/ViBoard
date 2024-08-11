<template>
    <div v-ripple="storeComputed.state.viRipple" :nni="props.nni" :style="{
        height: '100%',
        width: '100%',
    }">
        <audio :id="'audio' + props.nni"></audio>
        <div :style="{
            height: '100%',
            width: '100%',
            transform: `translate(-50%, -50%)`,
            position: 'absolute',
            left: '50%',
            top: '50%',

        }">
            <v-icon @click.stop="handlePlay" :size="size" :style="{
                height: '100%',
                width: '100%',
                position: 'absolute',
                left: '50%',
                top: '50%',
                transform: 'translate(-50%, -50%)',
                color: props.innerStyle.color,
            }">{{ SvgIconName[playicon as keyof typeof SvgIconName] }}</v-icon>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, computed, reactive, ref, nextTick, getCurrentInstance } from 'vue'
import { ViChartStyle, ViInnerStyle, ViOutStyle, store } from '../../../store';
import * as SvgIconName from '@mdi/js'
import { min } from 'mathjs';
import { EventsOn, EventsOff } from '../../../../wailsjs/runtime/runtime';
import { storeDisplay } from '../../../store/display';
import { watch } from 'vue';
let props = defineProps<{
    nni: string;
    name: string;
    innerStyle: ViInnerStyle;
    outStyle: ViOutStyle;
    chartStyle: ViChartStyle | undefined
}>()

// const instance = getCurrentInstance();

let storeComputed = computed(() => {
    // let fullpath = instance?.proxy?.$router.currentRoute.value.fullPath
    if (!props.outStyle.displayMode) {
        return store
    } else {
        return storeDisplay
    }
})


let size = reactive(computed(() => {
    let w = Number(props.outStyle.width)
    let h = Number(props.outStyle.height)
    return min(w, h) * storeComputed.value.state.canvas.scale / 100
}))


let playicon = ref("mdiMicrophone")
function handleAudioEnded() {
    playicon.value = "mdiMicrophone"
    audio.removeEventListener('ended', handleAudioEnded);
}
let audio: HTMLAudioElement
async function handlePlay() {
    audio = document.getElementById('audio' + props.nni) as HTMLAudioElement
    if (audio != null) {
        // if(audio.src ==""){
        //     return
        // }
        if (props.innerStyle.text == undefined || props.innerStyle.text == "") {
            return
        }
        await fetchAndDisplayAudio(props.innerStyle.text)
        audio.addEventListener('ended', handleAudioEnded);
        if (playicon.value != "mdiPause") {
            audio.play();
            playicon.value = "mdiPause"
        } else {
            audio.pause()
            playicon.value = "mdiMicrophone"
        }
    }
}

watch(() => props.innerStyle.text, (newValue) => {
    if (newValue != undefined && newValue != "") {
        // audio.pause()
        // audio.play();
        // playicon.value = "mdiPause"
        fetchAndDisplayAudio(newValue)
        setTimeout(() => {
            let audioElement = document.getElementById('audio' + props.nni) as HTMLAudioElement;
            audioElement.play();
            playicon.value = "mdiPause"
        }, 1000);
    }
})

onMounted(() => {
    // EventsOn("compo_" + props.nni, receiveCompoData)
})
onUnmounted(() => {
    // EventsOff("compo_" + props.nni)
})


async function fetchAndDisplayAudio(url: string) {
    try {
        let audioElement = document.getElementById('audio' + props.nni) as HTMLAudioElement;
        if (url.startsWith("http")) {
            audioElement.src = url
            return
        }

        if (props.innerStyle.text == undefined) {
            return ""
        }
        const response = await fetch(url);
        if (response.body == null) {
            return ""
        }
        const reader = response.body.getReader();

        const stream = new ReadableStream({
            async start(controller) {
                while (true) {
                    const { done, value } = await reader.read();

                    if (done) {
                        controller.close();
                        reader.releaseLock();
                        break;
                    }

                    controller.enqueue(value);

                    // 将读取到的数据追加到视频元素的缓冲区
                    const chunk = new Uint8Array(value);
                    audioElement.src = URL.createObjectURL(new Blob([chunk], { type: 'video/mp4' }));
                }
            }
        });

        // 创建一个包含可读流的 blob，并将其分配给 video 元素
        const blob = await new Response(stream).blob();
        audioElement.src = URL.createObjectURL(blob);
        return audioElement.src
    } catch (error) {
        console.error('Error fetching and displaying video:', error);
    }
    return ""
}
</script>

