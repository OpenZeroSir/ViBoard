<template>
    <div class="text-center">
        <v-dialog v-model="storeDisplay.state.viShowAsk" width="500px" activator="parent" persistent>
            <v-card>
                <v-card-title>
                    <div style="color:red; width: 100%;text-align:center;">请求输入视图密码</div>
                </v-card-title>

                <v-card-text>
                    <v-text-field v-model="passwd" type="password" label="视图密码" variant="underlined"
                        color="green"></v-text-field>
                </v-card-text>
                <v-card-subtitle>
                    <div style="padding-left: 15px;padding-right: 15px; padding-bottom: 15px;">
                        <div>
                            所属公司：{{ storeDisplay.state.preInfo.company }}
                        </div>
                        <div>
                            所属人员：{{ storeDisplay.state.preInfo.author }}
                        </div>
                        <div>
                            联系邮箱：{{ storeDisplay.state.preInfo.email }}
                        </div>
                        <div>
                            联系电话：{{ storeDisplay.state.preInfo.phone }}
                        </div>
                        <div>
                            服务网址：{{ storeDisplay.state.preInfo.website }}
                        </div>
                        <div style="white-space: nowrap;overflow: hidden;text-overflow: ellipsis;">
                            视图描述：{{ storeDisplay.state.preInfo.docinfo }}
                            <v-tooltip activator="parent" location="bottom">{{ storeDisplay.state.canvas.docinfo
                            }}</v-tooltip>
                        </div>
                    </div>
                </v-card-subtitle>
                <v-card-actions>
                    <div
                        style="display:flex; justify-content: space-between;width: 100%;padding-left: 100px;padding-right: 100px;">
                        <v-btn color="green" variant="tonal" @click="handleLoad">载入视图</v-btn>
                        <v-btn color="red" variant="tonal" @click="handleClose">取消载入</v-btn>
                    </div>
                </v-card-actions>

            </v-card>
        </v-dialog>
    </div>
</template>

<script lang="ts" setup>
import { computed, nextTick } from 'vue';
import { ref } from 'vue'
import { storeDisplay } from '../../../store/display';
import { ReqestImportViewCode, ReqestResetRemoteSession } from '../../../../wailsjs/go/main/App';
import { ViView } from '../../../store';
import { createRenderItemAPI } from '../../../utils/func';
// import {EventsEmit} from '../../../../wailsjs/runtime'
let passwd = computed({
    get: () => storeDisplay.state.preInfo.passwd,
    set: (val) => {
        storeDisplay.commit("set_pre_passwd", val)
    }
});

let handleLoad = () => {
    storeDisplay.commit("show_dialog_ask", false)
    storeDisplay.commit("show_loader", true)
    ReqestImportViewCode(storeDisplay.state.preInfo.filepath, storeDisplay.state.preInfo.passwd).then(v => {
        if (v.Error != "") {
            storeDisplay.commit("show_loader", false)
            storeDisplay.commit("show_dialog", { show: true, msg: "解密失败(" + v.Error + ")" })
        } else {
            ReqestResetRemoteSession().then(rsession => {
                // console.log("vv", v)
                let canvas = JSON.parse(v.DocInfo)
                canvas["passwd"] = storeDisplay.state.preInfo.passwd
                storeDisplay.commit("set_canvas", canvas)
                // let compons = JSON.parse(v.ObjData)
                let compons = new Function('return ' + v.ObjData)()
                // console.log("import compons", compons)
                storeDisplay.commit("set_baidu_map_key", compons.baidukey)
                let views = compons.views as ViView[]
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

                storeDisplay.commit("set_viViews", views)
                storeDisplay.commit("set_view_index")
                handleResize()
                storeDisplay.commit("set_FilePath", storeDisplay.state.preInfo.filepath)
                storeDisplay.commit("set_remote_passwd", rsession.Passwd)
                storeDisplay.commit("set_remote_token", rsession.Token)
                storeDisplay.commit("show_loader", false)
                // storeDisplay.commit("show_dialog", { show: true, msg: compons })
                // EventsEmit("PostTypeDisplayLoadedCompo",compons)
                storeDisplay.commit("emit_component_objects")
            }).catch(reason => {
                storeDisplay.commit("show_loader", false)
                storeDisplay.commit("show_dialog", { show: true, msg: "111" + reason })
            });
        }
    }).catch(reason => {
        storeDisplay.commit("show_loader", false)
        storeDisplay.commit("show_dialog", { show: true, msg: reason })
    })
}

function handleResize() {
    nextTick(() => {

        if (storeDisplay.state.display == null) {
            storeDisplay.commit("set_rect", { height: 0, width: 0 })
            return
        }
        let rect = storeDisplay.state.display.getBoundingClientRect()
        let width = storeDisplay.state.canvas.width
        let height = storeDisplay.state.canvas.height

        let width_new = 0
        let height_new = 0
        let src_cale = height / width
        let width_scale = width / rect.width
        let height_scale = height / rect.height
        if (height_scale < width_scale) {
            width_new = rect.width
            height_new = width_new * src_cale
        } else {
            height_new = rect.height
            width_new = height_new / src_cale
        }
        storeDisplay.commit("set_rect", { width: width_new, height: height_new })

    })
}


let handleClose = () => {
    storeDisplay.commit("show_dialog_ask", false)
}
</script>