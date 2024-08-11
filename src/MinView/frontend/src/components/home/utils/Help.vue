<template>
    <v-dialog v-model="store.state.viShowHelp" width="auto">
        <v-card>
            <div style="background-color: white;padding:10px;width: 100%;">
                <video src="../../../assets/videos/readme.webm" autoplay controls style="width: 100%;"></video>
            </div>
            <v-card-text>
                <!-- <div>
                    <div>
                        
                    </div>
                    <div>
                        <v-expansion-panels>
                            <v-expansion-panel v-for="(item, index) in lst" :key="index" :title="item.Q"
                                :text="item.A"></v-expansion-panel>
                        </v-expansion-panels>
                    </div>
                </div> -->
                <div id="helpdialog">
                    <div v-for="(item, index) in lst">
                        <div v-html="marked.parse(item.Q)"></div>
                        <pre>
                        <code v-html="marked.parse(item.A)"></code>
                    </pre>
                        <!-- <pre>
                        <code v-html=" marked.parse(item.A)"></code>
                    </pre> -->
                    </div>
                </div>
            </v-card-text>
            <div style="background-color: #01579B;padding:10px">
                <SmallInfo />
            </div>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import SmallInfo from './SmallInfo.vue';
import { store } from '../../../store/index';
import { nextTick, watch } from 'vue';
import * as marked from 'marked';
import hljs from 'highlight.js';
import 'highlight.js/styles/default.css';


watch(() => store.state.viShowHelp, async (newValue) => {
    if (newValue) {
        await nextTick(); // Wait for the next DOM update
        const container = document.getElementById("helpdialog");
        if (container) {
            // const codeBlocks = container.querySelectorAll('code');
            // codeBlocks.forEach((codeBlock) => {
            //     hljs.highlightBlock(codeBlock);
            // });
            // hljs.highlightBlock(container);
            // let a=  hljs.highlightAuto(store.state.viMessage,["markdown"])
            // marked(store.state.viMessage)
            // container.innerHTML = marked.parse(store.state.viMessage) as any
            const codeBlocks = container.querySelectorAll('code');
            codeBlocks.forEach((codeBlock) => {
                hljs.highlightBlock(codeBlock);
            });
            // hljs.highlightBlock(container);

        }
    }
})
//**视图工程文件(NNIE)**能增删和修改组件，但内容没有加密和压缩，也不能更新数据，主要用于设计；**视图编码文件(NNIC)**经过3重加密并压缩，组件无法增删组件，只能远程刷新数据，主要用于显示。
let lst = [
    {
        Q: '#### 视图工程文件(NNIE)和视图编码文件(NNIC)有什么区别？',
        A: '_视图工程文件(NNIE)_ 能导出 _视图编码文件(NNIC)_ 文件\n\n_视图工程文件(NNIE)_ 能增删和修改组件，但内容没有加密和压缩，也不能更新数据，主要用于设计；\n\n_视图编码文件(NNIC)_ 经过3重加密并压缩，组件无法增删组件，只能远程刷新数据，主要用于显示。'
    },
    {
        Q: '#### 视图工程文件(NNIE)设置的密码为什么无效？',
        A: '因为 _视图工程文件(NNIE)_ 中设置的密码是给 _视图编码文件(NNIC)_ 用的。'
    },
    {
        Q: '#### 视图编码文件(NNIC)密码忘记了怎么办？',
        A: '_视图编码文件(NNIC)_ 通过您输入的 _视图工程文件(NNIE)_ 密码采取一定算法派生出3个32字节的密码，\n\n分别对 _视图编码文件(NNIC)_ 的3部分内容进行对称加密，以我的技术并不能在没有密码的情况下进行解密，\n\n有技术和算力的同学或许可以试试。我皆以忘记密码后，直接从工程文件重新导出编码文件，并重新更新数据。'
    },
    {
        Q: '#### 视图编码文件(NNIC)支持什么方式更新数据？',
        A: '_视图编码文件(NNIC)_ 支持人工远程登陆更新数据和程序远程更新数据，人工登陆 _http://127.0.0.1:50520_ 远程更新，\n\n可以通过导入Excel表格完成更新数据，可以修改部分组件的部分属性（前提是组件开启 _远程更新_ 如果组件内容与表关联将不受该参数限制）；\n\n程序远程更新：\n\n```python\noptions = {\n\t"key":"e289ba5c-b692-4fd1-bcbd-7097e4734d48", #请登陆网页获取key，本key是表key，如果是组件就换成组件ID 例如：i0IV46AKSeIHhwf3bFeLF \n\t"data":json.dumps([["张三",69],["李四",88]]) , #原表格式是"姓名,分数"两列，如果更新的是组件，要json.dumps(["新内容"])，这个数组只有一个元素就是组件新内容。\n\t"replace":"true" #这个参数只有在覆盖更新的时候有效，增量更新或更新组件时无效\n}\nresponse = requests.post("http://127.0.0.1:50520/ai", data=options,headers={"token": "355458543157374159565337324e4547"}) #正常的话返回结果是"ok"\n```'
    },
]
</script>