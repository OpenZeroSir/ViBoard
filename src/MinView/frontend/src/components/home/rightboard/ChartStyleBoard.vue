<template>
    <div v-show="store.state.viBreadindex == 0" style=" height: 100%;  background-color: #003355;">
        <div style="margin: 10px;background-color: #003355; ">

            <v-row style="padding-left: 8px;padding-right: 24px;">
                <v-col cols="10">
                    <span style="font-weight: bold;color: green;">
                        坐标系统
                    </span>
                </v-col>
                <v-col cols="2">
                    <v-menu location="center">
                        <template v-slot:activator="{ props }">
                            <!-- <v-btn icon v-bind="props"> -->
                            <v-icon center color="green" size="large" v-bind="props">{{ mdiPlus }}</v-icon>
                            <!-- </v-btn> -->
                        </template>

                        <v-list>
                            <v-list-item v-for="(item, index) in coor" :key="index" @click="coorMenuItemClick(item, index)">
                                <v-list-item-title>{{ item }}</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-col>
            </v-row>
            <v-list>
                <v-list-item v-for="(item, index) in coor_selected" @dragover="dragover($event)"
                    @dragleave="dragleave($event)" @drop="coorHandleDrop($event, item, index)" :key="item.key"
                    :value="item.name" @click.stop="setupiscontrolcoor(item.key, index)"
                    style="color:white;padding-top: 0;padding-bottom: 0;border: 0;" density="compact" elevation="0" nav tile
                    flat>
                    <template v-slot:append>
                        <v-icon color="red" @click.stop="coorDeleteClick(item.key, index)">{{ mdiClose }}</v-icon>
                    </template>
                    <div v-ripple class="draggable" draggable="true" style="white-space: nowrap; /* 禁止换行 */
                                                overflow: hidden; /* 超出部分隐藏 */
                                                text-overflow: ellipsis; /* 超出部分显示省略号 */">
                        {{ item.name }}
                    </div>
                    <!-- <template v-slot:prepend>
                        <v-icon color="green">{{ mdiPencilPlusOutline }}</v-icon>
                    </template> -->

                </v-list-item>
            </v-list>
        </div>
        <div style="background-color: #01579B; width: 100%;height: 10px;"></div>
        <div style="margin: 10px;margin-top: 15px; background-color: #003355;">
            <v-row style="padding-left: 8px;padding-right: 24px;">
                <v-col cols="10">
                    <span style="font-weight: bold;color: green;">
                        数据系列
                    </span>
                </v-col>
                <v-col cols="2">
                    <v-menu location="center" height="700px" max-height="700px">
                        <template v-slot:activator="{ props }">
                            <!-- <v-btn icon v-bind="props"> -->
                            <v-icon center color="green" size="large" v-bind="props">{{ mdiPlus }}</v-icon>
                            <!-- </v-btn> -->
                        </template>

                        <v-list>
                            <v-list-item v-for="(item, index) in data" :key="index" @click="dataMenuItemClick(item, index)">
                                <v-list-item-title>{{ item }}</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-col>
            </v-row>
            <v-list>
                <v-list-item v-for="(item, index) in data_selected" :key="item.key" :value="item.name"
                    @dragover="dragover($event)" @drop="dataHandleDrop($event, item, index)" @dragleave="dragleave($event)"
                    @click.stop="setupiscontrolseries(item.key, index)"
                    style="color:white;padding-top: 0;padding-bottom: 0;border: 0;" density="compact" elevation="0" nav tile
                    flat>
                    <template v-slot:append>
                        <v-icon color="red" @click.stop="dataDeleteClick(item.name, index)">{{ mdiClose }}</v-icon>
                    </template>
                    <div v-ripple class="draggable" draggable="true" style="white-space: nowrap; /* 禁止换行 */
                                                overflow: hidden; /* 超出部分隐藏 */
                                                text-overflow: ellipsis; /* 超出部分显示省略号 */">
                        {{ item.name }}
                    </div>
                    <!-- <template v-slot:prepend>
                        <v-icon color="red">{{ mdiDragVariant }}</v-icon>
                    </template> -->

                </v-list-item>
            </v-list>
        </div>
        <div style="background-color: #01579B; width: 100%;height: 10px;"></div>
        <div style="margin: 10px;margin-top: 15px; background-color: #003355;">
            <v-row style="padding-left: 8px;padding-right: 24px;">
                <v-col cols="10">
                    <span style="font-weight: bold;color: green;">
                        辅助属性
                    </span>
                </v-col>
                <v-col cols="2">
                    <v-menu location="center">
                        <template v-slot:activator="{ props }">
                            <!-- <v-btn icon v-bind="props"> -->
                            <v-icon center color="green" size="large" v-bind="props">{{ mdiPlus }}</v-icon>
                            <!-- </v-btn> -->
                        </template>

                        <v-list>
                            <v-list-item v-for="(item, index) in other" :key="index"
                                @click="otherMenuItemClick(item, index)">
                                <v-list-item-title>{{ item }}</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-col>
            </v-row>
            <v-list>
                <v-list-item v-for="(item, index) in other_selected" :key="item.key" :value="item.name"
                    @click.stop="setupiscontrolother(item.key, index)"
                    style="color:white;padding-top: 0;padding-bottom: 0;border: 0;" density="compact" elevation="0" nav tile
                    flat>
                    <template v-slot:append>
                        <v-icon color="red" @click.stop="otherDeleteClick(item.key, index)">{{ mdiClose }}</v-icon>
                    </template>
                    <div v-ripple class="draggable" draggable="true" style="white-space: nowrap; /* 禁止换行 */
                                                overflow: hidden; /* 超出部分隐藏 */
                                                text-overflow: ellipsis; /* 超出部分显示省略号 */">
                        {{ item.name }}
                    </div>
                    <!-- <template v-slot:prepend>
                        <v-icon color="green">{{ mdiDragVariant }}</v-icon>
                    </template> -->

                </v-list-item>
            </v-list>
        </div>

        <!-- <div style="background-color: #01579B; width: 100%;height: 10px;"></div>
        <div style="margin: 10px;margin-top: 15px; background-color: #003355;">
            <v-row style="padding-left: 8px;padding-right: 24px;">
                <v-col cols="10">
                    <span style="font-weight: bold;color: green;">
                        条件过滤
                    </span>
                </v-col>
                <v-col cols="2">

                </v-col>
            </v-row>
            <div style="color: white;">
                <v-autocomplete chips v-model="condition_model" base-color="white" :items="condition" multiple
                    variant="underlined"></v-autocomplete>
            </div>
        </div> -->
        <div style="background-color: #01579B; width: 100%;height: 10px;"></div>
        <div v-show="chart_comp?.chartStyle?.style['bmap'] || has_geo"
            style="margin: 10px;margin-top: 15px;margin-right: 15px;margin-left: 15px; background-color: #003355; color: white;">
            <span style="font-weight: bold;color: green;">
                专用属性
            </span>
            <v-text-field v-show="chart_comp?.chartStyle?.style['bmap']" v-model="chart_text" flat rounded="false"
                :label="chart_text == '' ? '百度地图密钥' : ''" color="red" type="text" variant="underlined" density="compact"
                hide-details></v-text-field>

            <v-text-field v-show="has_geo" @click="selectGeoJson" color="red" @click:append-inner="selectGeoJson"
                type="text" :label="has_geo ? 'GeoJSON文件' : ''" readonly variant="underlined" hide-details
                append-inner-icon="md:attach_file">
            </v-text-field>

            <v-text-field v-show="chart_comp?.chartStyle?.style['bmap']" @click="selectKML" color="red"
                @click:append-inner="selectKML" type="text" label="KML图层" readonly variant="underlined" hide-details
                append-inner-icon="md:attach_file">
            </v-text-field>
        </div>
    </div>
    <div v-show="store.state.viBreadindex != 0" style=" height: 100%;  background-color: transparent;">
        <v-expansion-panels variant="Popout" v-model="selected_item" tile flat>
            <v-expansion-panel v-for="(item, key) in current_coor_object" :value="key" :style="{
                backgroundColor: 'transparent',
                color: item.is_object ? '#ff9900' : 'white'
            }" tile>
                <v-expansion-panel-title collapse-icon="md:keyboard_arrow_up" expand-icon="md:keyboard_arrow_down">
                    <span style="font-size:large;font-weight: bold; white-space: nowrap;
                                /* 防止内容换行 */
                                overflow: hidden;
                                /* 隐藏溢出部分的内容 */
                                text-overflow: ellipsis;
                                "> {{ key }}</span>
                    <v-tooltip activator="parent" location="left">{{ key }}</v-tooltip>
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                    <v-card flat rounded="0" style="background-color: #003355;">

                        <v-card-actions style="color: white;">
                            <div v-if="item.coor.uiControl" style="width: 100%;height: 100%;">

                                <v-text-field v-model="current_value_computed" color="green"
                                    v-if="item.coor.uiControl.type == 'number'" density="compact" elevation="0"
                                    type="number" :min="item.coor.uiControl.min" :max="item.coor.uiControl.max"
                                    :step="item.coor.uiControl.step" variant="underlined" hide-details>
                                </v-text-field>
                                <v-textarea v-model="current_value_computed"
                                    v-else-if="item.coor.uiControl.type == 'function'" variant="underlined"
                                    density="compact" elevation="0"></v-textarea>
                                <v-combobox v-model="current_value_computed" v-else-if="item.coor.uiControl.type == 'enum'"
                                    variant="underlined" :items="stringToArray(item.coor.uiControl.options)"
                                    density="compact" elevation="0"></v-combobox>
                                <v-combobox v-model="current_value_computed"
                                    v-else-if="item.coor.uiControl.type == 'boolean'" variant="underlined"
                                    :items="[true, false]" density="compact" elevation="0"></v-combobox>
                                <v-text-field v-model="current_value_computed"
                                    @click:append-inner="dialog_color_show = true" :modes="['hexa', 'rgba', 'hsla']"
                                    append-inner-icon="md:color_lens" :color="current_value_computed"
                                    v-else-if="item.coor.uiControl.type == 'color'" variant="underlined" density="compact"
                                    elevation="0"></v-text-field>
                                <v-text-field v-else v-model="current_value_computed" flat rounded="false" color="green"
                                    type="text" variant="underlined" density="compact" elevation="0"
                                    hide-details></v-text-field>

                            </div>
                            <v-text-field v-model="current_value_computed" @click:append-inner="dialog_color_show = true"
                                :modes="['hexa', 'rgba', 'hsla']" append-inner-icon="md:color_lens"
                                :color="current_value_computed"
                                v-else-if="key.toString().indexOf('Color') > -1 || key.toString().indexOf('color') > -1"
                                variant="underlined" density="compact" elevation="0"></v-text-field>
                            <v-textarea v-model="current_value_computed"
                                v-else-if="key.toString().indexOf('rich') > -1 || key.toString() == 'data' || key.toString() == 'indicator' || key.toString() == 'renderItem'"
                                variant="underlined" density="compact" elevation="0"></v-textarea>
                            <v-btn v-else-if="item.is_object == true" style="font-size:medium;" variant="plain"
                                color="green" @click="object_click(key)">
                                进入配置
                            </v-btn>
                            <v-text-field v-else v-model="current_value_computed" flat rounded="false" color="green"
                                type="text" variant="underlined" density="compact" elevation="0"
                                hide-details></v-text-field>
                        </v-card-actions>
                        <div style="display: flex;align-items: center;">
                            <!-- <v-btn @click.stop="functionIconClick" icon variant="plain" density="compact"
                                v-show="item.is_object != true" :color="is_function_computed ? '#ff9900' : 'white'"
                                elevation="0">
                                <v-icon>{{ mdiFunctionVariant }}</v-icon>
                                <v-tooltip activator="parent" location="left">编写函数</v-tooltip>
                            </v-btn> -->
                            <div style="padding-bottom: 10px;">
                                <v-icon @click.stop="functionIconClick"
                                    v-show="item.is_object != true || key.toString() == 'renderItem'"
                                    :color="is_function_computed ? 'green' : 'yellow'">{{ mdiFunctionVariant }}</v-icon>
                                <v-tooltip activator="parent" location="left">编写函数</v-tooltip>
                            </div>
                            <div :style="{ paddingLeft: item.is_object || key.toString().indexOf('rich') > -1 || key.toString() == 'data' || key.toString() == 'indicator' || key.toString() == 'renderItem' ? '15px' : '0px', width: '90%' }"
                                @click.stop="showInfo(key.toString() == 'renderItem' ? item.coor.desc + ' 详细参阅：https://echarts.apache.org/zh/option.html#series-custom' : item.coor.desc)"
                                class="hoverdiv">
                                <p>{{ formattedText(item.coor.desc) }}</p>
                            </div>
                        </div>
                    </v-card>

                </v-expansion-panel-text>

            </v-expansion-panel>
            <v-dialog v-model="dialog_color_show" width="auto">
                <v-color-picker v-model="current_value_computed" show-swatches></v-color-picker>
            </v-dialog>
        </v-expansion-panels>

    </div>
</template>


<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue';
import { OpenFileDialog, ReadKmlFile, RequestBaiduAPI, SendFileToFront } from '../../../../wailsjs/go/main/App';
import { mdiPlus, mdiClose, mdiDragVariant, mdiFunctionVariant, mdiChevronRight, mdiPencilCircle, mdiPencilCircleOutline, mdiDatabaseEdit, mdiDatabaseEditOutline, mdiPencilPlusOutline, mdiCodeJson } from '@mdi/js';
import { reactive } from 'vue';
import { ViChartOption, ViChartStyle, ViComponent, store, ViChartType, ViChartEnum } from '../../../store';
import { onUnmounted } from 'vue';
import { computed } from 'vue';
import { ComputedRef } from 'vue';
import { WritableComputedRef } from 'vue';
import { Ref } from 'vue';
import { cloneDeep, get, set, unset } from 'lodash';
import { watch } from 'vue';
import { watchEffect } from 'vue';
import { SeriesOption } from 'echarts';
import { QueryBySheetAndField } from '../../../../wailsjs/go/main/App';
import { WGS84ToBD09, convertARGBtoRGBA, createRenderItemAPI, drawPolygon, getElement, sleep } from '../../../utils/func';
import * as echarts from 'echarts'
import { NewDynamicFunction, IsFunction, get_current_path, setup_option } from '../../../utils/func';
import { EventsEmit } from '../../../../wailsjs/runtime/runtime';
import { storeDisplay } from '../../../store/display';
import { createRenderItem } from '../../../utils/func'
import { inject } from 'vue';

//条件过滤的组件名称和id,value是id
let condition = ref([] as { title: string, value: string }[]);

//选择的是哪个属性
let selected_item = ref('')
//当前属性的值
// let current_value = ref()
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
//对象属性模板列表
let object_template_list: Ref<{
    [key: string]: any;
}> = ref({})



let dialog_color_show = ref(false)

function selectKML() {
    let viewid = store.state.viViews[store.state.viCtrlIndex].id
    let compid = chart_comp.value?.id
    if (compid) {
        store.commit("show_loader_dialog", { show: true })
        ReadKmlFile(viewid, compid).then(geomsg => {
            if (geomsg.Error != "") {
                store.commit("show_loader_dialog", { show: false })
                console.log(geomsg.Error)
                return
            }
            let current_count = 0
            let bmap = chart_comp.value?.chartStyle?.style['bmap']
            data_selected.value?.forEach(v => {
                let a = Number(v.name[v.name.length - 1])
                if (a >= current_count) {
                    current_count = a + 1
                }

            })
            // console.log("geomsg", geomsg)
            if (geomsg.LinearRings) {
                if (geomsg.LinearRings.length > 0) {
                    var dataset = []
                    for (var i = 0; i < geomsg.LinearRings.length; i++) {
                        dataset.push(
                            [
                                //第二项应该是geomsg.LinearRings[i].name，但还不能在自定义的图上显示name,为了避免占用空间，先用空
                                geomsg.LinearRings[i].CoordString, "", geomsg.LinearRings[i].Color, geomsg.LinearRings[i].Fill
                            ]
                        )
                    }
                    data_selected.value = [
                        {
                            key: 'custom',
                            name: '自定义图层' + current_count,
                            count: current_count++,
                        }
                    ]
                    store.commit("add_echart_series", {
                        type: {
                            type: 'custom',
                            coordinateSystem: bmap ? 'bmap' : 'geo',
                            renderItem: createRenderItemAPI(bmap),
                            // dataset: dataset,
                            animation: false,
                            silent: true,
                            data: dataset,
                        }
                    })
                    nextTick(() => {
                        store.commit("snapshot")
                    })

                }
            }
            if (geomsg.LineStrings) {
                if (geomsg.LineStrings.length > 0) {
                    let data = []
                    for (var i = 0; i < geomsg.LineStrings.length; i++) {
                        let coords = geomsg.LineStrings[i].CoordString.split(" ")
                        let name = geomsg.LineStrings[i].name
                        let lst = []
                        for (var j = 0; j < coords.length; j++) {
                            let arr = coords[j].split(",")
                            if (arr.length >= 2) {
                                if (arr[0].length > 0 && arr[1].length > 0) {
                                    let lng = parseFloat(arr[0])
                                    let lat = parseFloat(arr[1])

                                    if (bmap) {
                                        let tmp = WGS84ToBD09(lng, lat)
                                        lng = tmp.lng
                                        lat = tmp.lat
                                    }
                                    lst.push([lng, lat])
                                }
                            }
                        }
                        data.push({ name: name, coords: lst })

                    }
                    data_selected.value = [
                        {
                            key: 'lines',
                            name: '路径图' + current_count,
                            count: current_count++,
                        }
                    ]
                    // console.log("data", data)
                    let obj = {
                        type: 'lines',
                        coordinateSystem: bmap ? 'bmap' : 'geo',
                        data: data,
                        polyline: true,
                        lineStyle: {
                            color: 'red', width: 3,
                        },
                    }
                    store.commit("add_echart_series", {
                        type: obj
                    })
                    nextTick(() => {
                        store.commit("snapshot")
                    })
                }
            }
            if (geomsg.Points) {
                if (geomsg.Points.length > 0) {
                    //

                    let data = []
                    for (var i = 0; i < geomsg.Points.length; i++) {
                        let arr = geomsg.Points[i].CoordString.split(",")
                        if (arr.length >= 2) {
                            if (arr[0].length > 0 && arr[1].length > 0) {
                                let lng = parseFloat(arr[0])
                                let lat = parseFloat(arr[1])
                                let name = geomsg.Points[i].name
                                if (bmap) {
                                    let tmp = WGS84ToBD09(lng, lat)
                                    lng = tmp.lng
                                    lat = tmp.lat
                                }
                                data.push([lng, lat, name])
                            }
                        }
                    }
                    data_selected.value = [
                        {
                            key: 'scatter',
                            name: '散点图' + current_count,
                            count: current_count++,
                        }
                    ]
                    let obj = {
                        type: 'scatter',
                        coordinateSystem: bmap ? 'bmap' : 'geo',
                        data: data,
                        label: {
                            formatter: (params: any) => {
                                return params.value[2]
                            },
                        },
                        tooltip: {
                            formatter: (params: any) => {
                                return params.value[2]
                            },
                        }
                    }
                    store.commit("add_echart_series", {
                        type: obj
                    })
                    nextTick(() => {
                        store.commit("snapshot")
                    })
                }
            }
            updateOption("selectKML")
            store.commit("show_loader_dialog", { show: false })

        }).catch(reason => {
            console.log("kml加载出错哦", reason)
            store.commit("show_loader_dialog", { show: false })
        })
        // store.commit("show_loader_dialog", { show: false })
    }

}

function selectGeoJson() {
    let viewid = store.state.viViews[store.state.viCtrlIndex].id
    let compid = chart_comp.value?.id
    if (compid) {
        OpenFileDialog("*.json", viewid, compid, false).then((result: string) => {
            if (result.indexOf("Error") != -1) {
                console.log(result)
                return
            }
            let urls = result.split(".")
            fetch(urls[0]).then(v => {
                // echarts.registerMap()
                v.json().then(value => {
                    // echarts.registerMap(compid as string,value)
                    store.commit("set_jsongeomap", value)
                    nextTick(() => {
                        store.commit("snapshot")
                    })
                    updateOption("selectGeoJson")
                })

            }).catch(reason => {
                console.log("GeoJson地图加载出错", reason)
            })

        })
    }

}

function loadBaiduMap(url: string) {
    window.BMAP_PROTOCOL = "https";
    window.BMap_loadScriptTime = (new Date).getTime();
    var scripts = document.getElementsByTagName('script');
    for (var i = 0; i < scripts.length; i++) {
        let src = scripts[i].getAttribute('src')
        if (src) {
            if (src?.indexOf("baidu.com") > -1) {
                scripts[i].parentNode?.removeChild(scripts[i])
            }
        }

    }
    var scriptElement = document.createElement('script');
    scriptElement.setAttribute('type', 'text/javascript');
    scriptElement.src = url;
    scriptElement.onload = function () {
        updateOption("loadBaiduMap")
    }

    document.head.append(scriptElement);
}

let chart_text = computed({
    get: () => {
        if (chart_comp.value?.style.displayMode) {
            return storeDisplay.state.viBaiduMapKey
        }
        return store.state.viBaiduMapKey
    },
    set: async (val: string) => {

        await RequestBaiduAPI(store.state.viBaiduMapKey).then(url => {
            if (url.indexOf("Error:") > -1) {
                console.log(url)
                return
            }
            loadBaiduMap(url)
            store.commit("set_baidu_map_key", val)
            nextTick(() => {
                store.commit("snapshot")
            })
        })
    }
})

//如果chart 属性有geo或者series有geo 是true
let has_geo = computed(() => {
    let result = false
    if (chart_comp.value) {
        let style = chart_comp.value?.chartStyle
        if (style) {
            if (style.style['geo']) {
                result = true
            }
            let series = style.style.series
            if (series) {
                let s = series as SeriesOption[]
                s.forEach(v => {
                    if (v.type == 'map') {
                        result = true
                    }
                })
            }
        }
    }

    return result
})


//可选坐标类型
const coor = {
    xAxis: '横坐标抽',
    yAxis: '纵坐标抽',
    singleAxis: '单坐标抽',
    geo: '地理坐标系',
    bmap: '百度地图',
    // calendar: '日历坐标系',
    parallel: '平行坐标系',
    parallelAxis: '平行坐标抽',
    polar: '极坐标系',
    radiusAxis: '极坐标径向抽',
    angleAxis: '极坐标角度抽',
    radar: '雷达图坐标',
}

//图表属性
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


// let condition_model = computed({
//     get: () => {
//         let res = chart_comp.value?.chartStyle?.chartConditionList
//         if (res == undefined) {
//             return []
//         } else {
//             return res
//         }
//     },
//     set: (val) => {

//         store.commit("set_chartConditionList", { value: val })

//     }
// })

let is_function_computed = computed(() => {
    return IsFunction(store_current_value.value)
})



//选择的是哪个属性 计算属性
let current_value_computed = computed({
    get: () => {
        if (typeof store_current_value.value == 'object') {
            let res = JSON.stringify(store_current_value.value)
            return res
        }
        let path = get_current_path(selected_item.value)
        if (path == "visualMap.type") {
            if (store_current_value.value == "") {
                return "continuous"
            }
        }
        return store_current_value.value
    },
    set: (val) => {
        store_current_value.value = val

        if (store_current_value.value == undefined) {
            return
        }
        // 匹配命名函数
        let str = store_current_value.value.toString()
        let path = get_current_path(selected_item.value)
        setup_option(str, path, chart_comp.value?.id)
        sleep(200).then(() => {
            nextTick(() => {
                updateOption("current_value_computed")
            })
        })
    }
});

watch(() => selected_item.value, (newValue) => {
    if (newValue != "") {
        let path = store.state.viBreaditems.join('.') + "." + selected_item.value
        if (store.state.ViChartRefIndex.type == ViChartEnum.Data) {
            let parts = path.split(".");
            path = "series[" + store.state.ViChartRefIndex.index + "]." + parts.slice(2).join(".");
        }
        path = path.replace("NNI.", "")
        let el = getElement(chart_comp.value?.id + "nni") as HTMLElement
        if (el != null) {
            let chart = echarts.getInstanceByDom(el)
            let option = chart?.getOption() as echarts.EChartsOption
            let value = get(option, path)
            if (value == undefined) {
                let p1 = path.substring(0, path.indexOf("."))
                let p2 = path.substring(path.indexOf(".") + 1)
                let newpath = p1 + "[0]." + p2
                // console.log(222,option,newpath, value)
                value = get(option, newpath)
            }
            // console.log(111,option,path, value)
            if (value) {
                store_current_value.value = value
            } else {
                store_current_value.value = ""
            }
        }
    }
})

//已选择的坐标名称
let coor_selected: WritableComputedRef<ViChartOption[] | undefined> = computed({
    get: () => {
        let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
        let name_list: ViChartOption[] | undefined = undefined
        if (selected.length == 1) {
            store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
                if (item.id == selected[0]) {
                    name_list = item.chartStyle?.coorNameList
                    return
                }
            })

        }
        //初始化条件过滤组件信息
        condition.value.length = 0
        store.state.viViews.forEach(v => {
            v.components.forEach(c => {
                if (c.typeName == "单选" || c.typeName == "多选") {
                    condition.value.push({ title: c.customName, value: c.id })
                }
                c.components.forEach(g => {
                    if (g.typeName == "单选" || g.typeName == "多选") {
                        condition.value.push({ title: g.customName, value: g.id })
                    }
                })
            })
        })
        return name_list
    },
    set: (val) => {

        store.commit("set_chart_style_name_list", { value: val })
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
                // if (option.series != undefined) {
                //     let series = option.series as SeriesOption[]
                //     for (var i = 0; i < series.length; i++) {
                //         if (series[i].type == 'custom') {
                //             let dataset = series[i] as any
                //             set((option.series as SeriesOption[])[i], "renderItem", createRenderItem(true, dataset.dataset as echarts.DatasetComponentOption[]))
                //         }

                //     }

                // }
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


function coorMenuItemClick(name: string, key: "parallel" | "angleAxis" | "radiusAxis" | "xAxis" | "yAxis" | "calendar" | "polar" | "singleAxis" | "parallelAxis" | "geo" | "radar" | "bmap") {
    // X坐标  Y坐标可以有两个
    if (name == "横坐标抽" || name == "纵坐标抽") {
        let sum = 0
        coor_selected.value?.forEach(v => {
            if (v.count > 0) {
                if (v.key == "calendar" || v.key == "singleAxis" || v.key == "polar" || v.key == "parallelAxis" || v.key == "geo") {
                    sum++
                }
            }
        })
        if (sum) {
            store.commit("show_dialog", { show: true, msg: "该类坐标与已存在的坐标抽存在冲突" })
            return
        }
        let exit = false
        let current_count = 0
        let max_count = 0
        coor_selected.value?.forEach(v => {
            if (v.key == "calendar" || v.key == "singleAxis" || v.key == "polar" || v.key == "geo") {
                if (v.count > 0) {
                    exit = true
                }
            }
            if (v.key == key) {
                current_count = v.count
                let n = Number(v.name[v.name.length - 1])
                if (n >= max_count) {
                    max_count = n + 1
                }
            }
        })
        if (exit) {
            return
        }
        if (current_count < 2) {
            // coor_selected.value.push(name + coor_count[key])
            // coor_count[key]++
            coor_selected.value = [
                {
                    key: key,
                    name: name + max_count,
                    count: current_count + 1,
                }
            ]
            store.commit("set_chart_style_coor", { key: key })
            nextTick(() => {
                store.commit("snapshot")
            })
            // updateCoorToChart(key)
        } else {
            store.commit("show_dialog", { show: true, msg: "该类坐标抽最多只能添加2个" })
            return
        }
    } else if (name.indexOf("极坐标") > -1) {
        let sum = 0
        let current_count = 0
        coor_selected.value?.forEach(v => {
            if (v.count > 0) {
                if (v.key == "polar" || v.key == "radiusAxis" || v.key == "angleAxis") {
                    //pass
                } else {
                    sum++
                }
            }
            if (v.key == key) {
                current_count = v.count
            }
        })
        if (sum) {
            store.commit("show_dialog", { show: true, msg: "该类坐标与已存在的坐标抽存在冲突" })
            return
        }
        if (current_count == 0) {
            coor_selected.value = [
                {
                    key: key,
                    name: name,
                    count: current_count + 1,
                }
            ]
            store.commit("set_chart_style_coor", { key: key })
            // nextTick(() => {
            //     store.commit("snapshot")
            // })
        } else {
            store.commit("show_dialog", { show: true, msg: "该类坐标与已存在的坐标抽存在冲突" })
            return
        }
    } else if (name.indexOf("平行坐标") > -1) {
        let sum = 0
        let current_count = 0
        coor_selected.value?.forEach(v => {
            if (v.count > 0) {
                if (v.key == "parallel" || v.key == "parallelAxis") {
                    //pass
                } else {
                    sum++
                }
            }
            if (v.key == key) {
                current_count = v.count
            }
        })
        if (sum) {
            store.commit("show_dialog", { show: true, msg: "该类坐标与已存在的坐标抽存在冲突" })
            return
        }
        if (current_count == 0) {
            coor_selected.value = [
                {
                    key: key,
                    name: name,
                    count: current_count + 1,
                }
            ]
            store.commit("set_chart_style_coor", { key: key })
            // nextTick(() => {
            //     store.commit("snapshot")
            // })
        } else {
            store.commit("show_dialog", { show: true, msg: "该类坐标与已存在的坐标抽存在冲突" })
            return
        }
    } else {
        let count = 0
        coor_selected.value?.forEach(v => {
            if (v.count > 0) {
                count++
                return
            }
        })
        if (count == 0) {
            coor_selected.value = [
                {
                    key: key,
                    name: name,
                    count: count + 1,
                }
            ]
            store.commit("set_chart_style_coor", { key: key })
            // nextTick(() => {
            //     store.commit("snapshot")
            // })
        } else {
            store.commit("show_dialog", { show: true, msg: "该类坐标与已存在的坐标抽存在冲突" })
            return
        }
    }

    updateOption("coorMenuItemClick")


}

//坐标抽删除按钮
async function coorDeleteClick(name: string, index: number) {
    store.commit("del_chart_style_name_list", { value: index })
    nextTick(() => {
        store.commit("snapshot")
    })
    updateOption("coorDeleteClick")
}

//数据删除按钮
function dataDeleteClick(name: string, index: number) {
    store.commit("del_chart_style_data_list", { value: index })
    nextTick(() => {
        store.commit("snapshot")
    })
    updateOption("dataDeleteClick")
}

//可选属性删除按钮
function otherDeleteClick(name: string, index: number) {
    store.commit("del_chart_style_other_list", { value: index })
    store.commit("del_chart_style_other", { value: name })
    nextTick(() => {
        store.commit("snapshot")
    })
    updateOption("otherDeleteClick")
}


let is_first_click = ref(true)

function updateMainOption(itemkey: string) {
    current_coor_object.value = {}
    let newkey = itemkey
    if (itemkey.endsWith("]")) {
        newkey = itemkey.slice(0, -3)
    }
    let obj = object_template_list.value[newkey]
    Object.keys(obj).forEach(k => {
        if (k.indexOf('.') == -1) {
            if (k != "id" && k != "data") {
                current_coor_object.value[k] = {
                    is_object: false,
                    coor: obj[k]
                }
            }
        }
    })
    Object.keys(obj).forEach(k => {
        if (k.indexOf('.') !== -1) {
            let arr = k.split(".")
            let item = current_coor_object.value[arr[0]]
            if (item === undefined) {
                if (arr[0] !== "data" && arr[0] !== "id") {
                    current_coor_object.value[arr[0]] = {
                        is_object: true,
                        coor: { desc: "" }
                    }
                }
            }
        }
    })
    Object.keys(current_coor_object.value).forEach(key => {
        if (Object.keys(obj).filter(it => {
            return it.startsWith(key + ".")
        }).length > 0) {
            current_coor_object.value[key].is_object = true
        }
    })
    // console.log("updateMainOption", current_coor_object.value)
}

function updateSubOption() {
    //第二个元素是主键，图xAxis
    let master_key = store.state.viBreaditems[1]
    master_key = master_key.replace("[0]", "").replace("[1]", "")

    //当前对象的属性路径前缀
    let current_key = store.state.viBreaditems.join(".").replace("NNI." + master_key + ".", "")
    current_key = current_key.replace("[0]", "").replace("[1]", "")
    if (current_key == "NNI." + master_key) {
        current_key = master_key
        updateMainOption(current_key)
        return
    }
    let obj = object_template_list.value[master_key]
    current_key = current_key.replace("NNI." + master_key + '.', "")
    // console.log(current_key, master_key)
    current_coor_object.value = {}
    Object.keys(obj).forEach(k => {
        if (k.startsWith(current_key + ".")) {

            let new_key = k.replace(current_key + ".", "")
            if (new_key.indexOf(".") == -1) {
                current_coor_object.value[new_key] = {
                    is_object: false,
                    coor: obj[k]
                }
            }
        }
    })
    Object.keys(obj).forEach(k => {
        if (k.startsWith(current_key + ".")) {
            let new_key = k.replace(current_key + ".", "")
            if (new_key.indexOf(".") != -1) {
                let arr = new_key.split(".")
                let item = current_coor_object.value[arr[0]]
                if (item === undefined) {
                    current_coor_object.value[arr[0]] = {
                        is_object: true,
                        coor: { desc: "" }
                    }
                }
            }
        }
    })

    // console.log("updateSubOption", current_coor_object.value)
    Object.keys(current_coor_object.value).forEach(key => {
        if (Object.keys(obj).filter(item => {
            // console.log(item, key)
            return item.startsWith(current_key + "." + key + ".")
        }).length > 0) {
            current_coor_object.value[key].is_object = true
        }
    })

}

//因为属性中有些属性是对象类型，需要点击对应属性进入对象里调整更多属性
function object_click(key: string | number) {
    store.commit("set_bread_items", { value: key })
    store.commit("set_bread_index", { value: store.state.viBreadindex + 1 })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })

    updateSubOption()
}

watch(() => store.state.viBreaditems.length, (newValue) => {
    if (is_first_click.value) {
        is_first_click.value = false
        return
    }
    if (newValue < 2) {
        store.commit("set_ViChartRefIndex", {
            value: {
                type: ViChartEnum.Undefined,
                index: -1
            }
        })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
        return
    }
    updateSubOption()

})






//已选择的数据名称
let data_selected: WritableComputedRef<ViChartOption[] | undefined> = computed({
    get: () => {
        let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
        let name_list: ViChartOption[] | undefined = undefined
        if (selected.length == 1) {
            store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
                if (item.id == selected[0]) {
                    name_list = item.chartStyle?.dataNameList
                    return
                }
            })

        }
        return name_list
    },
    set: (val) => {

        store.commit("set_chart_style_data_list", { value: val })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })

    }
});

function dataMenuItemClick(name: string, key: string) {
    // console.log(name, key)
    let current_count = 0
    data_selected.value?.forEach(v => {
        let a = Number(v.name[v.name.length - 1])
        if (a >= current_count) {
            current_count = a + 1
        }

    })
    data_selected.value = [
        {
            key: key,
            name: name + current_count,
            count: current_count + 1,
        }
    ]
    store.commit("add_echart_series", {
        type: key
    })
    nextTick(() => {
        store.commit("snapshot")
    })
}

function otherMenuItemClick(name: string, key: string) {
    let current_count = 0
    other_selected.value?.forEach(v => {
        if (v.key == key) {
            current_count = v.count
            return
        }
    })
    if (current_count == 0) {
        if (key == "piecewise" || key == "continuous") {
            other_selected.value?.forEach(v => {
                if (v.key == "piecewise" || v.key == "continuous") {
                    if (v.count > 0) {
                        current_count = v.count
                        return
                    }

                }
            })
        }
    }
    if (current_count > 0) {
        store.commit("show_dialog", { show: true, msg: "该属性类型与已存在的属性类型存在冲突" })
        return
    }
    other_selected.value = [
        {
            key: key,
            name: name,
            count: current_count + 1,
        }
    ]
    store.commit("add_chart_style_other", {
        value: key
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
    updateOption("otherMenuItemClick")
}

//当前选择的对象，比如xAxis对象,用来显示到组件上的
let current_coor_object: Ref<{
    [key: string]: {
        //是否是对象类型？
        is_object: boolean,
        //对象，可能是属性类型，也可能是对象
        coor: any,
    };
}> = ref({})

async function setupiscontrolcoor(itemkey: string, index: number) {

    is_first_click.value = true
    store_current_value.value = ""
    selected_item.value = ""
    let haskey = false
    Object.keys(object_template_list.value).forEach(key => {
        if (key == itemkey) {
            haskey = true
            return
        }
    })
    let error = false
    if (!haskey) {
        await SendFileToFront("templates/static/chartoptions/" + itemkey + ".json").then(v => {
            if (v.startsWith("Error:")) {
                error = true
                console.log(v)
                return
            }
            object_template_list.value[itemkey] = JSON.parse(v)
        }).catch(reason => {
            error = true
            console.log(reason)
            return
        })
    }
    let keyinfo = itemkey;
    if (itemkey == "xAxis" || itemkey == "yAxis") {
        let coorlist = coor_selected.value
        if (coorlist != undefined) {
            for (let i = 0; i < index; i++) {
                if (coorlist[i].key == itemkey) {
                    keyinfo = itemkey + "[1]"
                    break
                }
            }
        }
        if (itemkey == keyinfo) {
            keyinfo = itemkey + "[0]"
        }

    }
    updateMainOption(keyinfo)

    if (!error) {

        store.commit("set_bread_items", { value: keyinfo })
        store.commit("set_bread_index", { value: 1 })
        store.commit("set_ViChartRefIndex", {
            value: {
                type: ViChartEnum.Coor,
                index: index
            }
        })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }

}


async function setupiscontrolseries(itemkey: string, index: number) {
    is_first_click.value = true
    store_current_value.value = ""
    selected_item.value = ""
    let haskey = false
    Object.keys(object_template_list.value).forEach(key => {
        if (key == itemkey) {
            haskey = true
            return
        }
    })
    let error = false
    if (!haskey) {
        await SendFileToFront("templates/static/chartoptions/series/" + itemkey + ".json").then(v => {
            if (v.startsWith("Error:")) {
                error = true
                console.log(v)
                return
            }
            object_template_list.value[itemkey] = JSON.parse(v)
        }).catch(reason => {
            error = true
            console.log(reason)
            return
        })
    }
    updateMainOption(itemkey)
    if (!error) {
        store.commit("set_bread_items", { value: itemkey })
        store.commit("set_bread_index", { value: 1 })
        store.commit("set_ViChartRefIndex", {
            value: {
                type: ViChartEnum.Data,
                index: index
            }
        })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }

}


function stringToArray(text: string) {
    return text.split(',')
}

function formattedText(text: string) {
    // 使用DOMParser解析HTML并提取文本内容
    const parser = new DOMParser();
    const doc = parser.parseFromString(text, 'text/html');
    return doc.body.textContent || "";
}

let other_selected: WritableComputedRef<ViChartOption[] | undefined> = computed({
    get: () => {
        let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
        let name_list: ViChartOption[] | undefined = undefined
        if (selected.length == 1) {
            store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
                if (item.id == selected[0]) {
                    name_list = item.chartStyle?.otherNameList
                    return
                }
            })

        }
        return name_list
    },
    set: (val) => {

        store.commit("set_chart_style_other_list", { value: val })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});


async function setupiscontrolother(itemkey: string, index: number) {
    // console.log(itemkey, index)
    is_first_click.value = true
    let haskey = false
    let modkey = itemkey
    if (modkey == "piecewise" || modkey == "continuous") {
        modkey = "visualMap"
    }
    Object.keys(object_template_list.value).forEach(key => {
        if (key == modkey) {
            haskey = true
            return
        }
    })
    let error = false
    if (!haskey || modkey == "visualMap") {
        await SendFileToFront("templates/static/chartoptions/" + itemkey + ".json").then(v => {
            if (v.startsWith("Error:")) {
                error = true
                console.log(v)
                return
            }
            object_template_list.value[modkey] = JSON.parse(v)

        }).catch(reason => {
                error = true
                console.log(reason)
                return
            })
    }
    updateMainOption(modkey)
    if (!error) {
        store.commit("set_bread_items", { value: modkey })
        store.commit("set_bread_index", { value: 1 })
        store.commit("set_ViChartRefIndex", {
            value: {
                type: ViChartEnum.Other,
                index: index
            }
        })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }

}

// let current_chartstyle: ComputedRef<ViChartStyle | undefined> = computed(() => {
//     let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids
//     let style: ViChartStyle | undefined = undefined
//     if (selected.length == 1) {
//         store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
//             if (item.id == selected[0]) {
//                 style = item.chartStyle
//                 return
//             }
//         })
//     }
//     return style
// })


//可选数据类型
const data = {
    line: '折线图',
    bar: '柱状图',
    pie: '饼图',
    scatter: '散点图',
    effectScatter: '漪特散点',
    radar: '雷达图',
    // tree: '树图',
    // treemap: '树状数',
    // sunburst: '旭日图',
    boxplot: '箱形图',
    candlestick: 'K线图',
    heatmap: '热力图',
    map: '地图',
    // parallel: '平行坐标',
    lines: '路径图',
    // graph: '关系图',
    // sankey: '桑基图',
    funnel: '漏斗图',
    gauge: '仪表盘',
    // pictorialBar: '象形柱图',
    themeRiver: '主题河流',
    custom: '自定义系列'
}


//可选其他属性
const other = {
    title: '标题',
    legend: '图例',
    grid: "网格",
    tooltip: '提示框',
    // timeline: '时间线',
    // brush: '区域选择',
    axisPointer: '坐标指示',
    continuous: '连续映射',
    piecewise: '分段映射',
    textStyle: '全局字体',
}


watchEffect(() => {


    if (store.state.viViews[store.state.viCtrlIndex]) {

        if (coor_selected.value?.length == 0 && data_selected.value?.length == 0 && other_selected.value?.length == 0) {
            let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids

            if (selected.length == 1) {
                let item_count = ref(0)
                store.state.viViews[store.state.viCtrlIndex].components.forEach((item, index) => {
                    if (item.id == selected[0]) {
                        let style = item.chartStyle?.style
                        if (style != undefined) {
                            Object.keys(style).forEach((k, ind) => {
                                if (style != undefined) {
                                    let series: SeriesOption[] = style[k] as SeriesOption[]
                                    if (k == "series") {
                                        series.forEach((value) => {
                                            let series_type = value.type
                                            if (series_type != undefined) {
                                                let name = data[series_type as keyof typeof data]
                                                data_selected.value = [{
                                                    key: series_type,
                                                    name: name + item_count.value++,
                                                    count: 1,
                                                }]
                                            }

                                        })
                                    } else {
                                        if (coor[k as keyof typeof coor] != undefined) {
                                            coor_selected.value = [{
                                                key: k,
                                                name: coor[k as keyof typeof coor] + ind,
                                                count: 1,
                                            }]
                                        }
                                        // console.log("xxx", k)
                                        if (other[k as keyof typeof other] != undefined) {
                                            other_selected.value = [{
                                                key: k,
                                                name: other[k as keyof typeof other],
                                                count: 1,
                                            }]
                                        }
                                        if (k == "visualMap") {
                                            // console.log("yyy", style[k])
                                            let obj = style[k] as echarts.VisualMapComponentOption || undefined
                                            if (obj != undefined) {
                                                let visualtype = obj.type
                                                let name = "continuous"
                                                if (visualtype == undefined || visualtype == "continuous") {
                                                    name = "continuous"
                                                } else {
                                                    name = "piecewise"
                                                }
                                                other_selected.value = [{
                                                    key: name,
                                                    name: other[name as keyof typeof other],
                                                    count: 1,
                                                }]
                                            }

                                        }
                                    }
                                }
                            })

                        }
                        return
                    }
                })

            }

        }
    }
})

function dragover(event: any) {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'copy'
    event.target.style.color = "green";
}

function dragleave(event: any) {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'none'
    event.target.style.color = "white";
}



//点击的坐标对象item和在coor_selected的索引
function coorHandleDrop(event: any, item: ViChartOption, index: number) {
    event.target.style.color = "white";
    let data: string = event.dataTransfer.getData("application/json");
    if (data.length < 1) {
        return
    }
    let json = JSON.parse(data)
    // console.log(json.type, json.sheet_key, json.field_key)
    if (json.type == 'data') {
        QueryBySheetAndField(json.sheet_key, [json.field_key], false).then(v => {
            if (v.error != "") {
                console.log(v.error)
                return
            }

            let count = -1
            let coor_list = coor_selected.value
            if (coor_list != undefined) {
                for (let i = 0; i <= index; i++) {
                    if (coor_list[i].key == item.key) {
                        count++
                    }
                }
            } else {
                count = 0
            }
            let newdata = v.data.flat()
            let path = item.key + ".data"
            if (item.key == "radar") {
                path = item.key + ".indicator"
                newdata = []
                v.data.flat().forEach(value => {
                    newdata.push({
                        name: value
                    })
                })
            }
            // console.log("newdata", newdata)
            if (item.key == "xAxis" || item.key == "yAxis") {
                path = item.key + "[" + count + "].data"
            }
            store.commit("set_echart_option_by_path", {
                id: chart_comp.value?.id,
                path: path,
                value: newdata,
                name: json.field_key,
            })
            store.commit("set_echart_coor_map", {
                path: path,
                value: {
                    sheetkey: json.sheet_key,
                    fieldkey: json.field_key,
                }
            })
            // nextTick(() => {
            //     store.commit("snapshot")
            // })
            updateOption("coorHandleDrop")

        }).catch(reason => {
            console.log(reason)
        })
    }

}

//点击的数据对象item和在coor_selected的索引
function dataHandleDrop(event: any, item: ViChartOption, index: number) {
    event.target.style.color = "white";
    let data: string = event.dataTransfer.getData("application/json");
    if (data.length < 1) {
        return
    }
    let json = JSON.parse(data)
    // console.log(json.type, json.sheet_key, json.field_key)
    if (json.type == 'data') {
        // QueryBySheetAndField(json.sheet_key, [json.field_key], false).then(v => {
        //     if (v.error != "") {
        //         console.log(v.error)
        //         return
        //     }
        //     let path = "series[" + index + "].data"
        //     store.commit("set_echart_series_by_index", {
        //         id: chart_comp.value?.id,
        //         index: index,
        //         value: v.data.flat(),
        //         name: json.field_key 
        //     })
        //     store.commit("set_echart_data_map", {
        //         index: index,
        //         value: {
        //             sheetkey: json.sheet_key,
        //             fieldkey: json.field_key,
        //         }
        //     })
        //     updateOption("dataHandleDrop")
        // }).catch(reason => {
        //     console.log(reason)
        // })
        store.commit("set_echart_data_map", {
            index: index,
            value: {
                sheetkey: json.sheet_key,
                fieldkey: json.field_key,
            }
        })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
        EventsEmit("UIEmitSheetData", undefined)
    }
}

onMounted(() => {
    // coor_selected.value?.forEach(value => {
    //     SendFileToFront("templates/static/chartoptions/" + value.key + ".json").then(v => {
    //         if (v.startsWith("Error:")) {
    //             console.log(v)
    //             return
    //         }
    //         object_template_list.value[value.key] = JSON.parse(v)
    //     }).catch(reason => {
    //         console.log(reason)
    //         return
    //     })
    // })


})
onUnmounted(() => {

})


function showInfo(msg: string) {
    store.commit("set_show_html", {
        show: true,
        msg: msg
    })
}


function functionIconClick() {
    store.commit("set_ViShowFunction", { show: true, msg: store_current_value.value, path: get_current_path(selected_item.value) })
}



function coordragover(e: any) {
    e.preventDefault()
    console.log("coordragover")
    e.target.style.color = "red";
}

</script>

<style scoped>
.info-container {
    width: 10px;
    /* 设置内容容器的宽度 */
    white-space: nowrap;
    /* 防止内容换行 */
    overflow: hidden;
    /* 隐藏溢出部分的内容 */
    text-overflow: ellipsis;
    /* 显示省略号 */
}


/* 默认文字颜色 */
.hoverdiv p {
    color: white;
    /* 设置内容容器的宽度 */
    white-space: nowrap;
    /* 防止内容换行 */
    overflow: hidden;
    /* 隐藏溢出部分的内容 */
    text-overflow: ellipsis;
    /* 显示省略号 */
    background-color: #003355;
    padding-bottom: 10px;
    cursor: pointer;
}

/* 鼠标经过时改变文字颜色 */
.hoverdiv:hover p {
    color: #ff9900;
    /* 这里可以是任何你想要的颜色 */
    /* 设置内容容器的宽度 */
    white-space: nowrap;
    /* 防止内容换行 */
    overflow: hidden;
    /* 隐藏溢出部分的内容 */
    text-overflow: ellipsis;
    /* 显示省略号 */
    background-color: #003355;
    padding-bottom: 10px;
    cursor: pointer;
}



/* 默认文字颜色 */
.hoverenter p {
    color: green !important;
    cursor: pointer !important;
    padding-left: 10px !important;
}

/* 鼠标经过时改变文字颜色 */
.hoverenter:hover p {
    color: #ff9900 !important;
    cursor: pointer !important;
    padding-left: 10px !important;
}
</style>









