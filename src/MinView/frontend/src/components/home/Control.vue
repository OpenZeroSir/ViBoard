<template>
    <v-slide-group v-model="page_index" style="padding: 0px !important; " mandatory draggable center-active show-arrows>

        <v-slide-group-item v-for="n in store.state.viViews.length" :key="n" v-slot="{ isSelected, toggle }">
            <v-card @contextmenu.prevent="showControlMenu($event, n - 1)" :style="{
                border: store.state.viCtrlIndex == n - 1 ? '2px solid #ff9900' : '2px solid #003355',
                backgroundColor: `${store.state.viViews[n - 1].background_color}`,

            }" class="ma-4" height="62" :border="isSelected" :width="min_width" @click="toggle">
                <!-- :style="{
                    backgroundImage:`url(${store.state.viViews[store.state.viCtrlIndex].background_image})`,
                    
                }" -->
                <!-- <div :style="{
                    backgroundSize: 'cover',
                    backgroundColor: 'transparent',
                    backgroundImage: `url(${store.state.viViews[n - 1].background_image})`,
                }"> -->

                <v-img draggable @mousedown="mouseDown($event, n - 1)" @drop="handleDrop($event, n - 1)"
                    @dragover="dragover($event)" @dragstart="dragStart($event)" :src="shot_img(n - 1)" cover>
                </v-img>

                <!-- </div> -->
                <v-badge inline tile flat dot :color="store.state.viViews[n - 1].visible ? 'green' : 'red'">
                </v-badge>
                <v-menu v-model="openControlMenu" :style="{
                    top: `${controlMenuY}px !important`,
                    left: `${controlMenuX}px !important`,
                }" location="center">
                    <v-list>
                        <v-list-item v-for="(item, index) in controlMenuItems" :key="index">
                            <v-list-item-title @click="controlMenuItemClick(index, n - 1)">{{ item
                            }}</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </v-card>
        </v-slide-group-item>
    </v-slide-group>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, reactive, computed, getCurrentInstance } from 'vue'
import { store } from "../../store"



let min_width = computed(() => {
    return store.state.canvas.width * (50 / store.state.canvas.height)
})

let shot_img = computed(() => {
    return function (index: number) {
        if (store.state.viViews.length > 0) {
            return store.state.viViews[index].shot
        }
        return ""
    }
})

let openControlMenu = ref(false)
let controlMenuX = ref(0)
let controlMenuY = ref(0)
let page_index = computed({
    get: () => {
        return store.state.viCtrlIndex
    },
    set: (val) => {
        store.commit("select_View", val)
        nextTick(() => {
            store.commit("snapshot")
        })
    }
});

const controlMenuItems = ['插入视图', '删除视图'] //, '隐藏视图', '显示视图']
function showControlMenu(event: MouseEvent, page_id: number) {
    controlMenuX.value = event.clientX
    controlMenuY.value = event.clientY
    openControlMenu.value = true
    page_index.value = page_id
}
function mouseDown(event: MouseEvent, page_id: number) {
    page_index.value = page_id
}
function controlMenuItemClick(index: number, view_index: number) {
    if (index == 0) {
        store.commit("add_viView", page_index.value)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    } else if (index == 1) {
        store.commit("del_viView_ByIndex", page_index.value)
        nextTick(() => {
            store.commit("snapshot")
        })
    } else {
        let flag = false
        if (index == 3) {
            flag = true
        }
        store.commit("set_viView_visable",
            {
                view_index: page_index.value,
                visible: flag
            })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
    openControlMenu.value = false
}

function dragStart(event: any) {
    event.dataTransfer.setData("Text", page_index.value);
}
function dragover(event: any) {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
}
function handleDrop(event: any, current_index: number) {
    event.preventDefault()
    event.stopPropagation()
    let from_index: string = event.dataTransfer.getData("Text");
    let res = Number(from_index)
    if (isNaN(res)) {
        return
    }
    store.commit("swap_by_index", {
        from: res,
        to: current_index
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })

}


</script>

<style lang="scss" scoped></style>