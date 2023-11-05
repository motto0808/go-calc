# 表达式解析器
这是一个简单的表达式求值器，语法简单，支持最常规的整数四则运算和比较运算、以及类似C语言的逻辑运算。
## 如何使用
可以查看main.go文件
### 一步求值
```go
env := calc.Env{}
e := calc.NewEvaluator()
e.Eval("1+2>=1", env) //return 1
```

### 两步求值(先解析后求值)
```go
env := calc.Env{}
p := calc.NewParser()
stmts := p.Parse("1+2>=1")
e := calc.NewEvaluator()
e.EvaluateStmt(stmts[0], env) //return 1, nil
```

## 如何编写表达式
可以查看sample.calc文件以及unittest目录下的测试用例


### 支持的普通运算符
```js
1+1
1-1
1*2
4/2
5%2
a>b
a>=b
a<b
a<=b
a==b
a!=b
!foo
foo&&bar
foo||bar
a>0?a:0
```
### **in**关键字

用于判断数组中是否包含指定的值
```js
serverId in [1,2,3]
!(serverId in [1,2,3])
```

### 错误的用法
目前对于比较运算符禁止连续比较(为了避免不必要的错误)
例如下面的表达式会导致panic
```js
a<b<c
a<b>c
```
如果想判断在数学上a<b且b<c,应该写成
```js
a<b && b<c
```
如果确实想使用a<b的结果与c进行比较(比较运算符只会返回0或1),应该写成
```js
(a<b)<c
```


## 扩展
### 使用环境变量
可以将一些变量传入求值器
```go
env := calc.Env{"a": 123}
e := calc.NewEvaluator()
e.Eval("a>=100", env) //return 1
```

### 定义变量进行简单运算
例如，对如下文本进行求值将得到11
```javascript
var a=1
var b=a>10?a:10
a+b
```

### 使用外部条件求值
很多时候判断条件可能需要结合很多其他的信息进行判断，比如用户的等级、充值金额等，这时可以通过实现ICondHelper接口来实现

这样你只需要实现很多外部条件就可以直接在表达式中使用他们，而不需要把他们一一放入env中

表达式求值时优先在局部环境中查找变量，如果在局部环境中找不到变量，会尝试通过条件求值

具体使用例子可以参考unittest/cond_test.go
```go
func TestCondition(t *testing.T) {
	eva := calc.NewEvaluator()
	eva.SetCondHelper(&condHelper, nil)
	var v int = 0
	v = eva.Eval("charge>=200 && age<=30", calc.Env{})
	assert(t, v == 1)
}
```
## 如何编译
先安装goyacc
```
go install golang.org/x/tools/cmd/goyacc@latest
```
执行make
```
make all
```