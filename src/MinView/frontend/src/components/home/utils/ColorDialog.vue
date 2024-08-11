<template>
    <v-dialog  width="auto">
        <v-card>
            <v-card-text>
                <v-color-picker></v-color-picker>
            </v-card-text>
            <v-card-actions>
                <v-btn color="primary" block @click="$emit('update:visible', false)">Close Dialog</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script lang="ts" setup>
import { WritableComputedRef } from '@vue/runtime-core';
import { ref, computed,toRefs } from 'vue'

// let props = defineProps<{
//     show: boolean
// }>()

</script>