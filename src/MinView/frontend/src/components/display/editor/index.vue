<template>
    <div :style="{ height: '100%', width: '100%' }">
        <div class="bg-light-blue-darken-4" style="display:flex; justify-content:center; text-align: center; width: 100%;height: 100%;">
            <div class="bg-light-blue-darken-4" :style="{
                // backgroundColor: 'beige',
                // height: storeDisplay.state.controlHeight - 20 + 'px',
                width: '10%',
                textAlign: 'center',
                alignItems:'center',
            }">
                <Control />
            </div>

            <div id="displayEditorID"
                :style="{ backgroundColor: 'white', width: '90%', height: '100%', textAlign: 'center', display: 'flex', justifyContent: 'center', alignItems: 'center', color: 'white', paddingBottom: '25px' }">
                <div id="displayContextID" @dblclick.stop="setFullScreen" @click.stop="handleClick"
                    @contextmenu="handleClick" v-if="storeDisplay.state.viCtrlIndex >= 0" :style="{
                        width: storeDisplay.state.viViewWidth + 'px',
                        height: storeDisplay.state.viViewHeight + 'px',
                        backgroundColor: background_color,
                        backgroundImage: background_image,
                        backgroundRepeat: 'no-repeat',
                        backgroundSize: 'cover',
                        position: 'relative',

                    }
                        ">
                    <Shape v-for="( item ) in  componet_list " :outStyle="item.style" :inner-style="item.innerStyle"
                        :component_id="item.id" :nni="`nni${item.id}`" :components="item.components" :key="item.id"
                        :chartStyle="item.chartStyle">
                        <component :nni="`nni${item.id}`" :is="item.component" :name="item.name"
                            :innerStyle="item.innerStyle" :components="item.components" :outStyle="item.style"
                            :groupStyle="item.groupStyle" :chartStyle="item.chartStyle" />
                    </Shape>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, computed, onUnmounted } from 'vue';
import Control from '../Control.vue';
import { storeDisplay } from '../../../store/display';
import { nextTick } from 'vue';
import Shape from './Shape.vue';
import { watch } from 'vue';
import html2canvas from 'html2canvas';
import { onUpdated } from 'vue';
import { WindowFullscreen, WindowUnfullscreen } from '../../../../wailsjs/runtime/runtime';

// import { viFormat } from "../../../utils/func";

let componet_list = computed(() => {
    if (storeDisplay.state.viViews.length) {
        return storeDisplay.state.viViews[storeDisplay.state.viCtrlIndex].components
    }
})

let displayHeight = computed(() => {
    //document.body.scrollHeight
    let foot = document.getElementById('displayFoot')
    if (foot == null) {
        return '100%'
    }
    let menu = document.getElementById('displayMenu')
    if (menu == null) {
        return '100%'
    }
    let menu_height = menu.scrollHeight
    let foot_height = foot.scrollHeight
    let display_height = document.body.scrollHeight - menu_height - foot_height
    return display_height
})

function handleWheel(event: any) {
    if (storeDisplay.state.viViews.length > 0) {
        let index = storeDisplay.state.viCtrlIndex
        if (event.deltaY < 0) {
            //向上
            if (index == 0) {

            } else {
                index--
                storeDisplay.commit("set_view_index", index)
            }
        } else {
            if (index + 1 >= storeDisplay.state.viViews.length) {
                //index = 0
            } else {
                index++
                storeDisplay.commit("set_view_index", index)
            }
        }
    }
}

function handleKeydown(event: KeyboardEvent) {
    if (storeDisplay.state.viViews.length > 0) {
        let index = storeDisplay.state.viCtrlIndex
        if (event.key === 'ArrowUp' || event.key == 'ArrowLeft') {
            if (index == 0) {

            } else {
                index--
                storeDisplay.commit("set_view_index", index)
            }
        } else if (event.key === 'ArrowDown' || event.key == 'ArrowRight') {
            if (index + 1 >= storeDisplay.state.viViews.length) {
                //index = 0
            } else {
                index++
                storeDisplay.commit("set_view_index", index)
            }
        }
    }
}

function handleClick(event: MouseEvent) {

    if (storeDisplay.state.viViews.length > 0 && storeDisplay.state.viIsFullScreen) {
        let index = storeDisplay.state.viCtrlIndex
        if (event.button == 0) {
            if (index + 1 >= storeDisplay.state.viViews.length) {
                //index = 0
            } else {
                index++
                storeDisplay.commit("set_view_index", index)
            }
        } else if (event.button == 2) {
            if (index == 0) {

            } else {
                index--
                storeDisplay.commit("set_view_index", index)
            }
        }
    }
}


function setFullScreen() {
    var fullScreenDiv = document.getElementById("displayEditorID")
    var isFullScreen = document.fullscreenElement
    if (!isFullScreen) {
        if (fullScreenDiv) {
            if (fullScreenDiv.requestFullscreen) {
                fullScreenDiv.requestFullscreen();
                WindowFullscreen()
                storeDisplay.commit("set_viIsFullScreen", true)
            }
        }

    } else {
        if (document.exitFullscreen) {
            storeDisplay.commit("set_viIsFullScreen", false)
            document.exitFullscreen();
            WindowUnfullscreen()
        }
    }

}

var timrer: NodeJS.Timeout

function repeatFunction() {
    if (!storeDisplay.state.viIsFullScreen) {
        clearTimeout(timrer)
        return
    }
    if (storeDisplay.state.viViews.length == 1) {
        clearTimeout(timrer)
        return
    }
    if (storeDisplay.state.viPlayType == 'mdiRepeatVariant') {
        //循环播放
        // mdiGestureTap
        let index = storeDisplay.state.viCtrlIndex
        if (index + 1 >= storeDisplay.state.viViews.length) {
            index = 0
        } else {
            index++
        }
        nextTick(() => {
            storeDisplay.commit("set_view_index", index)
        })
        // timrer = setTimeout(repeatFunction, storeDisplay.state.viTimeDelay * 1000)

    } else if (storeDisplay.state.viPlayType == 'mdiShuffleVariant') {
        //随机播放
        let index = storeDisplay.state.viCtrlIndex
        var randomInt = 0;
        // do {
        //     randomInt = Math.floor(Math.random() * storeDisplay.state.viViews.length)

        // } while (randomInt != index);
        for (; ;) {
            randomInt = Math.floor(Math.random() * storeDisplay.state.viViews.length)
            if (randomInt != index) {
                break
            }
        }
        nextTick(() => {
            storeDisplay.commit("set_view_index", randomInt)
        })

    }
    clearTimeout(timrer)
    timrer = setTimeout(repeatFunction, storeDisplay.state.viTimeDelay * 1000)
}

watch(() => storeDisplay.state.viIsFullScreen, (newValue) => {
    if (newValue) {
        timrer = setTimeout(repeatFunction, storeDisplay.state.viTimeDelay * 1000)
    } else {
        clearTimeout(timrer)
    }
    handleResize()
})

function handleFullscreenChange() {
    var fullScreenDiv = document.getElementById("displayEditorID")
    if (fullScreenDiv) {
        if (document.fullscreenElement) {
            fullScreenDiv.style.backgroundColor = 'black';
            handleResize()
        } else {
            fullScreenDiv.style.backgroundColor = 'white';
            handleResize()
        }
    }

}

function asyncCanvasToBlob(canvas: HTMLCanvasElement, index: number) {
    return new Promise((resolve, reject) => {
        canvas.toBlob((blob) => {
            if (blob) {
                let url = URL.createObjectURL(blob)
                storeDisplay.state.viViews[index].shot = url
                resolve(url);
            }
            else {
                reject("")
            }
        }, 'image/png'); // 指定输出格式为PNG
    });
}


function getShot(index: number) {
    if (storeDisplay.state.viViews[index].shot != "") {
        return
    }
    let el = document.getElementById("displayContextID")
    if (el != null) {
        // 将 div 元素转换为 canvas 元素
        html2canvas(el, {
        }).then((canvas) => {
            // 在这里处理 canvas 元素
            const ctx = canvas.getContext('2d');
            if (ctx != null) {
                ctx.scale(0.005, 0.005)
                // let dataUrl = canvas.toDataURL("image/jpeg", 0.3)
                // WriteToCache(dataUrl).then(r => {
                //     State.viViews[State.viCtrlIndex].shot = r
                // })
                // canvas.toBlob(function (blob) {
                //     // 创建URL
                //     if (blob != null) {
                //         const url = URL.createObjectURL(blob);
                //         if (storeDisplay.state.viViews[index].shot != "") {
                //             URL.revokeObjectURL(storeDisplay.state.viViews[index].shot)
                //         }
                //         storeDisplay.state.viViews[index].shot = url
                //     }
                // }, "image/png");
                nextTick(() => {
                    asyncCanvasToBlob(canvas, index).then(v => {
                        console.log(v)
                    })
                })
            }
        });
    }
}
onUpdated(() => {

})

// watch(() => storeDisplay.state.viCtrlIndex, (newValue) => {
//     setTimeout(() => {
//         getShot(newValue)
//     }, 1000);
// })

function backgroundImageError(e: any) {
    // console.log("背景图片出错哦", e)
}

onMounted(() => {
    let el = document.getElementById("displayEditorID")
    if (el != null) {
        storeDisplay.commit("set_display", el)
    }
    window.addEventListener("resize", handleResize)
    document.addEventListener('fullscreenchange', handleFullscreenChange);
    window.addEventListener('keydown', handleKeydown);

    el = document.getElementById("displayContextID")
    if (el != null) {
        el.addEventListener("error", backgroundImageError)
    }
})

onUnmounted(() => {
    window.removeEventListener("resize", handleResize)
    window.removeEventListener('keydown', handleKeydown);
    document.removeEventListener('fullscreenchange', handleFullscreenChange);
    let el = document.getElementById("displayContextID")
    if (el != null) {
        el.removeEventListener("error", backgroundImageError)
    }
})

function handleResize() {
    // nextTick(() => {

    if (storeDisplay.state.display == null) {
        storeDisplay.commit("set_rect", { height: 0, width: 0 })
        return
    }
    let rect = storeDisplay.state.display.getBoundingClientRect()
    let width = storeDisplay.state.canvas.width
    let height = storeDisplay.state.canvas.height

    let width_new = 0
    let height_new = 0
    let src_cale = height / width
    let width_scale = width / rect.width
    let height_scale = height / rect.height
    if (height_scale < width_scale) {
        width_new = rect.width
        height_new = width_new * src_cale
    } else {
        height_new = rect.height
        width_new = height_new / src_cale
    }
    storeDisplay.commit("set_rect", { width: width_new, height: height_new })
    // })
}

let background_color = computed(() => {
    if (storeDisplay.state.viViews.length) {
        return storeDisplay.state.viViews[storeDisplay.state.viCtrlIndex].background_color
    } else {
        return ""
    }
})

let background_image = computed(() => {
    let res = "none"

    if (storeDisplay.state.viViews.length > 0) {
        let url = storeDisplay.state.viViews[storeDisplay.state.viCtrlIndex].background_image
        if (url == "") {
            res = "none"
        } else {
            res = "url(" + url + ")"
        }
    }
    return res
})



</script>