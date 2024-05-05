## 软件介绍

ViBoard是个轻量数据可视化软件，软件采用echarts提供大量图表组件用于设计，所有图表支持高度自定义，并且不依赖于服务器和开发人员，非开发人员也可以设计出高端图表数据可视化。软件分为`设计模式`和`显示模式`，`设计模式`用于设计界面样式和设计图表样式，`显示模式`用于显示已经设计好的界面和图表并提供更新数据接口服务，设计好的界面和图表可以导出给`显示模式`显示。

## 软件下载

| Windows | Linux | Mac  | 教程 | 案例 |
| ----    | ----  | ---- | ---- | ---- |
| [下载](https://github.com/nniai/ViBoard/releases/download/1.0.0/ViBoard_1.0.0_windows_amd64.rar) | [下载](https://github.com/nniai/ViBoard/releases/download/1.0.0/ViBoard_1.0.0_linux_amd64.xz) | - | [简易教程](https://cdn.jsdelivr.net/gh/nniai/ViBoard@main/assets/readme.webm)  | [案例参考](https://nniai.github.io/assets/style.webm) |

## 简易逻辑
```c
设计模式->导入excel->设计视图(*.NNIE)->导出视图编码(*.NNIC)
显示模式->打开视图编码(*.NNIC)->打开远程服务->人工或自动化定时更新数据
```

## 软件界面

1. 设计模式界面

![](https://github.com/nniai/ViBoard/blob/main/assets/screenshot.webp)

2. 显示模式界面

![](https://github.com/nniai/ViBoard/blob/main/assets/screenshot1.webp)

## 更新数据

视图编码文件(NNIC) 支持人工远程登陆更新数据和程序远程更新数据，人工登陆 http://127.0.0.1:50520 远程更新，
可以通过导入Excel表格完成更新数据，可以修改部分组件的部分属性（前提是组件开启 远程更新 如果组件内容与表关联将不受该参数限制）；
程序远程更新：
```Python
options = {
    "key":"e289ba5c-b692-4fd1-bcbd-7097e4734d48", #请登陆网页获取key，本key是表key，
                                                  #如果是组件就换成组件ID 例如：i0IV46AKSeIHhwf3bFeLF 
    "data":json.dumps([["张三",69],["李四",88]]) , #原表格式是"姓名,分数"两列，
                                                  #如果更新的是组件，要json.dumps(["新内容"])，这个数组只有一个元素就是组件新内容。
    "replace":"true" #这个参数只有在覆盖更新的时候有效，增量更新或更新组件时无效
}
response = requests.post("http://127.0.0.1:50520/ai", data=options,
headers={
    "token": "355458543157374159565337324e4547" #token可以在软件上获取到
}) #正常的话返回结果是"ok"
```

## 关于我们

我叫`李飞龙`，软件由我个人技术、爱好和资金在驱动，软件免费但无法保证任何包含但不仅限经济、安全和质量等问题，除了以下情况:
1. 软件不会植入广告
2. 软件不会留有后门
3. 软件不会收集信息

以上条件成立的前提是通过合规渠道获取到软件，任何形式交换获取软件或包含但不仅限参与软件官方或非官方活动等，均可以理解为骗局。

如您发现软件问题，可以把详细问题发送到我邮箱，当然也可以请我喝杯咖啡或者给我点个星星。

邮箱：`public@nniai.com`


![](https://github.com/nniai/ViBoard/blob/main/assets/wechat.webp)
