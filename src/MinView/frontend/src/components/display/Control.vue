<template>
    <v-list :style="{
        height: storeDisplay.state.controlHeight - 40 + 'px',
    }" class="bg-light-blue-darken-4" nav>
        <v-list-item v-for="(n, index ) in  storeDisplay.state.viViews " :key="n.id" :value="index"
            @click.stop="handleClick(index)" :v-show="n.visible" color="white" nav>
            <div :style="{

                width: '100%',
                // height: '100px',
                // backgroundColor: '#003355',
                // backgroundImage: 'url('+storeDisplay.state.viViews[index].shot+')',
                // backgroundSize:'cover',
            }" v-ripple>
                <!-- <template v-slot:title>
                    <div :style="{
                        color:storeDisplay.state.viCtrlIndex==index?'1px solid red':'1px solid #003355'
                    }">{{ index+1 }}</div>
                </template> -->
                <!-- <v-img :src="shot_img(index) ? img_url[index] : img_url[index]"></v-img> -->
                <!-- <v-avatar :color="storeDisplay.state.viCtrlIndex == index ? '#ff9900' : '#003355'" size="80"></v-avatar> -->

                <img :id="'img' + index" @error="handleImageError" :src="storeDisplay.state.viViews[index].shot" :style="{
                    width: '100%', height: '100%',
                    border: storeDisplay.state.viCtrlIndex == index ? '2px solid red' : '2px solid #003355'
                }">

            </div>
        </v-list-item>
    </v-list>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { storeDisplay } from '../../store/display';
import { onBeforeUnmount } from 'vue';
import { nextTick } from 'vue';
import { fetchEX } from '../../utils/func';

function handleClick(index: number) {
    storeDisplay.commit("set_view_index", index)
    fetch(storeDisplay.state.viViews[index].shot)
}

function handleResize() {
    nextTick(() => {
        let el = storeDisplay.state.display
        if (el != null) {
            let rect = el.getBoundingClientRect()
            storeDisplay.commit("set_control_height", rect.height)
        }
    })
}

function handleImageError(event: any) {
    var img = event.target;
    var retryCount = img.getAttribute('data-retry-count') || 0;

    if (retryCount < 5) { // 最多重试3次
        retryCount++;
        img.setAttribute('data-retry-count', retryCount);

        setTimeout(function () {
            img.src = img.src; // 重新加载图片
        }, 1000); // 延迟1秒后重试
    } else {
        img.src = ''; // 超过最大重试次数，显示默认图片
    }
}

let img_url = ref([] as string[])

let shot_img = computed(() => {
    return function (index: number) {
        if (storeDisplay.state.viViews.length > 0) {

            if (storeDisplay.state.viViews[index].shot.length == 0) {
                return false
            }
            fetch(storeDisplay.state.viViews[index].shot).then(v => {
                v.blob().then(b => {
                    let url = URL.createObjectURL(b)
                    img_url.value[index] = url

                    return true
                })
            })
            return false
        }
        return false
    }
})

onMounted(() => {
    window.addEventListener("resize", handleResize)
    //第一次初始化需要刷新
    handleResize()
})

onBeforeUnmount(() => {
    window.removeEventListener("resize", handleResize)
})

</script>