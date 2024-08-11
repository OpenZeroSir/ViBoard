<template>
    <!-- <div v-ripple="storeComputed.state.viRipple" :nni="props.nni"> -->
    <div v-ripple="storeComputed.state.viRipple" :nni="props.nni" @dragover="dragover($event)" @drop="drop($event)" :style="{
        height: '100%',
        width: '100%',
        // height: props.outStyle.height * storeComputed.state.canvas.scale / 100 + 'px',
        // width: props.outStyle.width * storeComputed.state.canvas.scale / 100 + 'px',
    }">
        <div :style="{
            height: '100%',
            width: '100%',
            transform: `translate(-50%, -50%) scale(${storeComputed.state.canvas.scale / 100}) `,
            position: 'absolute',
            left: '50%',
            top: '50%',

        }">
            <span :style="{
                fontFamily: props.innerStyle.fontfamily,
                fontSize: props.innerStyle.fontSize + 'px',
                fontWeight: props.innerStyle.fontWeight,
                fontStyle: props.innerStyle.fontStyle,
                position: 'absolute',
                left: '50%',
                top: '50%',
                transform: 'translate(-50%, -50%)',
                writingMode: props.innerStyle.vertical ? 'vertical-rl' : 'horizontal-tb',
                textOrientation: 'mixed',
                color: props.innerStyle.color,
                // whiteSpace:'normal',
                // overflow:'hidden',
                whiteSpace: 'nowrap',
            }">{{ value }}</span>

        </div>
    </div>
    <!-- </div> -->
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, computed, reactive, ref, watch, watchEffect, nextTick } from 'vue';
import { EventsOn, EventsOff, EventsEmit } from '../../../../wailsjs/runtime/runtime';
import { ViChartStyle, ViInnerStyle, ViOutStyle, store } from '../../../store';
import { getCurrentInstance } from 'vue';
import { storeDisplay } from '../../../store/display';
import { QueryBySheetAndField, RequestBaiduAPI } from '../../../../wailsjs/go/main/App';
let props = defineProps<{
    nni: string;
    name: string;
    innerStyle: ViInnerStyle;
    outStyle: ViOutStyle;
    chartStyle: ViChartStyle | undefined
}>()


let value = ref("")

watch(() => props.innerStyle.text, (newValue) => {
    if (newValue) {
        value.value = newValue
    }
})

const instance = getCurrentInstance();

// let storemode = computed(() => {
//     let fullpath = instance?.proxy?.$router.currentRoute.value.fullPath
//     if (fullpath == '/') {
//         return store.state.viRipple
//     } else {
//         return storeDisplay.state.viRipple
//     }
// })

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

async function UIEmitSheetDataFunction(display: boolean) {
    if (props.innerStyle.list.length > 0) {
        await QueryBySheetAndField(props.innerStyle.list[0].sheetkey, [props.innerStyle.list[0].fieldkey], display).then(v => {
            if (v.error != "") {
                console.log(v.error)
                return
            }
            value.value = v.data[0][0]
        }).catch(reason => {
            console.log(reason)
        })
    }
}

async function updatecompfromdropfield(data: any) {

    store.commit("set_innerStyle_list", {
        id: props.nni,
        sheetkey: data.sheet_key,
        fieldkey: data.field_key,
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
    await UIEmitSheetDataFunction(props.outStyle.displayMode)

}

async function UIEmitSheetDataCallback(data: any) {
    if(data === undefined){
        return
    }
    if (props.innerStyle.list.length > 0) {
        let item = props.innerStyle.list[0]
        let haskey = data[item.sheetkey]
        if (haskey) {
            await UIEmitSheetDataFunction(props.outStyle.displayMode)
        }
    }
}

onMounted(async () => {
    EventsOn("updatecompfromdropfield" + props.nni, updatecompfromdropfield)
    EventsOn("UIEmitSheetData", UIEmitSheetDataCallback)
    if (props.innerStyle.list.length > 0) {
        await UIEmitSheetDataFunction(props.outStyle.displayMode)
    } else if (props.innerStyle.text) {
        value.value = props.innerStyle.text
    }
})

function dragover(event: any) {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'copy'
}

function drop(event: any) {
    let data: string = event.dataTransfer.getData("application/json");
    if (data.length < 1) {
        return
    }
    let json = JSON.parse(data)
    if (json.type == "data") {
        event.preventDefault()
        event.stopPropagation()
        EventsEmit("updatecompfromdropfield" + props.nni, json)
    }

}

onUnmounted(() => {
    EventsOff("updatecompfromdropfield" + props.nni)
    EventsOff("UIEmitSheetData")
})


</script>