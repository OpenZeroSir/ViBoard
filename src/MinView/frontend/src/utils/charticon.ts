export const ChartIcons = [
    //空白
    {
        id: "741fef90-a46b-4cdc-9995-07515af1db31",
        message: "空白图 双击查看完整帮助",
        help: '本模板是空模板，理论上可以配置出很多形式的图表，但有很多内容需要了解，如果预设的模板有适合的，尽量使用预设模板，使用任何模板都应当遵守**尽量优先增加坐标系统再增加数据**'
    },
    // 折线图 start
    {
        id: "1E51ACC7-9E58-4704-BF1B-1142A07BAE3E",
        message: "基础折线图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**折线图**\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**折线图**属性上\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',
    },
    {
        id: "9cb49149-c520-4f53-a69e-4423e963f3d2",
        message: "堆积折线图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**两个折线图**\n\n- [x] 分别给两个折线图设置**stack**为"Total"\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**折线图**属性上\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',
    },
    {
        id: "9b1a3b1a-78ed-4308-9d7c-aae1e731248d",
        message: "基础面积图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**折线图** 给**areaStyle.color**设置一个颜色\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**折线图**属性上\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',
    },
    {
        id: "586e65fd-aed4-42ab-a797-271f578d7f5e",
        message: "堆积面积图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**两个折线图**\n\n- [x] 分别给两个折线图设置**stack**为"Total" 分别给**areaStyle.color**设置一个颜色\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**折线图**属性上\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',
    },
    {
        id: "f58cb5f2-ede3-4923-8fd2-65e5ef28641d",
        message: "基础柱状图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**柱状图**\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**柱状图**属性上\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',
    },
    {
        id: "1885c347-816b-4d9f-8f23-e80ad3e7c5b2",
        message: "堆积柱状图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**两个柱状图** 分别给两个柱状图设置**stack**为"Total"\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**柱状图**属性上\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',
    },
    {
        id: "db94cd1c-f811-494f-9803-fc0a57d5df0b",
        message: "水平柱状图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**柱状图**\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**柱状图**属性上\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',
    },
    {
        id: "5c1631ea-7831-4788-a132-c866f20c7548",
        message: "水平堆积柱状图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**两个柱状图** 分别给柱状图设置**total**为"total"\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**柱状图**属性上\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',

    },
    {
        id: "00de94ff-ecda-4231-acde-f49e8777933d",
        message: "饼图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**饼图**\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n- [ ] 从右上角的数据中拖一列数据到**饼图**属性上，如果需要饼图显示名称，可以再拖一列数据进去。'
    },
    {
        id: "ba2fafed-516f-48ce-b5e4-0925196edbf7",
        message: "环形饼图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**饼图** 设置饼图**radius**为["40%", "70%"]\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n- [ ] 从右上角的数据中拖一列数据到**饼图**属性上，如果需要饼图显示名称，可以再拖一列数据进去。'

    },
    // {
    //     id: "97854320-17ad-4c48-9dee-a560ecad171c",
    //     message: "半形饼图",
    // },
    {
        id: "3bb52213-654f-4bfc-b749-1bb7b698036a",
        message: "南丁格尔玫瑰图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**饼图** 设置饼图**radius**为"70%" 设置**center**为["50%", "50%"] 设置**roseType**为"area"\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n- [ ] 从右上角的数据中拖一列数据到**饼图**属性上，如果需要饼图显示名称，可以再拖一列数据进去。'

    },
    {
        id: "1782ebf9-589d-426c-b37f-9c4f5f8e7274",
        message: "极坐标柱状图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**极坐标系** 设置**radius**为[30, "80%"]\n\n- [x] 在图表**内部属性**添加**极坐标径向抽**\n\n- [x] 在图表**内部属性**添加**极坐标角度抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**柱状图** 设置**coordinateSystem**为"polar"\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n- [ ] 从右上角的数据中拖一列数据到**柱状图**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**极坐标角度抽**属性上'
    },
    {
        id: "643a5985-d595-4163-955d-b4cd1b212083",
        message: "极坐标柱状环形图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**极坐标系** 设置**radius**为[30, "80%"]\n\n- [x] 在图表**内部属性**添加**极坐标径向抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**极坐标角度抽**\n\n- [x] 在图表**内部属性**添加**柱状图** 设置**coordinateSystem**为"polar"\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n- [ ] 从右上角的数据中拖一列数据到**柱状图**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**极坐标径向抽**属性上'

    },
    {
        id: "3a2c1a57-66c8-4ae1-80b2-7866d68e415e",
        message: "极坐标柱状堆积图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**极坐标系**\n\n- [x] 在图表**内部属性**添加**极坐标径向抽**\n\n- [x] 在图表**内部属性**添加**极坐标角度抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加多个**柱状图** 设置**coordinateSystem**为"polar" 设置**stack**为"a"\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n- [ ] 从右上角的数据中分别拖一列数据到对应**柱状图**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**极坐标角度抽**属性上'
    },
    {
        id: "8a5dcc6b-0aa4-40f6-8323-d6ee3e18c59c",
        message: "极坐标柱状状环堆积图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**极坐标系**\n\n- [x] 在图表**内部属性**添加**极坐标径向抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**极坐标角度抽**\n\n- [x] 在图表**内部属性**添加多个**柱状图** 设置**coordinateSystem**为"polar" 设置**stack**为"a"\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n- [ ] 从右上角的数据中分别拖一列数据到对应**柱状图**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**极坐标径向抽**属性上'

    },
    {
        id: "f0dd12cd-0bcb-47ab-a94e-ade76cbb56af",
        message: "极坐标散点图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**极坐标系**\n\n- [x] 在图表**内部属性**添加**极坐标径向抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**极坐标角度抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**散点图** 设置**coordinateSystem**为"polar"设置**symbolSize**为\n\n```livecodeserver\n(value, params) => value[2]\n```\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n- [ ] 从右上角的数据中拖一列数据到**极坐标径向抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**极坐标角度抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**散点图**属性上\n\n第一列内容取值[0,n]是径向抽数据的序号\n\n第二列内容取值[0,m]是角度抽数据的序号\n\n如果有其他需要可以拖入更多列数据进来。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n'

    },
    {
        id: "b8563795-6f32-4f8c-a14a-06b93fa93435",
        message: "散点图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**散点图**\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖多列数据到**散点图**属性上 第一列是横坐标数据 第二列 纵坐标数据 第三列散点数据 如果还有其他用途可以拖其他数据进来\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',

    },


    {
        id: "8c798693-672f-4a5a-8cd1-107151ada9ba",
        message: "单抽坐标图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**单点坐标** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**散点图**\n\n 设置**coordinateSystem**为"singleAxis" 设置**symbolSize**为\n\n```livecodeserver\n(value)=>value[1]\n```\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**单点坐标抽**属性上\n\n- [ ] 从右上角的数据中拖多列数据到**散点图**属性上 第一列跟单点坐标一样的数据 第二列是散点数据\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',

    },
    // {
    //     id: "199c00b9-d3b9-4dd4-bfc0-8344089b1268",
    //     message: "平行坐标",
    // },

    {
        id: "d6f4f906-9539-4a28-b03e-486de2f0876c",
        message: "仪表盘 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**仪表盘** 设置**progress.color**为true\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**仪表盘**属性上 如果想显示名称 可以再拖一列数据进去'
    },
    {
        id: "de5415ca-2c0d-4700-87ce-12ef2955a344",
        message: "雷达图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**雷达坐标**\n\n- [x] 在图表**内部属性**添加**雷达图**\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**雷达坐标**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**雷达图**属性上\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```'
    },
    {
        id: "9d6749b8-6ac3-4f88-8cc8-e78952c67239",
        message: "盒须图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**箱形图**\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中分别拖5列数据到**箱形图**属性上\n\n```properties\n第一列:最小值\n第二列:下四分位数（第25百分位数）\n第三列:中位数（第50百分位数）\n第四列:上四分位数（第75百分位数）\n第五列:最大值\n```\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',

    },
    {
        id: "0c79ae07-91ef-4a17-adaa-6e70c7fbf152",
        message: "K线图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"value"\n\n- [x] 在图表**内部属性**添加**K线图**\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中分别拖4列数据到**K线图**属性上\n\n```properties\n第一列:开盘值\n第二列:收盘值\n第三列:最低值\n第四列:最高值\n```\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',

    },
    {
        id: "2caf407d-c219-4780-a66a-4f53ddd52421",
        message: "漏斗图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**漏斗图**\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**漏斗图**属性上 如果想显示名称 可以再拖一列数据进去'

    },
    {
        id: "48088e85-a412-4140-9bd9-56b2f7f93154",
        message: "主题河流图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**单点坐标** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**主题河流**\n\n 设置**coordinateSystem**为"singleAxis"\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**单点坐标抽**属性上\n\n- [ ] 从右上角的数据中分别拖3列数据到**主题河流**属性上 第一列取值是单点坐标数据里的内容 第二列是主题河流数据 第三列是类别 多个类别显示多条河流\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',

    },
    {
        id: "c9ea4933-c5a0-4432-9156-ca0cb7a8fd61",
        message: "笛卡尔坐标系上的热力图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**横坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**纵坐标抽** 设置**type**为"category"\n\n- [x] 在图表**内部属性**添加**热力图**\n\n- [x] 在图表**内部属性**添加**连续映射** 设置**min**为0 设置**max**为10\n\n- [x] 在图表**内部属性**添加**网格**\n\n- [x] 设置**网格**的**top**和**bottom**为"20%"\n\n这个设置可以防止缩放到很小的时候出现显示异常，大部分图表没有这种情况。\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中拖一列数据到**横坐标抽**属性上\n\n- [ ] 从右上角的数据中拖一列数据到**纵坐标抽**属性上\n\n- [ ] 从右上角的数据中分别拖3列数据到**热力图**属性上 第一列取值范围应该是**横坐标抽**数据里的内容 第二列取值应该是**纵坐标抽数据里的内容** 第三列是热力值\n\n采用本模板只需把本模板拖入到画布中，并拖入数据到对应属性上即可，也可以自行增加数据属性。\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n',
    },
    {
        id: "71a7a7a2-99c3-4575-bc03-4c1706a5f39b",
        message: "地图路径图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**百度地图** 设置**zoom** 设置**center** 设置**roam**为true\n\n- [x] 在图表**内部属性**添加**路径图** 设置**coordinateSystem**为"bmap" 设置**polyline**为true\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中分别拖2列数据到**路径图**属性上 第一列是经度 第二列是纬度 如果需要路径名称 需要拖入第三列 如果是多条路径 可以在数据中把经纬度设置为0用来分割\n\n```livecodeserver\n百度地图密钥申请 https://lbsyun.baidu.com 创建应用需选择为"浏览器端"\n```'
    },
    {
        id: "c163399d-06b8-436a-b620-070a87947201",
        message: "地图散点 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**百度地图** 设置**zoom** 设置**center** 设置**roam**为true\n\n- [x] 在图表**内部属性**添加**散点图** 设置**coordinateSystem**为"bmap"\n\n- [x] 在图表**内部属性**添加**涟漪图** 设置**coordinateSystem**为"bmap"\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中分别拖3列数据到**散点图/涟漪图**属性上 第一列是经度 第二列是纬度 第三列是散点/涟漪值 如需要设置散点大小颜色或名称等需要拖入更多数据并用函数来实现\n\n```livecodeserver\n百度地图密钥申请 https://lbsyun.baidu.com 创建应用需选择为"浏览器端"\n```\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n'

    },
    {
        id: "efcabf75-68a7-4941-8208-c0fd38a8aa64",
        message: "地图热力图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**百度地图** 设置**zoom** 设置**center** 设置**roam**为true\n\n- [x] 在图表**内部属性**添加**热力图** 设置**coordinateSystem**为"bmap" 设置**pointSize**为5 设置**blurSize**为6\n\n- [x] 在图表**内部属性**添加**连续映射** 设置**min**为0 设置**max**为5\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中分别拖3列数据到**热力图**属性上 第一列是经度 第二列是纬度 第三列数据 如果第三列没有可以用1填充第三列按密集热力渲染 否则按第三列数据大小渲染热力\n\n```livecodeserver\n百度地图密钥申请 https://lbsyun.baidu.com 创建应用需选择为"浏览器端"\n```'

    },
    {
        id: "68ddd844-d023-405e-94e2-88e55a135eb3",
        message: "GEOJSON热力图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**地图** 设置**roam**为true 选择一个GeoJson地图文件（可以在网上下载需要的GeoJson地图文件）\n\n- [x] 在图表**内部属性**添加**连续映射** 设置**min**和**max**\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中分别拖2列数据到**地图**属性上 第一列是地理位置名称 例如昆明市 第二列是数据'
    },
    {
        id: "472a9015-c131-44b5-88da-41dbc172a1c2",
        message: "GEOJSON散点图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**地理坐标系** 设置**roam**为true 选择一个GeoJson地图文件（可以在网上下载需要的GeoJson地图文件）\n\n- [x] 在图表**内部属性**添加**散点图/涟漪散点**\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中分别拖2列数据到**散点图/涟漪散点**属性上 第一列是经度 第二列是纬度 如需要设置散点大小颜色或名称等需要拖入更多数据并用函数来实现\n\n```livecodeserver\n"如果发现属性的字段拖入错误，可以按住ctrl键，并重新拖一个字段，该属性上的所有字段都被清空并添加刚刚拖入的字段。"\n"数据属性一般可以拖入多个字段，可以用来表达不同用途，比如设置波形图符号大小随数据大小呈现不同大小。"\n例如设置散点图符号大小函数声明：\n(value: Array|number, params: Object) => number|Array\n我们关注最多的是value参数，我们拖入一列数据value就是当前的数据，\n我们拖入多列数据可以用value[0]、value[1]...value[n]分别对应拖入顺序的数据列，函数可以如下：\n(value)=>value\n或是\n(value)=>value[n] //n=[0,n] 拖入多列数据取n列\n也可以进行一定运算如下，函数名可以随意,也可以用上面这样的匿名函数\nfunction nni(value){\n\treturn value[n]*5\n}\n```\n\n'

    },
    {
        id: "f37687e4-c20e-4bfc-8c5d-5c659154118d",
        message: "GEOJSON路径图 双击查看完整帮助",
        help:'1. ### 模板制作\n\n- [x] 从左侧图表组件内拖出**空白图表**到画布上\n\n- [x] 在图表**内部属性**添加**地理坐标系** 设置**roam**为true 选择一个GeoJson地图文件（可以在网上下载需要的GeoJson地图文件）\n\n- [x] 在图表**内部属性**添加**路径图**\n\n```livecodeserver\n"以上操作本模板已帮您完成"\n```\n\n2. ### 添加数据\n\n- [ ] 从右上角的数据中分别拖2列数据到**路径图**属性上 第一列是经度 第二列是纬度 如需显示名称需要拖入第三列名称数据 如果是多条路径需要把数据中的经度和纬度设置为0用来分割'
    },
]