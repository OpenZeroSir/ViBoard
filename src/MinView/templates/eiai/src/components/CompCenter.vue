<template>
    <v-pagination :length="viewidlist.length" v-model="pageindex"
        style="background-color: #003355;color: white;"></v-pagination>
    <form ref="input_file_form" v-show="false" style="display: none;">
        <v-file-input ref="input_file_object" v-model="files" variant="underlined"></v-file-input>
    </form>
    <v-expansion-panels v-model="panelmodel" @update:modelValue="onPanelsChange">
        <v-expansion-panel v-for="(data, index) in complist" :key="index" :value="index" v-show="data['name'] != 'ViGroup'">
            <v-expansion-panel-title>
                <v-row>
                    <v-col>
                        组件编号：{{ data["id"] }}
                    </v-col>
                    <v-col style="text-align: center;">
                        {{ data["innerStyle"]["selectList"] ? data["innerStyle"]["selectList"] : data["innerStyle"]["text"]
                        }}
                    </v-col>
                    <v-col style="text-align: right;">
                        组件类型：{{ data["typeName"] }}
                    </v-col>
                </v-row>
            </v-expansion-panel-title>
            <v-expansion-panel-text>
                <v-row style="color: red;">
                    <v-col>
                        组件编号
                    </v-col>
                    <v-col>
                        原始内容
                    </v-col>
                    <v-col>
                        目标内容
                    </v-col>
                    <v-col style="text-align: right;">

                    </v-col>
                </v-row>
                <v-row>
                    <v-col>
                        <div style="display: flex; height: 100%;align-items: center;">
                            {{ data["id"] }}
                        </div>
                    </v-col>
                    <v-col>
                        <div style="display: flex; height: 100%;align-items: center;">
                            {{ ((data['name'] == 'ViList') || (data['name'] == 'ViCheckBox') || (data['name'] == 'ViTimer'))
                                ? checkbox_list_old_select :
                                data["innerStyle"]["text"] }}
                        </div>
                    </v-col>
                    <v-col>
                        <v-chip-group v-if="data['name'] == 'ViCheckBox'" v-model="multi_select" multiple
                            selected-class="text-primary" variant="tonal">
                            <v-chip v-for="(chip, index) in checkbox_list" :key="index" :value="chip">
                                {{ chip }}
                            </v-chip>
                        </v-chip-group>
                        <v-chip-group v-else-if="data['name'] == 'ViList'" v-model="multi_select" multiple
                            selected-class="text-primary" variant="tonal">
                            <v-chip v-for="(chip, index) in checkbox_list" :key="index" :value="chip">
                                {{ chip }}
                            </v-chip>
                        </v-chip-group>
                        <v-chip-group v-else-if="data['name'] == 'ViRadio'" v-model="single_select"
                            selected-class="text-primary" variant="tonal">
                            <v-chip v-for="(chip, index) in checkbox_list" :key="index" :value="chip">
                                {{ chip }}
                            </v-chip>
                        </v-chip-group>
                        <div v-else>
                            <div v-if="data['name'] == 'ViTimer'">
                                <v-row>
                                    <v-text-field type="date" v-model="multi_select[0]" variant="underlined"
                                        density="compact" flat>
                                    </v-text-field>
                                    <v-text-field type="time" v-model="multi_select[1]" variant="underlined"
                                        density="compact" flat>
                                    </v-text-field>
                                </v-row>
                            </div>
                            <v-text-field v-else-if="data['name'] == 'ViFile'" variant="underlined" density="compact" flat
                                append-inner-icon="md:attach_file" @click:append-inner="file_select(data['id'], '*')"
                                v-model="single_select"></v-text-field>
                            <v-text-field v-else-if="data['name'] == 'ViImage'" variant="underlined" density="compact" flat
                                append-inner-icon="md:image_file" @click:append-inner="file_select(data['id'], '.jpg,.png')"
                                v-model="single_select"></v-text-field>
                            <v-text-field v-else-if="data['name'] == 'ViVideo'" variant="underlined" density="compact" flat
                                append-inner-icon="md:video_file"
                                @click:append-inner="file_select(data['id'], '.webm,.mp4')"
                                v-model="single_select"></v-text-field>
                            <v-text-field v-else-if="data['name'] == 'ViAudio'" variant="underlined" density="compact" flat
                                append-inner-icon="md:audio_file" @click:append-inner="file_select(data['id'], '.mp3')"
                                v-model="single_select"></v-text-field>
                            <v-text-field v-else variant="underlined" density="compact" flat v-model="single_select">
                            </v-text-field>
                        </div>
                    </v-col>
                    <v-col style="text-align: right;">
                        <v-btn variant="tonal" color="#ff9900" @click.stop="updateComponents(data['id'], data['name'])">
                            确认修改
                        </v-btn>
                    </v-col>
                </v-row>

            </v-expansion-panel-text>
        </v-expansion-panel>

    </v-expansion-panels>
    <!-- {{ multi_select }}{{ single_select }} {{ viewidlist }} -->
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, computed, watch, getCurrentInstance } from 'vue';
import { nanoid } from 'nanoid'

let files = ref([] as File[])

//当前页的激活组件索引
let panelmodel = ref(-1)
//单值组件的值
let single_select = ref("")
//多值组件的已经选择的内容
let multi_select = ref<string[]>([])

let pageindex = ref(1)
let viewlist: any[] = reactive([])

function nniid(){
    return "NIBoard" + nanoid(9)
}

let complist = computed((() => {
    let ret: any[] = []
    for (const key in viewlist[pageindex.value - 1]) {
        if (viewlist[pageindex.value - 1].hasOwnProperty(key)) {
            const value = viewlist[pageindex.value - 1][key];
            // console.log('Key:', key, 'Value:', value);
            if (key == "components") {
                // console.log('Key:', key, 'Value:', value);
                ret = value
            }
        }
    }
    // console.log("ret111",ret.filter(v=>v.innerStyle.styleFlag))
    return ret.filter(v=>v.innerStyle.styleFlag)
}))

let viewidlist = computed(() => {
    let ret: string[] = []
    for (const key in viewlist[pageindex.value - 1]) {
        if (viewlist[pageindex.value - 1].hasOwnProperty(key)) {
            const value = viewlist[pageindex.value - 1][key];
            // console.log('Key:', key, 'Value:', value);
            if (key == "id") {
                ret.push(value)
            }
        }
    }
    return ret
})

watch(files, (newValue) => {
    if (newValue.length) {
        handleFileChange()
    }
})


let checkbox_list = computed(() => {
    if (panelmodel.value > -1) {
        let data = complist.value[panelmodel.value]
        let stringArray: string[] = data['innerStyle']['textList'] as string[];
        // let innerStyle = data['innerStyle']
        if (stringArray == undefined) {
            return []
        }
        return stringArray
    }
    return []
})
let checkbox_list_old_select = computed(() => {
    if (panelmodel.value > -1) {
        let data = complist.value[panelmodel.value]
        let stringArray: string[] = data['innerStyle']['selectList'] as string[];
        // let innerStyle = data['innerStyle']
        if (stringArray == undefined) {
            stringArray = data['innerStyle']['textList'] as string[];
            if (stringArray == undefined) {
                return []
            }

        }
        return stringArray
    }

    return []
})

function onPanelsChange() {
    single_select.value = ""
    multi_select.value.length = 0
}


function GetComponents() {

    const formData = new FormData();
    formData.append("type", "PostTypeDisplayGetComponent");
    fetch("/api/nnipost", {
        method: 'POST',
        // headers: {
        //     'Content-Type': 'application/json',
        // },
        body: formData
    })
        .then(result => {
            result.json().then(v => {
                viewlist.length = 0
                // viewlist=v
                v.forEach((data: any) => {
                    viewlist.push(data)
                    // complist.push(data)
                })
            }).catch(reason => {
                alert(reason)
            })
        })
        .catch(error => {
            // 处理错误
            alert(error)
        });
}
function updateComponents(key: string, name: string) {

    const formData = new FormData();
    formData.append("type", "Component");
    formData.append("key", "nni" + key)
    if (name == 'ViCheckBox' || name == 'ViList' || name == 'ViTimer') {
        let select = JSON.stringify(multi_select.value)
        if (multi_select.value.length == 0) {
            select = JSON.stringify([""])
        }
        formData.append("data", select)
    } else {

        let arr: string[] = []
        arr.push(single_select.value)
        let select = JSON.stringify(arr)
        formData.append("data", select)
    }
    formData.append("comtype", name)
    // console.log(formData.get("data")?.toString())
    fetch("/api/nnipost", {
        method: 'POST',
        // headers: {
        //     'Content-Type': 'application/json',
        // },
        body: formData
    })
        .then(result => {
            result.text().then(v => {
                alert(v)
                GetComponents()
            }).catch(reason => {
                alert(reason)
            })
        })
        .catch(error => {
            // 处理错误
            alert(error)
        });
}
let instance = getCurrentInstance()
onMounted(() => {
    if (instance != null) {
        input_file = instance.refs.input_file_object as HTMLInputElement
    }
    GetComponents()
})


interface PostFile {
    type: string,
    //视图id
    viewid: string,
    //组件id
    compid: string,
    //文件数据
    data: File,
}

let input_file: HTMLInputElement
let current_compid = ref('')
// 提交文件给后端程序，后端程序写入数据库后，返回路径给前端，由前端更新路径给store
function file_select(compid: string, ext: string) {
    current_compid.value = compid
    input_file.setAttribute("accept", ext)
    files.value = []
    single_select.value = ""
    input_file.click()

}

function handleFileChange() {
    // if (input_file != null) {
    //     if (input_file.files != null) {
    let file = files.value[0]
    if (file != undefined) {
        let data: PostFile = {
            type: "File",
            viewid: viewidlist.value[pageindex.value - 1],
            compid: current_compid.value,
            data: file
        }
        PostTypeImportFile(data)
    }
    //     }
    // }
}

function PostTypeImportFile(data: PostFile) {
    const formData = new FormData();
    formData.append("type", data.type);
    formData.append("viewid", data.viewid);
    formData.append("compid", data.compid);
    formData.append("File", data.data);
    fetch("/api/nnipost", {
        method: 'POST',
        // headers: {
        //     'Content-Type': 'multipart/form-data',
        // },
        body: formData
    }).then(result => {
        result.text().then(v => {
            single_select.value = nniid()+v
        }).catch(error => {
            alert(error)
        });
    }).catch(error => {
        alert(error)
    });
}
</script>