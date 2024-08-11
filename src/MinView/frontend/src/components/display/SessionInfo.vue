<template>
    <div >
        <v-row class="bg-light-blue-darken-4"
            style="padding-left: 10px;padding-right: 10px;padding-bottom: 0px;padding-top: 0px;">
            <v-col style="text-align: left; width: 80%;" class="bg-light-blue-darken-4">
                <v-btn variant="tonal" density="compact" elevation="0" @click.stop="startRemote">
                    <template v-slot:prepend>
                        <v-icon :color="remote_icon_color">{{ SvgIconName[remote_icon as keyof typeof SvgIconName]
                        }}</v-icon>
                    </template>
                    <strong>{{ remote_label }}远程服务</strong>
                    <v-tooltip activator="parent" location="top">启动或停止远程更新数据功能，可以通过人工或程序远程定时推送数据到本程序显示。</v-tooltip>
                </v-btn>
                &nbsp;
                <v-chip variant="tonal" density="compact" color="white" @click.stop="copymsg(remote_ip+':'+remote_port)">
                    <strong>{{ remote_ip }}:{{ remote_port }}</strong>
                </v-chip>
                <v-btn variant="text" density="compact" color="green" elevation="0" :disabled="remote_label == '停止'" icon
                    @click.stop="refreshRemoteIPAddr">
                    <v-icon size="18px">{{ SvgIconName.mdiRefresh }}</v-icon>
                    <v-tooltip activator="parent" location="right">您电脑可能有多个IP地址，可以适用本功能切换不同IP地址。</v-tooltip>
                </v-btn>
            </v-col>
            <v-col style="text-align: center;" class="bg-light-blue-darken-4">
                <v-btn variant="text" density="compact" color="#ff9900" elevation="0" icon @click.stop="showCopyRight">
                    <v-icon size="18px">
                        {{ SvgIconName.mdiCopyright }}
                    </v-icon>
                </v-btn>
            </v-col>
            <v-col class="bg-light-blue-darken-4" style="text-align: right;margin-right: 0px;margin-left: 5px;">
                <v-row>
                    <v-col style="text-align: right;justify-items: right;align-items: right;">

                    </v-col>
                    <v-col>

                        <v-slider v-model="timer_slider" thumb-color="white" track-color="white" color="white" hide-details
                            thumb-size="10" min="1" max="99" elevation="0" density="compact" step="1">
                            <template v-slot:append>
                                <v-avatar color="white" size="26px" variant="tonal" pill density="compact">
                                    <div style="font-size:16px;">
                                        {{ timer_slider }}
                                    </div>
                                    <v-tooltip activator="parent" location="top">自动切换视图时间间隔，仅对循环切换或随机切换有效。</v-tooltip>
                                </v-avatar>
                            </template>
                            <template v-slot:prepend>
                                <v-btn variant="tonal" density="compact" size="26px" color="#ff9900" elevation="0" icon
                                    @click="playTypeClick">
                                    <v-icon size="18px">{{ SvgIconName[play_type as keyof typeof SvgIconName] }}</v-icon>
                                    <v-tooltip activator="parent" location="top">切换视图方式 手动 循环 随机</v-tooltip>
                                </v-btn>
                            </template>
                        </v-slider>
                    </v-col>
                </v-row>
            </v-col>
        </v-row>
        <div class="bg-light-blue-darken-4">

            <!-- <v-divider></v-divider> -->
            <div style="padding-left: 10px;padding-right: 10px;padding-bottom: 20px;padding-top: 0px;">
                <v-row>
                    <v-col>
                        <v-card v-ripple variant="tonal">
                            <v-card-text>
                                <v-row>
                                    <v-col cols="3">
                                        <div style="color: white;">
                                            <strong>账号：</strong>
                                            <span style="color: #ff9900;">
                                                <strong>nni</strong>
                                            </span>

                                        </div>
                                    </v-col>
                                    <v-col cols="5">
                                        <div style="color: white;">
                                            <strong>密码：</strong>
                                            <v-chip variant="text" density="compact" color="#ff9900"  @click.stop="copymsg(storeDisplay.state.viRemotePasswd)">
                                                <strong>{{ show_remote_passwd ? storeDisplay.state.viRemotePasswd :
                                                    '************' }}</strong>
                                            </v-chip>
                                        </div>
                                    </v-col>
                                    <v-col cols="2">
                                        <v-btn variant="text" density="compact" color="green" elevation="0" icon
                                            @click.stop="showRemotePasswd">
                                            <v-icon size="18px">{{ mdiEye }}</v-icon>
                                            <v-tooltip activator="parent" location="right">隐藏或显示远程登陆密码</v-tooltip>
                                        </v-btn>
                                        <v-btn variant="text" density="compact" color="green" elevation="0" icon
                                            @click.stop="refreshRemotePasswd">
                                            <v-icon size="18px">{{ SvgIconName.mdiRefresh }}</v-icon>
                                            <v-tooltip activator="parent" location="right">刷新远程登陆密码</v-tooltip>
                                        </v-btn>
                                    </v-col>
                                </v-row>
                            </v-card-text>
                            <v-card-subtitle>
                                <div style="color: white;padding-bottom: 10px;">适用人工远程定时刷新数据</div>
                            </v-card-subtitle>
                        </v-card>
                    </v-col>
                    <v-col>
                        <v-card v-ripple variant="tonal">
                            <v-card-text>
                                <div style="color: white;">
                                    <v-row>
                                        <v-col cols="8">
                                            <strong>TOKEN：</strong>
                                            <v-chip variant="text" density="compact" color="#ff9900" @click.stop="copymsg(storeDisplay.state.viRemoteToken)">
                                                <strong>{{ show_remete_token ? storeDisplay.state.viRemoteToken :
                                                    '********************************' }}</strong>
                                            </v-chip>
                                        </v-col>
                                        <v-col cols="2">
                                            <v-btn variant="text" density="compact" color="green" elevation="0" icon
                                                @click.stop="showRemoteToken">
                                                <v-icon size="18px">{{ mdiEye }}</v-icon>
                                                <v-tooltip activator="parent" location="right">显示或隐藏Token</v-tooltip>
                                            </v-btn>
                                            <v-btn variant="text" density="compact" color="green" elevation="0" icon
                                                @click.stop="refreshRemoteToken">
                                                <v-icon size="18px">{{ SvgIconName.mdiRefresh }}</v-icon>
                                                <v-tooltip activator="parent" location="right">刷新远程访问Token</v-tooltip>
                                            </v-btn>
                                        </v-col>
                                    </v-row>
                                </div>
                            </v-card-text>
                            <v-card-subtitle>
                                <div style="color: white;padding-bottom: 10px;">适用自动化远程定时刷新数据</div>
                            </v-card-subtitle>
                        </v-card>
                    </v-col>
                </v-row>
            </div>
            <!-- <v-divider ></v-divider> -->

        </div>
    </div>
</template>

<script lang="ts" setup>
import { mdiDotsCircle, mdiEye, } from '@mdi/js';
import * as SvgIconName from '@mdi/js'
import { ref } from 'vue';
import { storeDisplay } from '../../store/display';
import {
    RequestIPAddress, RequestGetRemotePasswd,
    RequestGetRemoteToken, ReqestResetRemoteSessionPasswd, ReqestResetRemoteSessionToken,
    StartRemoteWeb, StopRemoteWeb, StatusRemoteWeb
} from '../../../wailsjs/go/main/App';
import { onMounted } from 'vue';
import { computed,watch } from 'vue';
import { ClipboardSetText } from '../../../wailsjs/runtime/runtime';

function  copymsg(txt:string) {
    ClipboardSetText(txt).then(v=>{
        if(v){
            storeDisplay.commit("show_dialog",{show:true,msg:"复制完成"})
        }
    })
}

let play_types = [
    "mdiGestureTap",
    "mdiRepeatVariant",
    "mdiShuffleVariant",
]
let timer_slider = ref(5)
let play_type = ref("mdiGestureTap")


watch(() => timer_slider.value, (newValue) => {
    storeDisplay.commit("set_viTimeDelay",newValue)
})

function playTypeClick() {
    let index = play_types.indexOf(play_type.value) + 1
    if (index > 2) {
        index = 0
    }
    play_type.value = play_types[index]
    storeDisplay.commit("set_viPlayType",play_type.value)
}

function showCopyRight() {
    storeDisplay.commit("show_copyright", true)
}

let remote_icon = ref("mdiPlay")
let remote_label = ref("启动")
let remote_icon_color = ref("#ff0000")
function startRemote() {
    if (remote_icon.value == "mdiPlay") {
        StartRemoteWeb(remote_ip.value).then(v => {
            if (v.indexOf("启动成功") != -1) {
                let ress = v.split(":")
                remote_port.value = ress[1]
                remote_icon.value = "mdiStop"
                remote_label.value = "停止"
                remote_icon_color.value = "#00ff00"
                // storeDisplay.commit("show_dialog", { show: true, msg: v })
            } else {
                storeDisplay.commit("show_dialog", { show: true, msg: v })
            }
        }).catch(reason => {
            storeDisplay.commit("show_dialog", { show: true, msg: reason })
        })

    } else {
        StopRemoteWeb().then(v => {
            if (v == "停止成功") {
                remote_icon.value = "mdiPlay"
                remote_label.value = "启动"
                remote_icon_color.value = "#ff0000"
                // storeDisplay.commit("show_dialog", { show: true, msg: v })
            } else {
                storeDisplay.commit("show_dialog", { show: true, msg: v })
            }
        }).catch(reason => {
            storeDisplay.commit("show_dialog", { show: true, msg: reason })
        })

    }
}

let remote_ip = ref("127.0.0.1")
let remote_port = ref("50520")
function refreshRemoteIPAddr() {
    RequestIPAddress().then(v => {
        let result1 = v[0]
        if (result1.indexOf("Error:") != -1) {
            storeDisplay.commit("show_dialog", { show: true, msg: result1 })
        } else {
            let index = v.indexOf(remote_ip.value)
            if (index == -1) {
                remote_ip.value = v[0]
            } else {
                if (index + 1 > v.length - 1) {
                    index = 0
                } else {
                    index = index + 1
                }
                remote_ip.value = v[index]
            }

        }
    }).catch(reason => {
        storeDisplay.commit("show_dialog", { show: true, msg: reason })
    })
}

let show_remote_passwd = ref(false)
let show_remete_token = ref(false)

function showRemotePasswd() {
    show_remote_passwd.value = !show_remote_passwd.value
}
function showRemoteToken() {
    show_remete_token.value = !show_remete_token.value
}

function refreshRemotePasswd() {
    ReqestResetRemoteSessionPasswd().then(v => {
        storeDisplay.commit("set_remote_passwd", v)
    }).catch(reason => {
        storeDisplay.commit("show_dialog", { show: true, msg: reason })
    })
}
function refreshRemoteToken() {
    ReqestResetRemoteSessionToken().then(v => {
        storeDisplay.commit("set_remote_token", v)
    }).catch(reason => {
        storeDisplay.commit("show_dialog", { show: true, msg: reason })
    })
}

onMounted(() => {
    RequestGetRemotePasswd().then(v => {
        storeDisplay.commit("set_remote_passwd", v)
    }).catch(reason => {
        storeDisplay.commit("show_dialog", { show: true, msg: reason })
    });
    RequestGetRemoteToken().then(v => {
        storeDisplay.commit("set_remote_token", v)
    }).catch(reason => {
        storeDisplay.commit("show_dialog", { show: true, msg: reason })
    });

    StatusRemoteWeb().then(v => {
        if (v != -1) {
            remote_port.value = v.toString()
            remote_icon.value = "mdiStop"
            remote_label.value = "停止"
            remote_icon_color.value = "#00ff00"
        } else {
            remote_icon.value = "mdiPlay"
            remote_label.value = "启动"
            remote_icon_color.value = "#ff0000"
        }
    })
})

</script>