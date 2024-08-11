<script lang="ts" setup>
import { ref, computed, reactive } from 'vue'
import { store } from "../../store"
import { nextTick, watch, onMounted, onUnmounted } from 'vue';
import { fetchEX, isCeilAI, isUUID } from '../../utils/func';
import { publicStore } from '../../store/public';
import { OpenFileDialog, SelectFileDialog, ImportDataToDB, GetVOLUME, DeleteVOLUME } from "../../../wailsjs/go/main/App"
import { cloneDeep, keys } from 'lodash'
import { string } from 'mathjs';
import { isTSNumberKeyword } from '@babel/types';
import { mdiBackspace, mdiBackspaceOutline, mdiSelectDrag, mdiChevronRight, mdiClose, mdiImport, mdiInformation, mdiRecord, mdiReload } from '@mdi/js';
// import {
//     mdiContentSaveEdit,
//     mdiRecord
// } from '@mdi/js'

import { ViOutStyle } from '../../store/index'
import { OutStyleInfo } from '../../components/custom/option'
import OutStyleBoard from './rightboard/OutStyleBoard.vue'
import InnerStyleBoard from './rightboard/InnerStyleBoard.vue';
import { nanoid } from 'nanoid';
import { EventsOff, EventsOn } from '../../../wailsjs/runtime/runtime';

let right_drawer = ref(true)
let right_rail = ref(false)

let background_color_show = ref(false)

let background_selected = computed(() => {
    if (store.state.viViews.length > 0) {
        return store.state.viViews[store.state.viCtrlIndex].componet_ids.length < 1
    }
    return true
})

let editor_height = computed({
    get: () => store.state.canvas.height,
    set: (val) => {
        store.commit("set_editor_height", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

let editor_width = computed({
    get: () => store.state.canvas.width,
    set: (val) => {
        store.commit("set_editor_width", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

let background_color = computed({
    get: () => {
        let color = "#003355"
        if (store.state.viViews.length) {
            color = store.state.viViews[store.state.viCtrlIndex].background_color
        }
        return color
    },
    set: (val) => {
        store.commit("set_background_color", val)
        nextTick(() => {
            store.commit("snapshot")
        })
    }
});

// let background_image_files = ref<File[]>([])
// watch(background_image_files,
//     (newValue: File[], oldValue: File[]) => {
//         if (newValue.length > 0) {
//             fileToBase64(newValue[0]).then((base64) => {
//                 store.commit("set_background_image", base64)
//                 nextTick(() => {
//                     store.commit("snapshot")
//                 })
//                 background_label.value = "背景图片"
//             })
//         } else {
//             store.commit("set_background_image", "none")
//             nextTick(() => {
//                 store.commit("snapshot")
//             })
//             background_label.value = "背景图片"
//         }
//         background_image_files.value.length = 0
//     }
// );
function SetBackgroundImage() {
    let viewid = store.state.viViews[store.state.viCtrlIndex].id
    OpenFileDialog("*.jpg;*.png", viewid, "", false).then((result: string) => {
        if (result.indexOf("Error") != -1) {
            store.commit("set_background_image", "")
        } else {
            fetchEX(result).then((url) => {
                store.commit("set_background_image", url)
            })
        }
        // store.commit("snapshot")

    })
}
function clearBackgroundImage(event: any) {
    store.commit("set_background_image", "none")
    // nextTick(() => {
    //     store.commit("snapshot")
    // })
}
//保存当前工作表的 id key
let current_sheet_key = ref('')

// let background_label = ref("背景图片")

//导入输入到编辑器
//文件类型 xlsx csv
function handleImportToEditor() {
    SelectFileDialog("*.xlsx;*.csv").then((res) => {
        if (res.indexOf("Error:") != -1) {
            store.commit("show_dialog", { show: true, msg: res })
        } else {
            store.commit("show_loader_dialog", { show: true })
            ImportDataToDB(res, false).then((v) => {
                if (v == "") {
                    store.commit("show_loader_dialog", { show: false })
                    queryFiles("", "导入成功")
                } else {
                    store.commit("show_loader_dialog", { show: false })
                    store.commit("show_dialog", { show: true, msg: v })
                }
            }).catch((reason) => {
                store.commit("show_loader_dialog", { show: false })
                store.commit("show_dialog", { show: true, msg: reason })
            })
        }

    }).catch((reason) => {
        store.commit("show_dialog", { show: true, msg: reason })
    })
}
async function queryFiles(key: string, info: string) {
    await GetVOLUME(key).then((value) => {
        let keys = Object.keys(value)
        let err = value["Error"]
        if (err != "") {
            store.commit("show_dialog", { show: true, msg: err })
        } else {
            sheet_names.length = 0
            keys.forEach(key => {
                if (key != "Error") {
                    sheet_names.push({ key: key, name: value[key] })
                }
            })
            if (info != "") {
                store.commit("show_dialog", { show: true, msg: info })
                // console.log(sheet_names)
            }
        }
    }).catch((reason => {
        store.commit("show_dialog", { show: true, msg: reason })
    }))
}
let showImportButton = ref(true)
function handleFileNameClick(key: string) {
    showImportButton.value = false
    if (isUUID(key)) {
        current_sheet_key.value = key
        queryFiles(key, "")
    }
}
function handleToFileList() {
    queryFiles("", "")
    showImportButton.value = true
}
function handleDeleteFileFromList(key: string) {
    DeleteVOLUME(key).then((v) => {
        if (v == "") {
            queryFiles("", "")
            showImportButton.value = true
        } else {
            store.commit("show_dialog", { show: true, msg: v })
        }
    }).catch(reason => {
        store.commit("show_dialog", { show: true, msg: reason })
    })
}

function SelectFile() {
    let viewid = store.state.viViews[store.state.viCtrlIndex].id
    let comid = store.state.viViews[store.state.viCtrlIndex].componet_ids[0]
    OpenFileDialog("*.*", viewid, comid, true).then((result: string) => {
        let id = store.state.viViews[store.state.viCtrlIndex].componet_ids[0]
        let comps = store.state.viViews[store.state.viCtrlIndex].components.filter(v => {
            return v.id == id
        })
        let a = comps[0].innerStyle

        let style = a //cloneDeep(a)
        if (result.indexOf("Error") != -1) {
            style.text = ""
            store.commit("set_inner_style", {
                style: reactive(style),
                addr: {
                    viewid: store.state.viCtrlIndex,
                    comid: comps[0].id
                }
            })
        } else {
            //这里是获取图片视频之类的需要显示到视图中，所以需要请求文件。
            // fetchEX(result).then((url) => {
            //     style.text = url
            //     store.commit("set_inner_style", {
            //         style: reactive(style),
            //     })
            // })
            //这里是显示附件名字之类的，因为只要文件名，所以不用请求文件。
            style.text = result
            store.commit("set_inner_style", {
                style: reactive(style),
            })

        }
        store.commit("snapshot")

    })
}


let company = computed({
    get: () => store.state.canvas.company,
    set: (val) => {
        store.commit("set_company", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

let author = computed({
    get: () => store.state.canvas.author,
    set: (val) => {
        store.commit("set_author", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

let email = computed({
    get: () => store.state.canvas.email,
    set: (val) => {
        store.commit("set_email", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

let phone = computed({
    get: () => store.state.canvas.phone,
    set: (val) => {
        store.commit("set_phone", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

let website = computed({
    get: () => store.state.canvas.website,
    set: (val) => {
        store.commit("set_website", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});
let docinfo = computed({
    get: () => store.state.canvas.docinfo,
    set: (val) => {
        store.commit("set_docinfo", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

let passwd = computed({
    get: () => store.state.canvas.passwd,
    set: (val) => {
        store.commit("set_passwd", val)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

let show_passwd = ref(false)
interface Sheet {
    key: string;
    name: string;
}
const sheet_names = reactive<Sheet[]>([
    // { key: 'aaa1', name: '你好1bb1bb1bb1bb1bb1bb1bb1bb1bb1bb1bb1bb1' },
    // { key: 'aaa11', name: 'bb11' }
]);


let onboarding = computed({
    get: () => store.state.viBoarding,
    set: (val) => {
        store.commit("set_boarding", { value: val })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});


let current_comp_name = computed({
    get: () => {
        let name = ""
        let selected = store.state.viViews[store.state.viCtrlIndex].componet_ids[0]
        store.state.viViews[store.state.viCtrlIndex].components.forEach(c => {
            if (c.id == selected) {
                name = c.customName
                return
            }
            c.components.forEach(g => {
                if (g.id == selected) {
                    name = g.customName
                    return
                }
            })
        })
        return name
    },
    set: (val) => {
        store.commit("set_comp_name", { name: val })
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    }
});

let stylelist = ["外部属性", "内部属性"]//, "组件动作"]

function dragStart(event: DragEvent, name: any) {
    // 将拖拽元素添加一个类，以便应用拖拽时的样式
    const target = event.target as HTMLElement;
    if (target) {
        // 移除拖拽时添加的类
        target.classList.add('dragging');
    }
    if (event.dataTransfer) {
        event.dataTransfer.dropEffect = 'copy'
        event.dataTransfer?.setData("application/json", JSON.stringify({
            type: 'data',
            //表名称
            sheet_key: current_sheet_key.value,
            //字段名称
            field_key: name
        }))
    }
}

function dragEnd(event: DragEvent) {
    // 移除拖拽时添加的类
    const target = event.target as HTMLElement;
    if (target) {
        // 移除拖拽时添加的类
        target.classList.remove('dragging');
    }
    event.preventDefault()
}

function nniid() {
    return "NIBoard" + nanoid(9)
}

function file_select(filetype: string) {
    let viewid = store.state.viViews[store.state.viCtrlIndex].id
    OpenFileDialog(filetype, viewid, "", false).then((v) => {
        let url = nniid() + v
        // console.log("url", url)
        // fetchEX(result).then((url) => {
        //     store.commit("set_background_image", url)
        // })
        store.commit("set_background_image", url)
        // nextTick(() => {
        //     store.commit("snapshot")
        // })
    })
}

// function dragLeave(event: any) {
//     event.dataTransfer.dropEffect = 'copy'
// }


// let getInfoHeight = computed(() => {
//     return document.body.scrollHeight-300
// })

let getInfoHeight = ref(768 - 378)

function updateInfoHeight() {
    getInfoHeight.value = document.body.scrollHeight - 378
}

async function getSheetList() {
    await GetVOLUME("").then((value) => {
        let keys = Object.keys(value)
        let err = value["Error"]
        if (err != "") {
            store.commit("show_dialog", { show: true, msg: err })
        } else {
            sheet_names.length = 0
            keys.forEach(key => {
                if (key != "Error") {
                    sheet_names.push({ key: key, name: value[key] })
                }
            })
        }
    }).catch((reason => {
        store.commit("show_dialog", { show: true, msg: reason })
    }))
}

onMounted(async () => {
    window.addEventListener("resize", updateInfoHeight)
    getSheetList()
    updateInfoHeight()
    EventsOn("UIOpenProject", getSheetList)
})

onUnmounted(() => {
    window.removeEventListener("resize", updateInfoHeight)
    EventsOff("UIOpenProject")
})

</script>

<template>
    <v-navigation-drawer app class="drawerstyle" disable-resize-watcher color="light-blue-darken-4" permanent
        location="right" v-model="right_drawer" :rail="right_rail" @click="right_rail = false">
        <v-list-item density="compact" nav>
            <template v-if="!right_rail" v-slot:prepend>
                <v-btn variant="text" icon @click.stop="right_rail = !right_rail">
                    <span className="material-icons">arrow_forward_ios</span>
                    <v-tooltip activator="parent" location="left">隐藏</v-tooltip>
                </v-btn>
            </template>
            <template v-else v-slot:append>
                <v-btn variant="text" icon @click.stop="right_rail = !right_rail">
                    <span className="material-icons">arrow_back_ios_new</span>
                    <v-tooltip activator="parent" location="left">显示</v-tooltip>
                </v-btn>
            </template>
            <v-list-item-title>
                <div style="width: 100%; text-align: right;padding: 5px;">
                    <!-- <span style="margin-right: 5px;margin-top: 15px;font-weight: bold; font-size: large; color: #ff9900;">Dataset</span> -->
                    <v-btn v-if="showImportButton" icon variant="tonal" density="compact" @click="handleImportToEditor"
                        color="white" elevation="0">
                        <!-- <template v-slot:prepend>
                    <v-icon>{{ mdiImport }}</v-icon>
                </template> -->
                        <!-- 导入数据 -->
                        <v-icon v-if="showImportButton" color="#ff9900">
                            {{ mdiImport }}

                        </v-icon>
                        <v-tooltip activator="parent" location="left"> 导入数据</v-tooltip>
                    </v-btn>
                    <v-btn v-else variant="tonal" icon density="compact" @click="handleToFileList" color="white"
                        elevation="0">
                        <!-- <template v-slot:prepend>
                    <v-icon>{{ mdiBackspace }}</v-icon>
                </template>
                回到列表 -->
                        <v-icon>{{ mdiBackspaceOutline }}</v-icon>
                        <v-tooltip activator="parent" location="left"> 回到列表</v-tooltip>
                    </v-btn>
                </div>
            </v-list-item-title>
        </v-list-item>

        <!-- <v-divider></v-divider> -->
        <div v-if="right_rail == false" style="padding-left: 10px; padding-right: 10px;">
            <v-sheet height="200px" color="#003355" tile flat>
                <v-virtual-scroll height="200" :items="sheet_names">
                    <!-- <v-list lines="one" item-height="24"> -->
                    <template v-slot:default="{ item }">
                        <v-list-item style="color:white;" nav v-ripple>
                            <template v-slot:prepend>
                                <v-btn icon draggable="true" v-show="item.name.indexOf(item.key) > -1 ? true : false"
                                    @dragstart="dragStart($event, item.key)" @dragend="dragEnd($event)" variant="plant"
                                    density="compact" color="green" elevation="0">
                                    <v-icon>{{ mdiSelectDrag }}</v-icon>
                                    <v-tooltip activator="parent" location="left">拖拉本字段到图表对应属性上</v-tooltip>
                                </v-btn>
                            </template>

                            <template v-slot:append>
                                <v-btn v-if="showImportButton" @click.stop="handleDeleteFileFromList(item.key)" icon
                                    variant="plant" density="compact" color="red" elevation="0">
                                    <v-icon>{{ mdiClose }}</v-icon>
                                    <v-tooltip activator="parent" location="left"> 删除 {{ item.name }}</v-tooltip>
                                </v-btn>
                            </template>
                            <!-- <v-btn ripple variant="tonal" density="compact" color="white" elevation="0" width="100%">
                                {{ item.name }}
                            </v-btn> -->
                            <div @click.stop="handleFileNameClick(item.key)" style="white-space: nowrap; /* 禁止换行 */
                                                overflow: hidden; /* 超出部分隐藏 */
                                                text-overflow: ellipsis; /* 超出部分显示省略号 */">
                                {{ item.name }}
                                <v-tooltip activator="parent" location="left"> {{ item.name }}</v-tooltip>
                            </div>

                        </v-list-item>
                    </template>
                    <!-- </v-list> -->
                </v-virtual-scroll>
            </v-sheet>
        </div>
        <div v-if="right_rail == false">
            <div v-if="background_selected"
                :style="{ padding: '10px', height: getInfoHeight + 'px', overflow: 'auto !important', display: 'flex' }">

                <v-col>
                    <v-row>
                        <v-text-field v-model="editor_width" type="number" label="画布宽度" variant="underlined"
                            hide-details></v-text-field>
                    </v-row>
                    <v-row>
                        <v-text-field v-model="editor_height" type="number" label="画布高度" variant="underlined"
                            hide-details></v-text-field>
                    </v-row>
                    <v-row>
                        <v-text-field v-model="background_color" active variant="underlined" label="背景颜色" hide-details
                            @click:append-inner="background_color_show = true" :color="background_color" clearable
                            append-inner-icon="md:color_lens"></v-text-field>

                        <v-dialog v-model="background_color_show" width="auto">
                            <v-color-picker v-model="background_color" hide-inputs show-swatches></v-color-picker>
                        </v-dialog>
                    </v-row>
                    <v-row>
                        <!-- <v-file-input prepend-icon="" @click:control="clearBackgroundImage($event)" :label="background_label"
                        accept="image/*" persistent-clear v-model="background_image_files" variant="underlined"
                        append-inner-icon="md:image" single-line></v-file-input> -->
                        <!-- <v-icon>{{ mdiContentSaveEdit }}</v-icon> -->
                        <v-text-field @click="file_select('*.jpg;*.png')" @click:append-inner="file_select('*.jpg;*.png')"
                            type="text" readonly label="背景图片" variant="underlined" hide-details
                            append-inner-icon="md:image"></v-text-field>
                    </v-row>
                    <v-row>
                        <v-text-field v-model="passwd" :type="show_passwd ? 'text' : 'password'"
                            append-inner-icon="md:remove_red_eye" @click:append-inner="show_passwd = !show_passwd"
                            label="视图密码" variant="underlined" hide-details></v-text-field>
                    </v-row>
                    <v-row>
                        <v-text-field v-model="company" type="text" label="组织名称" variant="underlined"
                            hide-details></v-text-field>
                    </v-row>
                    <v-row>
                        <v-text-field v-model="author" type="text" label="人员姓名" variant="underlined"
                            hide-details></v-text-field>
                    </v-row>
                    <v-row>
                        <v-text-field v-model="phone" type="text" label="联系电话" variant="underlined"
                            hide-details></v-text-field>
                    </v-row>
                    <v-row>
                        <v-text-field v-model="email" type="text" label="联系邮箱" variant="underlined"
                            hide-details></v-text-field>
                    </v-row>
                    <v-row>
                        <v-text-field v-model="website" type="text" label="服务网址" variant="underlined"
                            hide-details></v-text-field>
                    </v-row>
                    <v-row>
                        <v-text-field v-model="docinfo" type="text" label="视图描述" variant="underlined"
                            hide-details></v-text-field>
                        <v-tooltip activator="parent"
                            location="left">您最好预留视图描述信息，方便您在不知道是什么视图编码文件时，可以通过视图描述回忆起视图密码。</v-tooltip>
                    </v-row>



                </v-col>

            </div>
            <div v-else>
                <!-- <v-expansion-panels variant="accordion">
                    <v-row>
                        <v-text-field @click="SelectFile" type="text" readonly label="选择文件" variant="underlined"
                            hide-details></v-text-field>
                    </v-row>

                </v-expansion-panels> -->
                <v-text-field style="padding-left: 10px;padding-right: 10px;" v-model="current_comp_name" type="text"
                    variant="underlined" hide-details></v-text-field>
                <v-card flat rounded="0"
                    style="margin-left: 0px;margin-right: 0px;background-color: transparent;padding: 0;margin-top: 0;">
                    <v-card-actions class="justify-space-between">
                        <v-item-group v-model="onboarding" class="text-center" mandatory>
                            <v-item v-for="n in 2" :key="`btn-${n}`" v-slot="{ isSelected, toggle }" :value="n">
                                <v-btn :variant="isSelected ? 'outlined' : 'text'" icon density="compact"
                                    style="color: white;" @click="toggle">
                                    <v-icon>{{ mdiRecord }}</v-icon>
                                    <v-tooltip activator="parent" location="top"> {{ stylelist[n - 1] }}</v-tooltip>
                                </v-btn>
                            </v-item>
                        </v-item-group>

                    </v-card-actions>
                    <v-window v-model="onboarding" continuous
                        :style="{ height: getInfoHeight - 50 + 'px', overflowY: 'auto !important', display: 'flex', width: '100%' }">
                        <OutStyleBoard />
                        <InnerStyleBoard />
                    </v-window>
                </v-card>
            </div>
        </div>


    </v-navigation-drawer>
</template>


<style>
.v-navigation-drawer__content {
    overflow-y: hidden !important;
    height: 100%;
}

.v-window__container {
    width: 100% !important;
    padding-left: 10px;
    padding-right: 10px;
}
</style>