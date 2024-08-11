<template>
    <div class="text-center">
        <v-dialog v-model="store.state.viShow" width="auto">
            <v-card>
                <v-card-text>
                    {{ store.state.viMessage }}
                </v-card-text>
                <!-- <v-card-actions>
                    <v-btn color="primary" block @click="handleClose">关闭</v-btn>
                </v-card-actions> -->
            </v-card>

        </v-dialog>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { ref } from 'vue'
import { store } from '../../../store/index';
let handleClose=()=>{
    store.commit("show_dialog",{show:false,msg:""})
}
</script>