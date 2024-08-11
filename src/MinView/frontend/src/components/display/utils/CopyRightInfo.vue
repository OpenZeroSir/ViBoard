<template>
    <v-dialog width="100%"  v-model="storeDisplay.state.viShowCopyRight">
        <v-card v-ripple style="height: 150px;" class="bg-light-blue-darken-4">
            <v-card-title>
                <div style="text-align: center;">
                    视图编码文件版权
                </div>
            </v-card-title>
            <v-card-text style="padding-top: 50px;padding-left: 40px;padding-right: 10px;">
                <v-row >
                    <!-- <v-col> -->
                    <v-row>
                        <v-col>
                            <v-row style="color: white;">
                                所属人员
                            </v-row>
                            <v-row style="color: #ff9900;">
                                {{ author }}
                            </v-row>
                        </v-col>
                        <v-col>
                            <v-row  style="color: white;">
                                联系电话
                            </v-row>
                            <v-row style="color: #ff9900;">
                                {{ phone }}
                            </v-row>
                        </v-col>
                        <v-col>
                            <v-row  style="color: white;">
                                联系邮箱
                            </v-row>
                            <v-row style="color: #ff9900;">
                                {{ email }}
                            </v-row>
                        </v-col>
                        <v-col>
                            <v-row  style="color: white;">
                                所属公司
                            </v-row>
                            <v-row style="color: #ff9900;">
                                {{ company }}
                            </v-row>
                        </v-col>
                        <v-col>
                            <v-row  style="color: white;">
                                服务网址
                            </v-row>
                            <v-row style="color: #ff9900;">
                                {{ website }}
                            </v-row>
                        </v-col>
                    </v-row>
                    <!-- </v-col> -->
                </v-row>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>
<script lang="ts" setup>
import { computed } from 'vue';
import { storeDisplay } from '../../../store/display';
let author = computed(() => {
    if (storeDisplay.state.canvas.author == "") {
        return "-"
    } else {
        return storeDisplay.state.canvas.author
    }
})
let phone = computed(() => {
    if (storeDisplay.state.canvas.phone == "") {
        return "-"
    } else {
        return storeDisplay.state.canvas.phone
    }
})
let email = computed(() => {
    if (storeDisplay.state.canvas.email == "") {
        return "-"
    } else {
        return storeDisplay.state.canvas.email
    }
})
let company = computed(() => {
    if (storeDisplay.state.canvas.company == "") {
        return "-"
    } else {
        return storeDisplay.state.canvas.company
    }
})
let website = computed(() => {
    if (storeDisplay.state.canvas.website == "") {
        return "-"
    } else {
        return storeDisplay.state.canvas.website
    }
})
</script>