import { divide, exp, multiply, number, string } from 'mathjs'
import { store, ViChartEnum, ViComponent, ViGroupStyle, ViOutStyle } from '../store'
import { angleToRadian } from './PositonAndSize'
// import { UploadFile } from '../../wailsjs/go/main/App'

// export function UploadFileToBackend(f:File) {
//     UploadFile(f).then(result => {
//         console.log("UploadFileToBackend result",result)
//     })
// }

export function sin(rotate: number) {
    return Math.abs(Math.sin(angleToRadian(rotate)))
}

export function cos(rotate: number) {
    return Math.abs(Math.cos(angleToRadian(rotate)))
}

export function viFormat(value: number, scale: number): number {
    return multiply(value, divide(scale, 100))
}

export function viDeFormat(value: number, scale: number): number {
    return divide(value, divide(scale, 100))
}

export function mod360(deg: number) {
    return (deg + 360) % 360
}

export function getNewValue(v: number) {
    let a = viDeFormat(v, store.state.canvas.scale_backup)
    let b = viFormat(a, store.state.canvas.scale)
    return b
}
export function getComponentRotatedStyle(style: {
    rotate: number,
    left: number,
    top: number,
    width: number,
    height: number,
    right: number,
    bottom: number,
}) {
    style = { ...style }
    if (style.rotate != 0) {
        const newWidth = style.width * cos(style.rotate) + style.height * sin(style.rotate)
        const diffX = (style.width - newWidth) / 2 // 旋转后范围变小是正值，变大是负值
        style.left += diffX
        style.right = style.left + newWidth

        const newHeight = style.height * cos(style.rotate) + style.width * sin(style.rotate)
        const diffY = (newHeight - style.height) / 2 // 始终是正
        style.top -= diffY
        style.bottom = style.top + newHeight

        style.width = newWidth
        style.height = newHeight
    } else {
        style.bottom = style.top + style.height
        style.right = style.left + style.width
    }

    return style
}

export function getElement(id: string): Element | null {
    // console.log("id", id)
    let els = document.querySelectorAll("[nni=nni" + id + "]")
    if (els.length < 1) {
        return null
    }
    return els.item(0)
}

// function getOriginalPosition(left: number, top: number, width: number, height: number, rotationAngle: number): { left: number, top: number } {
//     const cosTheta = Math.cos(rotationAngle * Math.PI / 180);
//     const sinTheta = Math.sin(rotationAngle * Math.PI / 180);
//     const rotationMatrix = [
//         [cosTheta, -sinTheta],
//         [sinTheta, cosTheta]
//     ];
//     const vector = [left, top];
//     const u = vector.map((v, i) => v + [width / 2, height / 2][i]);
//     const c = u.map((v, i) => v - rotationMatrix[i][0] * width / 2 - rotationMatrix[i][1] * height / 2);
//     const originalLeft = c[0] - width / 2;
//     const originalTop = c[1] - height / 2;
//     return { left: originalLeft, top: originalTop };
// }

export function uNGroupStyle(comRect: DOMRect, editorRect: DOMRect,
    parentStyle: ViOutStyle,
    groupStyle: ViGroupStyle) {

    const center = {
        x: comRect.left - editorRect.left + comRect.width / 2,
        y: comRect.top - editorRect.top + comRect.height / 2,
    }
    let width = groupStyle.width / 100 * parentStyle.width
    let height = groupStyle.height / 100 * parentStyle.height
    // 计算出元素新的 top left 坐标
    let left = center.x - width / 2
    let top = center.y - height / 2
    return { top, left }
}

// export function unGroupNewStyle(rect: DOMRect, r: number) {
//     let editorRect = store.state.viEditor!.getBoundingClientRect()
//     // let style = getComponentRotatedStyle({
//     //     rotate:-r,
//     //     left:rect.left,
//     //     top:rect.top,
//     //     width:rect.width,
//     //     height:rect.height,
//     //     right:0,
//     //     bottom:0
//     // })
//     let style = getOriginalPosition(rect.left, rect.top, rect.width, rect.height, r)
//     console.log(style)
//     // const center = {
//     //     x: rect.left - editorRect.left + rect.width / 2,
//     //     y: rect.top - editorRect.top + rect.height / 2,
//     // }
//     // let left = center.x - rect.width / 2
//     // let top = center.y - rect.height / 2
//     return { top: style.top - editorRect.top, left: style.left - editorRect.left }
// }
export function fileToBase64(file: File): Promise<string> {
    return new Promise<string>((resolve, reject) => {
        const reader = new FileReader()
        reader.readAsDataURL(file)
        reader.onload = () => {
            resolve(reader.result as string)
        }
        reader.onerror = (error) => {
            reject(error)
        }
    })
}

//下载特殊路径使用
export async function fetchEX(url: string) {
    //url.mime
    let urls = url.split(".")
    return await fetch(urls[0]).then(async (res) => {
        try {
            const bl = await res.blob()
            const videoBlob = new Blob([bl], { type: urls[1] });
            let blurl = URL.createObjectURL(videoBlob)
            return blurl
        } catch (reason) {
            console.log("fetchEX fetch blob", reason)
            return ""
        }
    }).catch((reason) => {
        console.log("fetchEX fetch ", reason)
        return ""
    })
}


//用来更新对象属性
// 使用示例
// updateObject(obj2, 'address', { ...obj2.address, zip: '90002' }, 'zip');
// updateObject(obj2, 'name', 'John');
function updateObject<T extends object, K extends keyof T>(obj: T, key: K, value: number | string): T {
    return {
        ...obj,
        [key]: value
    };
}

function updateNestedObject<T extends object, K1 extends keyof T, K2 extends keyof T[K1]>(obj: T, key1: K1, key2: K2, value: number | string): T {
    return {
        ...obj,
        [key1]: {
            ...obj[key1],
            [key2]: value
        }
    };
}

export function updateObjectWrapper<T extends object, K1 extends keyof T, K2 extends keyof T[K1] | undefined>(obj: T, key1: K1, key2: K2, value: number | string): T {
    if (key2 === undefined) {
        return updateObject(obj, key1, value);
    } else {
        return updateNestedObject(obj, key1, key2, value);
    }
}

const uuidPattern = /^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$/;

export function isUUID(input: string): boolean {
    return uuidPattern.test(input);
}
export function isCeilAI(input: string): boolean {
    return input.indexOf("CEILAI") != -1
}

//通过function string动态构建函数

// let x = NewDynamicFunction("function AA(a,b) { return a + b }")
// if(x){
//     console.log(x(1,2))
// }

import * as echarts from 'echarts'
import { nanoid } from 'nanoid'
export function NewDynamicFunction(functionstr: string) {

    const parameterPattern = /\(([^)]*)\)/;
    const match = functionstr.match(parameterPattern);
    if (match && match[1]) {
        const parameters = match[1].split(',').map(param => param.trim());
        function buildDynamicFunction(parameterNames: string[]): Function {
            return new Function(...parameterNames, 'return ' + functionstr)();
        }
        return buildDynamicFunction(parameters);
    } else {
        return undefined
    }
}


export function sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}

export function IsFunction(str: string) {
    const namedFunctionPattern = /\bfunction\s+(\w+)?\s*\([^)]*\)\s*{/;
    const anonymousFunctionPattern = /\(([^)]*)\)\s*=>\s*([^;]*)/;
    return namedFunctionPattern.test(str) || anonymousFunctionPattern.test(str)
}

//图表当前属性所在对象里的路径
export function get_current_path(item: string) {
    let path = store.state.viBreaditems.join('.') + "." + item
    if (store.state.ViChartRefIndex.type == ViChartEnum.Data) {
        let parts = path.split(".");
        path = "series[" + store.state.ViChartRefIndex.index + "]." + parts.slice(2).join(".");
    }
    path = path.replace("NNI.", "")
    return path
}

//给定的字符串 转换成对应值赋值给对象,path是图表当前属性所在对象里的路径
export function setup_option(str: string, path: string, id?: string) {
    if (IsFunction(str)) {
        if (isIllegalApiCall(str)) {
            //可能有风险的函数调用，用字符串原型
            let newstr = str.replace(/(['"])?([a-zA-Z0-9_]+)(['"])?:/g, '"$2": ').replaceAll("'", '"')
            store.commit("set_echart_option_by_path", {
                id: id,
                path: path,
                value: newstr
            })

        } else {
            //是函数需要构建函数，然后用函数赋值
            try {
                let func = NewDynamicFunction(str)
                if (func) {
                    store.commit("set_echart_option_by_path", {
                        id: id,
                        path: path,
                        value: func
                    })
                }
            } catch (error) {
                let newstr = str.replace(/(['"])?([a-zA-Z0-9_]+)(['"])?:/g, '"$2": ').replaceAll("'", '"')
                store.commit("set_echart_option_by_path", {
                    id: id,
                    path: path,
                    value: newstr
                })
            }
        }
    } else {
        //不是函数
        let newstr = str.replace(/(['"])?([a-zA-Z0-9_]+)(['"])?:/g, '"$2": ').replaceAll("'", '"')
        try {
            //先尝试转换成 对象、数组、布尔、数字、string 然后赋值
            let js = JSON.parse(newstr)
            store.commit("set_echart_option_by_path", {
                id: id,
                path: path,
                value: js
            })
        } catch (error) {
            //转换失败，比如一个没有引号的文字。。直接赋值，一般是比如线类型之乐的
            store.commit("set_echart_option_by_path", {
                id: id,
                path: path,
                value: newstr
            })
        }
    }
}



//下面是地理坐标转换


// 判断是否在国内，不在国内则不做偏移
function out_of_china(lng: number, lat: number) {
    return (lng < 72.004 || lng > 137.8347) || ((lat < 0.8293 || lat > 55.8271) || false)
}

export function WGS84ToBD09(lng: number, lat: number) {
    if (out_of_china(lng, lat)) {
        return { lng, lat }
    }
    const xPi: number = 3.14159265358979324 * 3000.0 / 180.0;
    const x: number = lng;
    const y: number = lat;
    const z: number = Math.sqrt(x * x + y * y) + 0.00002 * Math.sin(y * xPi);
    const theta: number = Math.atan2(y, x) + 0.000003 * Math.cos(x * xPi);
    const bd09Lon: number = z * Math.cos(theta) + 0.0065;
    const bd09Lat: number = z * Math.sin(theta) + 0.006;

    return { lng: bd09Lon, lat: bd09Lat };
}


//echarts自定义渲染函数
//coords坐标系列如[[x,y],[x1,y1]]
//baidumap是否是百度地图，如果是需要转为百度坐标
export function createRenderItem(baidumap: boolean, style: any[] | undefined) {
    if (style == undefined) {
        return
    }
    return function (params: any, api: any) {
        var points = [];
        for (var i = 0; i < style.length; i++) {
            let coords = style[i].source.CoordString as any[]
            // console.log("coords",coords)
            var lst = []
            for (var j = 0; j < coords.length; j++) {
                let coord = coords[j] as any[]
                if (baidumap) {
                    let tmp = WGS84ToBD09(coord[0], coord[1])
                    coord[0] = tmp.lng
                    coord[1] = tmp.lat
                }
                lst.push(api.coord([coord[0], coord[1]]))

            }
            points.push(lst);
            // break
        }
        // var color = api.visual('color');
        return {
            type: 'polygon',
            shape: {
                points: echarts.graphic.clipPointsByRect(points, {
                    x: params.coordSys.x,
                    y: params.coordSys.y,
                    width: params.coordSys.width,
                    height: params.coordSys.height
                })
            },
            style: api.style({
                fill: 'red',
                stroke: 'yellow' //echarts.color.lift(color, 5)
            })
        };
    };
}

//颜色转换AABBGGRR这种格式如何转换成RRGGBBAA
export function convertARGBtoRGBA(color: string): string {
    // 去除颜色值前的 '#' 符号
    color = color.slice(1);

    // 提取各个颜色通道的十六进制值
    const a = color.slice(0, 2);
    const b = color.slice(2, 4);
    const g = color.slice(4, 6);
    const r = color.slice(6, 8);

    // 重新组合为 RGBA 格式
    const rgba = r + g + b + a;

    // 添加 '#' 符号并返回结果
    return "#" + rgba;
}


//echarts 绘制多边形API

export function drawPolygon(points: any[], x: number, y: number, width: number, height: number) {
    return echarts.graphic.clipPointsByRect(points, {
        x: x,
        y: y,
        width: width,
        height: height
    })
}


export function createRenderItemAPI(bmap: any) {
    return function renderItem(params: any, api: any) {
        let coordsstr = api.value(0) as string
        var coor_points: [number, number][] = [];
        if (coordsstr) {
            const re = /(\d+\.\d+),(\d+\.\d+),\d+/g;
            let match: RegExpExecArray | null;
            while ((match = re.exec(coordsstr)) !== null) {
                let lng = parseFloat(match[1]);
                let lat = parseFloat(match[2]);
                if (bmap) {
                    let tmp = WGS84ToBD09(lng, lat)
                    lng = tmp.lng
                    lat = tmp.lat
                }
                coor_points.push(api.coord([lng, lat]));
            }

            // let arr = coordsstr.split(" ")
            // for (let i = 0; i < arr.length; i++) {
            //     let coord = arr[i].split(",")
            //     if (coord.length >= 2) {
            //         if (coord[0].length > 0 && coord[1].length > 0) {
            //             let lng = parseFloat(coord[0])
            //             let lat = parseFloat(coord[1])
            //             if (bmap) {
            //                 let tmp = WGS84ToBD09(lng, lat)
            //                 lng = tmp.lng
            //                 lat = tmp.lat
            //             }
            //             coor_points.push(api.coord([lng, lat]));
            //         }
            //     }
            // }
        }
        let color = api.value(2) as string
        if (color.length >= 6) {
            color = convertARGBtoRGBA("#" + color)
        } else {
            color = "#003355"
        }
        let fill = api.value(3) as string
        if (fill.length > 6) {
            fill = convertARGBtoRGBA("#" + fill)
        } else {
            fill = "#003355" //echarts.color.lift(api.visual('color'), 5)
        }
        // points: echarts.graphic.clipPointsByRect(points, {
        //     x: params.coordSys.x,
        //     y: params.coordSys.y,
        //     width: params.coordSys.width,
        //     height: params.coordSys.height
        // })
        // points: drawPolygon(points,params.coordSys.x,params.coordSys.y,params.coordSys.width,params.coordSys.height)
        // points:points
        if (coor_points.length > 0) {

            return {
                type: 'polygon',
                shape: {

                    points: echarts.graphic.clipPointsByRect(coor_points, {
                        x: params.coordSys.x,
                        y: params.coordSys.y,
                        width: params.coordSys.width,
                        height: params.coordSys.height
                    })
                },
                style: api.style({
                    fill: fill,
                    stroke: color,
                    name: api.value(1),
                })
            };
        }
    }
}

//检测文本内容是否违规调用
export function isIllegalApiCall(code: string) {
    const patterns = [/document/, /fetch/, /XMLHttpRequest/, /WebSocket/, /navigator/, /window/, /screen/, /location/, /localStorage/, /sessionStorage/,/process/,/exec/,/eval/,/new Function/];
    for (const pattern of patterns) {
        if (pattern.test(code)) {
            return true;
        }
    }
    return false;
}


export function nanoidEx(){
    let ret = nanoid()
    ret = ret.replaceAll("_","N")
    ret = ret.replaceAll("-","N")
    return ret
}