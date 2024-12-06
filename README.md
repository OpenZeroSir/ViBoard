## 软件介绍

ViBoard 数据大屏可视化软件，基于Excel数据设计可视化，像PPT一样设计、传播和显示，但也能远程更新数据。

## Features
- Powerpoint-like
- Designed based on Excel data
- Update data remotely
- Cross platform

## 特别声明

本人并非专业程序员，由于第一次使用这么多语言和框架来设计程序，可能确实比较吃力，主要语言有go、typescript、javascript、html、css和sql,主要框架有wails、vue3、ECharts和vuetify等，代码中有很多烂代码甚至严重设计缺陷，本想在第二个版本中改进和完善，但由于个人时间、精力、技术和资金等原因导致第二个版本遥遥无期，暂时先发布第一版本的代码。

## 简易逻辑
```c
设计模式->导入excel->设计视图(*.NNIE)->导出视图编码(*.NNIC)
显示模式->打开视图编码(*.NNIC)->打开远程服务->人工或自动化定时更新数据
```

## 简易教程
[demo.webm](https://github.com/OpenZeroSir/ViBoard/assets/121022414/511105f5-e774-48e5-8380-e692143aa5c3)


## 更新数据
```Python
启动远程服务 #启动之后就可以远程更新数据。

# 人工远程更新数据

登陆 http://127.0.0.1:50520 # 账号密码在显示模式下可以看到，该模式支持人工用Excel文件定时更新数据。


# 自动化更新数据

options = {
    "key":"e289ba5c-b692-4fd1-bcbd-7097e4734d48", # 请登陆网页获取key，本key是表key，
                                                  # 如果是组件就换成组件ID 例如：i0IV46AKSeIHhwf3bFeLF 
    "data":json.dumps([["张三",69],["李四",88]]) , # 原表格式是"姓名,分数"两列，
                                                  # 如果更新的是组件，要json.dumps(["新内容"])，
                                                  # 这个数组只有一个元素就是组件新内容。
    "replace":"true" # 这个参数只有在覆盖更新的时候有效，增量更新或更新组件时无效。
}


requests.post(
    "http://127.0.0.1:50520/ai", 
    data=options,
    headers={
        "token": "355458543157374159565337324e4547" # token可以在显示模式下获取。
    }
) 

# 正常的话返回结果是"ok"
```

邮箱：`open@zerosir.com`


