

## 配置 goLand快捷键总结


vscode配置的很多自定义快捷键，不想在 goLand 记另一套了，这里 也做个总结
参考官方文档 [idea 官网](https://www.jetbrains.com/help/idea/using-code-editor.html#editor_basic_usage)

1. ctrl e快捷键：
   - [x] `ctrl e .`  , 历史记录中打开上一个编辑器  【view > recent files 】
   - [x] `ctrl e [`, 左边 file列表可见性 【view > toolwindow > alt 1 】 
   - [x] `ctrl e ]` ,大纲可见性
   - [ ] `ctrl e ,hjkl` ,焦点移动 【无法完全解决】，这里用 switcher来代替，只能用 ctrl + tab , 
   - [x] `ctrl e m` 当前面板最大化，切换面板大小 【maximum toolwindow 】
   - [x] `ctrl e e` switer 面板 ,vscode 对应 ctrl  q

2. ctrl t 终端快捷键：
   - [x] `ctrl t t` ,切换终端【打开终端】alt f12
   - [x] `ctrl t M` 终端最大化
   - [ ] `ctrl t n` 下个终端 ,可以 ctrl  [] 代替
   - [ ] `ctrl t ([ or ])` 上一个或者下一个终端 
   - [ ] `ctrl t .` 拆分终端
   - [ ] `ctrl j ]` ，返回编辑器焦点 【终端和编辑器切换焦点】
   - [ ] `ctrl t g` 下个终端【下一组终端】

3. ctrl j 代码编辑快捷键

   - [x] `ctrl j j` 进入函数定义
   - [x] `ctrl j h` 触发参数提示
   - [x] `ctrl j k` 转到引用 find usage ,alt f7
   - [x] `ctrl j l` 转到实现
   - [ ] `ctrl j .` 触发面板可见性
   - [ ] `ctrl j u` 上一个标签面板
   - [ ] `ctrl j i` 下一个标签面板
   - [x] `ctrl j n` 下一个错误或者警告
   - [x] `ctrl j r` 运行代码
   - [ ] `ctrl j t` 在光标处运行测试
   - [ ] `ctrl j space` quick fix
   - [ ] `ctrl j f` 查找符号
   - [ ] `ctrl j ctrl f` 全文查找
   - [ ] `ctrl j ;` 跳到指定行

5. bookmarks
   - [x] ctrl b b  toggle bookmark
   - [x] ctrl b .   list all book marks\
   - [x] ctrl b d ,  打断点，debug break point

7. debug 快捷键
   - [x] ctrl numpad6  , step over
   - [x] ctrl numpad5 , run debug
   - [x] ctrl numpad8 ,step out
   - [x] ctrl numpad2, step into
   - [x] ctrl end    , end debug
   - [x] ctrl numpad4  , 重新debug


8. live template
   - [x] alt / , basic 代码提示




9.  同步goLand配置，方便复用


   - 点击 file,manager ide settings
   - export jar


10. 其他辅助
   - [ ] alt 1,alt 2, 跳转标签页



11. 安装vim插件， 尽量不使用鼠标操作，全部用快捷键操作。











 


