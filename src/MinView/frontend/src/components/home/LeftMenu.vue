<script lang="ts" setup>
import { ref, reactive, computed, ComputedRef, getCurrentInstance, onMounted, onUnmounted } from 'vue'
// import { VVirtualScroll } from 'vuetify/VVirtualScroll'
import * as SvgIconName from '@mdi/js'
import { store } from '../../store'
import '../../assets/images/charts/1E51ACC7-9E58-4704-BF1B-1142A07BAE3E.webp'
import '../../assets/images/charts/9b1a3b1a-78ed-4308-9d7c-aae1e731248d.webp'
import '../../assets/images/charts/9cb49149-c520-4f53-a69e-4423e963f3d2.webp'
import '../../assets/images/charts/586e65fd-aed4-42ab-a797-271f578d7f5e.webp'
import '../../assets/images/charts/f58cb5f2-ede3-4923-8fd2-65e5ef28641d.webp'
import '../../assets/images/charts/1885c347-816b-4d9f-8f23-e80ad3e7c5b2.webp'
import '../../assets/images/charts/db94cd1c-f811-494f-9803-fc0a57d5df0b.webp'
import '../../assets/images/charts/5c1631ea-7831-4788-a132-c866f20c7548.webp'

import '../../assets/images/charts/00de94ff-ecda-4231-acde-f49e8777933d.webp'
import '../../assets/images/charts/ba2fafed-516f-48ce-b5e4-0925196edbf7.webp'
import '../../assets/images/charts/97854320-17ad-4c48-9dee-a560ecad171c.webp'
import '../../assets/images/charts/3bb52213-654f-4bfc-b749-1bb7b698036a.webp'
import '../../assets/images/charts/b8563795-6f32-4f8c-a14a-06b93fa93435.webp'


import '../../assets/images/charts/0c79ae07-91ef-4a17-adaa-6e70c7fbf152.webp'
import '../../assets/images/charts/d6f4f906-9539-4a28-b03e-486de2f0876c.webp'
import '../../assets/images/charts/de5415ca-2c0d-4700-87ce-12ef2955a344.webp'
import '../../assets/images/charts/c9ea4933-c5a0-4432-9156-ca0cb7a8fd61.webp'
import '../../assets/images/charts/2caf407d-c219-4780-a66a-4f53ddd52421.webp'
import '../../assets/images/charts/9d6749b8-6ac3-4f88-8cc8-e78952c67239.webp'
import '../../assets/images/charts/8c798693-672f-4a5a-8cd1-107151ada9ba.webp'
import '../../assets/images/charts/199c00b9-d3b9-4dd4-bfc0-8344089b1268.webp'
import '../../assets/images/charts/1782ebf9-589d-426c-b37f-9c4f5f8e7274.webp'
import '../../assets/images/charts/643a5985-d595-4163-955d-b4cd1b212083.webp'
import '../../assets/images/charts/3a2c1a57-66c8-4ae1-80b2-7866d68e415e.webp'
import '../../assets/images/charts/8a5dcc6b-0aa4-40f6-8323-d6ee3e18c59c.webp'
import '../../assets/images/charts/f0dd12cd-0bcb-47ab-a94e-ade76cbb56af.webp'
import '../../assets/images/charts/48088e85-a412-4140-9bd9-56b2f7f93154.webp'

//空
import '../../assets/images/charts/741fef90-a46b-4cdc-9995-07515af1db31.webp'

//地图
import '../../assets/images/charts/71a7a7a2-99c3-4575-bc03-4c1706a5f39b.webp'
import '../../assets/images/charts/c163399d-06b8-436a-b620-070a87947201.webp'
import '../../assets/images/charts/efcabf75-68a7-4941-8208-c0fd38a8aa64.webp'
import '../../assets/images/charts/68ddd844-d023-405e-94e2-88e55a135eb3.webp'
import '../../assets/images/charts/472a9015-c131-44b5-88da-41dbc172a1c2.webp'
import '../../assets/images/charts/f37687e4-c20e-4bfc-8c5d-5c659154118d.webp'


// import { FontIconName } from '../../utils/data'
import { ChartIcons } from '../../utils/charticon'
import { Components } from '../custom/index'
import { nextTick } from 'process'

const instance = getCurrentInstance()
let left_drawer = ref(true)
let left_rail = ref(false)

let base_comp_names = Object.keys(Components).filter(v => {
    let name = Components[v as keyof typeof Components].typeName
    return (name != "图标") && (name != "组") && (name != "图表")
})

function baseCompDragStart(event: any) {
    let compname = event.target.getAttribute("data-icon")
    // event.dataTransfer.setData("Text", svgIconKey);
    event.dataTransfer?.setData("application/json", JSON.stringify({
        type: 'comp',
        //图标名称
        value: compname,
    }))
}



// let font_wheel_index = ref(0)
// let fontsearchkey = ref("")
// let reafontkeys = computed(
//     () => {
//         return reactive(FontIconName.filter(p => p.indexOf(fontsearchkey.value) > -1).slice(font_wheel_index.value, font_wheel_index.value + 60))
//     }
// )

let svgkeys = Object.keys(SvgIconName)
let svgsearchkey = ref("")
// let reasvgkeys = reactive(svgkeys.filter(p=>p.indexOf(svgsearchkey.value)>-1))
let reasvgkeys: ComputedRef<string[]>
let svg_wheel_index = ref(0)
reasvgkeys = computed(
    () => {
        return reactive(svgkeys.filter(p => p.toLowerCase().indexOf(svgsearchkey.value.toLowerCase()) > -1).slice(svg_wheel_index.value, Math.floor(((svg_wheel_index.value + (getInfoHeight.value / 35) * 6) + 5) / 6) * 6))
    }
)

let chartkey = ref("")
let chartkeys = computed(
    () => {
        return reactive(
            ChartIcons.filter(
                p => p.message.indexOf(chartkey.value) > -1
            )
        )
    }
)

function iconDragStart(event: any) {
    let svgIconKey = event.target.getAttribute("data-icon")
    // event.dataTransfer.setData("Text", svgIconKey);
    event.dataTransfer?.setData("application/json", JSON.stringify({
        type: 'icon',
        //图标名称
        value: svgIconKey,
    }))
}

function chartDragStart(event: any, id: string) {
    // let svgIconKey = event.target.getAttribute("data-icon")
    // event.dataTransfer.setData("Text", "nnichart" + id);
    event.dataTransfer?.setData("application/json", JSON.stringify({
        type: 'chart',
        //图表id
        value: id,
    }))
}

function handleSVGIconWheel(event: any) {
    if (event.deltaY < 0) {
        if (svg_wheel_index.value >= 60) {
            svg_wheel_index.value -= 60
        }
    } else {
        if (svg_wheel_index.value < svgkeys.length - 60) {
            svg_wheel_index.value += 60
        }
    }
}
// function handleFontIconWheel(event: any) {
//     if (event.deltaY < 0) {
//         if (font_wheel_index.value >= 60) {
//             font_wheel_index.value -= 60
//         }
//     } else {
//         if (font_wheel_index.value < FontIconName.length - 60) {
//             font_wheel_index.value += 60
//         }
//     }
// }


let getInfoHeight = ref(768 - 358)

function updateInfoHeight() {
    getInfoHeight.value = document.body.scrollHeight - 358
}

onMounted(() => {
    window.addEventListener("resize", updateInfoHeight)
    updateInfoHeight()
})

onUnmounted(() => {
    window.removeEventListener("resize", updateInfoHeight)
})

function charticondbclick(msg: string | undefined) {
    let message = ""
    if (msg) {
        message = msg
    }
    store.commit("set_viShowChartHelp",{
        show:true,
        msg:message
    })
}

</script>

<template>
    <v-navigation-drawer color="light-blue-darken-4" tile flat v-model="left_drawer" :rail="left_rail" permanent
        @click="left_rail = false">
        <v-list-item density="compact" nav>
            <v-list-item-title>
                <span style="font-size: large;margin-left: 15px;font-weight: bold;color: #ff9900;">ViBoard v1.0.0</span>
            </v-list-item-title>
            <template v-if="!left_rail" v-slot:append>
                <v-btn variant="text" icon @click.stop="left_rail = !left_rail">
                    <span className="material-icons">arrow_back_ios_new</span>
                    <v-tooltip activator="parent" location="left">隐藏</v-tooltip>
                </v-btn>
            </template>
            <template v-else v-slot:prepend>
                <v-btn variant="text" icon @click.stop="left_rail = !left_rail">
                    <span className="material-icons">arrow_forward_ios</span>
                    <v-tooltip activator="parent" location="left">显示</v-tooltip>
                </v-btn>
            </template>

        </v-list-item>

        <v-expansion-panels variant="accordion" tile flat v-if="left_rail == false">
            <v-expansion-panel style="background-color:transparent;color: white;" tile flat>
                <v-expansion-panel-title tile flat collapse-icon="md:keyboard_arrow_up"
                    expand-icon="md:keyboard_arrow_down">
                    <span style="font-size: small;">基础组件</span>
                </v-expansion-panel-title>
                <v-expansion-panel-text tile flat>
                    <!-- <v-card  tile flat v-if="left_rail == false" style="background-color:transparent;color: white;" >
                        <v-chip class="ma-2" color="cyan" label draggable v-for="(title, i) in base_comp_names"
                            :key="i">{{ Components[title as keyof typeof Components].friendlyName }}</v-chip>
                    </v-card> -->
                    <v-chip-group active-class="primary--text" column>
                        <v-chip v-bind:data-icon="title" class="ma-2" color="cyan" @dragstart="baseCompDragStart($event)"
                            label draggable v-for="(title, i) in base_comp_names" :key="i">{{
                                Components[title as keyof typeof Components].typeName }}</v-chip>
                    </v-chip-group>
                </v-expansion-panel-text>

            </v-expansion-panel>
            <v-expansion-panel style="background-color:transparent;color: white;" tile>
                <v-expansion-panel-title collapse-icon="md:keyboard_arrow_up" expand-icon="md:keyboard_arrow_down">
                    <span style="font-size: small;">图表组件</span>
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                    <v-text-field v-model="chartkey" variant="underlined" density="compact" hint="通过关键字搜索图表">
                        <template v-slot:append-inner>
                            <v-icon>
                                {{ SvgIconName["mdiMagnify" as keyof typeof SvgIconName] }}
                            </v-icon>
                        </template>
                    </v-text-field>
                    <!-- <v-virtual-scroll :height="350" :items="chartkeys">
                        <template v-slot:default="{ item }">
                            <v-list-item>
                                <v-img draggable :src="`${item.id}.webp`" />
                                </v-img>
                                <v-tooltip activator="parent" location="right"> {{ item.message }}</v-tooltip>
                                <template v-slot:prepend>
                                    <div style="margin-right: 4px;">
                                        <v-badge dot :color="item.vip ? 'red' : 'green'"></v-badge>
                                 </template>
                            </v-list-item>
                        </template>
                    </v-virtual-scroll> -->
                    <div :style="{
                        height: getInfoHeight + 'px', overflow: 'auto !important',

                    }">
                        <div v-for="item in chartkeys">
                            <div>
                                <img v-bind:data-icon="item.id" draggable :src="`webp/${item.id}.webp`"
                                    @dblclick="charticondbclick(item.help)" @dragstart="chartDragStart($event, item.id)"
                                    style="width: 100%;" />
                                <v-tooltip activator="parent" location="left">{{ item.message }}</v-tooltip>
                            </div>
                        </div>
                    </div>
                </v-expansion-panel-text>
            </v-expansion-panel>
            <!-- <v-expansion-panel style="background-color:transparent;color: white;" tile>
                <v-expansion-panel-title collapse-icon="md:keyboard_arrow_up" expand-icon="md:keyboard_arrow_down">
                    <span style="font-size: small;">字体图标</span>
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                    <v-text-field v-model="fontsearchkey" variant="underlined" density="compact" hint="通过关键字搜索字体图标">
                        <template v-slot:append-inner>
                            <v-icon>
                                {{ SvgIconName["mdiMagnify" as keyof typeof SvgIconName] }}
                            </v-icon>
                        </template>
                    </v-text-field>
                    <v-chip-group @wheel="handleFontIconWheel($event)" active-class="primary--text" column>
                        <v-chip v-bind:data-icon="`md:${tag}`" draggable label link ripple size="64px"
                            v-for="tag in reafontkeys" :key="tag" @dragstart="iconDragStart($event)">
                            <v-icon size="26px">
                                md:{{ tag }}
                            </v-icon>
                            <v-tooltip activator="parent" location="right"> {{ tag }}</v-tooltip>
                        </v-chip>
                    </v-chip-group>

                </v-expansion-panel-text>
            </v-expansion-panel> -->
            <v-expansion-panel style="background-color:transparent;color: white;" tile>
                <v-expansion-panel-title collapse-icon="md:keyboard_arrow_up" expand-icon="md:keyboard_arrow_down">
                    <span style="font-size: small;">矢量图标</span>
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                    <v-text-field v-model="svgsearchkey" single-line variant="underlined" density="compact"
                        hint="通过关键字搜索矢量图标" hide-details>
                        <template v-slot:append-inner>
                            <v-icon>
                                {{ SvgIconName["mdiMagnify" as keyof typeof SvgIconName] }}
                            </v-icon>
                        </template>
                    </v-text-field>

                    <v-chip-group @wheel="handleSVGIconWheel($event)" active-class="primary--text" column>
                        <v-chip v-bind:data-icon="tag" draggable label link ripple size="64px" v-for="tag in reasvgkeys"
                            :key="tag" @dragstart="iconDragStart($event)">
                            <v-icon size="26px">
                                {{ SvgIconName[tag as keyof typeof SvgIconName] }}
                            </v-icon>
                            <v-tooltip activator="parent" location="right"> {{ tag }}</v-tooltip>
                        </v-chip>
                    </v-chip-group>

                    <!-- <v-infinite-scroll height="400">
                        <template v-for="(item, index) in SvgIconName" :key="item">
                                Item #{{ item }}
                        </template>
                    </v-infinite-scroll> -->

                </v-expansion-panel-text>
            </v-expansion-panel>

        </v-expansion-panels>
    </v-navigation-drawer>
</template>

<style lang="scss">
.v-expansion-panel__shadow {
    box-shadow: none !important;
}
</style>
