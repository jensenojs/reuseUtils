
# README

对于一个对象如果想要使用reuse package做内存复用, 那么其需要满足几个条件

- 实现`ReusableObject`接口 : `TypeName`
- 实现对应的`new`, `reset`方法
- 在init函数中使用`CreatePool`进行注册

但问题是`parser-tree`或者`plan-tree`下有上百种类的对象需要做内存复用, 手动编写无异于愚公移山, 本脚本用于批量生成reuse做Parser内存复用的代码, 后续修改一下之后可以拿去给`Plan-tree`或者`Scope-tree`用, 主要需要修改的应该是`Reset`或者说`reset`相关的方法.

```go
func CreatePool[T ReusableObject](
	new func() *T,
	reset func(*T),
	opts *Options[T]) {

	}
```

关于像ast这样的树形结构想要使用`reuse`进行内存复用有需要递归Alloc和递归Free的场景, 递归`Alloc`的逻辑内嵌在了`mysql_sql.go`中, 因此`new`方法不需要特别讨论, 对于少部分的结构体(比如说里面包含了`map`, 或者原子变量的, 再讨论就行了), 但对于递归释放而言, 情况会更复杂一些, 具体看下面的例子.

使用本脚本前需要在main.go中设定需要批量生成的代码的路径, 以`alter.go`举例
```go
// 设定目录
srcFolderPath = "/Users/jensen/Projects/matrixorigin/matrixone/pkg/sql/parsers/tree/alter.go" 
```

假设该文件下的结构体有一个`AlterUser`, 多个结构体时也是类似的
```go
type AlterUser struct {
	statementImpl
	IfExists bool
	Users    []*User
	Role     *Role
	MiscOpt  UserMiscOption
	// comment or attribute
	CommentOrAttribute AccountCommentOrAttribute
}
```

在使用时需要指定`genType`, 一共有三个方法可供选择, 所有生成的代码都默认放在本目录下的`generate/`文件夹下, 然后可以再根据具体业务逻辑修改, 这里需要重点修改的可能是`GetReset`方法

- GenCreatePool
    批量生成`func init() {}`所需内容到对应文件, 示例如下
    ```go
	reuse.CreatePool[AlterUser](
		func() *AlterUser { return &AlterUser{} },
		func( a *AlterUser) { a.reset() },
		reuse.DefaultOptions[AlterUser]().
			WithEnableChecker())
    ```
	
- GenTypeName 
    
    批量生成reuse.ReusableObject所需接口到对应文件, 示例如下

	在使用前记得修改`gen_name_method`下的`packageName`字段
    ```go
    // generate/alter.go_
    func (node AlterUser) TypeName() string { return "tree.AlterUser" }
    ```
- GenReset

	这种情况比较复杂, 编写本文档时还没有实现相应脚本, 但Reset脚本大差不差应该如下所示

    ```go
    // generate/alter.go_
    func (node AlterUser) reset() { 
		...
		for _, u := range node.Users {
			u.Free[User](nil)
		}
		...
	}
    ```

	在`AlterUser.reset` 中需要调用`User.Free`, 是因为`User`也是`ast`中的一个节点, 且`User`也(会)做`reuse`的内存复用, 且在该结构体中`User`是指针, 不是结构体
	
	这样在`alter user`的SQL执行完之后, 需要释放相应地`stmt`时, 只需要调用`AlterUser.Free`即可, 这样就会从ast的头部开始进行递归地`Free`了

	```go
	func (node *AlterUser) Free() {
		reuse.Free[AlterUser](node, nil)
	}
	```

	如何自动生成相应地`reset`方法是比较麻烦的, 在`parser`下可以参考节点对应的`Format`方法进行改写, 逻辑会复杂一些, 目前已经初步完成了对Parser的重写


批量生成复制进去之后, 会有结构体和其对应的方法不在一起的问题, 可能可以通过一个format方法来加以修复, 它的主要作用就是将结构体以及其相关的方法聚集起来. Format方法的编写也比较麻烦, 还是有一些bug, 特别是有一些注释的时候, 但总体来说还在能接受的范围内.