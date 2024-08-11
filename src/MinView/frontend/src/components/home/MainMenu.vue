<script lang="ts" setup>
import { nextTick, ref, watchEffect } from 'vue'
import {
    mdiContentSave, mdiRecord, mdiGroup, mdiUngroup,
    mdiAccount, mdiAlignVerticalBottom, mdiAlignVerticalCenter,
    mdiAlignVerticalTop, mdiAlignHorizontalRight,
    mdiAlignHorizontalCenter, mdiAlignHorizontalLeft,
    mdiOpenInNew, mdiFolderOpenOutline, mdiExport,
    mdiContentCopy, mdiContentCut, mdiContentPaste, mdiContentDuplicate,
    mdiMagnify, mdiFindReplace, mdiRedo, mdiUndo, mdiChevronDown,
    mdiChevronUp, mdiFormatVerticalAlignTop, mdiFormatVerticalAlignBottom,
    mdiDotsVerticalCircleOutline, mdiRotate3dVariant, mdiContentSaveEdit, mdiRouter
} from '@mdi/js'
import { store, ViComponent, ViView } from '../../store';
import { computed, reactive, getCurrentInstance } from 'vue';
import { nanoid } from 'nanoid';
import { createRenderItemAPI, getComponentRotatedStyle, getElement, mod360, nanoidEx, sleep, uNGroupStyle } from '../../utils/func';
import { abs } from 'mathjs';
import { Components } from "../custom/index";
import { cloneDeep } from 'lodash'
import { RequestSavePath, RequestSaveProject, SelectFileDialog, ReqestOpenProjectData, ReqestNewProject, ReqestExportViewCode, AppUIWriteFile, GetVOLUME } from '../../../wailsjs/go/main/App';
import { useConfirm, useSnackbar } from 'vuetify-use-dialog'
import { SetupAppMode } from '../../../wailsjs/go/main/App'
import { watch } from 'vue';
import serialize from 'serialize-javascript';
import html2canvas from 'html2canvas';
import { onMounted } from 'vue';
import { EventsEmit, Quit } from '../../../wailsjs/runtime/runtime';

const instance = getCurrentInstance();
const createConfirm = useConfirm()

let tab = ref(0)

let active_group = computed(() => {
    return !(store.state.viViews[store.state.viCtrlIndex].componet_ids.length > 1)
})
let un_selected_shape = computed(() => {
    return !(store.state.viViews[store.state.viCtrlIndex].componet_ids.length > 0)
})
let active_ungroup = computed(() => {
    let res = false
    if (store.state.viViews[store.state.viCtrlIndex].componet_ids.length === 1) {
        let id = store.state.viViews[store.state.viCtrlIndex].componet_ids[0]
        store.state.viViews[store.state.viCtrlIndex].components.forEach(v => {
            if (v.id == id) {
                if (v.name == "ViGroup") {
                    res = true
                    return
                }
            }
        })
    }
    return !res
})
let active_undo = computed(() => {
    // if(isNaN(store.state.viCtrlIndex)){
    //     store.commit("set_viCtrlIndex",0)
    // }
    if (store.state.viViews.length > 0) {

        return store.state.viViews[store.state.viCtrlIndex].snapindex >= 1
    }
    return false
})
let active_redo = computed(() => {
    if (store.state.viViews.length > 0) {
        let id = store.state.viViews[store.state.viCtrlIndex].id
        if (Object.keys(store.state.viSnapshot).indexOf(id) != -1) {
            return store.state.viViews[store.state.viCtrlIndex].snapindex >= store.state.viSnapshot[id].length - 1
        }
    }
    return false
})
let active_save = computed(() => {
    let res = false
    Object.keys(store.state.viSnapshot).forEach(v => {
        if (store.state.viSnapshot[v].length > 1) {
            res = true
            return
        }
    })
    return res
})

async function handleNew() {
    const isConfirmed = await createConfirm({
        title: "新建工程",
        content: "该操作会导致当前工程被覆盖，确定吗？",
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
        ReqestNewProject().then(result => {
            if (result != "") {
                store.commit("show_dialog", { show: true, msg: "1" + result })
            } else {
                store.commit("clear_state")
                store.commit("snapshot")
            }
        }).catch(reason => {
            store.commit("show_dialog", { show: true, msg: "2" + reason })
        })
    }
}

async function handleOpen() {
    const isConfirmed = await createConfirm({
        title: "打开工程",
        content: "该操作会导致当前工程被覆盖，确定吗？",
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
        SelectFileDialog("*.NNIE").then(path => {
            if (path.indexOf("Error:") > -1) {
                store.commit("show_dialog", { show: true, msg: path })
            } else {
                store.commit("show_loader_dialog", { show: true })
                ReqestOpenProjectData(path).then(result => {
                    if (path.indexOf("Error:") > -1) {
                        store.commit("show_loader_dialog", { show: false })
                        store.commit("show_dialog", { show: true, msg: result })
                    } else {
                        //读取成功
                        // let res = JSON.parse(result)
                        let res = new Function('return ' + result)()
                        console.log("res", res)
                        let views = res.viViews as ViView[]
                        for (var i = 0; i < views.length; i++) {
                            for (var j = 0; j < views[i].components.length; j++) {

                                let obj = views[i].components[j].chartStyle?.style.series as echarts.SeriesOption[] | undefined
                                if (obj) {
                                    let bmap = views[i].components[j].chartStyle?.style['bmap']
                                    for (let index = 0; index < obj.length; index++) {
                                        const element = obj[index] as any;
                                        if (element.type == 'custom') {
                                            element.renderItem = createRenderItemAPI(bmap)
                                        }

                                    }
                                }
                                let length = views[i].components[j].components.length
                                for (var k = 0; k < length; k++) {
                                    let comp = views[i].components[j].components[k]
                                    if (comp.chartStyle) {
                                        let obj = comp.chartStyle.style.series as echarts.SeriesOption[] | undefined
                                        if (obj) {
                                            let bmap = comp.chartStyle.style['bmap']
                                            for (let index = 0; index < obj.length; index++) {
                                                const element = obj[index] as any;
                                                if (element.type == 'custom') {
                                                    element.renderItem = createRenderItemAPI(bmap)
                                                }
                                            }
                                        }
                                    }

                                }
                            }
                        }
                        store.commit("set_state", res)
                        store.commit("set_viFilePath", path)
                        store.commit("snapshot")
                        store.commit("show_loader_dialog", { show: false })
                        EventsEmit("UIOpenProject")
                    }
                }).catch(reason => {
                    store.commit("show_loader_dialog", { show: false })
                    store.commit("show_dialog", { show: true, msg: reason })
                })
            }
        }).catch(reason => {
            store.commit("show_loader_dialog", { show: false })
            store.commit("show_dialog", { show: true, msg: reason })
        })
    }
}

function handleSave() {

    if (store.state.viFilePath == "") {

        RequestSavePath("", "*.NNIE").then(v => {
            if (v.indexOf("Error:") != -1) {
                store.commit("show_dialog", { show: true, msg: v })
            } else {
                store.commit("clear_snapshot")
                // const objStr = JSON.stringify(store.state);
                const objStr = serialize(store.state)
                store.commit("show_loader_dialog", { show: true })
                RequestSaveProject(v, objStr).then(v => {
                    store.commit("show_loader_dialog", { show: false })
                    store.commit("show_dialog", { show: true, msg: v })
                }).catch(reason => {
                    store.commit("show_loader_dialog", { show: false })
                    store.commit("show_dialog", { show: true, msg: reason })
                })
                store.commit("set_viFilePath", v)
                // nextTick(() => {
                //     store.commit("snapshot")
                // })
            }
        }).catch(reason => {
            store.commit("show_dialog", { show: true, msg: reason })
        })
    } else {
        store.commit("clear_snapshot")
        // const objStr = JSON.stringify(store.state);
        const objStr = serialize(store.state)
        store.commit("show_loader_dialog", { show: true })
        RequestSaveProject(store.state.viFilePath, objStr).then(v => {
            store.commit("show_loader_dialog", { show: false })
            store.commit("show_dialog", { show: true, msg: v })
        }).catch(reason => {
            store.commit("show_loader_dialog", { show: false })
            store.commit("show_dialog", { show: true, msg: reason })
        })
    }

}

function handleSaveAs() {
    RequestSavePath("", "*.NNIE").then(v => {
        if (v.indexOf("Error:") != -1) {
            store.commit("show_dialog", { show: true, msg: v })
        } else {
            store.commit("clear_snapshot")
            // const objStr = JSON.stringify(store.state);
            const objStr = serialize(store.state)
            store.commit("show_loader_dialog", { show: true })
            RequestSaveProject(v, objStr).then(v => {
                store.commit("show_loader_dialog", { show: false })
                store.commit("show_dialog", { show: true, msg: v })
            }).catch(reason => {
                store.commit("show_loader_dialog", { show: false })
                store.commit("show_dialog", { show: true, msg: reason })
            })
            store.commit("set_viFilePath", v)
        }
    }).catch(reason => {
        store.commit("show_dialog", { show: true, msg: reason })
    })
}

function getShot(index: number) {
    let el = document.getElementById("displayContextID")
    if (el != null) {
        // 将 div 元素转换为 canvas 元素
        html2canvas(el, {
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
                        // const uint8Array = new Uint8Array(blob);
                        // 使用FileReader对象读取Blob数据
                        const fileReader = new FileReader();
                        fileReader.onload = function (event: any) {
                            // 获取读取到的Blob数据
                            const blobData = event.target.result;

                            // 将Blob数据转换为Uint8Array
                            const uint8Array = new Uint8Array(blobData);

                            // 将Uint8Array作为参数发送到后端
                            // AppUIWriteFile(uint8Array)
                        };

                        // 将Blob对象读取为ArrayBuffer
                        fileReader.readAsArrayBuffer(blob);

                        return "1x"
                    }
                }, "image/png");
            }
        });
    }
}

onMounted(() => {
    // readblob()
})

function readblob() {
    const blob = new Blob(["Hello, world!"], { type: "text/plain" });
    const fileReader = new FileReader();
    fileReader.onload = function (event: any) {
        // 获取读取到的Blob数据
        const blobData = event.target.result;

        // 将Blob数据转换为Uint8Array
        const uint8Array = new Uint8Array(blobData);

        // 将Uint8Array作为参数发送到后端
        AppUIWriteFile("sss", "ddd", `[${uint8Array.toString()}]`).then(v => {
            console.log("AppUIWriteFile", v)
        })
    };

    // 将Blob对象读取为ArrayBuffer
    fileReader.readAsArrayBuffer(blob);
}

function readFile(blob: Blob): Promise<Uint8Array> {
    return new Promise((resolve, reject) => {
        const fileReader = new FileReader();

        fileReader.onload = function (event: any) {
            const blobData = event.target.result;
            const uint8Array = new Uint8Array(blobData);
            resolve(uint8Array);
        };
        // fileReader.onerror = function (error) {
        //   reject(error);
        // };

        fileReader.readAsArrayBuffer(blob);
    });
}

function nniid() {
    return "NIBoard" + nanoid(9)
}

async function captureScreenshot(el: HTMLElement, v: ViView) {
    const canvas = await html2canvas(el, {});
    const ctx = canvas.getContext('2d');
    if (ctx != null) {
        ctx.scale(0.001, 0.001);
        const blob: Blob = await new Promise((resolve) => {
            canvas.toBlob(function (blob) {
                if (blob) {
                    resolve(blob);
                }
            }, "image/png");
        });
        if (blob != null) {
            const data = await readFile(blob);
            const url = await AppUIWriteFile(v.id, "shot", `[${data.toString()}]`);
            v.shot = nniid() + url;
            if (url.indexOf("Error") > -1) {
                v.shot = ""
            }
            // console.log("v.shot", v.shot);
        }
    }
}

function handleExportCode() {
    RequestSavePath("", "*.NNIC").then(async v => {
        if (v.indexOf("Error:") != -1) {
            store.commit("show_dialog", { show: true, msg: v })
        } else {
            store.commit("show_loader_dialog", { show: true })
            let views = cloneDeep(store.state.viViews)
            for (var i = 0; i < views.length; i++) {
                let v = views[i]
                for (var j = 0; j < v.components.length; j++) {
                    let c = v.components[j]
                    c.style.displayMode = true
                    for (var k = 0; k < c.components.length; k++) {
                        let g = c.components[k]
                        g.style.displayMode = true
                    }
                }
                v.shot = ""
                let el = document.getElementById("EditorID")
                if (el != null) {
                    store.commit("set_viCtrlIndex", i)
                    await sleep(3000)
                    await captureScreenshot(el, v)
                }
            }
            let obj = {
                baidukey: store.state.viBaiduMapKey,
                views: views,
            }
            // const objStr = JSON.stringify(obj);
            const docStr = JSON.stringify(store.state.canvas);
            const objStr = serialize(obj)
            // const deserializedContent = new Function('return ' + objStr)();
            ReqestExportViewCode(v, objStr, docStr, store.state.canvas.passwd).then(value => {
                store.commit("show_loader_dialog", { show: false })
                store.commit("show_dialog", { show: true, msg: value })
            }).catch(reason => {
                store.commit("show_loader_dialog", { show: false })
                store.commit("show_dialog", { show: true, msg: reason })
            })
        }
    }).catch(reason => {
        store.commit("show_dialog", { show: true, msg: reason })
    })
}

function handleGroup(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    let comps: ViComponent[] = []
    let indexs: number[] = []
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach((comp, index) => {
            let rotate = comp.style.rotate
            if (comp.id == id) {
                if (comp.components.length > 0) {
                    let edite_rect = store.state.viEditor!.getBoundingClientRect()
                    comp.components.forEach(v => {
                        let c = { ...v }
                        let el = getElement(c.id)
                        if (el == null) {
                            return
                        }
                        let r = mod360(rotate + c.style.rotate)
                        const { top, left } = uNGroupStyle(el.getBoundingClientRect(), edite_rect, comp.style, c.groupStyle)
                        c.style.top = top
                        c.style.left = left
                        c.style.width = c.groupStyle.width * comp.style.width / 100
                        c.style.height = c.groupStyle.height * comp.style.height / 100
                        c.style.rotate = r
                        comps.push(c)
                    })
                } else {
                    comps.push(comp)
                }
                indexs.push(index)
            }
        })
    })

    store.commit("del_components_by_indexs", indexs)
    let group: ViComponent = reactive(cloneDeep(Components.ViGroup))
    group.id = nanoidEx()
    group.style.top = comps[0].style.top
    group.style.left = comps[0].style.left
    group.style.width = comps[0].style.width
    group.style.height = comps[0].style.height
    group.style.zIndex = store.state.viViews[store.state.viCtrlIndex].z_index

    let width_arr: number[] = []
    let height_arr: number[] = []
    let top_arr: number[] = []
    let left_arr: number[] = []
    let temp_groups: ViComponent[] = []
    comps.forEach(v => {
        let c = { ...v };
        let cs = { ...c.style }
        const { left, top, width, height } = getComponentRotatedStyle({
            rotate: cs.rotate,
            left: cs.left,
            top: cs.top,
            width: cs.width,
            height: cs.height,
            right: 0,
            bottom: 0,
        })
        width_arr.push(left + width)
        height_arr.push(top + height)
        top_arr.push(top)
        left_arr.push(left)
        temp_groups.push(c)
    })
    group.style.width = Math.max(...width_arr) - Math.min(...left_arr)
    group.style.height = Math.max(...height_arr) - Math.min(...top_arr)

    group.style.top = Math.min(...top_arr)
    group.style.left = Math.min(...left_arr)
    temp_groups.forEach(V => {
        let c = { ...V }
        c.groupStyle.width = c.style.width / group.style.width * 100
        c.groupStyle.height = c.style.height / group.style.height * 100
        c.groupStyle.top = ((c.style.top - group.style.top) / group.style.height) * 100
        c.groupStyle.left = ((c.style.left - group.style.left) / group.style.width) * 100
        group.components?.push(c)
    })
    // console.log(222,group)
    group.customName = group.typeName + "-" + store.state.viCounter.toString()
    store.commit("unset_component_id")
    store.commit("add_component", group)
    store.commit("set_counter")
    nextTick(() => {
        store.commit("snapshot")
    })
}
function handleUnGroup(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    let id = store.state.viViews[store.state.viCtrlIndex].componet_ids[0]
    store.commit("unset_component_id")
    store.state.viViews[store.state.viCtrlIndex].components.forEach(v => {
        if (v.id == id) {
            let comp = { ...v }
            let edite_rect = store.state.viEditor!.getBoundingClientRect()
            comp.components.forEach(c => {
                let s = { ...c }
                let r = mod360(s.style.rotate + comp.style.rotate)
                let el = getElement(c.id)
                if (el == null) {
                    return
                }
                const { top, left } = uNGroupStyle(el.getBoundingClientRect(), edite_rect, comp.style, c.groupStyle)
                s.style.top = top
                s.style.left = left
                s.style.rotate = r
                c.style.width = c.groupStyle.width * comp.style.width / 100
                c.style.height = c.groupStyle.height * comp.style.height / 100
                store.commit("add_component", reactive({ ...s }))
            })
            store.commit("unset_component_id")
            store.commit("del_component", id)
            return
        }
    })
    nextTick(() => {
        store.commit("snapshot")
    })
}

function handleAlignHorizontalLeft(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()

    let edite_rect = store.state.viEditor!.getBoundingClientRect()
    let pos = 0
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach(comp => {
            if (comp.id == id) {
                let style = { ...comp.style }
                if (pos == 0) {
                    pos = style.left
                }
                let el = getElement(comp.id)
                let rect = el!.getBoundingClientRect()
                let offset = rect.left - edite_rect.left
                if (style.rotate != 0) {
                    style.left = pos + abs(style.left - offset)
                } else {
                    style.left = pos
                }
                store.commit("set_out_style", {
                    style: reactive(style),
                    addr: {
                        viewid: store.state.viCtrlIndex,
                        comid: comp.id
                    }
                })
            }
        })
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}

function handleAlignHorizontalCenter(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    let edite_rect = store.state.viEditor!.getBoundingClientRect()
    let pos = 0
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach(comp => {
            if (comp.id == id) {
                let style = { ...comp.style }
                if (pos == 0) {
                    pos = style.left
                }
                let el = getElement(comp.id)
                let rect = el!.getBoundingClientRect()
                let offset = rect.left - edite_rect.left
                if (style.rotate != 0) {
                    style.left = pos + abs(style.left - offset) - rect.width / 2
                } else {
                    style.left = pos - rect.width / 2
                }
                store.commit("set_out_style", {
                    style: reactive(style),
                    addr: {
                        viewid: store.state.viCtrlIndex,
                        comid: comp.id
                    }
                })
            }
        })
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}

function handleAlignHorizontalRight(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    let edite_rect = store.state.viEditor!.getBoundingClientRect()
    let pos = 0
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach(comp => {
            if (comp.id == id) {
                let style = { ...comp.style }
                if (pos == 0) {
                    pos = style.left
                }
                let el = getElement(comp.id)
                let rect = el!.getBoundingClientRect()
                let offset = rect.left - edite_rect.left
                if (style.rotate != 0) {
                    style.left = pos + abs(style.left - offset) - rect.width
                } else {
                    style.left = pos - rect.width
                }
                store.commit("set_out_style", {
                    style: reactive(style),
                    addr: {
                        viewid: store.state.viCtrlIndex,
                        comid: comp.id
                    }
                })
            }
        })
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}

function handleAlignVerticalTop(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    let pos = 0
    let edite_rect = store.state.viEditor!.getBoundingClientRect()
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach(comp => {
            if (comp.id == id) {
                let style = { ...comp.style }
                if (pos == 0) {
                    pos = style.top
                }
                let el = getElement(comp.id)
                let rect = el!.getBoundingClientRect()
                let offset = rect.top - edite_rect.top
                if (style.rotate != 0) {
                    style.top = pos + abs(style.top - offset)
                } else {
                    style.top = pos
                }
                store.commit("set_out_style", {
                    style: reactive(style),
                    addr: {
                        viewid: store.state.viCtrlIndex,
                        comid: comp.id
                    }
                })
            }
        })
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}

function handleAlignVerticalCenter(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    let edite_rect = store.state.viEditor!.getBoundingClientRect()
    let pos = 0
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach(comp => {
            if (comp.id == id) {
                let style = { ...comp.style }
                if (pos == 0) {
                    pos = style.top
                }
                let el = getElement(comp.id)
                let rect = el!.getBoundingClientRect()
                let offset = rect.top - edite_rect.top
                if (style.rotate != 0) {
                    style.top = pos + abs(style.top - offset) - rect.height / 2
                } else {
                    style.top = pos - rect.height / 2
                }
                store.commit("set_out_style", {
                    style: reactive(style),
                    addr: {
                        viewid: store.state.viCtrlIndex,
                        comid: comp.id
                    }
                })
            }
        })
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}

function handleAlignVerticalBottom(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    let edite_rect = store.state.viEditor!.getBoundingClientRect()
    let pos = 0
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach(comp => {
            if (comp.id == id) {
                let style = { ...comp.style }
                if (pos == 0) {
                    pos = style.top
                }
                let el = getElement(comp.id)
                let rect = el!.getBoundingClientRect()
                let offset = rect.top - edite_rect.top
                if (style.rotate != 0) {
                    style.top = pos + abs(style.top - offset) - rect.height
                } else {
                    style.top = pos - rect.height
                }
                store.commit("set_out_style", {
                    style: reactive(style),
                    addr: {
                        viewid: store.state.viCtrlIndex,
                        comid: comp.id
                    }
                })
            }
        })
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}

function moveToUp(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.commit("to_up", id)
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}
function moveToDown(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.commit("to_down", id)
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}
function moveToTop(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.commit("to_top", id)
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}
function moveToBottom(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(id => {
        store.commit("to_bottom", id)
    })
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}
function handleRedo(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    store.commit("redo_snapshot")
}
function handleUndo(e: MouseEvent | KeyboardEvent) {
    e.preventDefault()
    e.stopPropagation()
    store.commit("undo_snapshot")
}

function handleSwitchModel() {
    if (instance != null) {
        const proxy = instance.proxy
        if (proxy != null) {
            SetupAppMode("/display").then(v => {
                const router = proxy.$router
                router.replace("/display")
            }).catch(reason => {
                console.log(reason)
            })

        }
    }
}

// watch(() => store.state.viViews[store.state.viCtrlIndex].componet_ids, (newValue) => {
//     if (typeof newValue == 'undefined') {
//         return
//     }
//     if (newValue.length == 0) {
//         tab.value = 0
//     } else if (newValue.length > 0) {
//         if (store.state.viBoarding == 2) {
//             tab.value = 2
//         } else {
//             tab.value = 1
//         }
//     }
// })

watchEffect(() => {
    if (store.state.viViews[store.state.viCtrlIndex]) {
        if (store.state.viViews[store.state.viCtrlIndex].componet_ids.length == 0) {
            tab.value = 0
        }else if (store.state.viViews[store.state.viCtrlIndex].componet_ids.length == 1) {
            if (store.state.viBoarding == 2) {
                // tab.value = 2
            } else {
                // tab.value = 1
            }
        } else if (store.state.viViews[store.state.viCtrlIndex].componet_ids.length > 1) {
            if (store.state.viBoarding == 2) {
                tab.value = 2
            } else {
                tab.value = 1
            }
        }
    }
})

watch(() => store.state.viBoarding, (newValue) => {

    if (newValue == 2) {
        if (current_comptype.value == "ViChart") {
            tab.value = 2
        } else {
            tab.value = 1
        }

    } else {
        if (store.state.viViews[store.state.viCtrlIndex].componet_ids.length == 0) {
            tab.value = 0
        } else {
            tab.value = 1
        }

    }
})

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

let breaditems_model = computed({
    get: () => store.state.viBreadindex,
    set: (val) => {
        if (val == undefined) {
            return
        }
        store.commit("set_bread_index", { value: val })
    }
});


async function exit_fuction() {
    // const isConfirmed = await createConfirm({
    //     title: "退出程序",
    //     content: "确定要退出程序吗？",
    //     confirmationText: "确定",
    //     cancellationText: "取消",
    //     dialogProps: {
    //         width: 'auto',
    //     },
    //     confirmationButtonProps: {
    //         color: 'red',
    //     },
    // })
    // if (isConfirmed) {
    //     Quit()
    // }
    Quit()
}

</script>

<template>
    <v-app-bar dense color="light-blue-darken-4" tile flat>


        <v-menu>
            <template v-slot:activator="{ props }">
                <v-app-bar-nav-icon v-bind="props"></v-app-bar-nav-icon>
            </template>
            <v-list>
                <v-list-item key="exit" @click.stop="exit_fuction">
                    <v-list-item-title>退出程序</v-list-item-title>
                </v-list-item>
            </v-list>
        </v-menu>
        <v-card-actions class="justify-space-between">
            <v-item-group v-model="tab" class="text-center" mandatory>
                <v-item :key="`menu-1`" v-slot="{ active, toggle }">
                    <v-btn :input-value="active" icon @click="toggle" color="yellow">
                        <v-icon>{{ mdiRecord }}</v-icon>
                        <v-tooltip activator="parent" location="bottom">工程</v-tooltip>
                    </v-btn>
                </v-item>
                <v-item :key="`menu-2`" v-slot="{ active, toggle }">
                    <v-btn :input-value="active" icon @click="toggle" color="green">
                        <v-icon>{{ mdiRecord }}</v-icon>
                        <v-tooltip activator="parent" location="bottom">控制</v-tooltip>
                    </v-btn>
                </v-item>
            </v-item-group>
        </v-card-actions>
        <v-spacer></v-spacer>
        <v-window v-model="tab">
            <v-window-item>
                <v-chip label link ripple size="64px" style="padding: 2px;margin: 2px;" @click.stop="handleNew">
                    <v-icon size="20px" center>
                        {{ mdiOpenInNew }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">新建工程</v-tooltip>
                </v-chip>
                <v-chip label link ripple size="64px" style="padding: 2px;margin: 2px;" @click.stop="handleOpen">
                    <v-icon size="20px" center>
                        {{ mdiFolderOpenOutline }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">打开工程</v-tooltip>
                </v-chip>
                &nbsp;
                <v-chip label link ripple size="64px" style="padding: 2px;margin: 2px;" @click.stop="handleSave">
                    <v-icon size="20px" center>
                        {{ mdiContentSave }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">保存工程</v-tooltip>
                </v-chip>
                <v-chip label link ripple size="64px" style="padding: 2px;margin: 2px;" @click.stop="handleSaveAs">
                    <v-icon size="20px" center>
                        {{ mdiContentSaveEdit }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">另存工程</v-tooltip>
                </v-chip>
                <v-chip label link ripple size="64px" style="padding: 2px;margin: 2px;" @click.stop="handleExportCode">
                    <v-icon size="20px" center>
                        {{ mdiExport }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">导出视图</v-tooltip>
                </v-chip>
                &nbsp;
                <v-chip label link ripple size="64px" :disabled="!active_undo" @click="handleUndo($event)"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiUndo }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">撤销</v-tooltip>
                </v-chip>
                <v-chip label link ripple size="64px" :disabled="active_redo" @click="handleRedo($event)"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiRedo }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">重做</v-tooltip>
                </v-chip>
            </v-window-item>


            <v-window-item>
                <!-- 组和分组 -->
                <v-chip @click="handleGroup($event)" :disabled="active_group" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiGroup }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">组合</v-tooltip>
                </v-chip>
                <v-chip :disabled="active_ungroup" @click="handleUnGroup($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiUngroup }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">拆分</v-tooltip>
                </v-chip>
                &nbsp;
                <!-- 层级移动 -->
                <v-chip :disabled="un_selected_shape" @click="moveToUp($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiChevronUp }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">上移一层</v-tooltip>
                </v-chip>
                <v-chip :disabled="un_selected_shape" @click="moveToDown($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiChevronDown }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">下移一层</v-tooltip>
                </v-chip>
                <v-chip :disabled="un_selected_shape" @click="moveToTop($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiFormatVerticalAlignTop }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">移到等层</v-tooltip>
                </v-chip>
                <v-chip :disabled="un_selected_shape" @click="moveToBottom($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiFormatVerticalAlignBottom }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">移到底层</v-tooltip>
                </v-chip>
                &nbsp;
                <!-- 对齐方式 -->
                <v-chip :disabled="active_group" @click="handleAlignHorizontalLeft($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiAlignHorizontalLeft }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">左对齐</v-tooltip>
                </v-chip>
                <v-chip :disabled="active_group" @click="handleAlignHorizontalCenter($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiAlignHorizontalCenter }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">垂直对齐</v-tooltip>
                </v-chip>
                <v-chip :disabled="active_group" @click="handleAlignHorizontalRight($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiAlignHorizontalRight }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">右对齐</v-tooltip>
                </v-chip>
                &nbsp;
                <v-chip :disabled="active_group" @click="handleAlignVerticalTop($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiAlignVerticalTop }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">上对齐</v-tooltip>
                </v-chip>
                <v-chip :disabled="active_group" @click="handleAlignVerticalCenter($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiAlignVerticalCenter }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">水平对齐</v-tooltip>
                </v-chip>

                <v-chip :disabled="active_group" @click="handleAlignVerticalBottom($event)" label link ripple size="64px"
                    style="padding: 2px;margin: 2px;">
                    <v-icon size="20px" center>
                        {{ mdiAlignVerticalBottom }}
                    </v-icon>
                    <v-tooltip activator="parent" location="bottom">下对齐</v-tooltip>
                </v-chip>
            </v-window-item>
            <v-window-item>
                <v-chip-group mandatory='force' variant="text" v-model="breaditems_model">
                    <v-chip variant="text" elevation="0" density="compact" v-for="item in store.state.viBreaditems" pill>{{
                        item }}</v-chip>
                </v-chip-group>
            </v-window-item>
        </v-window>
        <v-spacer></v-spacer>

        <v-card-actions class="justify-space-between">
            <v-btn icon color="white" elevation="0">
                <v-icon color="red">{{ mdiRecord }}</v-icon>
                <v-tooltip activator="parent" location="bottom">正在开发...</v-tooltip>
            </v-btn>
            <v-btn variant="tonal" density="compact" color="white" elevation="0" @click="handleSwitchModel">
                <template v-slot:prepend>
                    <v-icon>{{ mdiRotate3dVariant }}</v-icon>
                </template>
                切换模式
                <v-tooltip activator="parent" location="left">编辑模式/显示模式</v-tooltip>
            </v-btn>

        </v-card-actions>


    </v-app-bar>
</template>

<style scoped>
.menu3style {
    color: red;
}
</style>