
import { createStore, Store } from 'vuex'
import { ref, Ref, reactive, nextTick } from 'vue'
// import { ref, InjectionKey } from 'vue'
import { nanoid } from 'nanoid'
import { cloneDeep, List, map, set, unset } from 'lodash'
// import html2canvas from 'html2canvas';
import { WriteToCache } from "../../wailsjs/go/main/App"
import { nanoidEx, updateObjectWrapper, WGS84ToBD09 } from '../utils/func'
import { EChartsOption } from 'echarts'
import * as echarts from 'echarts'

import html2canvas from 'html2canvas';
import { get } from 'lodash'
import { number } from 'mathjs'


// 快照对象类型，每个viewid对应一个视图的多次提交记录
export interface SnapObject { [viewid: string]: ViView[] }
// let oob:OOB={}
// let x = Object.keys(oob).indexOf("aaaa")


export interface ViGroupStyle {
    top: number,
    left: number,
    width: number,
    height: number,
}

export interface ViInnerStyle {
    //颜色
    color?: string,
    //已选择的文本或单一文本
    text?: string,
    //可选列表
    textList?: string[],
    //已选列表或是日期和时间
    selectList?: string[],
    //状态标识,比如下拉框单选或多选,
    styleFlag?: boolean,
    //水平或垂直。。。默认ffalse
    vertical?: boolean,
    fontfamily?: string,
    fontSize?: number,
    fontWeight?: any,
    fontStyle?: any,
    //保存组件所对应表和字段,一个组件对应和多个字段
    list: {
        sheetkey: string,
        fieldkey: string,
    }[]
}

export interface ViOutStyle {
    top: number,
    left: number,
    width: number,
    height: number,
    rotate: number,
    backgroundColor: string,

    borderTopColor: string,
    borderBottomColor: string,
    borderLeftColor: string,
    borderRightColor: string,

    borderTopStyle: string,
    borderBottomStyle: string,
    borderLeftStyle: string,
    borderRightStyle: string,

    borderTopWidth: string,
    borderBottomWidth: string,
    borderLeftWidth: string,
    borderRightWidth: string,

    borderTopLeftRadius: string,
    borderTopRightRadius: string,
    borderBottomLeftRadius: string,
    borderBottomRightRadius: string,

    opacity: number,
    zIndex: number,

    //是显示模式还是编辑模式,因为远程提交数据的时候无法区分组件是在显示模式还是编辑模式,如果是显示模式才会更新
    displayMode: boolean,
}

// @ts-ignore
export enum ViChartMainType {
    //指定图表默认类型，之所以要指定，比如一个图表有波形和柱状，那新建数据时候应当有一个默认了型
    //波形
    Line = 'line',
    //柱状
    Bar = 'bar'
}

//图表属性类型
export enum ViChartEnum {
    //坐标
    Coor = 'Coor',
    //数据
    Data = 'Data',
    //其他属性
    Other = 'Other',
    //未知
    Undefined = 'Undefined',
}

export interface ViChartType {
    type: ViChartEnum,
    index: number,
}

//图表已经选择的属性
export interface ViChartOption {
    //图表属性
    key: string,
    //属性名称
    name: string,
    //属性计数，因为有些属性只能有一个，比如极坐标只能有一个，X坐标可以有两个
    count: number,
}
export interface ViChartStyle {
    //图表类型
    type: ViChartMainType,
    //图表样式
    style: EChartsOption,
    //图表坐标列表
    coorNameList: ViChartOption[],
    //图表数据列表
    dataNameList: ViChartOption[],
    //图表其他属性
    otherNameList: ViChartOption[],
    //图表坐标数据对应表关系
    chartCoorList: {
        [key: string]: {
            sheetkey: string;
            fieldkey: string;
        }
    },
    //图表series数据对应表关系,之所以和chartCoorList分开，主要是chartSeriesList用索引删除更方便
    chartSeriesList: {//表key
        sheetkey: string,
        //字段key,列表
        fieldkey: string[]
    }[],
    //图表过滤器id列表
    // chartConditionList:string[],

}

export interface ViComponent {
    //组件id
    id: string,
    //组件名称，软件初始化时异步加载组件并注册组件名称 ViSVGIcon ViGroup ...
    //目前组件名称主要还是用于识别组件，但例如 ViSVGIcon ViFontIcon等组件就无法通过组件名称识别，
    //因为组件名称只能识别到组件是svg图标或是字体图标，无法识组件里是那个图标。
    //因此通过name属性指定，Font ICON名称或是SVG ICON名称用来加载对应ICON,
    //其它组件的name直接是component,例如：ViGroup、ViLabel ...
    name: string,
    // 组件对象 软件初始化时异步加载组件并注册组件对象 ViSVGIcon ViGroup ...
    component: string,
    //组件名称，例如图标、图表。。。。
    typeName: string,
    //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    customName: string,
    //组件内部样式
    innerStyle: ViInnerStyle,
    //组件外部样式
    style: ViOutStyle,
    //组件组样式，创建组后，组件需要有在组里的位置信息等。
    groupStyle: ViGroupStyle,
    //ViGroup类组件需要被所有子组件放到这里
    components: ViComponent[],
    //图表样式，不同图表用不同样式，
    chartStyle: ViChartStyle | undefined,
}

export interface ViView {
    //视图ID
    id: string,
    //视图里已经选择的组件IDs
    componet_ids: string[],
    //视图中组件当前最大层级，0-n
    z_index: number,
    //视图背景属性
    background_color: string,
    //视图背景图片
    background_image: string,
    //视图中所有组件
    components: ViComponent[],
    //视图缩略图，base64图片
    shot: string,
    //视图可见
    visible: boolean,
    //快照索引
    snapindex: number,
}

// 为 store state 声明类型
export interface State {
    control_pressed: boolean,
    shift_pressed: boolean,
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
    //所有视图
    viViews: ViView[],
    //当前视图索引
    viCtrlIndex: number,
    //编辑器引用
    viEditor: HTMLElement | null,
    //保存组件复制对象行，用于处理复制粘贴问题
    viCopy: ViComponent[],
    //保存快照记录
    viSnapshot: SnapObject,
    //显示提示对话框
    viShow: boolean,
    //显示对话框内容
    viMessage: string,
    //显示进度条
    viShowLoader: boolean,
    //显示提问问题
    viShowAsk: boolean,
    //文件路径，可能从文件读取的时候需要保存文件
    viFilePath: string,
    //组件显示波纹效果
    viRipple: boolean,
    //组件计数器,每产生一个组件，counter自增1
    viCounter: number,
    //保存右边面板的索引，因为面板的变化需要导致是否显示chartoption的bradcrumbs
    viBoarding: number,
    //顶部菜单抽屉面板内容，主要是chart组件用到，
    viBreaditems: string[],
    //顶部菜单抽屉面板索引，需要关联到chart组件的第n级options
    viBreadindex: number,
    //显示v-html内容窗口
    viShowHtml: boolean,
    //显示v-html内容
    viHtmlContext: string,
    //保存当前图标属性索引
    ViChartRefIndex: ViChartType,
    //显示编写函数框
    ViShowFunction: boolean,
    //编写的函数内容
    ViFunctionContext: string,
    //编写函数的对象属性路径
    ViFunctionPath: string,
    //保存当前图表属性值
    viChartCurrentValue: any,
    //百度地图密钥
    viBaiduMapKey: string,
    //显示帮助界面
    viShowHelp: boolean,
    //显示编辑器的程序版权
    viShowCopyRight: boolean,
    //显示图表的帮助
    viShowChartHelp: boolean,
    //显示更新窗口
    viShowUpdate: boolean,
    //显示license窗口
    viShowLicense: boolean,
}


// 定义 injection key
// export const key: InjectionKey<Store<State>> = Symbol()
export const store = createStore<State>({
    state: {
        control_pressed: false,
        shift_pressed: false,
        canvas: {
            scale_backup: 30,
            scale: 30,
            height: 1080,
            width: 1920,
            company: "",
            author: "",
            email: "",
            phone: "",
            website: "",
            passwd: "",
            docinfo: ""
        },
        viViews: [],
        viCtrlIndex: 0,
        viEditor: null,
        viCopy: [],
        viSnapshot: {},
        viShow: false,
        viMessage: "",
        viShowAsk: false,
        viShowLoader: false,
        viFilePath: "",
        viRipple: true,
        viCounter: 0,
        viBoarding: 1,
        viBreaditems: ["NNI"],
        viBreadindex: 0,
        viShowHtml: false,
        viHtmlContext: '',
        ViChartRefIndex: {
            type: ViChartEnum.Undefined,
            index: -1
        },
        ViShowFunction: false,
        ViFunctionContext: "",
        ViFunctionPath: "",
        viChartCurrentValue: "",
        viBaiduMapKey: "",
        viShowHelp: false,
        viShowCopyRight: false,
        viShowChartHelp: false,
        viShowUpdate: false,
        viShowLicense: false,
    },
    mutations: {
        //清空工程
        clear_state(state) {
            state.viCtrlIndex = 0
            state.viViews.length = 1
            // @ts-ignore
            state.viViews[0].componet_ids.length = 0
            state.viViews[0].components.length = 0
            state.viViews[0].snapindex = -1
            state.viCopy.length = 0
            state.viSnapshot = {}
        },
        //覆盖工程
        set_state(state, value: State) {
            state.viCtrlIndex = 0
            state.canvas = value.canvas
            let views = value.viViews
            views.forEach(v => {
                v.shot = ""
            })
            state.viViews = views
            state.viCopy.length = 0
            state.viSnapshot = {}
            state.viBaiduMapKey =value.viBaiduMapKey
        },
        set_control_pressed(state, val: boolean) {
            state.control_pressed = val
        },
        set_shift_pressed(state, val: boolean) {
            state.shift_pressed = val
        },
        set_editor_scale(state, value: number) {
            if (value > 0) {
                state.canvas.scale_backup = state.canvas.scale
                state.canvas.scale = value
            }
        },
        set_editor_height(state, value: number) {
            if (value > 0) {
                state.canvas.height = value
            }
        },
        set_editor_width(state, value: number) {
            if (value > 0) {
                state.canvas.width = value
            }
        },
        set_company(state, value: string) {
            state.canvas.company = value
        },
        set_author(state, value: string) {
            state.canvas.author = value
        },
        set_email(state, value: string) {
            state.canvas.email = value
        },
        set_phone(state, value: string) {
            state.canvas.phone = value
        },
        set_website(state, value: string) {
            state.canvas.website = value
        },
        set_passwd(state, value: string) {
            state.canvas.passwd = value
        },
        set_docinfo(state, value: string) {
            state.canvas.docinfo = value
        },
        set_viCtrlIndex(state, value: number) {
            state.viCtrlIndex = value
        },
        set_viFilePath(state, value: string) {
            state.viFilePath = value
        },
        //设置背景颜色
        set_background_color(state, value: string) {
            state.viViews[state.viCtrlIndex].background_color = value
        },
        //设置背景图片
        set_background_image(state, value: string) {
            if (state.viViews[state.viCtrlIndex].background_image != "") {
                URL.revokeObjectURL(state.viViews[state.viCtrlIndex].background_image)
            }
            state.viViews[state.viCtrlIndex].background_image = value
        },
        //通过视图索引选择视图
        select_View(state, controlId: number) {
            if (state.viCtrlIndex != controlId) {
                state.viCtrlIndex = controlId
            }
        },
        //增加视图
        add_viView(state, view_index: number): string {
            let view: ViView = {
                id: nanoidEx(),
                componet_ids: [],
                z_index: 0,
                background_color: '#003355',
                background_image: "",
                components: [],
                shot: "",
                visible: true,
                snapindex: -1,
            }
            // NewImg(state.canvas.width * 0.01, state.canvas.height * 0.01).then(r => {
            //     if (r.indexOf("wails") != -1) {
            //         view.shot = r
            //     }
            // })
            if (state.viViews.length == 0) {
                state.viViews.splice(view_index, 0, view)
            } else {
                state.viViews.splice(view_index + 1, 0, view)
                state.viCtrlIndex = view_index + 1
            }

            return view.id
        },
        //通过视图ID删除视图
        del_viView(state, view_id: string) {
            if (state.viViews.length == 1) {
                state.viViews[0].components.length = 0
                return
            }
            state.viViews.forEach((v, k) => {
                if (v.id == view_id) {
                    state.viViews.splice(k, 1)
                    if (state.viCtrlIndex == k) {
                        if (state.viCtrlIndex != 0) {
                            state.viCtrlIndex = k - 1
                        } else {
                            state.viCtrlIndex = 0
                        }
                        return
                    }
                }
            })
        },
        //通过视图索引删除视图
        del_viView_ByIndex(state, view_index: number) {
            if (state.viViews.length == 1) {
                state.viViews[view_index].components.length = 0
                return
            }
            state.viViews.splice(view_index, 1)
            if (state.viCtrlIndex == view_index) {
                if (state.viCtrlIndex != 0) {
                    state.viCtrlIndex = view_index - 1
                } else {
                    state.viCtrlIndex = 0
                }
            }
        },
        //设置视图可见
        set_viView_visable(State, payload: { view_index: number, visible: boolean }) {
            State.viViews[payload.view_index].visible = payload.visible
        },
        //给当前视图增加组件
        add_component(State, component: ViComponent) {
            State.viViews[State.viCtrlIndex].components.push(component)
            State.viViews[State.viCtrlIndex].componet_ids.length = 0
            State.viViews[State.viCtrlIndex].componet_ids.push(component.id)
        },
        //通过组件ID删除组件
        del_component(State, id: string) {
            // State.viViews[State.viCtrlIndex].components.forEach((v, index) => {
            //     if (v.id == id) {
            //         State.viViews[State.viCtrlIndex].components.splice(index, 1)
            //         return
            //     }
            // })
            for (var i = State.viViews[State.viCtrlIndex].components.length - 1; i >= 0; i--) {
                let v = State.viViews[State.viCtrlIndex].components[i]
                if (v.id == id) {
                    State.viViews[State.viCtrlIndex].components.splice(i, 1)
                    return
                }
            }
        },
        //通过组件ID移动zIndex 组件上移一层
        to_up(State, id: string) {

            State.viViews[State.viCtrlIndex].components.forEach((v, index) => {
                if (v.id == id) {
                    let zindex = v.style.zIndex
                    // zindex = zindex + 1
                    let arr: number[] = []
                    State.viViews[State.viCtrlIndex].components.forEach(item => {
                        if (item.style.zIndex > zindex && id != item.id) {
                            arr.push(item.style.zIndex)
                        }
                    })
                    if (arr.length == 0) {
                        v.style.zIndex = zindex + 1
                    } else {
                        let min_zindex = Math.min(...arr)
                        if (min_zindex >= zindex) {
                            v.style.zIndex = min_zindex + 1
                        }
                    }
                    return
                }
            })


        },
        //通过组件ID移动zIndex 组件下移一层
        to_down(State, id: string) {
            State.viViews[State.viCtrlIndex].components.forEach((v, index) => {
                if (v.id == id) {
                    let zindex = v.style.zIndex
                    // zindex = zindex + 1
                    let arr: number[] = []
                    State.viViews[State.viCtrlIndex].components.forEach(item => {
                        if (item.style.zIndex < zindex) {
                            arr.push(item.style.zIndex)
                        }
                    })
                    if (arr.length == 0) {
                        if (v.style.zIndex > 0) {
                            v.style.zIndex = zindex - 1
                        }
                    } else {
                        let max_zindex = Math.max(...arr)
                        if (max_zindex > 0 && max_zindex <= zindex) {
                            v.style.zIndex = max_zindex - 1
                        }
                    }

                    return
                }
            })
        },
        //通过组件ID移动zIndex 组件移到顶层
        to_top(State, id: string) {
            State.viViews[State.viCtrlIndex].components.forEach((v, index) => {
                if (v.id == id) {
                    let arr: number[] = []
                    State.viViews[State.viCtrlIndex].components.forEach(item => {
                        arr.push(item.style.zIndex)
                    })
                    let res = Math.max(...arr)
                    v.style.zIndex = res + 1
                    return
                }
            })
        },
        //通过组件ID移动zIndex 组件移到底层
        to_bottom(State, id: string) {
            State.viViews[State.viCtrlIndex].components.forEach((v, index) => {
                if (v.id == id) {
                    let arr: number[] = []
                    State.viViews[State.viCtrlIndex].components.forEach(item => {
                        arr.push(item.style.zIndex)
                    })
                    let res = Math.min(...arr)
                    if (res > 0) {
                        v.style.zIndex = res - 1
                    } else {
                        v.style.zIndex = 0
                    }

                    return
                }
            })
        },
        //通过批量组件ID选择组件
        set_component_id_by_area(State, ids: string[]) {
            State.viViews[State.viCtrlIndex].componet_ids.length = 0
            ids.forEach(id => {
                State.viViews[State.viCtrlIndex].componet_ids.push(id)
            })
        },
        //通过组件ID选择组件
        set_component_id(State, id: string) {
            let id_index = State.viViews[State.viCtrlIndex].componet_ids.indexOf(id)
            if (State.viViews[State.viCtrlIndex].componet_ids.length > 1 && !State.control_pressed) {
                if (id_index == -1) {
                    State.viViews[State.viCtrlIndex].componet_ids.push(id)
                }
                return
            }
            if (State.control_pressed) {
                if (id_index == -1) {
                    State.viViews[State.viCtrlIndex].componet_ids.push(id)

                } else {
                    State.viViews[State.viCtrlIndex].componet_ids.splice(id_index, 1)
                }
                return
            }
            State.viViews[State.viCtrlIndex].componet_ids.length = 0
            State.viViews[State.viCtrlIndex].componet_ids.push(id)

        },
        //让当前视图选择组件数量为0
        unset_component_id(State) {
            State.viViews[State.viCtrlIndex].componet_ids.length = 0
        },
        //初始化保存编辑器引用，方便各组件内获取。
        init_editor(State, e: HTMLElement | null) {
            if (e) {
                State.viEditor = e
            }
        },
        //设置组件外部样式
        set_out_style(State, payload: {
            style: ViOutStyle,
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
                    v.style = payload.style
                    return
                }
            })
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
        //改变组内组件的样式，但已经弃用，函数未删除。
        set_group_style(State, payload: {
            style: { outStyle: ViOutStyle, innerStyle: ViInnerStyle, groupStyle: ViGroupStyle },
            addr?: { viewid: number, comid: string, groupid: string }
        }) {
            let vid = State.viCtrlIndex
            let cid = State.viViews[State.viCtrlIndex].componet_ids[0]
            if (payload.addr != undefined) {
                vid = payload.addr.viewid
                cid = payload.addr.comid
            }
            State.viViews[vid].components.forEach(v => {
                if (v.id == cid) {
                    v.components.forEach(c => {
                        if (c.id == payload.addr?.groupid) {
                            c.innerStyle = payload.style.innerStyle
                            c.style = payload.style.outStyle
                            c.groupStyle = payload.style.groupStyle
                        }
                    })
                    return
                }
            })
        },
        //通过组件索引批量删除当前视图的组件，组件数组索引，需要从数组后面删除。
        del_components_by_indexs(State, indexs: number[]) {
            for (var i = indexs.length - 1; i >= 0; i--) {
                State.viViews[State.viCtrlIndex].components.splice(indexs[i], 1)
            }
        },
        //清空复制缓存
        clear_copy(State) {
            State.viCopy.length = 0
        },
        //复制到缓存
        copy_to(State, comp: ViComponent) {
            State.viCopy.push(comp)
        },
        //记录快照
        snapshot(State) {
            let data = cloneDeep(State.viViews[State.viCtrlIndex])
            if (Object.keys(State.viSnapshot).indexOf(data.id) == -1) {
                State.viSnapshot[data.id] = []
            }
            State.viSnapshot[data.id].push(data)
            State.viViews.forEach(v => {
                if (v.id == data.id) {
                    v.snapindex++
                    State.viSnapshot[data.id].splice(v.snapindex + 1)
                    nextTick(() => {
                        if (State.viEditor != null) {
                            // 将 div 元素转换为 canvas 元素
                            html2canvas(State.viEditor, {
                            }).then((canvas) => {
                                // 在这里处理 canvas 元素
                                const ctx = canvas.getContext('2d');
                                if (ctx != null) {
                                    ctx.scale(0.005, 0.005)
                                    // let dataUrl = canvas.toDataURL("image/jpeg", 0.3)
                                    // WriteToCache(dataUrl).then(r => {
                                    //     State.viViews[State.viCtrlIndex].shot = r
                                    // })
                                    canvas.toBlob(function (blob) {
                                        // 创建URL
                                        if (blob != null) {
                                            const url = URL.createObjectURL(blob);
                                            if (State.viViews[State.viCtrlIndex].shot != "") {
                                                URL.revokeObjectURL(State.viViews[State.viCtrlIndex].shot)
                                            }
                                            State.viViews[State.viCtrlIndex].shot = url
                                        }

                                    }, "image/png");
                                }
                            });
                        }
                    })

                    return
                }
            })

        },
        //重做快照，不做快照数量限制，每次保存文件需调用清空快照
        redo_snapshot(State) {
            if (State.viViews[State.viCtrlIndex].snapindex < State.viSnapshot[State.viViews[State.viCtrlIndex].id].length - 1) {
                State.viViews[State.viCtrlIndex].snapindex++
                let data = cloneDeep(State.viSnapshot[State.viViews[State.viCtrlIndex].id][State.viViews[State.viCtrlIndex].snapindex])
                data.snapindex = State.viViews[State.viCtrlIndex].snapindex
                State.viViews.forEach((value, index) => {
                    if (value.id == data.id) {
                        State.viViews[index] = data
                        return
                    }
                })
            }
        },
        //撤销快照
        undo_snapshot(State) {
            if (State.viViews[State.viCtrlIndex].snapindex > 0) {
                State.viViews[State.viCtrlIndex].snapindex--
                let data = cloneDeep(State.viSnapshot[State.viViews[State.viCtrlIndex].id][State.viViews[State.viCtrlIndex].snapindex])
                data.snapindex = State.viViews[State.viCtrlIndex].snapindex
                State.viViews.forEach((value, index) => {
                    if (value.id == data.id) {
                        State.viViews[index] = data
                        return
                    }
                })

            }
        },
        //清空快照，每次保存到文件后，要清空快照，避免长期操作导致快照数组过大,
        //保留第一个快照是因为第一次初始化完成，会记录一次初始化快照。
        clear_snapshot(State) {
            Object.keys(State.viSnapshot).forEach(v => {
                State.viSnapshot[v].length = 0
            })
            State.viViews.forEach(v => {
                v.snapindex = -1
            })
        },
        //通过控制面板交换视图位置
        swap_by_index(State, payload: {
            from: number,
            to: number
        }) {
            let view = cloneDeep(State.viViews[payload.from])
            let view_to = State.viViews[payload.to]
            State.viViews.splice(payload.from, 1)
            let index = State.viViews.indexOf(view_to)
            State.viViews.splice(index, 0, view)
        },
        //显示对话框提示
        show_dialog(state, payload: {
            msg: string,
            show: boolean,
        }) {
            state.viMessage = payload.msg
            state.viShow = payload.show
        },
        //显示问题提问
        show_ask_dialog(state, payload: {
            msg: string,
            show: boolean,
        }) {
            state.viMessage = payload.msg
            state.viShowAsk = payload.show
        },
        //显示或隐藏加载对话框
        show_loader_dialog(state, payload: { show: boolean }) {
            state.viShowLoader = payload.show
        },
        //设置组件数量计时器
        set_counter(state) {
            state.viCounter = state.viCounter + 1
        },
        //设置显示哪个属性标签
        set_boarding(state, payload: { value: number }) {
            state.viBoarding = payload.value
        },
        //设置顶部抽屉菜单内容
        set_bread_items(state, payload: { value: string }) {
            state.viBreaditems.push(payload.value)
        },
        //删除顶部抽屉菜单部分内容
        del_bread_items(state, payload: { value: number }) {
            state.viBreaditems.splice(payload.value)
        },
        //顶部菜单抽屉面板索引，需要关联到chart组件的第n级options
        set_bread_index(state, payload: { value: number }) {
            state.viBreadindex = payload.value
            state.viBreaditems.splice(state.viBreadindex + 1)
        },
        //设置chartstyle坐标系统
        set_chart_style_coor(state, payload: {
            key: string
        }) {
            let id = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == id) {
                    if (payload.key == "xAxis" || payload.key == "yAxis") {
                        let axis = item.chartStyle?.style[payload.key] as echarts.XAXisComponentOption[] | echarts.YAXisComponentOption[]
                        axis.push({})
                    } else if (payload.key == "bmap") {
                        let obj = item.chartStyle?.style
                        if (obj) {
                            set(obj, "bmap.zoom", 15)
                        }
                    } else {
                        let axis = item.chartStyle?.style[payload.key] as echarts.XAXisComponentOption | echarts.YAXisComponentOption
                        axis = {}
                    }
                }
            })
        },
        //设置chartstyle coorNamelist
        set_chart_style_name_list(state, payload: {
            value: [{
                key: string,
                name: string,
                count: number,
            }
            ]
        }) {
            let id = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == id) {
                    item.chartStyle?.coorNameList.push(payload.value[0])

                    if (payload.value[0].key == "xAxis" || payload.value[0].key == "yAxis") {
                        item.chartStyle?.coorNameList.forEach(v => {
                            if (payload.value[0].key == v.key) {
                                v.count = payload.value[0].count
                            }
                        })
                    }
                    // console.log("add", payload.value[0].key, item.chartStyle?.style)
                    return
                }

            })
        },
        //删除hartstyle coorNamelist对应Index对应的内容
        del_chart_style_name_list(state, payload: { value: number }) {
            let id = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == id) {
                    let coor = item.chartStyle?.coorNameList[payload.value]
                    if (coor != undefined) {
                        if (coor.key == "xAxis" || coor.key == "yAxis") {
                            let lst = item.chartStyle?.coorNameList
                            let count = -1
                            if (lst != undefined) {
                                for (let i = 0; i <= payload.value; i++) {
                                    if (lst[i].key == coor.key) {
                                        count++
                                    }
                                }
                            } else {
                                count = 0
                            }
                            if (count > 1) {
                                count = 1
                            }
                            let option = item.chartStyle?.style
                            if (option != undefined) {
                                // unset(option, coor.key + "[" + count + "]")
                                let axis = option[coor.key] as echarts.XAXisComponentOption[] | echarts.YAXisComponentOption[]
                                axis.splice(count, 1)
                                if (coor.key != undefined) {
                                    let k = coor.key
                                    item.chartStyle?.coorNameList.forEach(v => {
                                        if (v.key == k) {
                                            v.count--
                                        }
                                    })
                                }
                            }
                            let path = coor.key + "[" + count + "].data"
                            let obj = item.chartStyle?.chartCoorList
                            if (obj != undefined) {
                                unset(obj, path)
                                // console.log("del_chart_style_name_list 1", item)
                            }

                        } else {
                            let option = item.chartStyle?.style
                            if (option != undefined) {
                                unset(option, coor.key)
                            }
                            let obj = item.chartStyle?.chartCoorList
                            if (obj != undefined) {
                                unset(obj, coor.key + ".data")
                                // console.log("del_chart_style_name_list 2", item)
                            }
                        }
                        item.chartStyle?.coorNameList.splice(payload.value, 1)
                        // console.log("del", coor.key, item.chartStyle?.style)
                        return
                    }
                }
            })
        },
        //设置chartstyle dataNameList
        set_chart_style_data_list(state, payload: {
            value: [{
                key: string,
                name: string,
                count: number,
            }
            ]
        }) {
            let id = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == id) {
                    item.chartStyle?.dataNameList.push(payload.value[0])
                    return
                }
            })
        },
        //删除hartstyle dataNameList对应Index对应的内容
        del_chart_style_data_list(state, payload: { value: number }) {
            let id = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == id) {
                    item.chartStyle?.dataNameList.splice(payload.value, 1)
                    let series = item.chartStyle?.style.series as echarts.SeriesOption[] | undefined
                    if (series != undefined) {
                        series.splice(payload.value, 1)
                    }
                    let obj = item.chartStyle?.chartSeriesList
                    if (obj != undefined) {
                        // obj.delete("series[" + payload.value + "].data")
                        obj.splice(payload.value, 1)
                        // console.log("del_chart_style_data_list", item)
                    }
                    return
                }
            })
        },
        //设置chartstyle otherNameList
        set_chart_style_other_list(state, payload: {
            value: [{
                key: string,
                name: string,
                count: number,
            }
            ]
        }) {
            let id = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == id) {
                    item.chartStyle?.otherNameList.push(payload.value[0])
                    return
                }
            })
        },
        //删除chartstyle otherNameList对应Index对应的内容
        del_chart_style_other_list(state, payload: { value: number }) {
            let id = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == id) {
                    item.chartStyle?.otherNameList.splice(payload.value, 1)
                    return
                }
            })
        },
        //增加echarts option 的other属性 ,一般是给定 string类型value用来当成Key,增加一个空对象
        add_chart_style_other(state, payload: { value: string }) {
            let id = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == id) {
                    // item.chartStyle?.otherNameList.splice(payload.value, 1)
                    let style = item.chartStyle?.style
                    if (style != undefined) {
                        if (payload.value == "continuous" || payload.value == "piecewise") {
                            set(style, "visualMap[0]", { type: payload.value })
                        } else {
                            set(style, payload.value, {})
                        }
                    }
                    // console.log("add_chart_style_other", style)
                    return
                }
            })
        },
        //删除echarts option 的other属性 ,一般是给定 string类型value用来当成Key删除echats optian属性
        del_chart_style_other(state, payload: { value: string }) {
            let id = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == id) {
                    // item.chartStyle?.otherNameList.splice(payload.value, 1)
                    let style = item.chartStyle?.style
                    if (style != undefined) {
                        unset(style, payload.value)
                    }
                    // console.log("del_chart_style_other", style)
                    return
                }
            })
        },
        //通过路径更新图表的echart option，例如xAxis[0].data=...
        set_echart_option_by_path(state,
            payload: { id: string, path: string, value: any, name?: string }) {
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                let id = payload.id
                if (!id) {
                    id = state.viViews[state.viCtrlIndex].componet_ids[0]
                }

                if (item.id == id) {
                    let obj = item.chartStyle?.style
                    if (obj != undefined) {
                        set(obj, payload.path, payload.value)
                        // console.log("set_echart_option_by_path 2", obj, payload.path, payload.value)
                        if (payload.name) {
                            set(obj, payload.path.replace(".data", ".name").replace(".indicator", ".name"), payload.name)
                        }
                    }
                    return
                }
            })
        },
        //增加series空数据类型，一般是 例如 {type:'line'}
        add_echart_series(state, payload: {
            type: any
        }) {
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == state.viViews[state.viCtrlIndex].componet_ids[0]) {
                    let obj = item.chartStyle?.style.series as echarts.SeriesOption[] | undefined
                    if (obj != undefined) {
                        if (typeof payload.type === 'object') {
                            obj.push(
                                payload.type
                            )
                        } else {
                            obj.push({
                                type: payload.type
                            })
                        }
                        // console.log("add_echart_series", obj)
                    }
                    return
                }
            })
        },
        //通过索引更新图表的echart option的series，例如series[0]=...
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
                        } else if (series_type === 'radar') {
                            let data = payload.value as any[]
                            newdata = [{
                                value: data.flat()
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
        //保存echart对象属性的坐标是对应哪个表的哪个字段？
        set_echart_coor_map(state,
            payload: { path: string, value: any }) {
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == state.viViews[state.viCtrlIndex].componet_ids[0]) {
                    let obj = item.chartStyle?.chartCoorList
                    if (obj != undefined) {
                        // set(obj, payload.path, payload.value)
                        obj[payload.path] = payload.value
                        // obj.set(payload.path, payload.value)
                        // console.log("set_echart_coor_map", item)
                    }
                    return
                }
            })
        },
        //保存echart对象属性的数据是对应哪个表的哪个字段？
        set_echart_data_map(state,
            payload: {
                index: number, value: {
                    sheetkey: string;
                    fieldkey: string;
                }
            }) {
            state.viViews[state.viCtrlIndex].components.forEach(item => {
                if (item.id == state.viViews[state.viCtrlIndex].componet_ids[0]) {
                    let obj = item.chartStyle?.chartSeriesList
                    if (obj != undefined) {
                        // set(obj, payload.path, payload.value)
                        // obj[payload.path] = payload.value
                        // obj.set(payload.path, payload.value)
                        let value = payload.value
                        let idx = obj[payload.index]
                        if (idx === undefined) {
                            obj[payload.index] = {
                                sheetkey: value.sheetkey,
                                fieldkey: [value.fieldkey]
                            }
                        } else {
                            if (state.control_pressed) {
                                obj[payload.index].fieldkey.length = 0
                            }
                            if (obj[payload.index].fieldkey.indexOf(value.fieldkey) === -1) {
                                obj[payload.index].fieldkey.push(value.fieldkey)
                            }
                        }

                        // console.log("set_echart_data_map", obj)
                    }
                    return
                }
            })
        },
        //删除echart对象属性的data是对应哪个表的哪个字段？
        // del_echart_data_map(state,
        //     payload: { path: string }) {
        //     state.viViews[state.viCtrlIndex].components.forEach(item => {
        //         if (item.id == state.viViews[state.viCtrlIndex].componet_ids[0]) {
        //             let obj = item.chartStyle?.chartDataList
        //             if (obj != undefined) {
        //                 obj.delete(payload.path)
        //                 console.log("del_echart_data_map", item)
        //             }
        //             return
        //         }
        //     })
        // },

        //显示或隐藏v-html对话框
        set_show_html(state, payload: {
            show: boolean,
            msg: string
        }) {
            state.viHtmlContext = payload.msg
            state.viShowHtml = payload.show
        },
        //设置当前图表点击属性索引，可能是坐标属性、数据属性、其他属性
        set_ViChartRefIndex(state, payload: { value: ViChartType }) {
            state.ViChartRefIndex = payload.value
        },
        //设置函数对话框是否显示
        set_ViShowFunction(state, payload: {
            msg: string,
            show: boolean,
            path: string
        }) {
            state.ViShowFunction = payload.show
            state.ViFunctionContext = payload.msg
            state.ViFunctionPath = payload.path
        },
        //设置图表当前属性值
        set_viChartCurrentValue(state, payload: { value: any }) {
            state.viChartCurrentValue = payload.value
        },
        //修改组件名称
        set_comp_name(state, payload: { name: string }) {
            let selected = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(c => {
                if (c.id == selected) {
                    c.customName = payload.name
                    return
                }
                c.components.forEach(g => {
                    if (g.id == selected) {
                        g.customName = payload.name
                        return
                    }
                })
            })
        },
        //设置图表条件过滤id
        // set_chartConditionList(state,payload:{value:string[]}){
        //     let selected = state.viViews[state.viCtrlIndex].componet_ids[0]
        //     state.viViews[state.viCtrlIndex].components.forEach(c=>{
        //         if(c.id ==selected){
        //             if(c.chartStyle?.chartConditionList){
        //                 c.chartStyle?.chartConditionList.splice(0)
        //                 payload.value.forEach(v=>{
        //                     c.chartStyle?.chartConditionList.push(v)
        //                 })
        //                 console.log(111,c.chartStyle?.chartConditionList)
        //                 return
        //             }

        //         }
        //     })
        // }
        //设置innerstayle里的表字段属性
        set_innerStyle_list(state, payload: { id: string, sheetkey: string, fieldkey: string }) {
            // let selected = state.viViews[state.viCtrlIndex].componet_ids[0]
            state.viViews[state.viCtrlIndex].components.forEach(c => {
                if ("nni" + c.id == payload.id) {
                    if (c.name == 'ViData') {
                        if (c.innerStyle.list.length > 0) {
                            c.innerStyle.list.forEach((v, index) => {
                                if (v.sheetkey != payload.sheetkey) {
                                    c.innerStyle.list.length = 0
                                }
                                if (v.sheetkey == payload.sheetkey && v.fieldkey == payload.fieldkey) {
                                    c.innerStyle.list.splice(index, 1)
                                }
                            })
                        }
                        c.innerStyle.list.push({
                            sheetkey: payload.sheetkey,
                            fieldkey: payload.fieldkey,
                        })
                    } else {
                        c.innerStyle.list.length = 0
                        c.innerStyle.list.push({
                            sheetkey: payload.sheetkey,
                            fieldkey: payload.fieldkey,
                        })
                    }
                    // console.log(c.innerStyle.list)
                    return
                }
            })
        },
        //这个是专门给修改百度地图Key用的
        set_baidu_map_key(state, value: string) {
            state.viBaiduMapKey = value
        },
        //给geo组件注册地图
        set_jsongeomap(state, value) {
            let selected = state.viViews[state.viCtrlIndex].componet_ids[0]
            let comps = state.viViews[state.viCtrlIndex].components
            for (let index = 0; index < comps.length; index++) {
                const c = comps[index];
                if (c.id == selected) {
                    let series = c.chartStyle?.style.series as echarts.SeriesOption[]
                    for (let index = 0; index < series.length; index++) {
                        const obj = series[index] as echarts.MapSeriesOption;
                        if (obj.coordinateSystem == 'geo' || obj.type == 'map') {
                            echarts.registerMap(c.id, value)
                            obj.map = c.id
                            break
                        }
                    }
                    if (c.chartStyle?.style.geo) {
                        let obj = c.chartStyle?.style.geo as echarts.GeoComponentOption
                        echarts.registerMap(c.id, value)
                        obj.map = c.id
                    }
                    return
                }
                c.components.forEach(g => {
                    if (g.id == selected) {
                        let series = g.chartStyle?.style.series as echarts.SeriesOption[]
                        for (let index = 0; index < series.length; index++) {
                            const obj = series[index] as echarts.MapSeriesOption;
                            if (obj.coordinateSystem == 'geo' || obj.type == 'map') {
                                echarts.registerMap(g.id, value)
                                obj.map = g.id
                                break
                            }
                        }
                        if (g.chartStyle?.style.geo) {
                            let obj = g.chartStyle?.style.geo as echarts.GeoComponentOption
                            echarts.registerMap(g.id, value)
                            obj.map = g.id
                        }
                        return
                    }
                })

            }
        },
        //设置是否显示帮助
        set_viShowHelp(state, value: boolean) {
            state.viShowHelp = value
        },
        //设置是否显示版权
        set_viShowCopyRight(state, value: boolean) {
            state.viShowCopyRight = value
        },
        //设置是否显示图表帮助
        set_viShowChartHelp(state, payload: {
            show: true,
            msg: string
        }) {
            state.viShowChartHelp = payload.show
            state.viMessage = payload.msg
        },
        //设置是否显示更新对话框
        set_viShowUpdate(state, value: boolean) {
            state.viShowUpdate = value
        },
        //设置是否显示license窗口
        set_viShowLicense(state, value: boolean) {
            state.viShowLicense = value
        }
    }

})

