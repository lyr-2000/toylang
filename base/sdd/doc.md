---
title: "doc"
date: 2022-08-06T22:46:00+08:00
lastmod: 2022-05-11T20:02:48+08:00
categories: ["sdd"]
tags: ["sdd"]
author: "lyr"
draft: false
images: ["https://api.mtyqx.cn/api/random.php"]

---


## 构造三地址代码(中间代码)

这一步暂时没法实现，看了github很多的项目，都是之前做了个前端编译器，如果要做 Jit虚拟机的话，可能要学一些 汇编的代码，先放弃了。


`5*4 + 3*2 `

转化为 

```js
p0 = 5 * 3
p1 = 2 * 4
p2 = p0 * p1
```


## 符号表



```js
var a  =1
b = a+1

{
    var a = 10
    var b = a+1
}


```

符号表的设计


1. 静态符号表 
   - 存储常量
2.  符号表 (实现词法作用域)


```go
type Symbol struct {
    Address string //变量
    Label string //标签
    ConstNum string //常数

}

type SymbolTable struct {
    Symbols []*Symbol
    NextTabs []*SymbolTable
}

```


```js
var a = 0; //offset = 0
var b = 1 //offset=1
{
    c = b+1 //0，当前符号表找不到b，就在父符号找变量b,直到找到为止
    d = c+1 //1
}

//模拟压栈过程

```














