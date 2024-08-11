
export const OutStyleInfo = {
    top: '组件顶部距离',
    left: '组件左边距离',
    width: '组件宽度',
    height: '组件高度',
    rotate: '组件选择角度',
    backgroundColor: '组件背景颜色',

    borderTopColor: '组件上边框颜色',
    borderBottomColor: '组件下边框颜色',
    borderLeftColor: '组件左边框颜色',
    borderRightColor: '组件右边框颜色',

    borderTopStyle: '组件上边框样式',
    borderBottomStyle: '组件下边框样式',
    borderLeftStyle: '组件左边框样式',
    borderRightStyle: '组件右边框样式',

    borderTopWidth: '组件上边框宽度',
    borderBottomWidth: '组件下边框宽度',
    borderLeftWidth: '组件左边框宽度',
    borderRightWidth: '组件右边框宽度',

    borderTopLeftRadius: '组件左上角圆角',
    borderTopRightRadius: '组件右上角圆角',
    borderBottomLeftRadius: '组件左下角圆角',
    borderBottomRightRadius: '组件右上角圆角',

    opacity: '组件透明度',
    zIndex: '组件堆叠顺序',
}
/*
返回组件内部样式的中文内容，因为不同组件可能用同一个key名称来表示不同的内容
*/
export function InnerStyleInfo(key: string, comptype: string) {
    switch (key) {
        case "fontfamily":
            return "字体名称"
        case "fontSize":
            return "字体大小"
        case "color":
            return "前景颜色"
        case "fontStyle":
            return "字体风格"
        case "fontWeight":
            return "字体粗细"
        case "text":
            switch (comptype) {
                case "ViLabel":
                    return "文本内容"
                case "ViImage":
                    return "选择图片"
                case "ViButton":
                    return "按钮内容"
                case "ViAudio":
                    return "选择声音"
                case "ViVideo":
                    return "选择视频"
                case "ViFile":
                    return "选择文件"
                case "ViMath":
                    return "函数代码"
                case "ViCode":
                    return "代码内容"
                case "ViRadio":
                    return "当前选择"
                case "ViTimer":
                    return "显示格式"
                case "ViInput":
                    return "输入内容"
                default:
                    return ""
            }
        case "textList":
            switch (comptype) {
                case "ViRadio":
                    return "可选内容"
                case "ViCheckBox":
                    return "可选内容"
                case "ViList":
                    return "可选内容"
                case "ViTimer":
                    return "日期时间"
                default:
                    return ""
            }
        case "selectList":
            switch (comptype) {
                case "ViCheckBox":
                    return "已选内容"
                case "ViList":
                    return "已选内容"
                default:
                    return ""
            }
        case "styleFlag":
            return "远程更新"
        case "vertical":

            switch (comptype) {
                case "ViImage":
                    return "是否平埔"
                case "ViVideo":
                    return "自动播放"
                default:
                    return "垂直显示"
            }
        default:
            break;
    }
    return ""
}