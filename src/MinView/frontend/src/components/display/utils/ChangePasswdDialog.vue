<template>
    <div>
        <v-dialog v-model="storeDisplay.state.viShowChangePasswd" activator="parent" persistent>
            <v-card height="100%">
                <v-card-title style="text-align: center;">
                    重置密码
                </v-card-title>
                
                <v-card-text>
                    <v-row>
                        <v-col>
                            <v-text-field label="原密码" v-model="old_passwd" variant="underlined" clearable></v-text-field>
                        </v-col>
                        <v-col>
                            <v-text-field label="新密码" v-model="new_passwd" variant="underlined" clearable></v-text-field>
                        </v-col>
                        <v-col>
                            <v-text-field label="确认新密码" v-model="new_passwd_comf" variant="underlined"
                                clearable></v-text-field>
                        </v-col>
                    </v-row>
                    <v-row>
                        <div style="color: red;font-size:15px; padding: 10px;">本功能用于修改视图编码的密码，确认修改密码后，密码就已经修改成功，
                        但需要保存视图或另存视图到电脑上，下次加载该视图编码文件采用新密码，建议采用另存视图。</div>
                        <v-checkbox v-model="checked" color="red" :label="info_label"></v-checkbox>
                    </v-row>
                </v-card-text>
                
                <v-card-actions>
                    <v-row>
                        <v-col style="text-align: center;">
                            <v-btn variant="tonal" color="red" @click="changePasswdFunction">
                                确定修改密码
                            </v-btn>
                        </v-col>
                        <v-col style="text-align: center;">
                            <v-btn variant="tonal" color="green" @click="cancelChange">
                                放弃修改密码
                            </v-btn>
                        </v-col>
                    </v-row>
                </v-card-actions>
                <br/>
            </v-card>
        </v-dialog>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { storeDisplay } from '../../../store/display';
let checked = ref(false)
let old_passwd = ref("")
let new_passwd = ref("")
let new_passwd_comf = ref("")

let info_label = computed(()=>{
    return "我确认修改密码，忘记新密码再也没有办法加载该视图编码文件："+storeDisplay.state.viFilePath
})


function changePasswdFunction() {
    // if(old_passwd.value==""){
    //     storeDisplay.commit("show_dialog", { show: true, msg: "原密码不能为空" })
    //     return
    // }
    // if(new_passwd.value==""){
    //     storeDisplay.commit("show_dialog", { show: true, msg: "新密码不能为空" })
    //     return
    // }
    // if(new_passwd_comf.value ==""){
    //     storeDisplay.commit("show_dialog", { show: true, msg: "确认新密码不能为空" })
    //     return
    // }
    if (new_passwd_comf.value != new_passwd.value) {
        storeDisplay.commit("show_dialog", { show: true, msg: "新密码和确认新密码不一致" })
        return
    }
    if (old_passwd.value == "" && new_passwd.value == "") {
        storeDisplay.commit("show_dialog", { show: true, msg: "密码都为空" })
        return
    }
    if (old_passwd.value != storeDisplay.state.preInfo.passwd) {
        storeDisplay.commit("show_dialog", { show: true, msg: "旧密码无效" })
        return
    }
    if (!checked.value) {
        storeDisplay.commit("show_dialog", { show: true, msg: "请选择确认修改密码才能修改！" })
        return
    }
    storeDisplay.commit("set_pre_passwd", new_passwd.value)
    storeDisplay.commit("show_change_passwd_dialog", false)
}

function cancelChange() {
    storeDisplay.commit("show_change_passwd_dialog", false)
}

let show_dialog = computed({
    get: () => storeDisplay.state.viShowChangePasswd,
    set: (val) => {
        storeDisplay.commit("show_change_passwd_dialog", val)
    }
});

</script>