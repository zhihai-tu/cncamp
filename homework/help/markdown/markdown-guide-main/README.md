# markdown-gride

## 概述

### 1. 什么是Markdown

Markdown是以**易于阅读、创作和编辑**文档为目标的**轻量级标记语言**，是一种创作格式，最终以纯文本形式发布。

Markdown由**Aaron Swartz**和**John Gruber**共同设计。

### 2. 为什么是Markdown

Markdown视**可读性**为最高准则，**易于阅读、轻松创作、灵活编辑**，是文档编辑*极简主义*的代表。

主要优势有：
* 语法集非常小，完全由标点符号组成
* 兼容HTML，可用于创作web文档
* 轻量级，跨平台


## 常用语法规则
### 1. 标题
主要有两种形式：

* 使用`=`和`-`分别标记一级、二级标题
  #### 语法：
  > 一级标题   
  > `=====`   
  > 二级标题    
  > `-----`

  #### 效果：
  > 一级标题
  > ======
  > 二级标题
  > ----- 

* 使用`#`标记1~6级标题
  #### 语法：
  > \# 一级标题    
  > \## 二级标题  
  > \### 三级标题  
  > \#### 四级标题    
  > \##### 五级标题   
  > \###### 六级标题

  **Tips**: 建议在`#`后加一个空格

  #### 效果：
  > # 一级标题
  > ## 二级标题
  > ### 三级标题
  > #### 四级标题
  > ##### 五级标题
  > ###### 六级标题



### 2. 段落和换行
换行：只需在行末加两个空格键和一个回车键即可换行。

分段：段落之间空一行即可。

### 3. 块引用
#### 语法：
只需在段落首行前加上一个`>`，嵌套引用只需要加额外层级的`>`,如：
> \> 块引用   
> \>> 嵌套引用

#### 效果：
> 区块引用  
> 嵌套引用

### 4. 列表
Markdown支持*无序列表*和*有序列表*。

* 无序列表使用`*`,`-`,`+`标记
  #### 语法
  > \* (-\+) 第一项    
  > \* (-\+) 第二项   
  > \* (-\+) 第三项   

  #### 效果
  > - 第一项
  > - 第二项
  > - 第三项

  **Tips**: 标记后面最小加一个空格或制表符。

* 有序列表使用数字辅以`.`
  #### 语法:
  > 1. 第一项  
  > 2. 第二项  
  > 3. 第三项  

  #### 效果：
  > 1. 第一项  
  > 2. 第二项  
  > 3. 第三项

### 5. 代码块
插入代码块, 只需要将每一行都缩进 4 个空格或者 1 个水平制表符。

普通段落：  
public static void main(String[] args) {    
    System.out.println("Hello, Markdown");   
}

代码区块：

    public static void main(String[] args) {
        System.out.println("Hello, Markdown);
    }

也可使用如下方式格式化代码：
#### 语法：
> \```java  
>   This is a code block.    
> \```

**Tips**: \```后可添加编程语言类型，如java, c, c++, python, javascript等，将标识出语言关键字。
#### 效果：
> ```java
> public static void main(String[] args) {
>     System.out.println("Hello, Markdown");
> }
> ```

### 6. 分割线
可在一行中添加三个以上的`*`, `-`, 或者`_`添加分割线。
#### 语法：
> \***   
> \---   
> \___   

#### 效果：
> 以下是内容一、内容二分割线   
> ______
> 内容二

### 7. 文本样式
常用的文本样式主要有强调、倾斜、删除线和底纹。

#### 语法：
> \*\*加粗\*\*   
> \*斜体\*   
> \~\~删除线\~\~   
> \`底纹\`   

#### 效果：
> **加粗**   
> *斜体*   
> ~~删除线~~   
> `底纹`   

### 8. 链接
主要链接形式：*内联*、*引用*和自动链接。

* 内联
  #### 语法：
  > This is \[an example\]\(https://example.com\) inline link.    
  > \[szhxiao的Markdown手册\]\(https://github.com/szhxiao/markdown-guide\)

  #### 效果
  > This is [an example](http://example.com) inline link.    
  > [szhxiao的Markdown手册](https://github.com/szhxiao/markdown-guide)

* 引用
  #### 语法：
  > \[szhxiao_markdown-guide\]\[1\]    
  > \[szhxiao_self-code-lab]\[2\]    
  >\[1\]: https://github.com/szhxiao/markdown-guide   
  >\[2\]: https://github.com/szhxiao/self-code-lab

  #### 效果：
  > [szhxiao_markdown-guide][1]   
  > [szhxiao_self-code-lab][2]  

  > [1]: https://github.com/szhxiao/markdown-guide   
  > [2]: https://github.com/szhxiao/self-code-lab

* 自动链接
  
  自动创建URL和email地址链接的简短形式

  #### 语法：
  > \<http:// example.com\/\>   
  > \<http://https://github.com\/\>

  #### 效果：
  > <http://https://github.com/>

### 9. 图片
使用了类似链接的语法来插入图片, 有*内联*和*引用*两种形式。

* 内联
  #### 语法：
  > \!\[Alt text\]\(\/path\/to\/img.jpg\)   
  > \!\[阅读\]\(https://github.com/szhxiao/markdown-guide/blob/main/resources/beverage-book.jpg\)

  #### 效果：
  > ![Alt text](/path/to/img.jpg)   
  > ![阅读](/resources/"beverage-book.jpg")

* 引用
  #### 语法：
  > \!\[Alt text\]\[id\]
  > \!\[阅读\]\[beverage-book\]
  
  #### 效果：
  > [阅读][beverage-book]

  >[beverage-book]: https://images.gitee.com/uploads/images/2021/0425/095708_63965da9_1494433.jpeg "beverage-book.jpg"



### 10. 表格
制作表格使用 | 来分隔不同的单元格，使用 - 来分隔表头和其他行。

#### 语法：
> \| 表头 \| 表头 \|   
> \| ---- | ----- |   
> \| 单元格 \| 单元格 \|   
> \| 单元格 \| 单元格 \|

#### 效果：
> | 序号 | 语法 | 表格对齐方式 |
> | :---- | :----: | ----: |
> | 1 | :- | 左对齐 |
> | 2 | :-: | 居中对齐 |
> | 3 | -: | 右对齐 |

