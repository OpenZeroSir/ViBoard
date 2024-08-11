<template>
    <div v-ripple="storeComputed.state.viRipple" :nni="props.nni" :style="{
        // height: props.outStyle.height * storeComputed.state.canvas.scale / 100 + 'px',
        // width: props.outStyle.width * storeComputed.state.canvas.scale / 100 + 'px',
        height: '100%',
        width: '100%',
    }">
        <!-- <div v-if="props.chartStyle?.style['bmap'] && storeComputed.state.viBaiduMapKey === ''">请在内部属性中配置百度地图密钥</div> -->
        <div :style="{
            height: '100%',
            width: '100%',
            // transform: `translate(-50%, -50%)  `,
            // position: 'absolute',
            // left: '50%',
            // top: '50%',
        }">
            <span
                v-show="props.chartStyle?.style['bmap'] && storeComputed.state.viBaiduMapKey == ''">请在内部属性中配置百度地图密钥</span>
            <div :nni="getEl" style="height: 100%;width: 100%;">

            </div>

        </div>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, computed, reactive, getCurrentInstance, onUnmounted, ref } from 'vue'
import { ViChartStyle, ViInnerStyle, ViOutStyle, store } from '../../../store';
import { storeDisplay } from "../../../store/display"
import { cloneDeep, random, set, unset } from 'lodash'
import ResizeObserver from 'resize-observer-polyfill';
import * as echarts from 'echarts'
import { getElement, sleep } from "../../../utils/func";
import { nextTick } from 'process';
import { EChartsOption, XAXisComponentOption } from 'echarts';
import { time } from 'console';
import { EventsOn, EventsOff, EventsEmit } from '../../../../wailsjs/runtime/runtime';
import { CheckMapAPILoaded, QueryBySheetAndField, RequestBaiduAPI } from '../../../../wailsjs/go/main/App';

const instance = getCurrentInstance();

let props = defineProps<{
    nni: string;
    name: string;
    innerStyle: ViInnerStyle;
    outStyle: ViOutStyle;
    chartStyle: ViChartStyle | undefined
}>()

let getEl = computed((id: string) => {
    return props.nni + "nni"
})

let mapinfo = ref('请在内部属性中配置百度地图密钥')

let chart: echarts.ECharts
let el: HTMLElement
let resizeObserver: ResizeObserver
function chartResize() {
    nextTick(() => {
        chart.resize()

    })
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
        updateOption()
    }

    document.head.append(scriptElement);
}

let storeComputed = computed(() => {
    // let fullpath = instance?.proxy?.$router.currentRoute.value.fullPath
    if (!props.outStyle.displayMode) {
        return store
    } else {
        return storeDisplay
    }
})

function updateOption() {
    try {
        if (props.chartStyle?.style) {

            let ss = props.chartStyle?.style.series as echarts.SeriesOption[]
            let obj = ss[2] as any
            chart?.setOption(props.chartStyle?.style)
            // chart.resize()
        }
    } catch (error) {
        chart?.dispose()
        chart = echarts.init(el)
    }
}

async function UIEmitSheetDataCallback(data: any) {
    if (data == undefined) {
        let coor = props.chartStyle?.chartCoorList
        if (coor != undefined) {
            for (const key in coor) {
                if (coor != undefined) {
                    let value = coor[key]
                    await QueryBySheetAndField(value.sheetkey, [value.fieldkey], props.outStyle.displayMode).then(v => {
                        if (v.error != "") {
                            return
                        }
                        storeDisplay.commit("set_echart_option_by_path", {
                            id: props.nni,
                            path: key,
                            value: v.data.flat()
                        })
                        // set(option, key, v.data)
                    }).catch(reason => {
                        console.log(reason)
                    })

                }
            }

        }
        let series = props.chartStyle?.chartSeriesList
        if (series != undefined) {
            for (const [index, value] of series.entries()) {
                // let fields: string[] = []
                // props.innerStyle.list.forEach(v => {
                //     fields.push(v.fieldkey)
                // })
                await QueryBySheetAndField(value.sheetkey, value.fieldkey, props.outStyle.displayMode).then(v => {
                    if (v.error != "") {
                        console.log(v.error)
                        return
                    }
                    let test = v.data[0] as any[]
                    let newdata = v.data
                    if (test.length === 1) {
                        newdata = v.data.flat()
                    }
                    if (props.outStyle.displayMode) {
                        storeDisplay.commit("set_echart_series_by_index", {
                            id: props.nni,
                            index: index,
                            value: newdata
                        })
                    } else {
                        store.commit("set_echart_series_by_index", {
                            id: props.nni,
                            index: index,
                            value: newdata
                        })
                        // nextTick(() => {
                        //     store.commit("snapshot")
                        // })
                    }

                    // let path = "series[" + index + "].data"
                    // set(option, path, v.data)
                }).catch(reason => {
                    console.log(reason)
                })

            }
        }
    } else {
        Object.keys(data).forEach(async sheetkey => {
            let coor = props.chartStyle?.chartCoorList
            if (coor != undefined) {
                for (const key in coor) {
                    if (coor != undefined) {
                        let value = coor[key]
                        if (value.sheetkey == sheetkey) {
                            await QueryBySheetAndField(value.sheetkey, [value.fieldkey], props.outStyle.displayMode).then(v => {
                                if (v.error != "") {
                                    console.log(v.error)
                                    return
                                }
                                storeDisplay.commit("set_echart_option_by_path", {
                                    id: props.nni,
                                    path: key,
                                    value: v.data.flat()
                                })
                                // set(option, key, v.data)
                            }).catch(reason => {
                                console.log(reason)
                            })
                        }
                    }
                }

            }
            let series = props.chartStyle?.chartSeriesList
            if (series != undefined) {
                for (const [index, value] of series.entries()) {
                    if (value.sheetkey == sheetkey) {
                        await QueryBySheetAndField(value.sheetkey, value.fieldkey, props.outStyle.displayMode).then(v => {
                            if (v.error != "") {
                                console.log(v.error)
                                return
                            }
                            storeDisplay.commit("set_echart_series_by_index", {
                                id: props.nni,
                                index: index,
                                value: v.data.flat()
                            })
                            // let path = "series[" + index + "].data"
                            // set(option, path, v.data)
                        }).catch(reason => {
                            console.log(reason)
                        })
                    }
                }
            }
        })
    }

    sleep(200).then(() => {
        nextTick(() => {
            updateOption()
        })
    })

}




onMounted(async () => {
    console.log("storeComputed.value.state.viBaiduMapKey",storeComputed.value.state.viBaiduMapKey)
    if (props.name === "ViChart") {
        await CheckMapAPILoaded().then(async status => {
            if (!status) {
                if (storeComputed.value.state.viBaiduMapKey != "") {
                    await RequestBaiduAPI(storeComputed.value.state.viBaiduMapKey).then(url => {
                        if (url.indexOf("Error:") > -1) {
                            console.log(url)
                            return
                        }
                        loadBaiduMap(url)
                        console.log("map api loaded")
                    })
                }

            }
        }).catch(reason => {
            console.log("百度地图加载出错",reason)
            mapinfo.value = reason
        })
        let url = storeComputed.value.state.viViews[storeComputed.value.state.viCtrlIndex].id + "NNI" + props.nni.substring(3)
        await fetch(url).then(v => {
            // echarts.registerMap()
            v.json().then(value => {
                echarts.registerMap(props.nni.substring(3), value)
            })

        }).catch(reason => {
            console.log("GeoJson地图加载出错", reason)
        })

    }
    el = getElement(props.nni.substring(3) + "nni") as HTMLElement
    if (el != null) {
        chart = echarts.init(el)
        let option = props.chartStyle?.style
        if (option != undefined) {
            // try {
            //     console.log("init chart option",option)
            //     chart.setOption(option)
            // } catch (error) {
            //     console.log("chart mount chart.setOption(option)",option,error)
            //     chart?.dispose()
            //     chart = echarts.init(el)
            // }
            // updateOption()
        }
        // el.addEventListener('resize', chartResize)
        resizeObserver = new ResizeObserver(() => {
            chartResize()
        });
        resizeObserver.observe(el);
    }

    EventsOn("UIEmitSheetData", UIEmitSheetDataCallback)

    EventsEmit("UIEmitSheetData", undefined)

    // if (props.innerStyle.list.length > 0) {
    //     EventsEmit()
    // }

    // let x = chart.getOption()
    // console.log(x)
    // let x = chart.getOption()
    // let b:XAXisComponentOption={}
    // let option = {
    //     xAxis:b
    // }
    // chart.setOption(option)
})

onUnmounted(() => {
    // el.removeEventListener("resize", chartResize)
    resizeObserver.disconnect()
    chart.dispose()
    EventsOff("UIEmitSheetData")
})



</script>