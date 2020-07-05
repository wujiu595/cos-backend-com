# cos-backend-com

## 项目简介

该项目基于依赖注入的微服务项目，项目间使用rpc通信依据功能主要划分为 account、cores、notification、eth 4个服务

- account 提供用户相关的功能：用户信息、登陆鉴权等功能
- cores 项目的核心模块，提供startup、startupsetting、iro、bounty、hunter等功能
- notification 通知服务，主要提供通知功能
- eth 提供和以太坊交互的功能

## 开发环境准备

- 数据库：postgresql 10.3

- redis
- 文件存储minio
- go 1.10及以上

提供了docker-compose快速搭建开发环境

```
docker-compose -f docker-compose-devel.yaml up -d
```

## 目录介绍

```
├── README.md
├── bin
│   ├── strip
│   └── web3
├── docker-compose-devel.yaml
├── go.mod
├── go.sum
├── hack
│   ├── build （build 编译好的二进制文件目录）
│   ├── conf （项目的配置文件目录，以名称对应到服务，其中以 .dev. 为本地环境测试的配置环境）
│   ├── docker （提供把二进制打包成docker的及docker相关的文件目录）
│   ├── files （主要提供项目运行所需要的一些文件初始化目录）
│   ├── gnorm （生成代码的模版文件）
│   ├── migrate（数据库升级脚本目录）
│   │   └── 20200501_0_example
│   │       └── upgrade.sql
│   └── run （提供项目构建测试等脚本文件）
└── src
    ├── template
    │   ├── cmd
    │   │   └── template
    │   │       ├── app
    │   │       │   └── app.go （路由、应用配置、初始化环境等）
    │   │       └── main.go （入口文件）
    │   ├── env.go （导入的配置文件映射到 struct后到文件）
    │   ├── proto （不属于请求和响应结构的原型）
    │   └── routers
    │       └── template
    │           └── template.go （handle 文件）
    ├── common （一些通用的模块）
    ├── libs
    │   ├── apierror（自定义错误码）
    │   ├── auth （登陆鉴权）
    │   ├── filters （filters）
    │   ├── models（model 文件）
    └── └── sdk （对外暴露的proto、服务等）
```



## start project

- 环境变量配置

推荐使用goland开发

- 初始化环境

```shell
./hack/run database init
```

会执行数据库的启动，测试数据的导入等

- 启动指定服务

```
./hack/run start account
```

- 同步数据库结构到 `files database.sql`

```
./hack/run database dump.schema
```

默认会同步本地数据结构到files下的database.sql下，需要同步远程数据库结构时可以通过 `DATABASE_DUMP_HOST` 同步指定环境的数据库结构，在同步本地数据库结构的修改时应在migrate目录下新建文件目录并添加数据库升级脚本



## API规范

项目使用yapi做为api的管理工具，在api的设计上应该满足restful和cloud api的设计规范，并在yapi中做好单元测试

### 面向资源设计

2.1 集合与资源与属性

集合是具有相同属性的资源

举例： 书、某一本书、作者依次对应集合、资源、属性

#### 分类

- 命名：**必须**出现在资源定义的开头,作为某一类的接口的界定不能作为URI单独使用

- 格式：

  `[分类]/[集合/资源Id]...../集合`

- 注意：分类不一定出现在每个URI中，如果相同资源在不同功能中不同的展示，对于书而言，作者看到的内容与读者看到的内容不是完全相同,对书进行的操作也不完全相同，因此可把接口可分类为 reader和writer

- 举例：

  `/writer/books``/reader/books`

#### 集合

- 命名：**必须**使用有意义的`名词的复数形式`，**不得**使用values、items等

- 格式：

  `[分类]/[集合/资源Id]...../集合`

- 方法：`List集合`、`Create资源`

- 举例：

  `/writer/books`

#### 资源

- 命名：**必须**使用有意义的`名词的单数形式+Id`构成所指向的资源

- 格式

  `[分类]/[集合/资源Id]...../集合/资源Id`

- 方法：`Get资源`、`Delete资源`、`Update资源`

- 举例：

  `/writer/books/{bookId}`

#### 资源属性

- 命名：**必须**是名词的单数形式并属于当前资源的属性，如果属性为指向其他资源的id，那么属性的命名应该使用其他资源名称作为属性的获取内容,并且属性`必须`出现在URI的末尾

- 格式：

  `[分类]/[/集合/资源Id]...../集合/资源Id/属性`

- 方法：`Get资源属性`、`Update资源属性`

- 举例：

  `/reader/books/{bookId}/auther`

- 注意：在批量获取属性时需使用自定义方法

  `GET /reader/books:batchGetAuthers`

#### 资源属性的统计

- 命名：必须是名词的单数形式，统计内容属于资源属性，简单的方法

- 举例

  `/reader/books/{bookId}/favoriteCount`

- 方法：Get属性统计

### 方法

#### 标准方法

| 标准方法               | HTTP 映射 | HTTP 请求正文 | HTTP 响应正文 |
| :--------------------- | :-------- | :------------ | :------------ |
| List                   | GET       | 无            | 资源列表      |
| Get                    | GET       | 无            | 资源          |
| Create                 | POST      | 资源          | 资源          |
| Update                 | PUT       | 资源          | 资源          |
| Delete                 | DELETE    | 无            | 无            |
| Get*属性***            | GET       | 无            | 属性          |
| Get***属性统计***      | GET       | 无            | 统计结果      |
| Update***`资源`属性*** | PUT       | 属性          | 属性          |

注意：List(result item)、Get、Create、Update、Delete（标记）应返回相同的资源内容

### 自定义方法

自定义方法是指 8 个标准方法之外的 API 方法。这些方法**应该**仅用于标准方法不易表达的功能。通常情况下，API 设计者**应该**尽可能优先考虑使用标准方法，而不是自定义方法。

自定以方法的规则：

- 自定义方法**应该**使用 HTTP `POST` 动词，因为该动词具有最灵活的语义，但作为替代 get 或 list 的方法（如有可能，**可以**使用 `GET`）除外
- 请注意，使用 HTTP `GET` 的自定义方法**必须**具有幂等性并且无负面影响。例如，在资源上实现特殊视图的自定义方法**应该**使用 HTTP `GET`。
- 网址路径**必须**以包含冒号（后跟自定义动词）的后缀结尾。

常用的自定义方法：

| 方法名称            | 自定义动词  | HTTP 动词 | 备注                                                         |
| :------------------ | :---------- | :-------- | :----------------------------------------------------------- |
| 取消                | `:cancel`   | `POST`    | 取消一个未完成的操作（构建、计算等）                         |
| BatchGet <复数名词> | `:batchGet` | `GET`     | 批量获取多个资源。（详情请参阅[列表描述](https://cloud.google.com/apis/design/standard_methods#list)） |
| 移动                | `:move`     | `POST`    | 将资源从一个父级移动到另一个父级。                           |
| 搜索                | `:search`   | `GET`     | List 的替代方法，用于获取不符合 List 语义的数据。            |
| 逻辑删除            | `:delete`   | `POST`    | 逻辑删除资源，`并返回和Get方法相同的输出`                    |
| 恢复删除            | `:undelete` | `POST`    | 恢复之前删除的资源，`并返回和Get方法相同的输出`              |



### YApi 接口管理

#### 接口命名

接口名称是以分类，模块名称，对应方法组成

- 格式

  `[分类-]模块-方法`

- 方法和接口名称的对应

  | 标准方法     | 名称                           |
  | :----------- | :----------------------------- |
  | 标准方法     | 名称                           |
  | List         | [分类-]模块-列表               |
  | Get          | [分类-]模块-获取               |
  | Create       | [分类-]模块-创建               |
  | Update       | [分类-]模块-更新               |
  | Delete       | [分类-]模块-删除               |
  | Get`属性`    | [分类-]模块-获取`属性名`       |
  | Get属性统计  | [分类-]模块-属性统计           |
  | Update`属性` | [分类-]模块-修改`属性名`       |
  | 自定义方法   | [分类-]模块-方法实现功能的描述 |

### 接口定义

- 根据需求只返回前端需要展示和需要计算的字段
- 返回数据隐藏关键信息字段如：accessKeyId等
- 在返回的数据中，只要字段的tag中没有**‘validate:omitempty’**对应到yapi中都为必须字段



## 代码规范

### 路由

- 项目下的路由器都适用namespace的形式 在app.go下可以查看到代码示例
- 示例：

```go
func (p *appConfig) ConfigRoutes() {
	p.Routers(util.VersionRouter())
	p.Routers(
		s.Router("/user",
			s.Post(users.Users{}).Action("Get"),
		)
	)
}
```

通过请求 `/user`路径会执行`users.Users{}`的`Get`方法

### handle

- handle方法
- 示例：

```go
type Users struct {
	routers.Base
}

func (h *Users) Get() (res interface{}) {
	var input account.LoginInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.ErrBadRequest.WithData(err)
		return
	}

	var user account.UserResult
	if err := usermodels.Users.Get(h.Ctx, h.Uid, &user); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(user)
	return
}
```

`Users` struct继承 `routers.Base`

`Get`方法返回res 接口类型

### model

- models以功能划分例如`User`、`startups`
- 示例：

```go
var Users = &users{
	Connector: models.DefaultConnector,
}

type users struct {
	dbconn.Connector
}

func (p *users) Get(ctx context.Context, id flake.ID, output interface{}) (err error) {
	stmt := `
		SELECT *
		FROM users
		WHERE id = ${id};
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}": id,
	})

	return p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
}
```

`Users` 包涵 `Connector: models.DefaultConnector`并使用单例模式在项目运行时初始化



### 注释

- 命名及注释，应明确函数的功能，使用 动词+名词的形式命名函数，在命名上见名知意

- 对于复杂函数在命名上无法达到准确的含义时应标明注释描述函数功能

  ```go
  /*
  
  */
  ```

- todo

  对于暂时遗留未实现的内容应该使用todo进行注释

  ```go
  //todo
  ```

  

### 代码格式化

在goland 中应开启 file watch 进行代码格式化，值得注意的是在导入包时应去掉包与包之间的换行

```go
import (
	"context"
	accountEnv "cos-backend-com/src/account"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/sdk/account"
	"database/sql"
	"fmt"
	"math/rand"
	
	"time"

	"github.com/wujiu2020/strip/utils"
)
```

