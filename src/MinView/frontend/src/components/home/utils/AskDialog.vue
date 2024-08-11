<template>
    <v-row justify="center">
        <v-dialog v-model="store.state.viShowAsk" activator="parent" persistent width="auto">
            <v-card>
                <!-- <v-card-title class="text-h5">
                    Use Google's location service?
                </v-card-title> -->
                <v-card-text>{{ store.state.viMessage }}</v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="red" variant="tonal" @click="handleSubmit">
                        确定
                    </v-btn>
                    <v-btn color="green" variant="tonal" @click="handleClose">
                        取消
                    </v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </v-row>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { ref } from 'vue'
import { store } from '../../../store/index';
let handleSubmit = () => {
    store.commit("show_ask_dialog", { show: false, msg: "确定" })
}
let handleClose = () => {
    store.commit("show_ask_dialog", { show: false, msg: "" })
}
</script>