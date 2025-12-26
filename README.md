# dbx

基于 [sqlx](https://github.com/jmoiron/sqlx) 的半自动化 ORM 代码生成工具，灵感来自 MyBatis Generator。

## 什么是半自动化 ORM？

与 GORM 等全自动 ORM 不同，dbx 采用"半自动化"理念：

| 特性 | 全自动 ORM (GORM) | 半自动化 ORM (dbx) |
|------|------------------|-------------------|
| SQL 生成 | 框架自动生成，开发者不可见 | 开发者手写 SQL，完全可控 |
| 复杂查询 | 链式调用，学习成本高 | 原生 SQL，无学习成本 |
| 性能优化 | 难以优化生成的 SQL | 直接优化手写 SQL |
| 调试难度 | 需要开启日志查看生成的 SQL | SQL 就在代码里，一目了然 |

### 半自动化的含义

**自动化的部分（机器做）：**
- 根据数据库表结构自动生成 Model 结构体
- 自动生成单表的基础 CRUD 方法（Insert/Update/Delete/FindByID/FindByIds 等）
- 表结构变更后重新生成即可同步

**手动的部分（人来做）：**
- 复杂的业务查询 SQL 由开发者手写
- 多表关联、聚合统计等复杂场景完全可控
- SQL 写在代码里，便于 review 和优化

### 为什么选择半自动化？

```go
// 全自动 ORM：你需要学习框架的 DSL
db.Where("age > ?", 18).
   Joins("LEFT JOIN orders ON orders.user_id = users.id").
   Group("users.id").
   Having("count(orders.id) > ?", 5).
   Find(&users)

// 半自动化 ORM：直接写 SQL，所见即所得
userDAO.QueryRowsCtx(ctx, &users, `
    SELECT u.* FROM users u
    LEFT JOIN orders o ON o.user_id = u.id
    WHERE u.age > ?
    GROUP BY u.id
    HAVING count(o.id) > ?
`, 18, 5)
```

**优势：**
- SQL 能力不受框架限制，DBA 可以直接优化
- 没有 ORM 的 N+1 问题、预加载陷阱
- 代码即文档，SQL 逻辑清晰可见
- 单表 CRUD 自动生成，减少重复代码

## 安装

```bash
# 安装 CLI 工具
go install github.com/fanqingxuan/dbx/cmd/dbx@latest

# 项目中引入库
go get github.com/fanqingxuan/dbx
```

## 快速开始

### 1. 生成代码

在项目根目录（有 go.mod 的目录）执行：

```bash
dbx -dsn "user:pass@tcp(localhost:3306)/dbname"
```

会自动：
- 读取 go.mod 获取包名
- 生成数据库所有表的 model 和 dao

指定输出目录和表：
```bash
dbx -dsn "user:pass@tcp(localhost:3306)/dbname" -o ./internal user order
```

### 2. 使用配置文件

创建 `dbx.yml`：

```yaml
dsn: user:pass@tcp(localhost:3306)/dbname
output: ./internal
tables: user, order
```

然后执行 `dbx` 即可。

### 3. 生成的目录结构

```
├── model/
│   ├── user.go           # 表结构体
│   └── order.go
├── dao/
│   ├── gen/
│   │   ├── user.go       # 自动生成的 CRUD（可重新生成覆盖）
│   │   └── order.go
│   ├── user.go           # 业务 DAO（不会覆盖，可自定义扩展）
│   └── order.go
```

## 使用示例

### 初始化

```go
import (
    "github.com/fanqingxuan/dbx/pkg/dbx"
    _ "github.com/go-sql-driver/mysql"

    "yourproject/dao"
)

func main() {
    db, err := dbx.Open("mysql", "user:pass@tcp(localhost:3306)/dbname")
    if err != nil {
        panic(err)
    }

    userDAO := dao.NewUserDAO(db)
}
```

### 基础 CRUD

```go
ctx := context.Background()

// 插入
user := &model.User{Name: "test", Age: 18}
err := userDAO.Insert(ctx, user)

// 根据 ID 查询
user, err := userDAO.FindByID(ctx, 1)

// 根据多个 ID 查询
users, err := userDAO.FindByIds(ctx, []int64{1, 2, 3})

// 更新
user.Name = "updated"
err = userDAO.Update(ctx, user)

// 删除
err = userDAO.Delete(ctx, 1)

// 批量删除
affected, err := userDAO.DeleteByIds(ctx, []int64{1, 2, 3})

// 批量更新
affected, err := userDAO.UpdateByIds(ctx, []int64{1, 2, 3}, map[string]any{
    "status": 1,
    "name":   "batch",
})
```

### 自定义查询

```go
// 查询多行
var users []model.User
err := userDAO.QueryRowsCtx(ctx, &users, "SELECT * FROM user WHERE age > ?", 18)

// 查询单行
var user model.User
err := userDAO.QueryRowCtx(ctx, &user, "SELECT * FROM user WHERE id = ?", 1)

// 查询单个值
var count int
err := userDAO.QueryValueCtx(ctx, &count, "SELECT COUNT(*) FROM user")

// 查询单列
var names []string
err := userDAO.QueryColumnCtx(ctx, &names, "SELECT name FROM user")

// 执行语句
affected, err := userDAO.ExecCtx(ctx, "UPDATE user SET status = ? WHERE age < ?", 0, 18)
```

### 扩展 DAO

在 `dao/user.go` 中添加自定义方法（此文件不会被覆盖）：

```go
package dao

import (
    "context"
    "yourproject/model"
)

func (d *UserDAO) FindByName(ctx context.Context, name string) (*model.User, error) {
    var user model.User
    err := d.QueryRowCtx(ctx, &user, "SELECT * FROM user WHERE name = ?", name)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (d *UserDAO) FindActiveUsers(ctx context.Context) ([]model.User, error) {
    var users []model.User
    err := d.QueryRowsCtx(ctx, &users, "SELECT * FROM user WHERE status = 1")
    return users, err
}
```

## 日志

实现 Logger 接口即可：

```go
type MyLogger struct{}

func (MyLogger) Debug(format string, args ...any) { log.Printf("[DEBUG] "+format, args...) }
func (MyLogger) Info(format string, args ...any)  { log.Printf("[INFO] "+format, args...) }
func (MyLogger) Warn(format string, args ...any)  { log.Printf("[WARN] "+format, args...) }
func (MyLogger) Error(format string, args ...any) { log.Printf("[ERROR] "+format, args...) }

db.SetLogger(MyLogger{})
```

## CLI 参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| -dsn | MySQL 连接字符串 | - |
| -o | 输出目录 | . |
| -c | 配置文件路径 | dbx.yml |

## 依赖

- github.com/jmoiron/sqlx
- github.com/go-sql-driver/mysql
