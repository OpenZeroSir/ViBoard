<template>
    <div v-ripple="storeComputed.state.viRipple" data-vjs-player :nni="props.nni" :style="{
        height: '100%',
        width: '100%',
    }">
        <!-- <video :id="'video' + props.nni" class="video-js vjs-default-skin" controls loop
            :autoplay="props.innerStyle.vertical" :style="{
                height: '100%',
                width: '100%',
            }" >
        </video> -->
        <video :id="'video' + props.nni" class="video-js" :style="{
            height: '100%',
            width: '100%',
        }">
        </video>
        <!-- <VideoPlayer ref="videoPlayer" class="video-js" :loop="true" src="https://media.w3.org/2010/05/sintel/trailer.mp4"></VideoPlayer> -->
        <!-- <video-js ref="player" :options="playerOptions"></video-js> -->
    </div>
</template>

<script lang="ts" setup>
import { onMounted, computed, reactive, ref, nextTick, onBeforeUnmount } from 'vue'
import { ViChartStyle, ViInnerStyle, ViOutStyle, store } from '../../../store';
import videojs from "video.js";
// import 'video.js/dist/video-js.css'
// import 'video.js/src/css/video-js.scss';
import { getCurrentInstance } from 'vue';
import { storeDisplay } from '../../../store/display';
import { watch } from 'vue';
import 'video.js/dist/video-js.css'


let props = defineProps<{
    nni: string;
    name: string;
    innerStyle: ViInnerStyle;
    outStyle: ViOutStyle;
    chartStyle: ViChartStyle | undefined,
}>()

let player: videojs.Player

async function playVideo(url: string | undefined) {
    if (url == undefined || url == "") {
        return
    }
    // video.src = ""
    if (url.startsWith("NIBoard")) {

        await fetchAndDisplayVideo(url)
        // video.load();
    } else {
        player.src(url);
    }

}

onMounted(() => {
    // let video = document.getElementById('video' + props.nni) as HTMLVideoElement;
    // video.addEventListener("play", playVideo)
    // video.addEventListener("pause",pauseVideo)
    // 获取视频容器元素
    // let video = document.getElementById('video' + props.nni);

    // 配置视频播放器选项
    // const options = {
    //     // controls: true, // 显示控制栏
    //     // autoplay: true, // 是否自动播放
    //     // sources: [{
    //     //     src: 'https://media.w3.org/2010/05/sintel/trailer.mp4', // 视频文件的 URL
    //     //     type: 'video/mp4',
    //     // }],
    //     Preload:'auto',
    //     loop:false,
    // };

    // 初始化视频播放器
    // videojs()
    player = videojs('video' + props.nni);
    // // player.ready()
    // player.load()
    // player.play()
    player.on('click', function () {
        if (player.paused()) {
            player.play(); // 如果当前是暂停状态，点击后播放
        } else {
            player.pause(); // 如果当前是播放状态，点击后暂停
        }
    });
    player.on('dblclick', function () {
        if (player.isFullscreen()) {
            player.exitFullscreen(); // 如果当前是全屏状态，双击后退出全屏
        } else {
            player.requestFullscreen(); // 如果当前不是全屏状态，双击后进入全屏
        }
    });
    player.loop(true)
    player.preload(true)
    playVideo(props.innerStyle.text)

})

onBeforeUnmount(() => {
    if (player) {
        player.dispose()
    }
    // let video = document.getElementById('video' + props.nni) as HTMLVideoElement;
    // video.removeEventListener("play", playVideo)
    // video.removeEventListener("pause",pauseVideo)
})


// const instance = getCurrentInstance();

let storeComputed = computed(() => {
    // let fullpath = instance?.proxy?.$router.currentRoute.value.fullPath
    if (!props.outStyle.displayMode) {
        return store
    } else {
        return storeDisplay
    }
})

watch(() => props.innerStyle.text, (newValue) => {
    if (newValue != undefined &&  newValue != "") {
        playVideo(newValue)
        player.play()
    }
})


async function fetchAndDisplayVideo(url: string) {
    try {
        if (player.src().length > 0) {
            try {
                URL.revokeObjectURL(player.src())
                player.src("")
            } catch (error) {
                console.log(error)
            }
            // videoElement.play()
            // return ""
        }
        const mediaSource = new MediaSource();
        player.src(URL.createObjectURL(mediaSource));

        async function updateDate() {
            const sourceBuffer = mediaSource.addSourceBuffer("video/mp4");

            const response = await fetch(url);
            if (response.body != null) {
                const reader = response.body.getReader();

                async function appendData() {
                    const { done, value } = await reader.read();

                    if (done) {
                        mediaSource.endOfStream();
                        reader.releaseLock();
                        sourceBuffer.removeEventListener("updateend", appendData);
                        mediaSource.removeEventListener("sourceopen", updateDate);
                    } else {
                        sourceBuffer.addEventListener("updateend", appendData);
                        sourceBuffer.appendBuffer(value);
                    }
                }

                appendData();
            }

        }

        mediaSource.addEventListener("sourceopen", updateDate);



        // 设置 video 元素的其他属性
        // videoElement.controls = true;
        // videoElement.autoplay = true;
        // videoElement.load()
        // videoElement.play()
    } catch (error) {
        console.error('Error fetching and displaying video:', error);
    }

}





</script>