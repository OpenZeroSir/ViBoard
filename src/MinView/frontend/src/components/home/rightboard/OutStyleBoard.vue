<template>
    <v-window-item key="board1" :value="1" style="background-color: #01579B; color: green;">
        <div v-if="store.state.viViews[store.state.viCtrlIndex].componet_ids.length > 1">
            <div style="color: white;">
                Copyright © 2024 zerosir.com
            </div>
        </div>
        <v-expansion-panels v-else variant="Popout" v-model="selected_item" tile flat>
            <v-expansion-panel v-for="(item, index) in out_style_option.keys" :value="item"
                v-show="out_style_option.keys[index] !== 'displayMode'" style="background-color:transparent;color: white;"
                tile>
                <v-expansion-panel-title collapse-icon="md:keyboard_arrow_up" expand-icon="md:keyboard_arrow_down">
                    <span style="font-size:small">{{ out_style_option.infos[index] }}</span>
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                    <v-card flat rounded="0" style="background-color: #003355;">
                        <!-- <v-card-subtitle style="font-size: small; color: white;margin-top: 10px;">
                            {{ out_style_option.infos[index] }}
                        </v-card-subtitle> -->
                        <!-- {{ out_style_option[1][index] }} {{ out_style_option[2][index] }} -->
                        <v-card-actions style="color: white;">
                            <div style="padding-left: 10px;padding-right: 10px;width: 100%;height: 100%;">
                                <v-text-field v-model="value" color="green"
                                    v-if="out_style_option.keys[index] === 'opacity'" density="compact" elevation="0"
                                    type="number" variant="underlined" min="0" max="1" step="0.1" hide-details>
                                </v-text-field>
                                <v-text-field v-model="value" color="green"
                                    v-else-if="out_style_option.types[index] == 'number'" density="compact" elevation="0"
                                    type="number" variant="underlined" hide-details>
                                </v-text-field>
                                <v-combobox v-model="value"
                                    v-else-if="out_style_option.types[index] == 'string' && item.includes('Style')"
                                    variant="underlined" :items="border_style" density="compact" elevation="0"></v-combobox>
                                <v-text-field v-model="value" @click:append-inner="dialog_color_show_func"
                                    :modes="['hexa', 'rgba', 'hsla']" append-inner-icon="md:color_lens" :color="value"
                                    v-else-if="out_style_option.types[index] == 'string' && (item.indexOf('Color') > -1 || item.indexOf('color') > -1)"
                                    variant="underlined" density="compact" elevation="0"></v-text-field>
                                <v-btn v-else-if="out_style_option.types[index] == 'object'" color="green" density="compact"
                                    elevation="0" variant="plain">
                                    进入配置>>
                                </v-btn>
                                <v-text-field v-else v-model="value" flat rounded="false" color="green" type="text"
                                    variant="underlined" density="compact" elevation="0" hide-details></v-text-field>
                            </div>
                        </v-card-actions>
                    </v-card>
                </v-expansion-panel-text>

            </v-expansion-panel>
        </v-expansion-panels>
        <v-dialog v-model="dialog_color_show" width="auto">
            <v-color-picker v-model="value" show-swatches></v-color-picker>
        </v-dialog>
    </v-window-item>
</template>

<script setup lang="ts">
import { computed, nextTick, reactive, ref, watch } from 'vue';
import { ViOutStyle, store } from '../../../store'
import { OutStyleInfo } from '../../custom/option'
import { cloneDeep, set } from 'lodash';
import { number } from 'mathjs';
import { onMounted } from 'vue';
import SmallInfo from '../utils/SmallInfo.vue';

let dialog_color_show = ref(false)

let dialog_color_show_func = (() => {
    dialog_color_show.value = true
})

//选择的是哪个属性
let selected_item = ref('')
let out_style_index = ref(0)
let out_style_option = computed(() => {
    let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
    let keys: string[] = []
    let values: any[] = []
    let types: string[] = []
    let infos: string[] = []
    if (selected.length == 1) {
        store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
            if (item.id == selected[0]) {
                let style = item.style
                keys = Object.keys(style)
                values = Object.values(style)
                keys.forEach((v) => {
                    let reftype = typeof style[v as keyof ViOutStyle];
                    types.push(reftype)
                    let info = OutStyleInfo[v as keyof typeof OutStyleInfo]
                    infos.push(info)
                })
                return
            }
        })

    }
    return { keys, values, types, infos }
})

let border_style = ['none', 'dotted', 'dashed', 'solid', 'double']

// let need_update = 0
// watch(selected_item, (newValue) => {
//     need_update = 0
//     let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
//     if (selected.length == 1) {
//         store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
//             if (item.id == selected[0]) {
//                 let style = item.style
//                 let value = style[newValue as keyof ViOutStyle]
//                 return
//             }
//         })

//     }
// })



let value = computed({
    get: () => {
        let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
        if (selected.length == 1) {
            store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
                if (item.id == selected[0]) {
                    out_style_index.value = index
                    return
                }
            })

        }
        return store.state.viViews[store.state.viCtrlIndex].components[out_style_index.value].style[selected_item.value as keyof ViOutStyle]
    },
    set: (val) => {
        let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
        if (selected.length == 1) {
            store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
                if (item.id == selected[0]) {
                    out_style_index.value = index
                    let style = cloneDeep(item.style)
                    let key1 = selected_item.value as keyof ViOutStyle
                    let reftype = typeof style[key1 as keyof ViOutStyle];
                    let newValue1 = val
                    if (reftype == "number") {
                        newValue1 = number(newValue1)
                    }
                    set(style, key1, newValue1)
                    store.commit("set_out_style", {
                        style: reactive(cloneDeep(style)), addr: {
                            viewid: store.state.viCtrlIndex,
                            comid: store.state.viViews[store.state.viCtrlIndex].componet_ids[0]
                        }
                    })
                    // nextTick(() => {
                    //     store.commit("snapshot")
                    // })
                    return
                }
            })

        }


    }
});

onMounted(() => {

})

</script>