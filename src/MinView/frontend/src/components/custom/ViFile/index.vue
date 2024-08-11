<template>
    <div v-ripple="storeComputed.state.viRipple" :nni="props.nni" :style="{
        height: '100%',
        width: '100%',
    }" @dblclick.stop="handelDbClick">
        <!-- <div :style="{
            height: props.outStyle.height + 'px',
            width: props.outStyle.width + 'px',
            transform: `translate(-50%, -50%) scale(${store.state.canvas.scale / 100}) `,
            position: 'absolute',
            left: '50%',
            top: '50%',
        }"> -->
        <span :nni="props.nni" :style="{
            fontFamily: props.innerStyle.fontfamily,
            fontSize: fontsize,
            fontWeight: props.innerStyle.fontWeight,
            fontStyle:props.innerStyle.fontStyle,
            position: 'absolute',
            left: '0%',
            top: '50%',
            transform: 'translate(0, -50%)',
            color: props.innerStyle.color,
            whiteSpace: 'nowrap',  
        }">{{ filename }}</span>

        <!-- </div> -->
    </div>
</template>

<script lang="ts" setup>

import { onMounted, computed, reactive, ref, nextTick } from 'vue'
import { ViChartStyle, ViInnerStyle, ViOutStyle, store } from '../../../store';
import * as SvgIconName from '@mdi/js'
import { min } from 'mathjs';
import { OpenFile } from '../../../../wailsjs/go/main/App';
import { getCurrentInstance } from 'vue';
import { storeDisplay } from '../../../store/display';

let props = defineProps<{
    nni: string;
    name: string;
    innerStyle: ViInnerStyle;
    outStyle: ViOutStyle;
    chartStyle: ViChartStyle | undefined,
}>()

let filename = computed(() => {
    let names = props.innerStyle.text!.split('.')
    return names[names.length - 2] + "." + names[names.length - 1]
})

function handelDbClick() {
    let viewid: string
    let comid: string
    if (storeComputed.value.state.viRipple) {
        viewid = store.state.viViews[store.state.viCtrlIndex].id
        comid = store.state.viViews[store.state.viCtrlIndex].componet_ids[0]
    } else {
        viewid = storeDisplay.state.viViews[store.state.viCtrlIndex].id
        comid = storeDisplay.state.viViews[store.state.viCtrlIndex].componet_ids[0]
    }
    OpenFile(viewid, comid).then(v => {
        if (v != "") {
            storeComputed.value.commit("show_dialog", { show: true, msg: v })
        }
    })
}

// const instance = getCurrentInstance();

let storeComputed = computed(() => {
    // let fullpath = instance?.proxy?.$router.currentRoute.value.fullPath
    if (!props.outStyle.displayMode) {
        return store
    } else {
        return storeDisplay
    }
})


let fontsize = computed(() => {
    if (props.innerStyle.fontSize) {
        return (props.innerStyle.fontSize * storeComputed.value.state.canvas.scale / 100) + 'px'
    }
    return (storeComputed.value.state.canvas.scale / 100) + 'px'

})

</script>

<style lang="scss" scoped>
.filecontainer {
    display: flex;
    justify-content: center;
    width: 100%;
    height: 100%;
}

.lefticon {
    /* 可选：设置左侧内容的宽度 */
    width: 30%;
    height: 100%;
    align-items: center;
    justify-content: center;
    display: flex;
}

.rightname {
    /* 可选：设置右侧内容的宽度 */
    width: 70%;
    height: 100%;
    align-items: center;
    justify-content: left;
    display: flex;
}
</style>
