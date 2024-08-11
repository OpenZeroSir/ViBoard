<template>
    <v-dialog v-model="store.state.viShowLicense" width="auto">
        <v-card>
            <v-card-text>
                <pre>
                    {{ info }}
                </pre>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { ReadLicense } from '../../../../wailsjs/go/main/App';
import { store } from '../../../store/index';
import { ref } from 'vue';

let info = ref('')

onMounted(()=>{
    ReadLicense().then(v=>{
        info.value = v
    })
})

</script>