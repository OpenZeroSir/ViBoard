<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount, nextTick, reactive, computed, getCurrentInstance } from 'vue'
import MainMenu from './MainMenu.vue'
import LeftMenu from './LeftMenu.vue'
import RightMenu from './RightMenu.vue'
import Editor from "./editor/index.vue"
import Foot from './Foot.vue'
import Control from './Control.vue'
import InfoDialog from './utils/InfoDialog.vue'

// import { VVirtualScroll } from 'vuetify/labs/VVirtualScroll'
import { mapState } from 'vuex'
import { store, ViChartMainType, ViComponent, ViView } from "../../store"
import { viFormat, viDeFormat, getNewValue, WGS84ToBD09, nanoidEx } from '../../utils/func'
import { nanoid } from 'nanoid'
import { GetSocketPort, GetSocketToken, RequestBaiduAPI } from "../../../wailsjs/go/main/App"
// import { EventsOn } from '../../../wailsjs/runtime/runtime'
import { publicStore } from '../../store/public'
import { Components } from "../custom/index"
import { cloneDeep } from 'lodash'
// import { ChartIcons } from "../../utils/charticon"
import LoaderDialog from './utils/LoaderDialog.vue'
import AskDialog from './utils/AskDialog.vue'
import { EChartsOption } from 'echarts'
import HtmlInfo from './utils/HtmlInfo.vue'
import FunctionDialog from './utils/FunctionDialog.vue'
import { CHINA } from '../../utils/china'
import * as echarts from 'echarts'
import Help from './utils/Help.vue'
import CopyRight from './utils/CopyRight.vue'
import ChartHelp from './utils/ChartHelp.vue'
import { Quit } from '../../../wailsjs/runtime/runtime'
import { useConfirm, useSnackbar } from 'vuetify-use-dialog'
import UpdateInfo from './utils/UpdateInfo.vue'
import LicenseDialog from './utils/LicenseDialog.vue'


const createConfirm = useConfirm()

const ViContentRef = ref<HTMLElement | null>(null);



let view_height = ref(500)
let editor_scale_min = ref(0)
let editor_scale = computed({
    get: () => store.state.canvas.scale,
    set: (val) => {
        store.commit("set_editor_scale", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
        updateScale()
    }
});



function updateScale(index?: number) {
    if (index == undefined) {
        store.state.viViews[store.state.viCtrlIndex].components.forEach(v => {


            let style = { ...v.style }
            style.top = getNewValue(style.top)
            style.left = getNewValue(style.left)
            style.width = getNewValue(style.width)
            style.height = getNewValue(style.height)
            store.commit("set_out_style", {
                style: reactive(style), addr: {
                    viewid: store.state.viCtrlIndex,
                    comid: v.id
                }
            })
            // nextTick(() => {
            //     store.commit("snapshot")
            // })

        })
    }
}



let main_control: HTMLElement | null = null;
let main_center: HTMLElement | null = null;
let main_container: HTMLElement | null = null;
let main_vicontent: HTMLElement | null = null;
let EditorID: HTMLElement | null;


function reSizeed() {
    nextTick(() => {
        let main_center_padding = Number(window.getComputedStyle(
            main_center!).padding.replace("px", ""));
        let main_container_padding = Number(window.getComputedStyle(
            main_container!).padding.replace("px", ""));
        let offset_vmain_center = window.innerHeight -
            main_center!.offsetHeight - main_center!.offsetTop;
        let new_view_height = main_center!.offsetHeight +
            offset_vmain_center - main_control!.offsetHeight - (
                2 * main_center_padding + main_container_padding / 2 / 2
            );
        view_height.value = new_view_height;
        // if(new_view_height> EditorID!.clientHeight){
        //     view_height.value = EditorID!.clientHeight+12
        // }

        // EditorID = document.getElementById('EditorID');
    })
}

function initWebSocket() {
    GetSocketPort().then((result: number) => {
        // setTimeout(() => {
        const socket = new WebSocket('ws://192.168.50.40:' + result.toString() + '/nni');
        socket.onopen = () => {
            // 将字节数组发送给后端
            const data = { token: publicStore.state.wstoken }; // 要发送的 JSON 数据
            const jsonData = JSON.stringify(data); // 将 JSON 对象转换为字符串
            socket.send(jsonData); // 发送 JSON 数据到服务器
            publicStore.commit("set_wssock", socket)
        };
        socket.onerror = (e: any) => {
            console.log("socket error", e)
        }
        // }, 3000);

    }).catch((reason) => {
        console.log("前端初始化WebSocket连接失败", reason)
    })
}

async function exit() {
    const isConfirmed = await createConfirm({
        title: "退出程序",
        content: "确定要退出程序吗？",
        confirmationText: "确定",
        cancellationText: "取消",
        dialogProps: {
            width: 'auto',
        },
        confirmationButtonProps: {
            color: 'red',
        },
    })
    if (isConfirmed) {
        Quit()
    }
}

onMounted(() => {

    main_control = document.getElementById('main_control');
    main_center = document.getElementById('main_center');
    main_container = document.getElementById('main_container');
    main_vicontent = document.getElementById('vicontent');
    EditorID = document.getElementById('EditorID');
    store.commit("init_editor", EditorID)
    // nextTick(() => {
    //     store.commit("snapshot")
    // })

    reSizeed();
    // nextTick(() => {
    //     let main_center_padding = Number(window.getComputedStyle(
    //         main_center!).padding.replace("px", ""));
    //     editor_scale_min.value = ((main_center!.offsetWidth - 2 * main_center_padding) /
    //         store.state.canvas.width) * 100
    //     store.commit("set_editor_scale", editor_scale_min.value);
    // })

    window.addEventListener('keydown', handleKeyDown);
    window.addEventListener('keyup', handleKeyUp);
    window.addEventListener('beforeunload', exit)

    if (store.state.viViews.length == 0) {
        store.commit("add_viView")
        nextTick(() => {
            store.commit("snapshot")
        })
    }
    // let functionstr = "function AA(a,b) { return a + b }"
    // function buildDynamicFunction(parameterNames: string[]): Function {
    //     return new Function(...parameterNames, 'return ' + functionstr)();
    // }

    // // 使用动态构建的函数，传递参数
    // const dynamicFunction = buildDynamicFunction(['param1', 'param2']);

    // // 调用动态构建的函数，并传递参数
    // console.log(dynamicFunction('Hello', 'world'))

    // // 使用正则表达式匹配参数
    // const parameterPattern = /\(([^)]*)\)/;
    // const match = functionstr.match(parameterPattern);

    // if (match && match[1]) {
    //     // 分割参数字符串，并去除可能存在的空格
    //     const parameters = match[1].split(',').map(param => param.trim());

    //     console.log("函数参数:", parameters);
    // } else {
    //     console.log("未找到函数参数");
    // }

    // RequestBaiduAPI("xxx").then(value=>{
    //     console.log("baidu api result",value)
    //     var scriptElement = document.createElement('script');
    //     scriptElement.text = value;
    //     document.head.append(scriptElement);
    // })
    echarts.registerMap("map", CHINA as any)

})

// function testEvent(aa:string){
//     console.log("testEvent",aa)
// }

onBeforeUnmount(() => {
    window.removeEventListener('keydown', handleKeyDown);
    window.removeEventListener('keyup', handleKeyUp);
    window.removeEventListener('beforeunload', exit)
})

window.onresize = () => {
    // console.log("onresize")
    reSizeed();
}


function dragover(event: any) {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'copy'
}
function handleDrop(event: any) {

    // console.log(event,EditorID?.getBoundingClientRect())
    event.preventDefault()
    event.stopPropagation()
    let dropdata: string = event.dataTransfer.getData("application/json");
    if (dropdata.length < 1) {
        return
    }
    let json = JSON.parse(dropdata)

    let rect = EditorID?.getBoundingClientRect()
    let comp: ViComponent | null = null
    if (json.type == "icon") {
        comp = reactive(cloneDeep(Components.ViIcon))
        comp.name = json.value
        comp.component = "ViSVGIcon"

        // } else if (data.startsWith("md:")) {
        //     comp = reactive(cloneDeep(Components.ViIcon))
        //     comp.name = data.substring(3)
        //     comp.component = "ViFontIcon"
    } else if (json.type == "chart") {
        let chartid = json.value
        comp = reactive(cloneDeep(Components.ViChart))
        comp.name = "ViChart"
        comp.component = "ViChart"
        let chartOption: EChartsOption
        chartOption = {
            xAxis: [{
                type: 'category',
                data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
            }],
            yAxis: [{
                type: 'value'
            }],
            grid: {
                top: '20%',
                bottom: '20%',
            },
            // tooltip: {
            //     trigger: 'item',
            //     axisPointer: {
            //         type:'cross'
            //     }
            // },
        };
        if (chartid == "1E51ACC7-9E58-4704-BF1B-1142A07BAE3E") {
            //基础折线图
            // if (chartOption.series != undefined) {
            chartOption.series = [
                {
                    data: [150, 230, 224, 218, 135, 147, 260],
                    type: 'line'
                }
            ]
            // }
        } else if (chartid == "9cb49149-c520-4f53-a69e-4423e963f3d2") {
            //堆积折线图
            chartOption.series = [
                {
                    type: 'line',
                    stack: 'Total',
                    data: [120, 132, 101, 134, 90, 230, 210]
                },
                {
                    type: 'line',
                    stack: 'Total',
                    data: [220, 182, 191, 234, 290, 330, 310]
                }
            ]
        } else if (chartid == "9b1a3b1a-78ed-4308-9d7c-aae1e731248d") {
            //基础面积图
            chartOption.series = [
                {
                    data: [820, 932, 901, 934, 1290, 1330, 1320],
                    type: 'line',
                    areaStyle: {}
                }
            ]
        } else if (chartid == "586e65fd-aed4-42ab-a797-271f578d7f5e") {
            //堆积面积图
            chartOption.series = [
                {
                    type: 'line',
                    stack: 'Total',
                    areaStyle: {},
                    data: [120, 132, 101, 134, 90, 230, 210]
                },
                {
                    type: 'line',
                    stack: 'Total',
                    areaStyle: {},
                    data: [220, 182, 191, 234, 290, 330, 310]
                },
            ]
        } else if (chartid == "f58cb5f2-ede3-4923-8fd2-65e5ef28641d") {
            //基础柱状图
            chartOption.series = [
                {
                    type: 'bar',
                    data: [120, 132, 101, 134, 90, 230, 210]
                }
            ]
        } else if (chartid == "1885c347-816b-4d9f-8f23-e80ad3e7c5b2") {
            //堆积柱状图
            chartOption.series = [
                {
                    data: [120, 200, 150, 80, 70, 110, 130],
                    type: 'bar',
                    stack: 'total',
                },
                {
                    data: [220, 200, 100, 80, 70, 100, 150],
                    type: 'bar',
                    stack: 'total',
                }
            ]
        } else if (chartid == "db94cd1c-f811-494f-9803-fc0a57d5df0b") {
            //水平柱状图
            chartOption.series = [
                {
                    type: 'bar',
                    data: [220, 182, 191, 234, 290, 330, 310]
                }
            ]
            chartOption.xAxis = [{
                type: 'value',
            }]
            chartOption.yAxis = [{
                type: 'category',
                data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
            }]
        } else if (chartid == "5c1631ea-7831-4788-a132-c866f20c7548") {
            //水平柱状图
            chartOption.series = [
                {
                    type: 'bar',
                    stack: 'total',
                    data: [220, 182, 191, 234, 290, 330, 310]
                },
                {
                    type: 'bar',
                    stack: 'total',
                    data: [220, 182, 191, 234, 290, 330, 310]
                }
            ]
            chartOption.xAxis = [{
                type: 'value',
            }]
            chartOption.yAxis = [{
                type: 'category',
                data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
            }]
        } else if (chartid == "00de94ff-ecda-4231-acde-f49e8777933d") {
            chartOption = {
                series: [
                    {
                        type: 'pie',
                        data: [
                            { value: 1048, name: 'Search Engine' },
                            { value: 735, name: 'Direct' },
                            { value: 580, name: 'Email' },
                            { value: 484, name: 'Union Ads' },
                            { value: 300, name: 'Video Ads' }
                        ]
                    }
                ]
            }
        }
        else if (chartid == "ba2fafed-516f-48ce-b5e4-0925196edbf7") {
            chartOption = {
                series: [
                    {
                        type: 'pie',
                        radius: ['40%', '70%'],
                        // itemStyle: {
                        //     borderRadius: 10,
                        //     borderColor: '#fff',
                        //     borderWidth: 2
                        // },
                        data: [
                            { value: 1048, name: 'Search Engine' },
                            { value: 735, name: 'Direct' },
                            { value: 580, name: 'Email' },
                            { value: 484, name: 'Union Ads' },
                            { value: 300, name: 'Video Ads' }
                        ]
                    }
                ]
            }
        }
        // else if (chartid == "97854320-17ad-4c48-9dee-a560ecad171c") {
        //     chartOption = {
        //         series: [
        //             {
        //                 name: 'Access From',
        //                 type: 'pie',
        //                 radius: ['40%', '70%'],
        //                 center: ['50%', '70%'],
        //                 // adjust the start angle
        //                 startAngle: 180,
        //                 data: [
        //                     { value: 1048, name: 'Search Engine' },
        //                     { value: 735, name: 'Direct' },
        //                     { value: 580, name: 'Email' },
        //                     { value: 484, name: 'Union Ads' },
        //                     { value: 300, name: 'Video Ads' },
        //                     {
        //                         // make an record to fill the bottom 50%
        //                         value: 1048 + 735 + 580 + 484 + 300,
        //                         itemStyle: {
        //                             // stop the chart from rendering this piece
        //                             color: 'none',
        //                             decal: {
        //                                 symbol: 'none'
        //                             }
        //                         },
        //                         label: {
        //                             show: false
        //                         }
        //                     }
        //                 ]
        //             }
        //         ]
        //     }
        // } 
        else if (chartid == "3bb52213-654f-4bfc-b749-1bb7b698036a") {
            chartOption = {
                series: [
                    {
                        type: 'pie',
                        radius: '70%',
                        center: ['50%', '50%'],
                        roseType: 'area',
                        data: [
                            { value: 40, name: 'rose 1' },
                            { value: 38, name: 'rose 2' },
                            { value: 32, name: 'rose 3' },
                            { value: 30, name: 'rose 4' },
                            { value: 28, name: 'rose 5' },
                            { value: 26, name: 'rose 6' },
                            { value: 22, name: 'rose 7' },
                            { value: 18, name: 'rose 8' }
                        ]
                    }
                ]
            }
        } else if (chartid == "b8563795-6f32-4f8c-a14a-06b93fa93435") {
            chartOption = {
                xAxis: [{
                    type: 'category'
                }],
                yAxis: [{
                    type: 'value'
                }],
                grid: {
                    top: '20%',
                    bottom: '20%',
                },
                series: [
                    {
                        // symbolSize: 10,
                        data: [
                            [10.0, 8.04, 15],
                            [8.07, 6.95, 1],
                            [13.0, 7.58, 3],
                            [9.05, 8.81, 18],
                            [11.0, 8.33, 22]
                        ],
                        type: 'scatter'
                    }
                ]
            }
        } else if (chartid == "0c79ae07-91ef-4a17-adaa-6e70c7fbf152") {
            chartOption = {
                xAxis: [{
                    type: 'category',
                    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
                }],
                yAxis: [{
                    type: 'value'
                }],
                grid: {
                    top: '20%',
                    bottom: '20%',
                },
                series: [
                    {
                        type: 'candlestick',
                        data: [
                            [20, 34, 10, 50],
                            [40, 35, 30, 50],
                            [31, 38, 33, 44],
                            [38, 15, 5, 42],
                            [25, 35, 30, 19],
                            [11, 23, 50, 44],
                            [30, 14, 5, 39]
                        ]
                    }
                ]
            }
        } else if (chartid == "d6f4f906-9539-4a28-b03e-486de2f0876c") {
            chartOption = {

                series: [
                    {
                        progress: {
                            show: true
                        },
                        type: 'gauge',
                        data: [
                            50
                        ]
                    }
                ]
            }
        } else if (chartid == "c9ea4933-c5a0-4432-9156-ca0cb7a8fd61") {
            chartOption = {
                xAxis: [{
                    type: 'category',
                    data: ['a', 'b', 'c']
                }],
                yAxis: [{
                    type: 'category',
                    data: [1, 2, 3],
                }],
                grid: {
                    top: '20%',
                    bottom: '20%',
                },
                visualMap: {
                    min: 0,
                    max: 10,
                    calculable: true,
                    orient: 'horizontal',
                    left: 'center',
                },
                series: [
                    {
                        type: 'heatmap',
                        data: [[0, 0, 5], [0, 1, 1], [0, 2, 0], [1, 0, 1], [1, 1, 3], [1, 2, 6], [2, 0, 2], [2, 1, 4], [2, 2, 5]],

                    }
                ]
            }
        } else if (chartid == "de5415ca-2c0d-4700-87ce-12ef2955a344") {
            chartOption = {
                radar: {
                    indicator: [
                        { name: 'A' },
                        { name: 'B' },
                        { name: 'C' },
                        { name: 'D' },
                        { name: 'E' },
                        { name: 'F' }
                    ]
                },
                series: [
                    {
                        type: 'radar',
                        data: [
                            {
                                value: [4200, 3000, 20000, 35000, 50000, 18000],
                            }
                        ]
                    }
                ]
            }
        } else if (chartid == "2caf407d-c219-4780-a66a-4f53ddd52421") {
            chartOption = {

                series: [
                    {
                        type: 'funnel',
                        data: [
                            { value: 60, },
                            { value: 40, name: 'Inquiry' },
                            { value: 20, name: 'Order' },
                            { value: 80, name: 'Click' },
                            { value: 100, name: 'Show' }
                        ]
                    }
                ]
            }
        } else if (chartid == "9d6749b8-6ac3-4f88-8cc8-e78952c67239") {
            chartOption = {
                xAxis: [{
                    type: 'category',
                    data: ['Category 1', 'Category 2', 'Category 3']
                }],
                yAxis: [{
                    type: 'value',
                }],
                grid: {
                    top: '20%',
                    bottom: '20%',
                },
                series: [
                    {
                        type: 'boxplot',
                        /*
                            a：最小值
                            b：下四分位数（第25百分位数）
                            c：中位数（第50百分位数）
                            d：上四分位数（第75百分位数）
                            e：最大值
                        */
                        data: [
                            [850, 960, 940, 980, 1000],  // 类别1的数据
                            [890, 900, 920, 930, 950],   // 类别2的数据
                            [850, 880, 870, 910, 950],   // 类别3的数据
                        ],
                    }
                ]
            }
        } else if (chartid == "8c798693-672f-4a5a-8cd1-107151ada9ba") {
            chartOption = {
                singleAxis: {
                    type: 'category',
                    data: ['a', 'b', 'c', 'd', 'e']
                },
                series: [{
                    type: 'scatter',
                    coordinateSystem: 'singleAxis',
                    data: [
                        ['a', 10],
                        ['b', 20],
                        ['c', 30],
                        ['d', 25],
                        ['e', 5],
                    ],
                    symbolSize: function (dataItem, params) {
                        return dataItem[1];
                    },
                    itemStyle: {
                        color: function (params: any) { return params['value'][1] as number > 10 ? 'red' : 'green' }
                    }
                }]
            };
        } else if (chartid == "199c00b9-d3b9-4dd4-bfc0-8344089b1268") {
            chartOption = {
                parallelAxis: [
                    { dim: 0, name: 'Price' },
                    { dim: 1, name: 'Net Weight' },
                    { dim: 2, name: 'Amount' },
                ],
                series: [{
                    type: 'parallel',
                    data: [
                        [12.99, 100, 82],
                        [9.99, 80, 77],
                        [5, 120, 60]
                    ]
                }]
            };
        } else if (chartid == "1782ebf9-589d-426c-b37f-9c4f5f8e7274") {
            chartOption = {
                polar: {
                    radius: [30, '80%']
                },
                radiusAxis: {
                },
                angleAxis: {
                    type: 'category',
                    data: ['a', 'b', 'c', 'd'],
                },
                series: [{
                    type: 'bar',
                    data: [2, 1.2, 2.4, 3.6],
                    coordinateSystem: 'polar',
                }],
            };
        } else if (chartid == "643a5985-d595-4163-955d-b4cd1b212083") {
            chartOption = {
                polar: {
                    radius: [30, '80%']
                },
                angleAxis: {
                },
                radiusAxis: {
                    type: 'category',
                    data: ['a', 'b', 'c', 'd']
                },
                series: [{
                    type: 'bar',
                    data: [2, 1.2, 2.4, 3.6],
                    coordinateSystem: 'polar'
                }]
            };
        } else if (chartid == "3a2c1a57-66c8-4ae1-80b2-7866d68e415e") {
            chartOption = {
                angleAxis: {
                    type: 'category',
                    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
                },
                radiusAxis: {},
                polar: {},
                series: [
                    {
                        type: 'bar',
                        data: [1, 2, 3, 4, 3, 5, 1],
                        coordinateSystem: 'polar',
                        name: 'A',
                        stack: 'a',
                        emphasis: {
                            focus: 'series'
                        }
                    },
                    {
                        type: 'bar',
                        data: [2, 4, 6, 1, 3, 2, 1],
                        coordinateSystem: 'polar',
                        name: 'B',
                        stack: 'a',
                        emphasis: {
                            focus: 'series'
                        }
                    },
                    {
                        type: 'bar',
                        data: [1, 2, 3, 4, 1, 2, 5],
                        coordinateSystem: 'polar',
                        name: 'C',
                        stack: 'a',
                        emphasis: {
                            focus: 'series'
                        }
                    }
                ],
            };
        } else if (chartid == "8a5dcc6b-0aa4-40f6-8323-d6ee3e18c59c") {
            chartOption = {
                angleAxis: {},
                radiusAxis: {
                    type: 'category',
                    data: ['Mon', 'Tue', 'Wed', 'Thu'],
                    z: 10
                },
                polar: {},
                series: [
                    {
                        type: 'bar',
                        data: [1, 2, 3, 4],
                        coordinateSystem: 'polar',
                        name: 'A',
                        stack: 'a',
                        emphasis: {
                            focus: 'series'
                        }
                    },
                    {
                        type: 'bar',
                        data: [2, 4, 6, 8],
                        coordinateSystem: 'polar',
                        name: 'B',
                        stack: 'a',
                    },
                    {
                        type: 'bar',
                        data: [1, 2, 3, 4],
                        coordinateSystem: 'polar',
                        name: 'C',
                        stack: 'a',
                    }
                ],
            };
        } else if (chartid == "f0dd12cd-0bcb-47ab-a94e-ade76cbb56af") {
            chartOption = {

                polar: {},
                angleAxis: {
                    type: 'category',
                    data: [
                        '12a', '1a', '2a', '3a', '4a', '5a', '6a',
                        '7a', '8a', '9a', '10a', '11a',
                        '12p', '1p', '2p', '3p', '4p', '5p',
                        '6p', '7p', '8p', '9p', '10p', '11p'
                    ],
                    splitLine: {
                        show: true
                    },
                },
                radiusAxis: {
                    type: 'category',
                    data: [
                        'Saturday', 'Friday', 'Thursday',
                        'Wednesday', 'Tuesday', 'Monday', 'Sunday'
                    ],
                },
                series: [
                    {
                        type: 'scatter',
                        coordinateSystem: 'polar',
                        symbolSize: (value, params) => value[2],
                        data: [[0, 0, 5], [0, 1, 1], [0, 2, 0], [0, 3, 0], [0, 4, 0], [0, 5, 0], [0, 6, 0], [0, 7, 0], [0, 8, 0], [0, 9, 0], [0, 10, 0], [0, 11, 2], [0, 12, 4], [0, 13, 1], [0, 14, 1], [0, 15, 3], [0, 16, 4], [0, 17, 6], [0, 18, 4], [0, 19, 4], [0, 20, 3], [0, 21, 3], [0, 22, 2], [0, 23, 5], [1, 0, 7], [1, 1, 0], [1, 2, 0], [1, 3, 0], [1, 4, 0], [1, 5, 0], [1, 6, 0], [1, 7, 0], [1, 8, 0], [1, 9, 0], [1, 10, 5], [1, 11, 2], [1, 12, 2], [1, 13, 6], [1, 14, 9], [1, 15, 11], [1, 16, 6], [1, 17, 7], [1, 18, 8], [1, 19, 12], [1, 20, 5], [1, 21, 5], [1, 22, 7], [1, 23, 2], [2, 0, 1], [2, 1, 1], [2, 2, 0], [2, 3, 0], [2, 4, 0], [2, 5, 0], [2, 6, 0], [2, 7, 0], [2, 8, 0], [2, 9, 0], [2, 10, 3], [2, 11, 2], [2, 12, 1], [2, 13, 9], [2, 14, 8], [2, 15, 10], [2, 16, 6], [2, 17, 5], [2, 18, 5], [2, 19, 5], [2, 20, 7], [2, 21, 4], [2, 22, 2], [2, 23, 4], [3, 0, 7], [3, 1, 3], [3, 2, 0], [3, 3, 0], [3, 4, 0], [3, 5, 0], [3, 6, 0], [3, 7, 0], [3, 8, 1], [3, 9, 0], [3, 10, 5], [3, 11, 4], [3, 12, 7], [3, 13, 14], [3, 14, 13], [3, 15, 12], [3, 16, 9], [3, 17, 5], [3, 18, 5], [3, 19, 10], [3, 20, 6], [3, 21, 4], [3, 22, 4], [3, 23, 1], [4, 0, 1], [4, 1, 3], [4, 2, 0], [4, 3, 0], [4, 4, 0], [4, 5, 1], [4, 6, 0], [4, 7, 0], [4, 8, 0], [4, 9, 2], [4, 10, 4], [4, 11, 4], [4, 12, 2], [4, 13, 4], [4, 14, 4], [4, 15, 14], [4, 16, 12], [4, 17, 1], [4, 18, 8], [4, 19, 5], [4, 20, 3], [4, 21, 7], [4, 22, 3], [4, 23, 0], [5, 0, 2], [5, 1, 1], [5, 2, 0], [5, 3, 3], [5, 4, 0], [5, 5, 0], [5, 6, 0], [5, 7, 0], [5, 8, 2], [5, 9, 0], [5, 10, 4], [5, 11, 1], [5, 12, 5], [5, 13, 10], [5, 14, 5], [5, 15, 7], [5, 16, 11], [5, 17, 6], [5, 18, 0], [5, 19, 5], [5, 20, 3], [5, 21, 4], [5, 22, 2], [5, 23, 0], [6, 0, 1], [6, 1, 0], [6, 2, 0], [6, 3, 0], [6, 4, 0], [6, 5, 0], [6, 6, 0], [6, 7, 0], [6, 8, 0], [6, 9, 0], [6, 10, 1], [6, 11, 0], [6, 12, 2], [6, 13, 1], [6, 14, 3], [6, 15, 4], [6, 16, 0], [6, 17, 0], [6, 18, 0], [6, 19, 0], [6, 20, 1], [6, 21, 2], [6, 22, 2], [6, 23, 6]],
                    }
                ]
            };
        } else if (chartid == "48088e85-a412-4140-9bd9-56b2f7f93154") {
            chartOption = {

                singleAxis: {
                    type: 'category',
                    data: ['1', '2', '3', '4', '5', '6']
                },
                series: [
                    {
                        type: 'themeRiver',
                        data: [
                            ['1', 10, 'b'],
                            ['2', 15, 'b'],
                            ['3', 35, 'b'],
                            ['4', 38, 'b'],
                            ['5', 22, 'b'],
                            ['1', 8, 'a'],
                            ['2', 13, 'a'],
                            ['3', 2, 'a'],
                            ['4', 17, 'a'],
                            ['5', 6, 'a'],
                        ]
                    }
                ]
            };
        } else if (chartid == "71a7a7a2-99c3-4575-bc03-4c1706a5f39b") {
            chartOption = {
                bmap: {
                    center: [116.46, 39.92],
                    zoom: 10,
                    roam: true,
                },
                series: [
                    {
                        type: 'lines',
                        coordinateSystem: 'bmap',
                        polyline: true,
                        // silent: true,
                        lineStyle: {
                            color: 'red', width: 3,
                        },
                        data: [
                            {
                                name: "线路1",
                                coords: [[116.4696, 39.9173], [116.4696, 39.9231], [116.478, 39.9228], [116.5187, 39.923], [116.5224, 39.9237], [116.5441, 39.9234], [116.5567, 39.9238], [116.5754, 39.9239], [116.5902, 39.9251], [116.6055, 39.9258], [116.6104, 39.924], [116.6133, 39.9238], [116.6242, 39.9261], [116.6153, 39.9394], [116.6048, 39.9542], [116.603, 39.9616], [116.5997, 39.9699], [116.5967, 39.973], [116.5933, 39.9785], [116.6091, 39.9814], [116.6161, 39.9855], [116.625, 39.9883], [116.6265, 39.9874], [116.6281, 39.981], [116.6316, 39.9637]]

                            }, {
                                name: "线路2",
                                coords: [[116.4383, 40.1471], [116.4372, 40.1458], [116.4373, 40.1409], [116.4347, 40.1395], [116.4446, 40.1225], [116.445, 40.1189], [116.4537, 40.1187], [116.4553, 40.1046], [116.4551, 40.1031], [116.4504, 40.1025], [116.4336, 40.1016], [116.4334, 40.1038], [116.426, 40.1034], [116.4122, 40.1044], [116.4134, 40.0892], [116.4143, 40.0837], [116.4126, 40.0726], [116.4139, 40.055], [116.4119, 40.0512], [116.412, 40.0455], [116.4151, 40.0401], [116.4179, 40.0316], [116.4174, 40.019], [116.4161, 40.0128], [116.4162, 40.0094], [116.4078, 40.0095], [116.4075, 39.9877], [116.409, 39.9883], [116.4168, 39.9885], [116.4172, 39.9833], [116.4242, 39.9834], [116.4249, 39.9692], [116.4348, 39.9694], [116.4369, 39.9689], [116.4598, 39.9539], [116.4614, 39.9516], [116.4614, 39.9336], [116.4782, 39.9337], [116.4778, 39.908], [116.486, 39.908], [116.4893, 39.9058], [116.4971, 39.9078]]
                            }
                        ],
                    }
                ]
            }
        } else if (chartid == "c163399d-06b8-436a-b620-070a87947201") {
            // let bd09 = WGS84ToBD09(102.73286745432786, 25.03032166639342)
            // function renderItem(params: any, api: any) {
            //     var coords = [
            //         [100.03962, 27.24345],
            //         [100.039895980774, 27.2441277564757],
            //         [100.039882704364, 27.2441319554595],
            //         [100.039869348122, 27.2441359469013],
            //         [100.039855916101, 27.2441397295809],
            //         [100.039842412377, 27.2441433023412],
            //         [100.039828841047, 27.2441466640892],
            //         [100.039815206231, 27.2441498137956],
            //         [100.039801512071, 27.2441527504963],
            //         [100.039787762724, 27.2441554732914],
            //         [100.039773962367, 27.2441579813469],
            //         [100.039760115194, 27.2441602738939],
            //         [100.039746225412, 27.2441623502296],
            //         [100.039732297243, 27.2441642097173],
            //         [100.03971833492, 27.2441658517867],
            //         [100.03970434269, 27.2441672759339],
            //         [100.039690324806, 27.2441684817221],
            //         [100.039676285532, 27.2441694687811],
            //         [100.039662229137, 27.2441702368081],
            //         [100.039648159896, 27.2441707855672],
            //         [100.03963408209, 27.2441711148901],
            //         [100.03962, 27.2441712246757],
            //         [100.03960591791, 27.2441711148901],
            //         [100.039591840104, 27.2441707855672],
            //         [100.039577770863, 27.2441702368081],
            //         [100.039563714468, 27.2441694687811],
            //         [100.039549675194, 27.2441684817221],
            //         [100.03953565731, 27.2441672759339],
            //         [100.03952166508, 27.2441658517867],
            //         [100.039507702757, 27.2441642097173],
            //         [100.039493774588, 27.2441623502296],
            //         [100.039479884806, 27.2441602738939],
            //         [100.039466037633, 27.2441579813469],
            //         [100.039452237276, 27.2441554732914],
            //         [100.039438487929, 27.2441527504963],
            //         [100.039424793769, 27.2441498137956],
            //         [100.039411158953, 27.2441466640892],
            //         [100.039397587623, 27.2441433023412],
            //         [100.039384083899, 27.2441397295809],
            //         [100.039370651878, 27.2441359469013],
            //         [100.039357295636, 27.2441319554595],
            //         [100.039344019226, 27.2441277564757],
            //         [100.03962, 27.24345]

            //     ];

            //     var points = [];
            //     for (var i = 0; i < coords.length; i++) {
            //         points.push(api.coord(coords[i]));
            //     }
            //     console.log("renderItem1", points)
            //     var color = api.visual('color');
            //     return {
            //         type: 'polygon',
            //         shape: {
            //             points: echarts.graphic.clipPointsByRect(points, {
            //                 x: params.coordSys.x,
            //                 y: params.coordSys.y,
            //                 width: params.coordSys.width,
            //                 height: params.coordSys.height
            //             })
            //         },
            //         style: api.style({
            //             fill: color,
            //             stroke: echarts.color.lift(color, 5)
            //         })
            //     };
            // }
            // var coords = [
            //     [100.03962, 27.24345],
            //     [100.039895980774, 27.2441277564757],
            //     [100.039882704364, 27.2441319554595],
            //     [100.039869348122, 27.2441359469013],
            //     [100.039855916101, 27.2441397295809],
            //     [100.039842412377, 27.2441433023412],
            //     [100.039828841047, 27.2441466640892],
            //     [100.039815206231, 27.2441498137956],
            //     [100.039801512071, 27.2441527504963],
            //     [100.039787762724, 27.2441554732914],
            //     [100.039773962367, 27.2441579813469],
            //     [100.039760115194, 27.2441602738939],
            //     [100.039746225412, 27.2441623502296],
            //     [100.039732297243, 27.2441642097173],
            //     [100.03971833492, 27.2441658517867],
            //     [100.03970434269, 27.2441672759339],
            //     [100.039690324806, 27.2441684817221],
            //     [100.039676285532, 27.2441694687811],
            //     [100.039662229137, 27.2441702368081],
            //     [100.039648159896, 27.2441707855672],
            //     [100.03963408209, 27.2441711148901],
            //     [100.03962, 27.2441712246757],
            //     [100.03960591791, 27.2441711148901],
            //     [100.039591840104, 27.2441707855672],
            //     [100.039577770863, 27.2441702368081],
            //     [100.039563714468, 27.2441694687811],
            //     [100.039549675194, 27.2441684817221],
            //     [100.03953565731, 27.2441672759339],
            //     [100.03952166508, 27.2441658517867],
            //     [100.039507702757, 27.2441642097173],
            //     [100.039493774588, 27.2441623502296],
            //     [100.039479884806, 27.2441602738939],
            //     [100.039466037633, 27.2441579813469],
            //     [100.039452237276, 27.2441554732914],
            //     [100.039438487929, 27.2441527504963],
            //     [100.039424793769, 27.2441498137956],
            //     [100.039411158953, 27.2441466640892],
            //     [100.039397587623, 27.2441433023412],
            //     [100.039384083899, 27.2441397295809],
            //     [100.039370651878, 27.2441359469013],
            //     [100.039357295636, 27.2441319554595],
            //     [100.039344019226, 27.2441277564757],
            //     [100.03962, 27.24345]

            // ];

            chartOption = {
                bmap: {
                    center: [116.46, 39.92],
                    zoom: 10,
                    roam: true,
                },
                series: [
                    {
                        type: 'scatter',
                        coordinateSystem: 'bmap',
                        data: [
                            [116.4696, 39.9173, 100], [116.4383, 40.1471, 50]
                            // [bd09.lon,bd09.lat]
                            // [102.73286745432786, 25.03032166639342]
                        ],
                    },
                    {
                        type: 'effectScatter',
                        coordinateSystem: 'bmap',
                        data: [
                            [116.4695, 39.9170, 100], [116.4971, 39.9078, 40]
                        ],
                    },
                    {
                        type: 'custom',
                        coordinateSystem: 'bmap',
                        renderItem: function renderItem(params: any, api: any) {
                            var coor_points = [
                                [100.03962, 27.24345],
                                [100.039895980774, 27.2441277564757],
                                [100.039882704364, 27.2441319554595],
                                [100.039869348122, 27.2441359469013],
                                [100.039855916101, 27.2441397295809],
                                [100.039842412377, 27.2441433023412],
                                [100.039828841047, 27.2441466640892],
                                [100.039815206231, 27.2441498137956],
                                [100.039801512071, 27.2441527504963],
                                [100.039787762724, 27.2441554732914],
                                [100.039773962367, 27.2441579813469],
                                [100.039760115194, 27.2441602738939],
                                [100.039746225412, 27.2441623502296],
                                [100.039732297243, 27.2441642097173],
                                [100.03971833492, 27.2441658517867],
                                [100.03970434269, 27.2441672759339],
                                [100.039690324806, 27.2441684817221],
                                [100.039676285532, 27.2441694687811],
                                [100.039662229137, 27.2441702368081],
                                [100.039648159896, 27.2441707855672],
                                [100.03963408209, 27.2441711148901],
                                [100.03962, 27.2441712246757],
                                [100.03960591791, 27.2441711148901],
                                [100.039591840104, 27.2441707855672],
                                [100.039577770863, 27.2441702368081],
                                [100.039563714468, 27.2441694687811],
                                [100.039549675194, 27.2441684817221],
                                [100.03953565731, 27.2441672759339],
                                [100.03952166508, 27.2441658517867],
                                [100.039507702757, 27.2441642097173],
                                [100.039493774588, 27.2441623502296],
                                [100.039479884806, 27.2441602738939],
                                [100.039466037633, 27.2441579813469],
                                [100.039452237276, 27.2441554732914],
                                [100.039438487929, 27.2441527504963],
                                [100.039424793769, 27.2441498137956],
                                [100.039411158953, 27.2441466640892],
                                [100.039397587623, 27.2441433023412],
                                [100.039384083899, 27.2441397295809],
                                [100.039370651878, 27.2441359469013],
                                [100.039357295636, 27.2441319554595],
                                [100.039344019226, 27.2441277564757],
                                [100.03962, 27.24345]

                            ];

                            var points = [];
                            for (var i = 0; i < coor_points.length; i++) {
                                points.push(api.coord(coor_points[i]));
                            }
                            var color = api.visual('color');
                            return {
                                type: 'polygon',
                                shape: {
                                    points: echarts.graphic.clipPointsByRect(coor_points, {
                                        x: params.coordSys.x,
                                        y: params.coordSys.y,
                                        width: params.coordSys.width,
                                        height: params.coordSys.height
                                    })
                                },
                                style: api.style({
                                    fill: color,
                                    stroke: echarts.color.lift(color, 5)
                                })
                            };
                        },

                        animation: false,
                        silent: true,
                        data: [
                            0
                        ],

                    } as echarts.CustomSeriesOption
                ]

            }
        } else if (chartid == "efcabf75-68a7-4941-8208-c0fd38a8aa64") {
            chartOption = {
                bmap: {
                    center: [120.13066322374, 30.240018034923],
                    zoom: 14,
                    roam: true
                },
                visualMap: {
                    // type: "continuous",
                    show: false,
                    top: 'top',
                    min: 0,
                    max: 5,
                    // seriesIndex: 0,
                    calculable: true,
                    inRange: {
                        color: ['blue', 'blue', 'green', 'yellow', 'red']
                    }
                },
                series: [
                    {
                        type: 'heatmap',
                        coordinateSystem: 'bmap',
                        // roam: true,
                        data: [
                            [120.12623751461, 30.225514547091, 0],
                            [120.12569412003, 30.225289138515, 3],
                            [120.1260081246, 30.225109979145, 0],
                            [120.12428900347, 30.224907917069, 4],
                            [120.1233608862, 30.224531990576, 0],
                            [120.12328968155, 30.225342953599, 3],
                            [120.12289821299, 30.22630750923, 0],
                            [120.1226090834, 30.226410111043, 0],
                            [120.12245512244, 30.226878337044, 3],
                            [120.1221672377, 30.227216311881, 1],
                            [120.12164622224, 30.227314672485, 4],
                            [120.12131843541, 30.2278850071, 2],
                            [120.12014167505, 30.227860326084, 1],
                            [120.11982765304, 30.227656247151, 0],
                            [120.11973240019, 30.227795357688, 5],
                            [120.11925363536, 30.227845331234, 3],
                            [120.11885762159, 30.227618695459, 2],
                            [120.11832083554, 30.227828718428, 1],
                            [120.11764254714, 30.2275760858, 2],
                            [120.11660933416, 30.227423442401, 0],
                            [120.11656948113, 30.227779565953, 2],
                            [120.11622329162, 30.228135226062, 3],
                            [120.11555815443, 30.228552151927, 3],
                            [120.11535925104, 30.228933739657, 4],
                            [120.11508923593, 30.228929830582, 5],
                            [120.1146001061, 30.22934780379, 0],
                            [120.11445996865, 30.229331065139, 1],
                            [120.11431067403, 30.22978764355, 0],
                            [120.11302670426, 30.23089800754, 0],
                            [120.11221369625, 30.230690045417, 4],
                            [120.11084513509, 30.230688382099, 0],
                            [120.11075387695, 30.231230018241, 1],
                            [120.11047529772, 30.231392448881, 3],
                            [120.11044468237, 30.231294617088, 4],
                            [120.11025496346, 30.231342673512, 0],
                            [120.11030754358, 30.231163266196, 4],
                            [120.1100134074, 30.231086634903, 2],
                            [120.10942045927, 30.231514831682, 0],
                            [120.10871903343, 30.231333695816, 4],
                            [120.10807570214, 30.231423668372, 0],
                            [120.10763082528, 30.231618832285, 4],
                            [120.1076172214, 30.231857645696, 5],
                            [120.10689398326, 30.232254414884, 3],
                            [120.10654498471, 30.232314939116, 0],
                            [120.10625143658, 30.23209170759, 2],
                            [120.10561817759, 30.232356818249, 0],
                            [120.1052679884, 30.232200325878, 1],
                            [120.10469210597, 30.232873117082, 2],
                            [120.10389149564, 30.23276397585, 1],
                            [120.10356823004, 30.233116970493, 4],
                            [120.10317755244, 30.232897606201, 1],
                            [120.10252540126, 30.232913879755, 1],
                            [120.10279292553, 30.233485195807, 0],
                            [120.10276069499, 30.233702758115, 3],
                            [120.10215452954, 30.233790321439, 2],
                            [120.1009310109, 30.234705361223, 4],
                            [120.10070707153, 30.235535810923, 4],
                            [120.10094839226, 30.236688233814, 3],
                            [120.10068527845, 30.237148154301, 4],
                            [120.0991389033, 30.237324990834, 5],
                            [120.09909068745, 30.237447538545, 3],
                            [120.09889701846, 30.237349749781, 4],
                            [120.09909203766, 30.237447835223, 5],
                            [120.09885018257, 30.23750460963, 0],
                            [120.0985803733, 30.237399683068, 2],
                            [120.09846495357, 30.237544566136, 1],
                            [120.09691421859, 30.237265892435, 2],
                            [120.09640608591, 30.237517005461, 0],
                            [120.09592462988, 30.237417530439, 4],
                            [120.09578540325, 30.237505515712, 5],
                            [120.09585978634, 30.237657781556, 4],
                            [120.09513916535, 30.237598831505, 4],
                            [120.09537128694, 30.237827879177, 5],
                            [120.09463588238, 30.237553240794, 1],
                            [120.09434036469, 30.237675194088, 0],
                            [120.09435981434, 30.237393110666, 0],
                            [120.09416203126, 30.237518267502, 1],
                            [120.09442420859, 30.237544095207, 5],
                            [120.09406873401, 30.237645935348, 4],
                            [120.09421500188, 30.23757678504, 2],
                            [120.09417631939, 30.237788505192, 0],
                            [120.09401715446, 30.23763248493, 5],
                            [120.09346056016, 30.237763679213, 4],
                            [120.0931873532, 30.237498100783, 3],
                            [120.09307923181, 30.237657251, 3],
                            [120.09326434256, 30.237542950418, 2],
                            [120.09322848324, 30.2373027005, 4],
                            [120.09301091533, 30.237355077676, 1],
                            [120.09300231512, 30.237102585171, 3],
                            [120.09273277667, 30.237014001726, 3],
                            [120.09235667235, 30.237023849668, 5],
                            [120.09247034134, 30.237020799186, 2],
                            [120.09235415838, 30.236931281436, 3],
                            [120.09249523307, 30.237077424673, 3],
                            [120.09227857692, 30.236953488801, 0],
                            [120.09236889214, 30.236802403617, 4],
                            [120.09228327523, 30.236774501485, 0],
                            [120.0915479786, 30.237013030973, 3],
                            [120.09170832511, 30.236965232782, 4],
                            [120.09165932082, 30.236819989034, 1],
                            [120.09141982762, 30.236834325375, 5],
                            [120.09105978327, 30.237018774407, 5],
                            [120.09122130957, 30.236776979775, 0],
                            [120.09096738259, 30.23666355983, 1],
                            [120.09029968116, 30.236770616438, 2],
                            [120.08983820472, 30.237106950056, 1],
                            [120.08945597294, 30.23699652666, 0],

                        ],
                        pointSize: 5,
                        blurSize: 6
                    } as any
                ]

            }
        } else if (chartid == "68ddd844-d023-405e-94e2-88e55a135eb3") {

            chartOption = {
                visualMap: {
                    show: false,
                    min: 800,
                    max: 50000,
                    type: "continuous",
                    realtime: false,
                    calculable: true,
                    inRange: {
                        color: ['lightskyblue', 'yellow', 'orangered']
                    }
                },
                series: [
                    {
                        type: 'map',
                        map: 'map',
                        roam: true,
                        label: {
                            show: true
                        },
                        data: [
                            { name: '中西区', value: [20057.34] },
                            { name: '湾仔', value: [15477.48] },
                            { name: '东区', value: 31686.1 },
                            { name: '南区', value: 6992.6 },
                            { name: '油尖旺', value: 44045.49 },
                            { name: '深水埗', value: 40689.64 },
                            { name: '九龙城', value: 37659.78 },
                            { name: '黄大仙', value: 45180.97 },
                            { name: '观塘', value: 55204.26 },
                            { name: '葵青', value: 21900.9 },
                            { name: '荃湾', value: 4918.26 },
                            { name: '屯门', value: 5881.84 },
                            { name: '元朗', value: 4178.01 },
                            { name: '北区', value: 2227.92 },
                            { name: '大埔', value: 2180.98 },
                            { name: '沙田', value: 9172.94 },
                            { name: '西贡', value: 3368 },
                            { name: '离岛', value: 806.98 }
                        ],
                    }
                ]
            }
        } else if (chartid == "472a9015-c131-44b5-88da-41dbc172a1c2") {

            chartOption = {
                geo: {
                    map: 'map',
                    roam: true,
                },
                series: [
                    {
                        type: 'scatter',
                        coordinateSystem: 'geo',
                        data: [
                            [114.113747, 22.285694, 100], [114.223022, 22.33391, 50]
                        ],
                    },
                    {
                        type: 'effectScatter',
                        coordinateSystem: 'geo',
                        data: [
                            [114.19825, 22.258787, 100], [113.924026, 22.157084, 40]
                        ],
                    }
                ]
            }
        } else if (chartid == "f37687e4-c20e-4bfc-8c5d-5c659154118d") {

            chartOption = {
                geo: {
                    map: 'map',
                    roam: true,
                },
                series: [
                    {
                        type: 'lines',
                        coordinateSystem: 'geo',
                        polyline: true,
                        lineStyle: {
                            color: 'red', width: 3,
                        },
                        data: [
                            {
                                name: "线路1",
                                coords: [[114.19825, 22.258787], [114.179657, 22.349068]],
                            }
                        ],
                    }
                ]
            }
        } else if (chartid == "741fef90-a46b-4cdc-9995-07515af1db31") {
            //空白
            chartOption = {
                series: []
            }
        }
        comp.chartStyle = {
            type: ViChartMainType.Line,
            style: chartOption,
            coorNameList: [],
            dataNameList: [],
            otherNameList: [],
            chartCoorList: {},
            chartSeriesList: [],
            // chartConditionList:[],
        }
    } else if (json.type == "comp") {
        comp = reactive(cloneDeep(Components[json.value as keyof typeof Components]))
        comp.name = json.value
        comp.component = json.value
    } else {
        console.log("拖出信息内容", json)
    }
    if (comp != null) {
        comp.id = nanoidEx()
        comp.customName = comp.typeName + "-" + store.state.viCounter.toString()
        store.commit("set_counter")
        comp.style.top = event.clientY - rect!.top
        comp.style.left = event.clientX - rect!.left
        comp.style.width = comp.style.width + comp.style.width * store.state.canvas.scale / 100
        comp.style.height = comp.style.height + comp.style.height * store.state.canvas.scale / 100
        comp.style.zIndex = store.state.viViews[store.state.viCtrlIndex].z_index
        store.commit("add_component", comp)
        store.commit("set_boarding", { value: 1 })
        store.commit("set_bread_index", { value: 0 })
        nextTick(() => {
            store.commit("snapshot")
        })
    }

}

function handleMouseDown(event: any) {
    event.stopPropagation()
    store.commit("unset_component_id")
}


function handleKeyDown(event: KeyboardEvent) {
    // console.log("DOWN", event.key, event.key == "Control")
    if (event.key == "Control") {
        store.commit("set_control_pressed", true)
    }
    if (event.key == "Shift") {
        store.commit("set_shift_pressed", true)
    }
}
function handleKeyUp(event: KeyboardEvent) {
    // console.log("UP", event.key, event.key == "Control")
    if (event.key == "Control") {
        store.commit("set_control_pressed", false)
    }
    if (event.key == "Shift") {
        store.commit("set_shift_pressed", false)
    }
}
function handleWheel(event: any) {
    if (store.state.control_pressed) {
        event.preventDefault()
        event.stopPropagation()
        // console.log(event.deltaY)
        if (ViContentRef.value) {
            if (event.deltaY < 0) {
                if (ViContentRef.value.scrollLeft > 0) {
                    ViContentRef.value.scrollLeft = ViContentRef.value.scrollLeft - 100;
                }
            } else if (event.deltaY > 0) {
                if (ViContentRef.value.scrollLeft < ViContentRef.value.scrollWidth) {
                    ViContentRef.value.scrollLeft = ViContentRef.value.scrollLeft + 100;
                }
            }
        }
    }
}

// const ready = ({BMap,map})=>{
//   //  对地图进行自定义操作
// };


</script>

<template>
    <v-main>
        <v-container id="main_container">
            <v-row>
                <MainMenu />
                <v-col id="main_center">
                    <!-- <baidu-map class="map" :center="{ lng: 118.454, lat: 32.955 }" :zoom="5" >
                    </baidu-map> -->
                    <v-list :height="view_height" id="main_list" :style="{
                        display: 'flex',
                        justifyContent: 'center',
                    }">
                        <section class="vicenter">
                            <!-- <div class="vicontent" id="vicontent" ref="ViContentRef" @dragover="dragover($event)"
                                @drop="handleDrop($event)" @mousedown="handleMouseDown" @wheel="handleWheel($event)">
                                <div id="PreEditorID">
                                    <Editor id="EditorID" />
                                </div>
                            </div> -->
                            <!-- <div class="vicontent" id="vicontent" ref="ViContentRef" @dragover="dragover($event)"
                                @drop="handleDrop($event)" @mousedown="handleMouseDown" @wheel="handleWheel($event)">
                                <div id="PreEditorID"> -->
                            <Editor class="vicontent" id="EditorID" @dragover="dragover($event)" @drop="handleDrop($event)"
                                @mousedown="handleMouseDown" @wheel="handleWheel($event)" />
                            <!-- </div>
                            </div> -->
                        </section>
                    </v-list>

                    <div id="main_control">

                        <Control />
                        <InfoDialog></InfoDialog>
                        <AskDialog></AskDialog>
                        <LoaderDialog></LoaderDialog>
                        <HtmlInfo></HtmlInfo>
                        <FunctionDialog></FunctionDialog>
                        <Help></Help>
                        <CopyRight></CopyRight>
                        <ChartHelp></ChartHelp>
                        <UpdateInfo></UpdateInfo>
                        <LicenseDialog></LicenseDialog>
                        <v-slider v-model="editor_scale" color="orange" density="comfortable" elevation="0" tick-size="0"
                            track-size="1" thumb-size="10" min="30" max="100"></v-slider>

                        <Foot />
                    </div>
                </v-col>

                <LeftMenu />
                <RightMenu />
            </v-row>
        </v-container>
        <!-- <p style="font-family: Source Han Sans Simplified Chinese Bold;font-size:xx-large;">xd就这测试</p>
        <p style="font-family: Source Han Sans Simplified Chinese Bold,sans-serif;font-size:xx-large;">xd就这测试</p>
        <p style="font-family: Source Han Sans Simplified Chinese Bold,sans-serif;font-size:xx-large;font-weight: 900;color: black;">xd就这测试</p>
        <p style="font-family: Noto Sans Simplified Chinese Bold,sans-serif;font-size:xx-large;">xd就这测试</p>
        <p style="font-family: Noto Sans Simplified Chinese Bold,sans-serif;font-size:xx-large;font-weight:900;color: black;">xd就这测试</p> -->
    </v-main>
</template>

<style lang="scss">
.v-input__details {
    height: 0px !important;
    min-height: 0px !important;
    padding-top: 0px !important;
}


.vicenter {
    height: 100%;
    overflow: auto;
    /* 让 maindiv 显示滚动条 */
    // border: 1px solid #000; /* 为了可视化效果添加边框 */
}

.vicontent {
    // width: 500px; /* 设置 subdiv 的宽度 */
    // height: 500px; /* 设置 subdiv 的高度 */
    background-color: #ccc;
    /* 为了可视化效果添加背景颜色 */
}

// .vicenter {
//     // margin-left: 200px;
//     // margin-right: 288px;
//     // background: #c67f7f;
//     height: 100%;
//     overflow: scroll;
//     scrollbar-width: thin;
//     // padding: 20px;

//     .vicontent {
//         justify-content: center;
//         width: 100%;
//         height: 100%;
//         overflow: auto;
//         scrollbar-width: thin;
//     }
// }


/* 自定义滚动条的样式 */
.vicenter::-webkit-scrollbar {
    width: 14px;
    /* 设置滚动条的宽度 */
}

.vicenter::-webkit-scrollbar-thumb {
    background-color: #00335522;
    /* 设置滚动条滑块的颜色 */
    border-radius: 6px;
    /* 设置滚动条滑块的边角半径 */
}


// /* 自定义滚动条的样式 */
// .vicontent::-webkit-scrollbar {
//     width: 14px;
//     /* 设置滚动条的宽度 */
// }

// .vicontent::-webkit-scrollbar-thumb {
//     background-color: #00335522;
//     /* 设置滚动条滑块的颜色 */
//     border-radius: 6px;
//     /* 设置滚动条滑块的边角半径 */
// }

// .vicontent::-webkit-scrollbar-track {
//   background-color: transparent; /* 设置滚动条轨道的颜色 */
// }

// /* 自定义水平滚动条的样式 */
// .vicontent::-webkit-scrollbar-horizontal {
//   border-bottom: 12px solid #f1f1f1; /* 设置水平滚动条的高度和颜色 */
// }

// .vicontent::-webkit-scrollbar-thumb:horizontal {
//   background-color: #888; /* 设置水平滚动条滑块的颜色 */
//   border-radius: 6px; /* 设置水平滚动条滑块的边角半径 */
// }

// .vicontent::-webkit-scrollbar{
//     width: 10px;  
// }

.map {
    width: 400px;
    height: 300px;
}
</style>