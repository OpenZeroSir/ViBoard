<template>
    <div v-ripple="storeComputed.state.viRipple" :style="{
        height: '100%',
        width: '100%',
    }" :nni="props.nni">
        <img :src="props.innerStyle.text" :style="{
            // height: props.outStyle.height + 'px',
            // width: props.outStyle.width + 'px',
            height: '100%',
            width: '100%',
            objectFit: props.innerStyle.vertical ? 'cover' : 'contain',
        }">
    </div>
</template>

<script lang="ts" setup>
import { onMounted, computed, reactive } from 'vue'
import { ViChartStyle, ViInnerStyle, ViOutStyle, store } from '../../../store';
import { getCurrentInstance } from 'vue';
import { storeDisplay } from '../../../store/display';
let props = defineProps<{
    nni: string;
    name: string;
    innerStyle: ViInnerStyle;
    outStyle: ViOutStyle;
    chartStyle: ViChartStyle | undefined
}>()


// const instance = getCurrentInstance();

let storeComputed = computed(() => {
    // let fullpath = instance?.proxy?.$router.currentRoute.value.fullPath
    if (!props.outStyle.displayMode) {
        return store
    } else {
        return storeDisplay
    }
})
</script>