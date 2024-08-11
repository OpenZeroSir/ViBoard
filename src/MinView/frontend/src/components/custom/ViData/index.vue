<script lang="ts" setup>
import { computed, getCurrentInstance, nextTick, onMounted, ref } from 'vue';
import { VDataTable } from 'vuetify/lib/labs/components.mjs';
import { ViChartStyle, ViInnerStyle, ViOutStyle, store } from '../../../store';
import { EventsEmit, EventsOff, EventsOn } from '../../../../wailsjs/runtime/runtime';
import { storeDisplay } from '../../../store/display';
import { QueryBySheetAndField } from '../../../../wailsjs/go/main/App';
import { onUnmounted } from 'vue';

let props = defineProps<{
    nni: string;
    name: string;
    innerStyle: ViInnerStyle;
    outStyle: ViOutStyle;
    chartStyle: ViChartStyle | undefined
}>()

let headers = ref([
    { title: 'Plant', align: 'start', key: 'name' },
    { title: 'Light', align: 'end', key: 'light' },
    { title: 'Height', align: 'end', key: 'height' },
    { title: 'Pet Friendly', align: 'end', key: 'petFriendly' },
    { title: 'Price ($)', align: 'end', key: 'price' },
])
let plants = ref([
    {
        name: 'Fern',
        light: 'Low',
        height: '20cm',
        petFriendly: 'Yes',
        price: 20,
    },
    {
        name: 'Snake Plant',
        light: 'Low',
        height: '50cm',
        petFriendly: 'No',
        price: 35,
    },
    {
        name: 'Monstera',
        light: 'Medium',
        height: '60cm',
        petFriendly: 'No',
        price: 50,
    },
    {
        name: 'Pothos',
        light: 'Low to medium',
        height: '40cm',
        petFriendly: 'Yes',
        price: 25,
    },
    {
        name: 'ZZ Plant',
        light: 'Low to medium',
        height: '90cm',
        petFriendly: 'Yes',
        price: 30,
    },
])

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
    // var fullScreenDiv = document.getElementById("displayEditorID")
    // if (fullScreenDiv == null) {
    //     return
    // }

    if (props.innerStyle.list.length <= 0) {
        return
    }
    let fields: string[] = []
    props.innerStyle.list.forEach(v => {
        fields.push(v.fieldkey)
    })
    await QueryBySheetAndField(props.innerStyle.list[0].sheetkey, fields, display).then(v => {
        if (v.error != "") {
            console.log(v.error)
            return
        }
        headers.value.length = 0
        fields.forEach((field, index) => {
            headers.value.push({
                title: field,
                key: field,
                align: (index == 0 ? 'start' : 'end')
            })
        })
        plants.value.length = 0
        v.data.forEach((row) => {
            let res: any = {}
            let temp = row as any[]
            temp.forEach((d, index) => {
                res[fields[index]] = d
            })
            plants.value.push(res)
        })
    }).catch(reason => {
        console.log(reason)
    })
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
    UIEmitSheetDataFunction(props.outStyle.displayMode)

}




async function UIEmitSheetDataCallback(data: any) {
    if(data === undefined){
        return
    }
    if (props.innerStyle.list.length > 0) {
        let item = props.innerStyle.list[0]
        let haskey = data[item.sheetkey]
        if (haskey) {
            if (props.outStyle.displayMode) {
                await UIEmitSheetDataFunction(props.outStyle.displayMode)
            }
        }
    }
}



onMounted(async () => {
    EventsOn("updatecompfromdropfield" + props.nni, updatecompfromdropfield)
    EventsOn("UIEmitSheetData", UIEmitSheetDataCallback)
    if (props.innerStyle.list.length > 0) {
        await UIEmitSheetDataFunction(props.outStyle.displayMode)
    }

})

onUnmounted(() => {
    EventsOff("updatecompfromdropfield" + props.nni)
    EventsOff("UIEmitSheetData")
})

let fontsize = computed(() => {
    if (props.innerStyle.fontSize) {
        return (props.innerStyle.fontSize * storeComputed.value.state.canvas.scale / 100) + 'px'
    }
    return (storeComputed.value.state.canvas.scale / 100) + 'px'

})


</script>

<template>
    <div v-ripple="storeComputed.state.viRipple" :style="{
        height: '100%',
        width: '100%',
    }" :nni="props.nni" @dragover="dragover($event)" @drop="drop($event)">
        <v-data-table multi-sort density="comfortable" items-per-page="32" :style="{
            backgroundColor: 'transparent', fontSize: fontsize,
            fontFamily: props.innerStyle.fontfamily,
            fontWeight: props.innerStyle.fontWeight,
            fontStyle: props.innerStyle.fontStyle,
            position: 'absolute',
            // left: '50%',
            top: '50%',
            transform: `translate(0, -50%) scaleY(${storeComputed.state.canvas.scale / 100}) `,
        }" :headers="headers" :items="plants" item-key="name">
            <template v-slot:item="{ item, index }">
                <tr>
                    <td v-for="(value, key) in headers" :key="key" :style="{
                        borderWidth: '0px !important',
                        fontFamily: props.innerStyle.fontfamily,
                        fontWeight: props.innerStyle.fontWeight,
                        fontStyle: props.innerStyle.fontStyle,
                        textAlign: key == 0 ? 'left' : 'right',
                        color: props.innerStyle.color,
                        backgroundColor: index % 2 !== 0 ? '#ffffff11 !important' : 'transparent !important'
                    }">
                        {{ item.columns[value.key] }}</td>
                </tr>
            </template>
        </v-data-table>
    </div>
</template>


<style>
/* td {  */
/* background-color: transparent !important;  */
/* color: yellow !important; */
/* } */

th {
    background-color: #ffffff11 !important;
    color: white !important;
    font-weight: bolder !important;
}

.v-data-table-footer {
    display: none !important;
}


.even-row {
    color: red;
    /* 偶数行的背景颜色 */
}



.odd-row {
    color: white;
    /* 奇数行的背景颜色 */

}
</style>