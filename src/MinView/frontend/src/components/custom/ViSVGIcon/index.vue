
<template>
    <div v-ripple="storeComputed.state.viRipple" :nni="props.nni" @dragover="dragover($event)" @drop="drop($event)" :style="{
        height: '100%',
        width: '100%',

    }">
        <v-icon :style="{
            // height: props.outStyle.height + 'px',
            // width: props.outStyle.width + 'px',
            // // transform: `translate(-50%, -50%) scale(${storeComputed.state.canvas.scale / 100}) `,
            // // position: 'absolute',
            // // left: '50%',
            // // top: '50%',
            // position: 'absolute',
            // left: '50%',
            // top: '50%',
            // transform: 'translate(-50%, -50%)',
            height: '100%',
            width: '100%',
            color: value,

        }">{{ SvgIconName[props.name as keyof typeof SvgIconName] }}</v-icon>
    </div>
</template>


<script lang="ts" setup>
import { computed, ref, reactive, watch, onMounted, nextTick } from 'vue'
import { min } from 'mathjs';
import * as SvgIconName from '@mdi/js'
import { ViComponent, ViInnerStyle, ViOutStyle, store } from '../../../store';
import { getCurrentInstance } from 'vue';
import { storeDisplay } from '../../../store/display';
import { EventsEmit, EventsOff, EventsOn } from '../../../../wailsjs/runtime/runtime';
import { QueryBySheetAndField } from '../../../../wailsjs/go/main/App';
import { onUnmounted } from 'vue';
let props = defineProps<{
    nni: string;
    name: string;
    outStyle: ViOutStyle;
    innerStyle: ViInnerStyle;
}>()
// let size = reactive(computed(() => {
//     let w = Number(props.outStyle.width) * storeComputed.value.state.canvas.scale/100
//     let h = Number(props.outStyle.height) * storeComputed.value.state.canvas.scale/100
//     console.log("storeComputed.value.state.canvas.scale",storeComputed.value.state.canvas.scale)
//     return min(w, h)
// }))

let value = ref("")

watch(() => props.innerStyle.color, (newValue) => {
    if (newValue) {
        value.value = newValue
    }
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
    if (props.innerStyle.list.length > 0) {
        if (props.innerStyle.list[0].sheetkey != data.sheet_key) {
            return
        }
    }
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
    } else if (props.innerStyle.color) {
        value.value = props.innerStyle.color
    }
})

onUnmounted(() => {
    EventsOff("updatecompfromdropfield" + props.nni)
    EventsOff("UIEmitSheetData")
})

</script>

<style lang="scss">
// div[data-v-c9ed1c95] {
//   height: 0px;
// }
</style>