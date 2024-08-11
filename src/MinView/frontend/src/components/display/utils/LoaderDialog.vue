<template>
    <div class="text-center">
        <!-- <v-btn :disabled="dialog" :loading="dialog" color="purple-darken-2" @click="dialog = true">
            Start loading
        </v-btn> -->
        <v-dialog v-model="enable" activator="parent" persistent width="auto">
            <v-card >
                <v-card-text>
                    正在处理，请稍后...
                    <v-progress-linear indeterminate color="red" class="mb-0"></v-progress-linear>
                </v-card-text>
            </v-card>
        </v-dialog>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { ref } from 'vue'
import { storeDisplay } from '../../../store/display';
let enable = computed({
    get: () => storeDisplay.state.viShowLoader,
    set: (val) => {
        storeDisplay.commit("show_loader", val)
    }
});
</script>