import { createStore, Store } from 'vuex'
import { ViComponent, ViInnerStyle, ViView } from './index'
import { EventsEmit } from '../../wailsjs/runtime'
import { get, set } from 'lodash'
import { WGS84ToBD09 } from '../utils/func'
export interface ViInfomation {
    //组织名称
    company: string,
    //人员姓名
    author: string,
    //邮箱地址
    email: string,
    //联系电话
    phone: string,
    //官方网址
    website: string,
    //未解码的密码
    passwd: string,
    //文档描述
    docinfo: string,
    //请求解密的文件路径
    filepath: string,
}
// 为 store state 声明类型
export interface DisplayState {
    //预加载信息
    preInfo: ViInfomation,
    canvas: {
        //保存缩放前的缩放值
        scale_backup: number,
        //保存当前缩放值
        scale: number,
        //布画高度
        height: number,
        //布画宽度
        width: number,
        //组织名称
        company: string,
        //人员姓名
        author: string,
        //邮箱地址
        email: string,
        //联系电话
        phone: string,
        //官方网址
        website: string,
        //文档密码，导出的小视图文件需要密码才能打开。
        passwd: string,
        //文档描述
        docinfo: string,
    },
    //原始项目的缩放,需要记录原始项目导出文件时候的缩放,方便到显示器显示的时候可以计算新的位置
    viInitScale: number,
    //视图数组
    viViews: ViView[],
    //视图索引
    viCtrlIndex: number,
    //视图显示宽度
    viViewWidth: number,
    //视图显示高度
    viViewHeight: number,
    //显示提示对话框
    viShow: boolean,
    //显示对话框内容
    viMessage: string,
    //显示进度条
    viShowLoader: boolean,
    //显示提问问题
    viShowAsk: boolean,
    //显示版权对话框
    viShowCopyRight: boolean,
    //显示修改密码
    viShowChangePasswd: boolean,
    //文件路径，可能从文件读取的时候需要保存文件
    viFilePath: string,
    //远程登陆密码
    viRemotePasswd: string,
    //远程Token
    viRemoteToken: string,
    //显示器对象
    display: HTMLElement | null,
    //控制条高度
    controlHeight: number,
    //组件显示波纹效果，display的时候应该是false
    viRipple: boolean,
    //百度地图密钥
    viBaiduMapKey: string,
    //是否是全屏
    viIsFullScreen:boolean,
    //自动切换时间
    viTimeDelay:number,
    //视图切换类型
    viPlayType:string,
}

export const storeDisplay = createStore<DisplayState>({
    state: {
        preInfo: {
            author: "",
            email: "",
            company: "",
            phone: "",
            website: "",
            passwd: "",
            docinfo: "",
            filepath: "",
        },
        canvas: {
            scale_backup: 50,
            scale: 50,
            height: 1080,
            width: 1920,
            company: "",
            author: "",
            email: "",
            phone: "",
            website: "",
            passwd: "",
            docinfo: "",
        },
        viInitScale: 100,
        viViews: [],
        viCtrlIndex: -1,
        viViewWidth: 0,
        viViewHeight: 0,
        viShow: false,
        viMessage: "",
        viShowAsk: false,
        viShowLoader: false,
        viShowCopyRight: false,
        viShowChangePasswd: false,
        viFilePath: "",
        viRemotePasswd: "",
        viRemoteToken: "",
        display: null,
        controlHeight: 200,
        viRipple: false,
        viBaiduMapKey: "",
        viIsFullScreen:false,
        viTimeDelay:5,
        viPlayType:"mdiGestureTap",
    },
    mutations: {
        //设置预加载信息
        set_preinfo(state, payload: ViInfomation) {
            state.preInfo = payload
        },
        //设置布画属性
        set_canvas(state, value) {
            state.canvas = value
        },
        set_viInitScale(state, value: number) {
            state.viInitScale = value
        },
        //设置组件数组
        set_viViews(state, views: ViView[]) {
            state.viViews = views
        },
        //设置视图索引
        set_view_index(state, index?: number) {
            if (index == undefined) {
                state.viViews.forEach((v, i) => {
                    if (v.visible) {
                        state.viCtrlIndex = i
                        return
                    }
                })
            } else {
                state.viCtrlIndex = index
                // console.log("state.viCtrlIndex",state.viCtrlIndex)
            }
        },
        //设置视图矩形
        set_rect(state, payload: { width: number, height: number }) {
            state.viViewHeight = payload.height
            state.viViewWidth = payload.width
            //更新缩放比例
            state.canvas.scale = 100 * payload.width / state.canvas.width
        },
        //设置密码
        set_pre_passwd(state, value: string) {
            state.preInfo.passwd = value
        },
        //显示对话框提示
        show_dialog(state, payload: {
            msg: string,
            show: boolean,
        }) {
            state.viMessage = payload.msg
            state.viShow = payload.show
        },
        //显示输入密码框
        show_dialog_ask(state, value: boolean) {
            state.viShowAsk = value
        },
        //显示进度条
        show_loader(state, value: boolean) {
            state.viShowLoader = value
        },
        //显示版权
        show_copyright(state, value: boolean) {
            state.viShowCopyRight = value
        },
        //显示修改密码
        show_change_passwd_dialog(state, value: boolean) {
            state.viShowChangePasswd = value
        },
        //设置文件路径
        set_FilePath(state, file_path: string) {
            state.viFilePath = file_path
        },
        //甚至远程密码
        set_remote_passwd(state, passwd: string) {
            state.viRemotePasswd = passwd
        },
        //设置远程token
        set_remote_token(state, rtoken: string) {
            state.viRemoteToken = rtoken
        },
        //设置显示器对象
        set_display(state, value: HTMLElement) {
            state.display = value
        },
        //设置控制条高度
        set_control_height(state, value: number) {
            state.controlHeight = value
        },
        //设置组件内部样式
        set_inner_style(State, payload: {
            style: ViInnerStyle,
            addr?: { viewid: number, comid: string }
        }) {

            let vid = State.viCtrlIndex
            let cid = State.viViews[State.viCtrlIndex].componet_ids[0]
            if (payload.addr != undefined) {
                vid = payload.addr.viewid
                cid = payload.addr.comid
            }
            State.viViews[vid].components.forEach(v => {
                if (v.id == cid) {
                    v.innerStyle = payload.style
                    return
                }
            })
        },
        //通过组件ID选择组件
        set_component_id(State, id: string) {
            //let id_index = State.viViews[State.viCtrlIndex].componet_ids.indexOf(id)
            State.viViews[State.viCtrlIndex].componet_ids.length = 0
            State.viViews[State.viCtrlIndex].componet_ids.push(id)
        },
        //触发组件数据给ui前端
        emit_component_objects(state) {
            // let data: ViComponent[][] = []
            // state.viViews.forEach((v, i) => {
            //     data.push(v.components)
            // })
            EventsEmit("UIEmitComponentsData", JSON.stringify(state.viViews))
        },
        //更新显示器视图值，由原车提交修改组件值。
        //data应该是数组类数据的string,例如："[选项1,选项2]",由后端go发送给前端之组件，组件接收到后，调用该函数来更新
        //map_data := make(map[string]string)
        //map_data["key"] = post_key
        //map_data["data"] = data
        update_compo_value(state, value: any) {
            // let info = JSON.parse(value)
            let key = value["key"] as string
            let data = value["data"] as string
            let comptype = value["comtype"] as string
            if (data.length < 2) {
                return
            }
            state.viViews.forEach((view, vindex) => {
                view.components.forEach((compo, cindex) => {
                    if ("nni" + compo.id == key) {
                        if (compo.innerStyle.styleFlag) {
                            // console.log(compo.id, key, compo.name, comptype)
                            try {
                                let obj = JSON.parse(data) as string[]
                                switch (compo.name) {
                                    case "ViCheckBox":
                                        if (compo.innerStyle.selectList) {
                                            compo.innerStyle.selectList.length = 0
                                            obj.forEach(d => {
                                                if (compo.innerStyle.textList?.indexOf(d) != -1) {
                                                    compo.innerStyle.selectList?.push(d)
                                                }
                                            })
                                        }
                                        break;
                                    case "ViTimer":
                                        if (compo.innerStyle.textList) {
                                            compo.innerStyle.textList.length = 0
                                            // console.log("obj:", obj)
                                            obj.forEach(d => {
                                                compo.innerStyle.textList?.push(d)
                                                if (d.indexOf("-") > -1 && compo.innerStyle.textList) {
                                                    compo.innerStyle.textList[0] = d
                                                }
                                                if (d.indexOf(":") > -1 && compo.innerStyle.textList) {
                                                    compo.innerStyle.textList[1] = d
                                                }
                                            })

                                        }
                                        break;
                                    case "ViIcon":
                                        compo.innerStyle.color = obj[0]
                                        break;
                                    default:
                                        compo.innerStyle.text = obj[0]
                                        break;
                                }
                            } catch (error) {
                                console.log(error)
                            }
                        }
                        return
                    }
                })
            })
        },
        //通过路径更新图表的echart option，例如xAxis[0].data=...
        set_echart_option_by_path(state,
            payload: { id: string, path: string, value: any }) {
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                // console.log("set_echart_option_by_path", item.id, payload.id)
                if ("nni" + item.id == payload.id) {
                    let obj = item.chartStyle?.style
                    if (obj != undefined) {
                        set(obj, payload.path, payload.value)
                    }
                    // console.log("set_echart_option_by_path", obj)
                    return
                }
            })
        },
        //通过索引更新图表的echart option的series，例如series[0]=...
        // set_echart_series_by_index(state,
        //     payload: { id: string, index: number, value: any }) {
        //     state.viViews[state.viCtrlIndex].components.forEach(item => {
        //         console.log(item.id, payload.id)
        //         if ("nni" + item.id == payload.id) {
        //             let obj = item.chartStyle?.style
        //             if (obj != undefined) {
        //                 let path = "series[" + payload.index + "].data"
        //                 set(obj, path, payload.value)
        //             }
        //             console.log("update", obj)
        //             return
        //         }
        //     })
        // },
        set_echart_series_by_index(state,
            payload: { id: string, index: number, value: any, name?: string }) {
            // console.log("set_echart_series_by_index", payload)
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                // console.log(item.id, payload.id)
                if (item.id == payload.id || "nni" + item.id == payload.id) {
                    let obj = item.chartStyle?.style
                    if (obj != undefined) {
                        let path = "series[" + payload.index + "].data"
                        let series_type = get(obj, "series[" + payload.index + "].type")
                        let series_coordinatesystem = get(obj, "series[" + payload.index + "].coordinateSystem")
                        let bmap_option = get(obj, "bmap")
                        // console.log("series_type", series_type, series_coordinatesystem)
                        let newdata = payload.value
                        if (Array.isArray(payload.value[0]) && payload.value[0].length >= 2 && (series_type === 'pie' || series_type === 'funnel')) {
                            // console.log("payload.value[0]", payload.value[0])
                            let lst = payload.value as any[]
                            newdata = []
                            lst.forEach(v => {
                                let tmp = v as any[]

                                newdata.push({
                                    value: tmp[0],
                                    name: tmp[1]
                                })
                            })
                        } else if (series_type === 'gauge') {
                            newdata = [payload.value[0]]
                            if (Array.isArray(payload.value[0]) && payload.value[0].length >= 2) {
                                let data = payload.value[0] as any[]
                                newdata = [{
                                    value: data[0],
                                    name: data[1]
                                }]
                            }
                        }else if (series_type === 'radar') {
                            let data = payload.value as any[]
                            newdata = [{
                                value:data.flat()
                            }]
                            // if (Array.isArray(payload.value[0]) && payload.value[0].length >= 2) {
                            //     let data = payload.value[0] as any[]
                            //     console
                            //     newdata = [{
                            //         value: data.flat(),
                            //     }]
                            // }
                        } else if (series_type === 'scatter' || series_type === 'effectScatter') {
                            newdata = payload.value
                            if (Array.isArray(payload.value[0]) && payload.value[0].length >= 2) {
                                let lst = payload.value as any[]
                                newdata = []
                                lst.forEach(v => {
                                    let tmp = v as any[]
                                    let lng = tmp[0]
                                    let lat = tmp[1]
                                    if (lng == null || lat == null) {
                                        return
                                    }
                                    let coor = { lng, lat }
                                    if (series_coordinatesystem == 'bmap' || bmap_option) {
                                        coor = WGS84ToBD09(coor.lng, coor.lat)
                                        // console.log('WGS84ToBD09', lng, lat, coor)
                                    }
                                    tmp[0] = coor.lng
                                    tmp[1] = coor.lat
                                    newdata.push(tmp)
                                })
                            }
                        } else if (series_type === 'map') {
                            //需要给map类series填充名字，也就name属性，用来对应geojson上的地理名称，正常显示对应位置的数据
                            // console.log("series map", payload.value)
                            newdata = payload.value
                            if (Array.isArray(payload.value[0]) && payload.value[0].length >= 2) {
                                let lst = payload.value as any[]
                                newdata = []
                                lst.forEach(v => {
                                    let tmp = v as any[]

                                    newdata.push({
                                        value: tmp[1],
                                        name: tmp[0]
                                    })
                                })
                            }

                        } else if (series_type === 'lines') {
                            // console.log("series lines", payload.value)
                            newdata = payload.value
                            if (Array.isArray(payload.value[0]) && payload.value[0].length >= 2) {
                                let lst = payload.value as any[]
                                newdata = []
                                let temp_obj = { name: "", coords: [] as any[] }
                                lst.forEach(v => {
                                    let tmp = v as any[]
                                    if (tmp[0] === 0 && tmp[1] === 0) {
                                        if (temp_obj.coords.length > 0) {
                                            newdata.push(temp_obj)
                                        }
                                        temp_obj = {
                                            name: "",
                                            coords: []
                                        }
                                    } else {
                                        let lng = tmp[0]
                                        let lat = tmp[1]
                                        if (lng == null || lat == null) {
                                            return
                                        }
                                        let coor = { lng, lat }
                                        if (series_coordinatesystem == 'bmap' || bmap_option) {
                                            coor = WGS84ToBD09(coor.lng, coor.lat)
                                            // console.log('WGS84ToBD09', lng, lat, coor)
                                        }
                                        temp_obj.coords.push([coor.lng, coor.lat])
                                        if (tmp[2] != undefined) {
                                            temp_obj.name = tmp[2]
                                        }
                                    }
                                })
                                if (temp_obj.coords.length > 0) {
                                    newdata.push(temp_obj)
                                }
                                // console.log("series lines newdata", newdata)
                            }
                        }
                        set(obj, path, newdata)
                        if (payload.name) {
                            path = "series[" + payload.index + "].name"
                            set(obj, path, payload.name)
                        }
                    }
                    // console.log("update", obj)
                    return
                }
            })
        },
        //这个是专门给修改百度地图Key用的
        set_baidu_map_key(state, value: string) {
            state.viBaiduMapKey = value
        },
        //设置是否是全屏
        set_viIsFullScreen(state,value:boolean){
            state.viIsFullScreen = value
        },
        //设置自动切换时间
        set_viTimeDelay(state,value:number){
            state.viTimeDelay = value
        },
        //设置切换类型
        set_viPlayType(state,value:string){
            state.viPlayType = value
        }

    },
})