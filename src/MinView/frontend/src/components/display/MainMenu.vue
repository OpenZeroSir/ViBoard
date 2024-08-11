<script lang="ts" setup>
import { nextTick, ref, watch } from 'vue'
import {
    mdiAlarm,
    mdiContentSave,
    mdiContentSaveAll,
    mdiFullscreen,
    mdiOnepassword,
    mdiOpenInApp,
    mdiOpenInNew,
    mdiRecord, mdiRotate3dVariant, mdiShieldLock,
} from '@mdi/js'
import { storeDisplay } from '../../store/display';
import { computed, reactive, getCurrentInstance } from 'vue';
import { cloneDeep } from 'lodash'
import { SelectFileDialog, RequestViewCodeInfo, RequestSavePath, ReqestSaveViewCode, SetupAppMode } from '../../../wailsjs/go/main/App';
import { WindowFullscreen, WindowUnfullscreen } from '../../../wailsjs/runtime/runtime';
import serialize from 'serialize-javascript';
const instance = getCurrentInstance();

function handleOpen() {
    let author = {
        author: "",
        email: "",
        company: "",
        phone: "",
        website: "",
        docinfo: "",
        passwd: "",
        filepath: "",
    }
    storeDisplay.commit("set_preinfo", author)
    SelectFileDialog("*.NNIC").then(path => {
        if (path.indexOf("Error:") != -1) {
            storeDisplay.commit("show_dialog", { show: true, msg: path })
        } else {

            RequestViewCodeInfo(path).then(value => {
                if (value.indexOf("Error:") != -1) {
                    storeDisplay.commit("show_dialog", { show: true, msg: value })
                } else {
                    let info = JSON.parse(value)
                    let author = {
                        author: info["author"],
                        email: info["email"],
                        company: info["company"],
                        phone: info["phone"],
                        website: info["website"],
                        docinfo: info["docinfo"],
                        passwd: "",
                        filepath: path,
                    }
                    storeDisplay.commit("set_preinfo", author)
                    storeDisplay.commit("set_viInitScale", info["scale"])
                    storeDisplay.commit("show_dialog_ask", true)
                }
            }).catch(reason => {
                storeDisplay.commit("show_dialog", { show: true, msg: reason })
            })
        }
    }).catch(reason => {
        storeDisplay.commit("show_dialog", { show: true, msg: reason })
    })
}

function handleSaveAs() {
    RequestSavePath("", "*.NNIC").then(v => {
        if (v.indexOf("Error:") != -1) {
            storeDisplay.commit("show_dialog", { show: true, msg: v })
        } else {
            storeDisplay.commit("show_loader", true)
            // const objStr = JSON.stringify(storeDisplay.state.viViews);
            let obj = {
                baidukey: storeDisplay.state.viBaiduMapKey,
                views: storeDisplay.state.viViews,
            }
            const objStr = serialize(obj)
            ReqestSaveViewCode(v, objStr, storeDisplay.state.preInfo.passwd).then(value => {
                storeDisplay.commit("show_loader", false)
                storeDisplay.commit("set_FilePath", v)
                storeDisplay.commit("show_dialog", { show: true, msg: value })
            }).catch(reason => {
                storeDisplay.commit("show_loader", false)
                storeDisplay.commit("show_dialog", { show: true, msg: reason })
            })
        }
    }).catch(reason => {
        storeDisplay.commit("show_dialog", { show: true, msg: reason })
    })
}

function handleSave() {
    storeDisplay.commit("show_loader", true)
    // const objStr = JSON.stringify(storeDisplay.state.viViews);
    let obj = {
        baidukey: storeDisplay.state.viBaiduMapKey,
        views: storeDisplay.state.viViews,
    }
    const objStr = serialize(obj)
    ReqestSaveViewCode(storeDisplay.state.viFilePath, objStr, storeDisplay.state.preInfo.passwd).then(value => {
        storeDisplay.commit("show_loader", false)
        storeDisplay.commit("show_dialog", { show: true, msg: value })
    }).catch(reason => {
        storeDisplay.commit("show_loader", false)
        storeDisplay.commit("show_dialog", { show: true, msg: reason })
    })
}

function handleChangePasswd() {
    if (storeDisplay.state.viFilePath == "") {
        storeDisplay.commit("show_dialog", { show: true, msg: "请先加载视图编码文件，才能修改对应文件密码。" })
        return
    }
    storeDisplay.commit("show_change_passwd_dialog", true)
}

function handleSwitchModel() {
    if (instance != null) {
        const proxy = instance.proxy
        if (proxy != null) {
            // const router = proxy.$router
            // router.replace("/")
            SetupAppMode("/").then(v => {
                const router = proxy.$router
                router.replace("/")
            }).catch(reason => {
                console.log(reason)
            })
        }
    }
}



function setFullScreen() {
    var fullScreenDiv = document.getElementById("displayEditorID")
    var isFullScreen = document.fullscreenElement
    if (!isFullScreen) {
        if (fullScreenDiv) {
            if (fullScreenDiv.requestFullscreen) {
                fullScreenDiv.requestFullscreen();
                WindowFullscreen()
                storeDisplay.commit("set_viIsFullScreen", true)
            }
        }

    } else {
        if (document.exitFullscreen) {
            storeDisplay.commit("set_viIsFullScreen", false)
            WindowUnfullscreen()
            document.exitFullscreen();
        }
    }
}

</script>

<template>
    <v-main>
        <v-app-bar dense color="light-blue-darken-4" tile flat>

            <v-row>
                <v-col style="margin-left: 10px;">
                    <v-btn variant="tonal" density="compact" color="white" elevation="0" @click="handleOpen">
                        <template v-slot:prepend>
                            <v-icon color="green">{{ mdiOpenInNew }}</v-icon>
                        </template>
                        <strong>载入视图</strong>
                        <v-tooltip activator="parent" location="bottom">打开视图编码文件</v-tooltip>
                    </v-btn>

                    &nbsp;
                    <v-btn variant="tonal" density="compact" color="white" elevation="0" @click.stop="handleSaveAs"
                        :disabled="storeDisplay.state.viFilePath == ''">
                        <template v-slot:prepend>
                            <v-icon color="#ff9900">{{ mdiContentSaveAll }}</v-icon>
                        </template>
                        <strong>另存视图</strong>
                        <v-tooltip activator="parent" location="bottom">另保存视图编码文件</v-tooltip>
                    </v-btn>
                    &nbsp;
                    <v-btn variant="tonal" density="compact" color="white" elevation="0" @click.stop="handleSave"
                        :disabled="storeDisplay.state.viFilePath == ''">
                        <template v-slot:prepend>
                            <v-icon color="red">{{ mdiContentSave }}</v-icon>
                        </template>
                        <strong>保存视图</strong>
                        <v-tooltip activator="parent" location="bottom">保存视图编码文件</v-tooltip>
                    </v-btn>
                </v-col>

                <v-col style="text-align: right;margin-right: 10px;">
                    <v-btn variant="tonal" density="compact" elevation="0" @click="handleChangePasswd">
                        <template v-slot:prepend>
                            <v-icon color="red">{{ mdiShieldLock }}</v-icon>
                        </template>
                        <strong>重置密码</strong>
                        <v-tooltip activator="parent" location="top">修改视图编码文件的加载密码，如果修改后忘记了密码，基本不可能恢复。</v-tooltip>
                    </v-btn>
                    &nbsp;
                    <v-btn variant="tonal" density="compact" color="white" elevation="0" @click="setFullScreen">
                        <template v-slot:prepend>
                            <v-icon>{{ mdiFullscreen }}</v-icon>
                        </template>
                        <strong>全屏显示</strong>
                        <v-tooltip activator="parent" location="bottom">全屏显示视图</v-tooltip>
                    </v-btn>
                    &nbsp;
                    <v-btn variant="tonal" density="compact" color="white" elevation="0" @click="handleSwitchModel">
                        <template v-slot:prepend>
                            <v-icon>{{ mdiRotate3dVariant }}</v-icon>
                        </template>
                        <strong>切换模式</strong>
                        <v-tooltip activator="parent" location="bottom">编辑模式/显示模式</v-tooltip>
                    </v-btn>
                </v-col>
            </v-row>


        </v-app-bar>
    </v-main>
</template>

