import { ViChartMainType, ViComponent } from "../../store";

export const ComponentsName = [
    "ViAudio",
    // "ViButton",
    // "ViProgress",
    // "ViCheckBox",
    // "ViCode",
    "ViData",
    "ViFile",
    // "ViFontIcon",
    "ViEmpty",
    'ViImage',
    // "ViInput",
    "ViLabel",
    // "ViList",
    // "ViMath",
    // "ViRadio",
    "ViSVGIcon",
    // "ViTimer",
    // "ViVariable",
    'ViVideo',
    "ViGroup",
    "ViChart"
]

const OutStyle = {
    top: 0,
    left: 0,
    width: 80,
    height: 20,

    backgroundColor: 'transparent',//'rgba(255,255,255,0.5)',
    // color: 'white',

    borderTopColor: 'transparent',
    borderBottomColor: 'transparent',
    borderLeftColor: 'transparent',
    borderRightColor: 'transparent',

    borderTopStyle: 'none',
    borderBottomStyle: 'none',
    borderLeftStyle: 'none',
    borderRightStyle: 'none',

    borderTopWidth: '0px',
    borderBottomWidth: '0px',
    borderLeftWidth: '0px',
    borderRightWidth: '0px',

    borderTopLeftRadius: '0px',
    borderTopRightRadius: '0px',
    borderBottomLeftRadius: '0px',
    borderBottomRightRadius: '0px',

    rotate: 0,
    opacity: 1,
    zIndex: 0,
    displayMode: false,
}

const GroupStyle = {
    top: 0,
    left: 0,
    width: 0,
    height: 0
}
export const Components = {
    //标签组件属性
    "ViLabel": {
        id: "",
        name: "ViLabel",
        component: "ViLabel",
        typeName: "文本",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "文本-",
        innerStyle: {
            color: 'white',
            text: "文本",
            vertical: false,
            //是否支持远程更新
            styleFlag: false,
            fontfamily: "sans-serif",
            fontSize: 24,
            fontWeight:'normal',
            fontStyle:'normal',
            list: [],
        },
        style: {
            ...OutStyle
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },
    //图片标签
    "ViImage": {
        id: "",
        name: "ViImage",
        component: "ViImage",
        typeName: "图片",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "图片-",
        innerStyle: {
            // color:'white',
            //图片url
            text: "",
            //是否铺满，本参数在v-img中，暂时没有起到作用，所以是预留到后面用的
            vertical: false,
            //是否支持远程更新
            styleFlag: false,
            list: [],
        },
        style: {
            ...OutStyle,
            height: 80,
            width: 200,
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },
    //按钮组件属性
    // "ViButton": {
    //     id: "",
    //     name: "ViButton",
    //     component: "ViButton",
    //     typeName: "按钮",
    //     //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    //     customName: "按钮-",
    //     innerStyle: {
    //         text: "按钮",
    //         fontfamily: "sans-serif",
    //         fontSize: 24,
    //     },
    //     style: {
    //         ...OutStyle,
    //         width: 128,
    //         height: 48,
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: [],
    //     chartStyle: undefined,
    // },


    //进度属性
    // "ViProgress": {
    //     id: "",
    //     name: "ViProgress",
    //     component: "ViProgress",
    //     typeName: "进度",
    //     //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    //     customName: "进度-",
    //     innerStyle: {
    //         color:'white',
    //         // text: "进度",
    //     },
    //     style: {
    //         ...OutStyle,
    //         width: 128,
    //         height: 48,
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: [],
    //     chartStyle: undefined,
    // },

    //声音标签
    "ViAudio": {
        id: "",
        name: "ViAudio",
        component: "ViAudio",
        typeName: "声音",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "声音-",
        innerStyle: {
            color: 'white',
            //声音url
            text: "",
            //是否支持远程更新
            styleFlag: false,
            list: [],
        },
        style: {
            ...OutStyle,
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },
    //视频标签
    "ViVideo": {
        id: "",
        name: "ViVideo",
        component: "ViVideo",
        typeName: "视频",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "视频-",
        innerStyle: {
            //url
            text: "http://0.0.0.0:8080/flower.webm",
            //是否自动播放
            vertical: false,
            //是否支持远程更新
            styleFlag: false,
            list: [],
        },
        style: {
            ...OutStyle,
            height: 80,
            width: 200,
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },
    //文件标签
    "ViFile": {
        id: "",
        name: "ViFile",
        component: "ViFile",
        typeName: "文件",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "文件-",
        innerStyle: {
            color: 'white',
            //文件名称，但显示的时候不会显示后缀名，
            text: "文件名称.txt",
            //是否支持远程更新
            styleFlag: false,
            fontfamily: "sans-serif",
            fontSize: 24,
            fontWeight:'normal',
            fontStyle:'normal',
            list: [],
        },
        style: {
            ...OutStyle,
            width: 80,
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },
    //函数组件属性
    // "ViMath": {
    //     id: "",
    //     name: "ViMath",
    //     component: "ViMath",
    //     typeName: "函数",
    //     //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    //     customName: "函数-",
    //     innerStyle: {
    //         fontSize: 24,
    //         text: `f\\lparen\\relax{\\red{x}}\\rparen=\\int_{-\\infty}^\\infty\\hat{f}\\lparen\\xi\\rparen\\,e^{2 \\pi i \\xi \\red{x}}\\,d\\xi`
    //     },
    //     style: {
    //         ...OutStyle,
    //         width: 420,
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: [],
    //     chartStyle: undefined,
    // },
    //代码组件属性
    // "ViCode": {
    //     id: "",
    //     name: "ViCode",
    //     component: "ViCode",
    //     typeName: "代码",
    //     //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    //     customName: "代码-",
    //     innerStyle: {
    //         fontSize: 24,
    //         text: `
    //     func f_x(x float64, xi []float64) float64 {
    //         // 计算函数在x处的取值
    //         var res complex128
    //         for i, xi_i := range xi {
    //             res += cmplx.Conj(complex(math.Exp(
    //                 -math.Pow(xi_i, 2)), 0)) 
    //             * cmplx.Exp(
    //                 2*math.Pi*1i*xi_i*x/float64(len(xi)))
    //         }
    //         return real(res)
    //     }
    //     `,
    //     },
    //     style: {
    //         ...OutStyle,
    //         width: 400,
    //         height: 220,
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: [],
    //     chartStyle: undefined,
    // },
    //白板标签
    "ViEmpty": {
        id: "",
        name: "ViEmpty",
        component: "ViEmpty",
        typeName: "白板",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "白板-",
        innerStyle: {
            //白板内部什么东西都没有，
            //主要用于制作特殊背景
            list: [],
        },
        style: {
            ...OutStyle,
            height: 80,
            width: 200,
            backgroundColor:'rgba(255,255,255,0.05)'
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },
    //单选标签
    // "ViRadio": {
    //     id: "",
    //     name: "ViRadio",
    //     component: "ViRadio",
    //     typeName: "单选",
    //     //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    //     customName: "单选-",
    //     innerStyle: {
    //         text: "选项1",
    //         textList: ["选项1", "选项2", "选项3"],
    //         vertical:false,
    //         styleFlag: false,
    //         fontfamily: "sans-serif",
    //         fontSize: 24,
    //     },
    //     style: {
    //         ...OutStyle,
    //         width: 400,
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: [],
    //     chartStyle: undefined,
    // },
    // //多选标签
    // "ViCheckBox": {
    //     id: "",
    //     name: "ViCheckBox",
    //     component: "ViCheckBox",
    //     typeName: "多选",
    //     //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    //     customName: "多选-",
    //     innerStyle: {
    //         textList: ["选项1", "选项2", "选项3"],
    //         selectList: ["选项1", "选项2"],
    //         vertical:false,
    //         styleFlag: true,
    //         fontfamily: "sans-serif",
    //         fontSize: 24,
    //     },
    //     style: {
    //         ...OutStyle,
    //         width: 400,
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: [],
    //     chartStyle: undefined,
    // },
    //列表标签
    // "ViList": {
    //     id: "",
    //     name: "ViList",
    //     component: "ViList",
    //     typeName: "列表",
    //     //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    //     customName: "列表-",
    //     innerStyle: {
    //         fontfamily: "sans-serif",
    //         fontSize: 24,
    //         textList: ["选项1", "选项2", "选项3", "选项4", "选项5"],
    //         selectList: ["选项1", "选项3"],
    //         styleFlag: true,
    //     },
    //     style: {
    //         ...OutStyle,
    //         // width:400,
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: [],
    //     chartStyle: undefined,
    // },
    //定时组件属性
    // "ViTimer": {
    //     id: "",
    //     name: "ViTimer",
    //     component: "ViTimer",
    //     typeName: "定时",
    //     //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    //     customName: "定时-",
    //     innerStyle: {
    //         color: 'white',
    //         //日期格式
    //         text: "[DD] HH:mm:ss",
    //         //2022-05-20T10:00:00
    //         //日期内容，第一个是日期，第二个是时间
    //         textList: ["2022-05-20", "10:00:00"],
    //         //是否支持远程更新
    //         styleFlag: false,
    //         fontfamily: "sans-serif",
    //         fontSize: 24,
    //         fontWeight:'normal',
    //         fontStyle:'normal',
    //         list: [],
    //     },
    //     style: {
    //         ...OutStyle
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: [],
    //     chartStyle: undefined,
    // },
    //输入组件属性
    // "ViInput": {
    //     id: "",
    //     name: "ViInput",
    //     component: "ViInput",
    //     typeName: "输入",
    //     //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
    //     customName: "输入-",
    //     innerStyle: {
    //         fontfamily: "sans-serif",
    //         fontSize: 24,
    //         text: "输入框",
    //         //宽度
    //         // text2:"200",
    //     },
    //     style: {
    //         ...OutStyle,
    //         width: 100,
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: [],
    //     chartStyle: undefined,
    // },
    //表格组件属性
    "ViData": {
        id: "",
        name: "ViData",
        component: "ViData",
        typeName: "表格",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "表格-",
        innerStyle: {
            color: 'white',
            fontfamily: "sans-serif",
            fontSize: 24,
            fontWeight:'normal',
            fontStyle:'normal',
            list: [],
        },
        style: {
            ...OutStyle,
            height: 160,
            width: 500,
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },
    // //可变标签组件属性
    // "ViVariable": {
    //     id: "",
    //     name: "ViVariable",
    //     component: "ViVariable",
    //     typeName: "变量",
    //     innerStyle: {
    //         fontfamily: "sans-serif",
    //         fontSize: 24,
    //         text: "支持远程更新的文本"
    //     },
    //     style: {
    //         ...OutStyle
    //     },
    //     groupStyle: {
    //         ...GroupStyle
    //     },
    //     components: []
    // },

    //图标组件属性，有svg icon 和 font icon两种组件会使用这个属性，
    //也就是 ViFontIcon ViSVGIcon
    //component是ViFontIcon时：
    //name = event.dataTransfer.getData("Text").substring(3)
    //因为字体图标是默认的图标集，但传过来是md:这3个字节开头，所以要把这3个字节去掉，之后剩下的字符是对应相应图标名称
    //component是ViSVGIcon时：
    //name=event.dataTransfer.getData("Text")
    //用前三个字节区分 ViSVGIcon 的图标是mdi开头的字符串，mdi后面对应相应图标
    /*
        识别前三个字母：
        ViFontIcon md:
        ViSVGIcon mdi
     */
    //虽然ViFontIcon已经被弃用，但还是在软件中注册
    "ViIcon": {
        id: "",
        name: "",
        component: "",
        typeName: "图标",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "图标-",
        innerStyle: {
            color: 'white',
            styleFlag: false,
            list: [],
        },
        style: {
            ...OutStyle,
            width: 26,
            height: 26,
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },
    //组组件属性
    /*
        不像正常组件一样，因为组组件不需要拖拽创建，
        有自己的创建按钮，所以不需要识别event.dataTransfer.getData("Text")
        当然也就不需要识别的关键字了。
        
    */
    "ViGroup": {
        id: "",
        name: "ViGroup",
        component: "ViGroup",
        typeName: "组",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "组-",
        innerStyle: {
            list: [],
        },
        style: {
            ...OutStyle,
            width: 0,
            height: 0,
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },

    //图表组件，显示echarts用的
    "ViChart": {
        id: "",
        name: "ViChart",
        component: "ViChart",
        typeName: "图表",
        //自定义名称，默认组件typeName开头，初始化时应当给组件初始化一个序号，方便编辑人员修改。
        customName: "图表-",
        innerStyle: {
            list: [],
        },
        style: {
            ...OutStyle,
            height: 100,
            width: 200,
        },
        groupStyle: {
            ...GroupStyle
        },
        components: [],
        chartStyle: undefined,
    },
}

