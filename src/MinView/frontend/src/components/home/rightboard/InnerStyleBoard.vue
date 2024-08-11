<template>
    <v-window-item key="board2" :value="2" style="background-color: #01579B; color: green;">
        <!-- {{ current_comptype }} -->
        <div
            v-if="current_comptype == '' || current_comptype == 'ViGroup' || store.state.viViews[store.state.viCtrlIndex].componet_ids.length > 1">
            <div style="color: white;">
                Copyright © 2024 zerosir.com
            </div>
        </div>
        <ChartStyleBoard v-else-if="current_comptype == 'ViChart'" />
        <v-expansion-panels v-else variant="Popout" v-model="selected_item" tile flat>
            <v-expansion-panel v-for="(item, index) in inner_style_option.keys" v-show="item != 'list'"
                :disabled="item == 'list'" :value="item" style="background-color:transparent;color: white;" tile>
                <v-expansion-panel-title collapse-icon="md:keyboard_arrow_up" expand-icon="md:keyboard_arrow_down"
                    :v-show="item != 'list'">
                    <span style="font-size:small;">{{ inner_style_option.infos[index] }}</span>
                </v-expansion-panel-title>
                <v-expansion-panel-text :v-show="item != 'list'">
                    <v-card flat rounded="0" style="background-color: #003355;">
                        <!-- <v-card-subtitle style="font-size: small; color: white;margin-top: 10px;">
                            {{ inner_style_option.infos[index] }}
                        </v-card-subtitle> -->
                        <v-card-actions style="color: white;">
                            <div style="padding-left: 10px;padding-right: 10px;width: 100%;height: 100%;">
                                <v-text-field v-model="value" color="green"
                                    v-if="inner_style_option.types[index] == 'number'" density="compact" elevation="0"
                                    type="number" variant="underlined" hide-details>
                                </v-text-field>
                                <v-text-field v-model="value" @click:append-inner="dialog_color_show_func"
                                    :modes="['hexa', 'rgba', 'hsla']" append-inner-icon="md:color_lens" :color="value"
                                    v-else-if="inner_style_option.types[index] == 'string' && (item.includes('Color') || item.includes('color'))"
                                    variant="underlined" density="compact" elevation="0"></v-text-field>
                                <v-combobox v-model="value" v-else-if="inner_style_option.types[index] == 'boolean'"
                                    variant="underlined" :items="[true, false]" density="compact"
                                    elevation="0"></v-combobox>
                                <v-combobox v-model="value" v-else-if="inner_style_option.keys[index] == 'fontfamily'"
                                    variant="underlined"
                                    :items="['sans-serif', 'serif', 'monospace', 'Arial', 'Courier New', 'Microsoft YaHei']"
                                    density="compact" elevation="0"></v-combobox>
                                <v-combobox v-model="value" v-else-if="inner_style_option.keys[index] == 'fontStyle'"
                                    variant="underlined" :items="['normal', 'italic', 'oblique']" density="compact"
                                    elevation="0"></v-combobox>
                                <v-combobox v-model="value" v-else-if="inner_style_option.keys[index] == 'fontWeight'"
                                    variant="underlined"
                                    :items="['normal', 'bold', 'bolder', 'lighter', '100', '200', '300', '400', '500', '600', '700', '800', '900', '1000']"
                                    density="compact" elevation="0"></v-combobox>
                                <v-text-field v-model="value" color="green" disabled
                                    v-else-if="inner_style_option.keys[index] == 'selectList'" density="compact"
                                    elevation="0" type="text" variant="underlined" hide-details>
                                </v-text-field>
                                <v-text-field v-model="value" color="green" v-else-if="comptype == 'ViImage'"
                                    append-inner-icon="md:image_file" @click:append-inner="file_select('*.jpg;*.png')"
                                    density="compact" elevation="0" type="text" variant="underlined" hide-details>
                                </v-text-field>
                                <v-text-field v-model="value" color="green" v-else-if="comptype == 'ViVideo'"
                                    append-inner-icon="md:video_file" @click:append-inner="file_select('*.webm;*.mp4')"
                                    density="compact" elevation="0" type="text" variant="underlined" hide-details>
                                </v-text-field>
                                <v-text-field v-model="value" color="green" v-else-if="comptype == 'ViAudio'"
                                    append-inner-icon="md:audio_file" @click:append-inner="file_select('*.mp3')"
                                    density="compact" elevation="0" type="text" variant="underlined" hide-details>
                                </v-text-field>
                                <v-text-field v-model="value" color="green" v-else-if="comptype == 'ViFile'"
                                    append-inner-icon="md:attach_file" @click:append-inner="file_select('*.*')"
                                    density="compact" elevation="0" type="text" variant="underlined" hide-details>
                                </v-text-field>
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
import { ViInnerStyle, store } from '../../../store'
import { InnerStyleInfo } from '../../custom/option'
import { cloneDeep, set } from 'lodash';
import { number } from 'mathjs';
import { OpenFileDialog } from '../../../../wailsjs/go/main/App';
import { nanoid } from 'nanoid';
import ChartStyleBoard from './ChartStyleBoard.vue';
import SmallInfo from '../utils/SmallInfo.vue';

let dialog_color_show = ref(false)

let dialog_color_show_func = (() => {
    dialog_color_show.value = true
})

//选择的是哪个属性
let selected_item = ref('')
let inner_style_index = ref(0)
let inner_style_option = computed(() => {
    let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
    let keys: string[] = []
    let values: any[] = []
    let types: string[] = []
    let infos: string[] = []
    if (selected.length == 1) {
        store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
            if (item.id == selected[0]) {
                let style = item.innerStyle
                keys = Object.keys(style)
                values = Object.values(style)
                keys.forEach((v) => {
                    let reftype = typeof style[v as keyof ViInnerStyle];
                    types.push(reftype)
                    // let info = OutStyleInfo[v as keyof typeof OutStyleInfo]
                    let info = InnerStyleInfo(v, item.component)
                    infos.push(info)
                })
                return
            }
        })

    }
    return { keys, values, types, infos }
})

let comptype = computed(() => {
    let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
    let comp = ""
    if (selected.length == 1) {
        store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
            if (item.id == selected[0]) {
                comp = item.component
                return
            }
        })
    }
    return comp
})

let value = computed({
    get: () => {
        let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
        if (selected.length == 1) {
            store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
                if (item.id == selected[0]) {
                    inner_style_index.value = index
                    return
                }
            })

        }
        return store.state.viViews[store.state.viCtrlIndex].components[inner_style_index.value].innerStyle[selected_item.value as keyof ViInnerStyle]
    },
    set: (val) => {
        let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
        if (selected.length == 1) {
            store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
                if (item.id == selected[0]) {
                    inner_style_index.value = index
                    let style = cloneDeep(item.innerStyle)
                    let key1 = selected_item.value as keyof ViInnerStyle
                    let reftype = typeof style[key1 as keyof ViInnerStyle];
                    let newValue1 = val
                    if (reftype == "number") {
                        set(style, key1, number(newValue1))
                    } else if (reftype == "object" && newValue1 != undefined) {
                        set(style, key1, newValue1.toString().split(','))
                    } else {
                        set(style, key1, newValue1)
                    }
                    store.commit("set_inner_style", {
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

let current_comptype = computed(() => {
    let comptype = ""
    let comid = store.state.viViews[store.state.viCtrlIndex].componet_ids[0]
    store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
        if (item.id == comid) {
            comptype = item.name
            if (item.name.startsWith('chart')) {
                comptype = "ViChart"
            }
            return
        }
    })
    return comptype
})

function nniid() {
    return "NIBoard" + nanoid(9)
}

function file_select(filetype: string) {
    let viewid = store.state.viViews[store.state.viCtrlIndex].id
    let comid = store.state.viViews[store.state.viCtrlIndex].componet_ids[0]
    OpenFileDialog(filetype, viewid, comid, true).then((v) => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
            if (item.id == comid) {
                let style = cloneDeep(item.innerStyle)
                if (v.indexOf("Error:") > -1) {
                    return
                }
                style.text = nniid() + v
                store.commit("set_inner_style", {
                    style: reactive(style), addr: {
                        viewid: store.state.viCtrlIndex,
                        comid: comid
                    }
                })
                // nextTick(() => {
                //     store.commit("snapshot")
                // })
                return
            }
        })
    })
}



</script>