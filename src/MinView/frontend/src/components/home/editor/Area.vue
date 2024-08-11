<template>
    <div :style="{
        left: start.x + 'px',
        top: start.y + 'px',
        width: width + 'px',
        height: height + 'px',
    }" class="area"></div>
</template>

<script lang="ts" setup>

defineProps({
    start: {
        type: Object,
        default: () => { },
    },
    width: {
        type: Number,
        default: 0,
    },
    height: {
        type: Number,
        default: 0,
    },
})

</script>

<style lang="scss" scoped>
.area {
    border: 1px dotted #ff9900;
    position: absolute;
}
</style>