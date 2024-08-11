<template>
    <v-dialog v-model="store.state.ViShowFunction" width="600px">
        <v-card>
            <v-card-text>
                <v-textarea v-model="value" label="Javascript Function"></v-textarea>
            </v-card-text>
            <v-card-actions>
                <div style="display: flex;">
                    <v-btn color="green" @click="onok">确定</v-btn>
                    <v-btn color="red" @click="oncancel">取消</v-btn>
                </div>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>
  
<script lang="ts" setup>
import { WritableComputedRef, computed, ref } from 'vue';
import { ViComponent, store } from '../../../store/index';
import { nextTick, onBeforeMount, watch } from 'vue';
import { onMounted } from 'vue';
import { getElement, setup_option, sleep } from '../../../utils/func';

import * as echarts from 'echarts'

let value = ref('')

let chart_comp: WritableComputedRef<ViComponent | undefined> = computed({
    get: () => {
        let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
        let name_list: ViComponent | undefined = undefined
        if (selected.length == 1) {
            store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
                if (item.id == selected[0]) {
                    name_list = item
                    return
                }
            })

        }
        return name_list
    },
    set: (val) => {

        // store.commit("set_chart_style_other_list", { value: val })

    }
});

let store_current_value = computed({
    get: () => {
        return store.state.viChartCurrentValue
    },
    set: (val) => {
        store.commit("set_viChartCurrentValue", { value: val })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

function updateOption(function_name: string) {
    let el = getElement(chart_comp.value?.id + "nni") as HTMLElement
    if (el != null) {
        let chart = echarts.getInstanceByDom(el)
        let option = chart_comp.value?.chartStyle?.style
        if (option != undefined) {
            try {
                chart?.clear()
                chart?.setOption(option)
                // console.log(function_name, 1, option)
            } catch (error) {
                // console.log(function_name, 2, error)
                chart?.dispose()
                chart = echarts.init(el)
            }
        }
    }
}

function onok() {
    setup_option(value.value, store.state.ViFunctionPath, undefined)
    nextTick(() => {
        store.commit("snapshot")
    })
    sleep(200).then(() => {
        nextTick(() => {
            updateOption("current_value_computed")
        })
    })
    store_current_value.value = value.value
    store.commit("set_ViShowFunction", { show: false, msg: "", path: "" })
}

function oncancel() {
    store.commit("set_ViShowFunction", { show: false, msg: "", path: "" })
}

watch(() => store.state.ViShowFunction, async (newValue) => {
    if (newValue) {
        value.value = store.state.ViFunctionContext
    }
})

</script>
