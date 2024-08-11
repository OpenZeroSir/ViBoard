<template>
    <div ref="maindiv" class="shape" :class="{ active }" :style="{
        width: props.outStyle.width + 'px',
        height: props.outStyle.height + 'px',
        top: props.outStyle.top + 'px',
        left: props.outStyle.left + 'px',
        // rotate: props.outStyle.rotate + 'deg',
        transform: `rotate(${props.outStyle.rotate}deg)`,
        position: 'absolute',
        backgroundColor: props.outStyle.backgroundColor,
        // display:'flex',
        zIndex: props.outStyle.zIndex,

        borderTopColor: props.outStyle.borderTopColor,
        borderBottomColor: props.outStyle.borderBottomColor,
        borderLeftColor: props.outStyle.borderLeftColor,
        borderRightColor:  props.outStyle.borderRightColor,

        borderTopStyle: props.outStyle.borderTopStyle  as any,
        borderBottomStyle:props.outStyle.borderBottomStyle  as any,
        borderLeftStyle: props.outStyle.borderLeftStyle  as any,
        borderRightStyle:props.outStyle.borderRightStyle  as any,

        borderTopWidth: props.outStyle.borderTopWidth,
        borderBottomWidth: props.outStyle.borderBottomWidth,
        borderLeftWidth: props.outStyle.borderLeftWidth,
        borderRightWidth: props.outStyle.borderRightWidth,

        borderTopLeftRadius: props.outStyle.borderTopLeftRadius,
        borderTopRightRadius: props.outStyle.borderTopRightRadius,
        borderBottomLeftRadius: props.outStyle.borderBottomLeftRadius,
        borderBottomRightRadius: props.outStyle.borderBottomRightRadius,

        opacity:props.outStyle.opacity,

    }" @mousedown="handleMouseDownOnShape($event)">
        <span v-show="props.active" class="material-icons icon-rotation"
            @mousedown="handleRotate($event)">rotate_right</span>
        <!-- <span v-show="element.isLock" class="iconfont icon-suo"></span> -->
        <div v-for="item in (active ? getPointList() : [])" :key="item" class="shape-point" :style="getPointStyle(item)"
            @mousedown="handleMouseDownOnPoint(item, $event)">
        </div>
        <slot></slot>
    </div>
</template>
<script setup lang="ts">
import { string } from "mathjs";
import { getComponentRotatedStyle, mod360 } from "../../../utils/func"
import { ref, reactive, computed, nextTick, ComputedRef } from "vue"
import { store, ViChartStyle, ViComponent, ViGroupStyle, ViInnerStyle, ViOutStyle } from "../../../store";
import positonAndSize from "../../../utils/PositonAndSize"
import { cloneDeep } from "lodash";

let props = defineProps<{
    active: boolean;
    outStyle: ViOutStyle;
    innerStyle: ViInnerStyle,
    components: ViComponent[];
    component_id: string;
    chartStyle: ViChartStyle | undefined,
}>()


let maindiv = ref<HTMLElement>()
let cursors: {
    [key: string]: any;
} = computed(() => getCursor())

const pointList = ['lt', 't', 'rt', 'r', 'rb', 'b', 'lb', 'l']
const initialAngle = { // 每个点对应的初始角度
    lt: 0,
    t: 45,
    rt: 90,
    r: 135,
    rb: 180,
    b: 225,
    lb: 270,
    l: 315,
}
const angleToCursor = [ // 每个范围的角度对应的光标
    { start: 338, end: 23, cursor: 'nw' },
    { start: 23, end: 68, cursor: 'n' },
    { start: 68, end: 113, cursor: 'ne' },
    { start: 113, end: 158, cursor: 'e' },
    { start: 158, end: 203, cursor: 'se' },
    { start: 203, end: 248, cursor: 's' },
    { start: 248, end: 293, cursor: 'sw' },
    { start: 293, end: 338, cursor: 'w' },
]

function getPointList() {
    return pointList
}
// let styleReact = computed(() => {
//     return props.outStyle
// })
function getPointStyle(point: string) {
    // const { width, height } = this.defaultStyle
    const width = Number(props.outStyle.width)
    const height = Number(props.outStyle.height)
    const hasT = /t/.test(point)
    const hasB = /b/.test(point)
    const hasL = /l/.test(point)
    const hasR = /r/.test(point)
    let newLeft = 0
    let newTop = 0

    // 四个角的点
    if (point.length === 2) {
        newLeft = hasL ? 0 : width
        newTop = hasT ? 0 : height
    } else {
        // 上下两点的点，宽度居中
        if (hasT || hasB) {
            newLeft = width / 2
            newTop = hasT ? 0 : height
        }

        // 左右两边的点，高度居中
        if (hasL || hasR) {
            newLeft = hasL ? 0 : width
            newTop = Math.floor(height / 2)
        }
    }
    cursors = getCursor()
    let style = {
        marginLeft: '-4px',
        marginTop: '-4px',
        left: `${newLeft}px`,
        top: `${newTop}px`,
        // cursor: cursors[point as keyof typeof string],
        cursor: cursors[point as keyof typeof string]
    }
    // console.log(point, cursors[point as keyof typeof string], cursors)
    // console.log("style", style)
    return style
}

function getCursor() {
    const rotate = mod360(Number(props.outStyle.rotate)) // 取余 360
    const result: { [key: string]: any } = {};
    let lastMatchIndex = -1 // 从上一个命中的角度的索引开始匹配下一个，降低时间复杂度

    pointList.forEach(point => {
        const angle = mod360(initialAngle[point as keyof typeof initialAngle] + rotate)
        const len = angleToCursor.length
        // eslint-disable-next-line no-constant-condition
        while (true) {
            lastMatchIndex = (lastMatchIndex + 1) % len
            const angleLimit = angleToCursor[lastMatchIndex]
            if (angle < 23 || angle >= 338) {
                result[point] = 'nw-resize'

                return
            }

            if (angleLimit.start <= angle && angle < angleLimit.end) {
                result[point] = angleLimit.cursor + '-resize'

                return
            }
        }
    })

    return result
}

function handleMouseDownOnShape(event: MouseEvent) {
    event.preventDefault()
    event.stopPropagation()

    if (props.chartStyle == undefined) {
        store.commit("set_boarding", { value: 1 })
    }
    store.commit("set_bread_index", { value: 0 })
    store.commit("set_component_id", props.component_id)
    cursors = getCursor()
    let top_left_objects: { id: string, top: number, left: number }[] = []
    let pos_list: ViComponent[] = []
    store.state.viViews[store.state.viCtrlIndex].componet_ids.forEach(v => {
        store.state.viViews[store.state.viCtrlIndex].components.forEach(c => {
            if (v == c.id) {
                pos_list.push(c)
                top_left_objects.push({
                    id: c.id,
                    top: c.style.top,
                    left: c.style.left
                })
                return
            }
        })

    })
    const startY = event.clientY
    const startX = event.clientX
    function getTopLeft(id: string) {
        let top = -1
        let left = -1
        top_left_objects.forEach(v => {
            if (v.id == id) {
                top = v.top
                left = v.left
                return
            }
        })
        return { top: top, left: left }
    }

    let hasMove = false
    const move = (moveEvent: MouseEvent) => {
        hasMove = true
        if (moveEvent.buttons != 1) {
            return
        }
        const curX = moveEvent.clientX
        const curY = moveEvent.clientY
        pos_list.forEach(c => {
            let pos = cloneDeep(c.style)
            let top_left = getTopLeft(c.id)
            if (top_left.top == -1 || top_left.left == -1) {
                return
            }
            let offset_y = curY - startY
            let offset_x = curX - startX
            if (offset_y == 0 && offset_x == 0) {
                return
            }
            pos.top = offset_y + top_left.top
            pos.left = offset_x + top_left.left
            store.commit("set_out_style", {
                style: reactive(pos),
                addr: {
                    viewid: store.state.viCtrlIndex,
                    comid: c.id
                }
            })
        })

    }

    const up = () => {
        document.removeEventListener('mousemove', move)
        document.removeEventListener('mouseup', up)
        if (hasMove) {
            nextTick(() => {
                store.commit("snapshot")
            })
        }
    }

    document.addEventListener('mousemove', move)
    document.addEventListener('mouseup', up)
}

function handleMouseDownOnPoint(item: any, event: any) {
    event.preventDefault()
    event.stopPropagation()

    const style = { ...props.outStyle }
    style.top = style.top
    style.left = style.left
    cursors = getCursor()

    // 组件宽高比

    const proportion = style.width / style.height

    // 组件中心点
    const center = {
        x: style.left + style.width / 2,
        y: style.top + style.height / 2,
    }

    // 获取画布位移信息
    const editorRectInfo = store.state.viEditor!.getBoundingClientRect()

    // 获取 point 与实际拖动基准点的差值 @justJokee
    const pointRect = event.target.getBoundingClientRect()
    // 当前点击圆点相对于画布的中心坐标
    const curPoint = {
        x: Math.round(pointRect.left - editorRectInfo.left + event.target.offsetWidth / 2),
        y: Math.round(pointRect.top - editorRectInfo.top + event.target.offsetHeight / 2),
    }

    // 获取对称点的坐标
    const symmetricPoint = {
        x: center.x - (curPoint.x - center.x),
        y: center.y - (curPoint.y - center.y),
    }

    // 是否需要保存快照
    let needSave = false
    let isFirst = true

    // const needLockProportion = this.isNeedLockProportion()
    const move = (moveEvent: MouseEvent) => {
        // 第一次点击时也会触发 move，所以会有“刚点击组件但未移动，组件的大小却改变了”的情况发生
        // 因此第一次点击时不触发 move 事件
        if (isFirst) {
            isFirst = false
            return
        }

        needSave = true
        const curPositon = {
            x: moveEvent.clientX - Math.round(editorRectInfo.left),
            y: moveEvent.clientY - Math.round(editorRectInfo.top),
        }
        const begin_width = style.width
        const begin_height = style.height
        positonAndSize(item, style, curPositon, proportion, false, {
            center,
            curPoint,
            symmetricPoint,
        })
        const end_width = style.width
        const end_height = style.height
        store.commit("set_out_style", {
            style: reactive(style),
            addr: {
                viewid: store.state.viCtrlIndex,
                comid: props.component_id
            }
        })
        let newStyle = cloneDeep(props.innerStyle)
        store.commit("set_inner_style", {
            style: reactive(//{
                //     width: style.width,
                //     height: style.height,
                //     fontfamily:props.innerStyle.fontfamily,
                //     fontSize:props.innerStyle.fontSize,
                //     text:props.innerStyle.text,
                //     list:props.innerStyle.textList,
                // }
                newStyle
            ),
            addr: {
                viewid: store.state.viCtrlIndex,
                comid: props.component_id
            }
        })
    }
    const up = () => {
        document.removeEventListener('mousemove', move)
        document.removeEventListener('mouseup', up)
        if (needSave) {
            nextTick(() => {
                store.commit("snapshot")
            })
        }
    }

    document.addEventListener('mousemove', move)
    document.addEventListener('mouseup', up)
}

function handleRotate(e: MouseEvent) {
    e.preventDefault()
    e.stopPropagation()
    // 初始坐标和初始角度
    const pos = { ...props.outStyle }
    const startY = e.clientY
    const startX = e.clientX
    const startRotate = Number(pos.rotate)

    // 获取元素中心点位置
    const rect = store.state.viEditor!.getBoundingClientRect()
    const centerX = rect.left + rect.width / 2
    const centerY = rect.top + rect.height / 2

    // 旋转前的角度
    const rotateDegreeBefore = Math.atan2(startY - centerY, startX - centerX) / (Math.PI / 180)

    // 如果元素没有移动，则不保存快照
    let hasMove = false
    const move = (moveEvent: MouseEvent) => {
        hasMove = true
        const curX = moveEvent.clientX
        const curY = moveEvent.clientY
        // 旋转后的角度
        const rotateDegreeAfter = Math.atan2(curY - centerY, curX - centerX) / (Math.PI / 180)
        // 获取旋转的角度值
        pos.rotate = (startRotate + rotateDegreeAfter - rotateDegreeBefore)
        // 修改当前组件样式

        store.commit('set_out_style', {
            style: reactive(pos), addr: {
                viewid: store.state.viCtrlIndex,
                comid: props.component_id
            }
        })
        let newStyle = cloneDeep(props.innerStyle)
        store.commit("set_inner_style", {
            style: reactive(//{
                // width: pos.width,
                // height: pos.height,
                // fontfamily: props.innerStyle.fontfamily,
                // fontSize: props.innerStyle.fontSize,
                // text: props.innerStyle.text,
                // list: props.innerStyle.textList,
                //}
                newStyle
            ), addr: {
                viewid: store.state.viCtrlIndex,
                comid: props.component_id
            }
        })
    }

    const up = () => {
        if (hasMove) {
            nextTick(() => {
                store.commit("snapshot")
            })
        }
        document.removeEventListener('mousemove', move)
        document.removeEventListener('mouseup', up)
        cursors = getCursor() // 根据旋转角度获取光标位置
    }

    document.addEventListener('mousemove', move)
    document.addEventListener('mouseup', up)
}
</script>

<style lang="scss" scoped>
.shape {
    position: absolute;

    &:hover {
        cursor: move;
    }
}

.active {
    outline: 1px solid #ff9900;
    // user-select: none;
}

.shape-point {
    position: absolute;
    background: #fff;
    border: 1px solid #ff9900;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    z-index: 1;
}

.icon-rotation {
    position: absolute;
    top: -25px;
    left: 50%;
    transform: translateX(-50%);
    cursor: grab;
    color: #ff9900;
    font-size: 18px;
    font-weight: 100;

    &:active {
        cursor: grabbing;
    }
}

.icon-suo {
    position: absolute;
    top: 0;
    right: 0;
}
</style>