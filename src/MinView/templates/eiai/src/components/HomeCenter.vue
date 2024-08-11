<template >
    <v-layout>
        <v-navigation-drawer style="color:lightgray;background-color: #003355; text-align: center;">
            <InfoCenter />
            <v-list elevation="0">
                <v-list-item>
                    <v-btn variant="text" width="100%" flat nav @click.stop="handleToHome">
                        更新数据表格
                    </v-btn>
                </v-list-item>
                <v-list-item>
                    <v-btn variant="text" width="100%" flat nav @click.stop="handleToCompo">
                        更新组件数据
                    </v-btn>
                </v-list-item>
            </v-list>
            <div
                style="width: 100%;background-color: #003344; color: lightgray;font-size:medium;position: absolute;bottom: 5px;text-align: center;">
                <span>Feedback:open@zerosir.com</span>
            </div>
        </v-navigation-drawer>
        <v-main>
            <router-view></router-view>
        </v-main>
    </v-layout>
</template>

<script setup lang="ts">
import InfoCenter from './InfoCenter.vue';
import { getCurrentInstance } from 'vue';

const instance = getCurrentInstance();

function handleToCompo() {
    if (instance != null) {
        const proxy = instance.proxy
        if (proxy != null) {
            const router = proxy.$router
            router.push("/comp")
        }
    }
}
function handleToHome() {
    if (instance != null) {
        const proxy = instance.proxy
        if (proxy != null) {
            const router = proxy.$router
            router.push("/")
        }
    }
}

</script>