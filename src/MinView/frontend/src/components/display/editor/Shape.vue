<template>
    <div ref="maindiv" :style="{
        width: (storeDisplay.state.viViewWidth * (props.outStyle.width * storeDisplay.state.canvas.width / (storeDisplay.state.canvas.width * (storeDisplay.state.viInitScale / 100)))) / storeDisplay.state.canvas.width + 'px',
        height: (storeDisplay.state.viViewHeight * (props.outStyle.height * storeDisplay.state.canvas.height / (storeDisplay.state.canvas.height * (storeDisplay.state.viInitScale / 100)))) / storeDisplay.state.canvas.height + 'px',
        top: (storeDisplay.state.viViewHeight * (props.outStyle.top * storeDisplay.state.canvas.height / (storeDisplay.state.canvas.height * (storeDisplay.state.viInitScale / 100)))) / storeDisplay.state.canvas.height + 'px',
        left: (storeDisplay.state.viViewWidth * (props.outStyle.left * storeDisplay.state.canvas.width / (storeDisplay.state.canvas.width * (storeDisplay.state.viInitScale / 100)))) / storeDisplay.state.canvas.width + 'px',
        // rotate: props.outStyle.rotate + 'deg',
        transform: `rotate(${props.outStyle.rotate}deg)`,
        position: 'absolute',
        backgroundColor: props.outStyle.backgroundColor,
        // display:'flex',
        zIndex: props.outStyle.zIndex,

        borderTopColor: props.outStyle.borderTopColor,
        borderBottomColor: props.outStyle.borderBottomColor,
        borderLeftColor: props.outStyle.borderLeftColor,
        borderRightColor: props.outStyle.borderRightColor,

        borderTopStyle: props.outStyle.borderTopStyle as any,
        borderBottomStyle: props.outStyle.borderBottomStyle as any,
        borderLeftStyle: props.outStyle.borderLeftStyle as any,
        borderRightStyle: props.outStyle.borderRightStyle as any,

        borderTopWidth: props.outStyle.borderTopWidth,
        borderBottomWidth: props.outStyle.borderBottomWidth,
        borderLeftWidth: props.outStyle.borderLeftWidth,
        borderRightWidth: props.outStyle.borderRightWidth,

        borderTopLeftRadius: props.outStyle.borderTopLeftRadius,
        borderTopRightRadius: props.outStyle.borderTopRightRadius,
        borderBottomLeftRadius: props.outStyle.borderBottomLeftRadius,
        borderBottomRightRadius: props.outStyle.borderBottomRightRadius,

        opacity: props.outStyle.opacity,

    }" @mousedown="handleMouseDownOnShape($event)">
        <slot></slot>
    </div>
</template>
<script setup lang="ts">

import { ref, reactive, computed, nextTick, ComputedRef, onMounted, onUnmounted } from "vue"
import { ViChartStyle, ViComponent, ViGroupStyle, ViInnerStyle, ViOutStyle } from "../../../store";
import { storeDisplay } from "../../../store/display";
import { EventsOff, EventsOn } from "../../../../wailsjs/runtime/runtime";

let props = defineProps<{
    outStyle: ViOutStyle;
    innerStyle: ViInnerStyle,
    components: ViComponent[];
    component_id: string;
    nni: string,
    chartStyle: ViChartStyle | undefined,
}>()

let scale = computed(() => {
    let width_scale = storeDisplay.state.viViewWidth / storeDisplay.state.canvas.width
    let height_scale = storeDisplay.state.viViewHeight / storeDisplay.state.canvas.height
    return { width_scale, height_scale }
})

//接收组件数据，用于远程提交修改组件。any应该是数组类数据的string,例如："[选项1,选项2]"
function receiveCompoData(data: any) {
    storeDisplay.commit("update_compo_value", data)
    storeDisplay.commit("emit_component_objects")
}

onMounted(() => {
    EventsOn(props.nni, receiveCompoData)
    props.components.forEach(v => {
        EventsOn("nni" + v.id, receiveCompoData)
    })
    // let newwidth = storeDisplay.state.canvas.width * (storeDisplay.state.viInitScale / 100)
    // let x = props.outStyle.left * storeDisplay.state.canvas.width / newwidth
    // let x1 = storeDisplay.state.viViewWidth * x / storeDisplay.state.canvas.width
    // let newx = (storeDisplay.state.viViewWidth * (props.outStyle.left * storeDisplay.state.canvas.width / (storeDisplay.state.canvas.width * (storeDisplay.state.viInitScale / 100)))) / storeDisplay.state.canvas.width
    // let newy = (storeDisplay.state.viViewHeight * (props.outStyle.top * storeDisplay.state.canvas.height / (storeDisplay.state.canvas.height * (storeDisplay.state.viInitScale / 100)))) / storeDisplay.state.canvas.height
    // console.log("xxx", x, x1, newx, newy, storeDisplay.state.canvas)
})
onUnmounted(() => {
    EventsOff(props.nni)
    props.components.forEach(v => {
        EventsOff("nni" + v.id)
    })
})

function handleMouseDownOnShape(event: MouseEvent) {
    storeDisplay.commit("set_component_id", props.component_id)
}
</script>

<style lang="scss" scoped></style>