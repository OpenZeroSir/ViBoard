options = {
    "key":"e289ba5c-b692-4fd1-bcbd-7097e4734d48", #请登陆网页获取key，本key是表key，如果是组件就换成组件ID 例如：i0IV46AKSeIHhwf3bFeLF 
    "data":json.dumps([["张三",69],["李四",88]]) , #原表格式是"姓名,分数"两列，如果更新的是组件，要json.dumps(["新内容"])，这个数组只有一个元素就是组件新内容。
    "replace":"true" #这个参数只有在覆盖更新的时候有效，增量更新或更新组件时无效
}
response = requests.post("http://127.0.0.1:50520/ai", data=options,headers={"token": "355458543157374159565337324e4547"}) #正常的话返回结果是"ok"

