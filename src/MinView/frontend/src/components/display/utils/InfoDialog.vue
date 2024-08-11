<template>
    <div class="text-center">
        <v-dialog v-model="storeDisplay.state.viShow" width="auto">
            <v-card>
                <v-card-text>
                    {{ storeDisplay.state.viMessage }}
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
import { storeDisplay } from '../../../store/display';
let handleClose=()=>{
    storeDisplay.commit("show_dialog",{show:false,msg:""})
}
</script>