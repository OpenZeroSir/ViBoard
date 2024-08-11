<template>
    <!-- 用来显示需要更新的数据表或组件 -->

    <v-expansion-panels>
        <v-expansion-panel v-for="key in sheetkeys" :key="key">
            <v-expansion-panel-title>
                <v-row>
                    <v-col>
                        文件名称：{{ sheetobjs[key][0].FileName }}
                    </v-col>
                    <!-- <v-col>
                        <div style="text-align: center;">
                            {{ sheetobjs[key][0].PriKey }}
                        </div>
                    </v-col> -->
                    <v-col>
                        <div style="text-align: right;">表格数量：{{ sheetobjs[key].length }}</div>
                    </v-col>
                </v-row>
            </v-expansion-panel-title>
            <v-expansion-panel-text>
                <v-row style="color: #ff0000;font-weight: bold; display: flex;align-items: center;justify-content: center;">
                    <!-- <v-col>
                        文件编号
                    </v-col> -->
                    <v-col>
                        工作表编号
                    </v-col>
                    <v-col style="text-align: left;">
                        工作表名称
                    </v-col>
                    <v-col cols="auto">
                        <v-row>
                            <v-col>
                                <v-btn :disabled="sheetobjs[key].length < 2" variant="tonal" color="red"
                                    @click.stop="handleImportFileClick(key, false)">增量更新</v-btn>
                                <v-tooltip activator="parent"
                                    location="bottom">增量更新整个工作簿的所有工作表到数据库，需要确保您导入的文件内包含和数据库表同样的表数量、表名和结构。这会增加更新失败的风险，如果更新失败，可能部分表内容无法回退，需要放弃保存文件并重新加载文件，才能达到回退效果。</v-tooltip>
                            </v-col>
                            <v-col style="text-align: right;" cols="auto">
                                <v-btn :disabled="sheetobjs[key].length < 2" variant="tonal" color="red"
                                    @click.stop="handleImportFileClick(key, true)">覆盖更新</v-btn>
                                <v-tooltip activator="parent"
                                    location="bottom">覆盖更新整个工作簿的所有工作表到数据库，需要确保您导入的文件内包含和数据库表同样的表数量、表名和结构。这会增加更新失败的风险，如果更新失败，可能部分表内容无法回退，需要放弃保存文件并重新加载文件，才能达到回退效果。</v-tooltip>
                            </v-col>
                        </v-row>
                    </v-col>


                </v-row>

                <v-row v-for="sheet in sheetobjs[key]" style=" display: flex;align-items: center;justify-content: center;">
                    <v-divider></v-divider>
                    <!-- <v-col>
                        {{ sheet.PriKey }}
                    </v-col> -->
                    <v-col>
                        {{ sheet.SecKey }}
                    </v-col>
                    <v-col>
                        {{ sheet.SheetName }}
                    </v-col>
                    <v-col cols="auto">
                        <v-row>
                            <v-col>
                                <v-btn variant="tonal" color="#ff9900"
                                    @click.stop="handleImportFileClick(sheet.SecKey, false)">增量更新</v-btn>
                                <v-tooltip activator="parent" location="bottom">保留留原有数据基础上增加数据</v-tooltip>
                            </v-col>
                            <v-col style="text-align: right;" cols="auto">

                                <v-btn variant="tonal" color="#ff9900"
                                    @click.stop="handleImportFileClick(sheet.SecKey, true)">覆盖更新</v-btn>
                                <v-tooltip activator="parent" location="bottom">删除留原有数据并新增数据</v-tooltip>
                            </v-col>
                        </v-row>
                    </v-col>
                </v-row>
            </v-expansion-panel-text>
        </v-expansion-panel>
    </v-expansion-panels>
    <form ref="input_file_form" v-show="false" style="display: none;">
        <v-file-input ref="input_file_object" accept=".xlsx,*" @change="handleFileChange"
            variant="underlined"></v-file-input>
    </form>
</template>

<script setup lang="ts">
import { onMounted, reactive, getCurrentInstance, ref } from "vue";
let input_file: HTMLInputElement
let input_form: HTMLFormElement
let instance = getCurrentInstance()
onMounted(() => {
    if (instance != null) {
        input_file = instance.refs.input_file_object as HTMLInputElement
        input_form = instance.refs.input_file_form as HTMLFormElement
    }
    GetVolume()
})

function handleFileChange() {
    if (input_file != null) {
        if (input_file.files != null) {
            let file = input_file.files[0]
            if (file != undefined) {
                let data: PostSheet = {
                    type: "Sheet",
                    key: current_sheet_key.value,
                    replace: current_sheet_replace.value,
                    data: file
                }
                PostTypeImportData(data)
            }
        }
    }
}



interface SheetObject {
    PriKey: string;
    FileName: string;
    SecKey: string;
    SheetName: string;
}
let sheetkeys: string[] = reactive([])
let sheetobjs: { [key: string]: SheetObject[] } = reactive({})

function GetVolume() {
    // const data = {
    //     type: "PostTypeDisplayGetVolume",
    // };
    const formData = new FormData();
    formData.append("type", "PostTypeDisplayGetVolume");
    fetch("/api/nnipost", {
        method: 'POST',
        // headers: {
        //     'Content-Type': 'multipart/form-data',
        // },
        body: formData
    })
        .then(result => {
            if (result.status == 200) {
                result.json().then((obj: { [key: string]: SheetObject[] }) => {
                    // console.log(obj)
                    sheetobjs = reactive(obj)
                    let keys = Object.keys(obj)
                    keys.forEach(v => {
                        sheetkeys.push(v)
                    })
                    // sheetkeys = reactive(keys)
                }).catch(() => {
                    result.text().then(v => {
                        alert(v)
                    }).catch(error => {
                        alert(error)
                    });
                })
            } else {
                alert(result.statusText)
            }
        })
        .catch(error => {
            // 处理错误
            alert(error)
        });
}

interface PostSheet {
    type: string,
    //表的key
    key: string,
    //是否替换,覆盖和增量更新，默认增量更新
    replace: boolean,
    //文件数据
    data: File,
}

//当前表的key
let current_sheet_key = ref("")
let current_sheet_replace = ref(false)
function handleImportFileClick(key: string, replace: boolean) {
    current_sheet_key.value = key
    current_sheet_replace.value = replace
    input_form.reset()
    input_file.click()
}

function PostTypeImportData(data: PostSheet) {
    const formData = new FormData();
    formData.append("type", data.type);
    formData.append("key", data.key);
    formData.append("replace", data.replace.toString());
    formData.append("File", data.data);
    fetch("/api/nnipost", {
        method: 'POST',
        headers: {
            // 'Content-Type': 'multipart/form-data',
        },
        body: formData
    }).then(result => {
        result.text().then(v => {
            alert(v)
        }).catch(error => {
            alert(error)
        });
    }).catch(error => {
        alert(error)
    });
}

// function PostTypeImport() {
//     let lst = [
//         ["xxx", 11, "yyyy"],
//         ["mm", 23, "nnn"],
//         ["mm", , "nnn"],
//         ["mm", "", "nnn"],
//     ]
//     const formData = new FormData();
//     formData.append("key", "bb496da1-f87f-456c-b060-5c6b2305d242");
//     formData.append("replace", false.toString());
//     formData.append("data", JSON.stringify(lst));
//     fetch("/api/ai", {
//         method: 'POST',
//         // headers: {
//         //     'Content-Type': 'multipart/form-data',
//         // },
//         body: formData
//     }).then(result => {
//         result.text().then(v => {
//             alert(v)
//         }).catch(error => {
//             alert(error)
//         });
//     }).catch(error => {
//         alert(error)
//     });
// }
</script>