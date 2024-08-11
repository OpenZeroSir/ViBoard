<template>
    <div id="editor" class="editor" :style="{
        height: `${viFormat(store.state.canvas.height, store.state.canvas.scale)}px`,
        width: `${viFormat(store.state.canvas.width, store.state.canvas.scale)}px`,
        backgroundColor: background_color,
        backgroundImage: background_image,
        backgroundRepeat: 'no-repeat',

        // backgroundPosition: 'center center',
        backgroundSize: 'cover',
        minWidth: `${viFormat(store.state.canvas.width, store.state.canvas.scale)}px`,
        color: 'white',
        position: 'relative',
        borderWidth: '1px 1px 1px 1px',
        borderColor: '#70c0ff',
        borderStyle: 'solid',
    }" @mousedown="handleMouseDown" @contextmenu.prevent.capture="showEditorMenu($event)">
        <!-- 网格线 -->
        <!-- <ViGrid /> -->
        <Shape v-for="(item) in componet_list" :outStyle="item.style" :inner-style="item.innerStyle"
            :active="activeShape(item.id)" :component_id="item.id" :components="item.components"
            :chartStyle="item.chartStyle">
            <component :nni="`nni${item.id}`" :is="item.component" :name="item.name" :key="item.id"
                :innerStyle="item.innerStyle" :components="item.components" :outStyle="item.style"
                :groupStyle="item.groupStyle" :chartStyle="item.chartStyle" />
        </Shape>
        <Area v-show="isShowArea" :start="area_start" :width="area_width" :height="area_height" />
        <v-menu v-model="isShowEditorMenu" :style="{
            top: `${editorMenuY}px !important`,
            left: `${editorMenuX}px !important`,
        }">
            <v-list>
                <v-list-item v-for="item in editorMenuItems" :disabled="!item.show" :key="item.title">
                    <v-list-item-title v-text="item.title"
                        @click="editorMenuItemClick($event, item.title)"></v-list-item-title>
                </v-list-item>
            </v-list>
        </v-menu>
    </div>
</template>

<script lang="ts" setup>

import ViGrid from './ViGrid.vue'
import Shape from './Shape.vue';
import Area from './Area.vue';
import { store, ViComponent } from "../../../store"
import { viFormat, getComponentRotatedStyle, nanoidEx } from '../../../utils/func';

import { computed, ref, toRaw, reactive, nextTick, defineComponent } from 'vue';
import { nanoid } from 'nanoid';
import { cloneDeep } from 'lodash'

let editorMenuItems = reactive([
    { title: '复制组件', show: computed(() => store.state.viViews[store.state.viCtrlIndex].componet_ids.length > 0) },
    { title: '粘贴组件', show: computed(() => store.state.viCopy.length > 0) },
    { title: '剪切组件', show: computed(() => store.state.viViews[store.state.viCtrlIndex].componet_ids.length > 0) },
    { title: '删除组件', show: computed(() => store.state.viViews[store.state.viCtrlIndex].componet_ids.length > 0) },
])
let isShowEditorMenu = ref(false)
let editorMenuX = ref(0)
let editorMenuY = ref(0)

let isShowArea = ref(false)
let area_width = ref(0)
let area_height = ref(0)
let area_start = reactive({ x: 0, y: 0 })

let componet_list = computed(() => {
    if (store.state.viViews.length) {
        return store.state.viViews[store.state.viCtrlIndex].components
    }
})
let background_color = computed(() => {
    let color = "#003355"
    if (store.state.viViews.length) {
        color = store.state.viViews[store.state.viCtrlIndex].background_color
    }
    return color
})
let background_image = computed(() => {
    // let color="url(src/assets/images/logo-universal.png)"
    let color = "none"

    if (store.state.viViews.length > 0) {
        let img = store.state.viViews[store.state.viCtrlIndex].background_image
        if (img == "") {
            color = "none"
        } else {
            color = "url(" + img + ")"
        }
    }
    return color
})
function activeShape(shapeId: string) {
    let result = false
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(v => {
        if (v == shapeId) {
            result = true
            return
        }
    })
    return result
}

function showEditorMenu(e: MouseEvent) {
    e.preventDefault()
    if (e.buttons == 2) {
        e.stopPropagation()
    }
    editorMenuX.value = e.clientX
    editorMenuY.value = e.clientY
    isShowEditorMenu.value = true
}

function copyComp() {
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach(comp => {
            if (id == comp.id) {
                let new_comp = cloneDeep(comp)
                store.commit("copy_to", new_comp)
            }
        })
    })
}

function pasteComp(x: number, y: number) {
    let rect = store.state.viEditor!.getBoundingClientRect()
    store.state.viCopy.forEach((comp, index) => {
        let c = cloneDeep(comp)
        c.id = nanoidEx()
        c.style.top = y - rect.top + index * 10
        c.style.left = x - rect.left + index * 10
        store.commit("add_component", reactive(c))
    })
}

function deleteComp() {
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.commit("del_component", id)
    })
}

function editorMenuItemClick(e: MouseEvent, title: string) {
    e.preventDefault()
    e.stopPropagation()
    if (title == "复制组件") {
        store.commit("clear_copy")
        copyComp()
    } else if (title == "粘贴组件") {
        pasteComp(e.clientX, e.clientY)
        store.commit("snapshot")
    } else if (title == "剪切组件") {
        store.commit("clear_copy")
        copyComp()
        deleteComp()
        store.commit("snapshot")
    } else if (title == "删除组件") {
        deleteComp()
        store.commit("snapshot")
    }
    
    isShowEditorMenu.value = false
}

function hideArea() {
    isShowArea.value = false
    area_width.value = 0
    area_height.value = 0
}
function selectAreaComponent() {
    const result: string[] = []
    // 区域起点坐标
    const { x, y } = area_start
    // 计算所有的组件数据，判断是否在选中区域内
    store.state.viViews[store.state.viCtrlIndex].components.forEach(c => {
        const { left, top, width, height } = getComponentRotatedStyle({
            rotate: c.style.rotate,
            left: c.style.left,
            top: c.style.top,
            width: c.style.width,
            height: c.style.height,
            right: 0,
            bottom: 0,
        })
        if (x <= left && y <= top && (left + width <= x + area_width.value) && (top + height <= y + area_height.value)) {
            result.push(c.id)
        }
    })
    store.commit("set_component_id_by_area", result)
}

function handleMouseDown(e: MouseEvent) {
    e.preventDefault()
    hideArea()
    if (e.buttons == 2) {
        return
    }
    // 获取编辑器的位移信息，每次点击时都需要获取一次。主要是为了方便开发时调试用。
    const rectInfo = store.state.viEditor!.getBoundingClientRect()
    let editorX = rectInfo.x
    let editorY = rectInfo.y

    const startX = e.clientX
    const startY = e.clientY
    area_start.x = startX - editorX
    area_start.y = startY - editorY
    // 展示选中区域
    isShowArea.value = true
    let isFirst = true
    const move = (moveEvent: MouseEvent) => {
        if (isFirst) {
            isFirst = false
            return
        }
        area_width.value = Math.abs(moveEvent.clientX - startX)
        area_height.value = Math.abs(moveEvent.clientY - startY)
        if (moveEvent.clientX < startX) {
            area_start.x = moveEvent.clientX - editorX
        }

        if (moveEvent.clientY < startY) {
            area_start.y = moveEvent.clientY - editorY
        }
    }

    const up = (event: MouseEvent) => {
        document.removeEventListener('mousemove', move)
        document.removeEventListener('mouseup', up)
        selectAreaComponent()
        hideArea()
    }

    document.addEventListener('mousemove', move)
    document.addEventListener('mouseup', up)

}

</script>

<style lang="scss">
.aa {
    position: relative;
}

.editor-active-class {
    background-color: rgba(0, 0, 255, 0.453) !important;
    color: white !important;
}
</style>