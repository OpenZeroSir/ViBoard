<template>
    <div v-ripple="storeComputed.state.viRipple">
        <div :nni="props.nni" :style="{
            height: '100%',
            width:'100%',
        }">
        </div>
    </div>
</template>
<script lang="ts" setup>
import { computed } from 'vue';
import { store, ViChartStyle, ViInnerStyle, ViOutStyle } from '../../../store';
import { getCurrentInstance } from "vue";
import { storeDisplay } from '../../../store/display';

let props = defineProps<{
    nni: string;
    name: string;
    innerStyle: ViInnerStyle;
    outStyle: ViOutStyle;
    chartStyle: ViChartStyle | undefined
}>()

let storeComputed = computed(() => {
    // let fullpath = instance?.proxy?.$router.currentRoute.value.fullPath
    if (!props.outStyle.displayMode) {
        return store
    } else {
        return storeDisplay
    }
})

// const instance = getCurrentInstance();

</script>