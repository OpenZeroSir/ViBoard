<template>
    <v-dialog v-model="store.state.viShowUpdate" width="500px">
        <v-card>
            <!-- <div style="background-color: #01579B;padding:10px">
                <SmallInfo />
            </div> -->
            <v-card-text>
                <div>
                    <div>
                        <v-row>
                            <v-col>
                                <span style="font-weight: bold;">当前版本:</span>
                            </v-col>
                            <v-col>
                                v1.0.0
                            </v-col>
                            <v-col>
                                <!-- <v-btn icon @click="update">
                                    <v-icon style="color: red;">{{ mdiUpdate }}</v-icon>
                                </v-btn> -->

                            </v-col>
                        </v-row>
                    </div>
                    <div style="margin-top: 20px;">
                        <v-row v-for="(key, index) in keys">
                            <v-col>
                                <span style="font-weight: bold;">{{ key }}</span>
                            </v-col>
                            <v-col>
                                {{ values[index] }}
                            </v-col>
                        </v-row>
                    </div>
                    <div style="text-align: center;margin-top: 20px;width: 100%;">
                        <v-btn @click="update">
                            检查更新
                        </v-btn>
                    </div>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import SmallInfo from './SmallInfo.vue';
import { store } from '../../../store/index';
import { CheckUpdate } from '../../../../wailsjs/go/main/App';
import { onMounted, ref } from 'vue';
import { mdiUpdate } from '@mdi/js';

let values = ref([] as string[])
let keys = ref([] as string[])

function update() {
    values.value.length = 0
    keys.value.length = 0
    CheckUpdate().then(v => {
        if (v.indexOf("Error:") > -1) {
            values.value.push(v)
            keys.value.push("错误提示")
        } else {
            let json = JSON.parse(v)
            let tags = json["tags"] as string
            let tagarr = tags.split(" ")
            tagarr.forEach(k => {
                keys.value.push(k)
                values.value.push(json[k])
            })
        }
    }).catch(reason => {
        values.value.push(reason)
        keys.value.push("错误提示")
    })
}

onMounted(() => {

})
</script>