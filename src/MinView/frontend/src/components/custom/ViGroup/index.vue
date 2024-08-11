<template>
    <div v-ripple="storeComputed.state.viRipple" :style="{
        height: '100%',
        width: '100%',

    }">

        <component v-for="item in style.components" :style="{
            // width: item.style.width + 'px',
            // height: item.style.height + 'px',
            width: item.groupStyle.width + '%',
            height: item.groupStyle.height + '%',
            top: item.groupStyle.top + '%',
            left: item.groupStyle.left + '%',
            // rotate: item.style.rotate + 'deg',
            transform: `rotate(${item.style.rotate}deg)`,
            position: 'absolute',
            backgroundColor: item.style.backgroundColor,
        }" :nni="`nni${item.id}`" :is="item.component" :name="item.name" :key="item.id"
            :innerStyle="item.innerStyle" :components="item.components" :outStyle="item.style" :groupStyle="item.groupStyle"
            :chartStyle="item.chartStyle" />
    </div>
</template>
<script lang="ts" setup>
import { computed, ref, onMounted } from 'vue'
import { ViChartStyle, ViComponent, ViGroupStyle, ViInnerStyle, ViOutStyle, store } from '../../../store';
import { getCurrentInstance } from 'vue';
import { storeDisplay } from '../../../store/display';
let props = defineProps<{
    name: string;
    nni: string
    innerStyle: ViInnerStyle;
    outStyle: ViOutStyle;
    groupStyle: ViGroupStyle;
    components: ViComponent[];
    chartStyle: ViChartStyle | undefined,
}>()

let style = computed(() => {
    return props
})
onMounted(() => {

})

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

<style lang="scss" scoped>
.group {
    // & > div {
    //     position: relative;
    //     width: 100%;
    //     height: 100%;

    //     .component {
    //         position: absolute;
    //     }
    // }
    position: relative;
    width: 100%;
    height: 100%;
}
</style>