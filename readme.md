# 目录结构说明

## cmd 在此包内组织和调用pkg包内代码，没有其他包依赖此包。

    - app 科脉小云后台管理

    - report 云鼎报表系统

## pkg 应用包，定义业务模型和接口

    - config 系统设置

    - http 实现http服务

    - jwt 实现json web token 验证

    - mock 实现数据模拟

    - mssql 实现mssql服务

    - mysql 实现mysql服务

    - util 实现工具类

## 第三方包依赖

    - [sqlx](go get github.com/jmoiron/sqlx)
    - [sqlmock](go get github.com/DATA-DOG/go-sqlmock)
    - [logrus](go get github.com/sirupsen/logrus)
    - [go-jwt](go get github.com/dgrijalva/jwt-go)
    - [httprouter](go get github.com/julienschmidt/httprouter)
    - [go-mssqldb](go get github.com/denisenkom/go-mssqldb)
    - [go-mysql](go get github.com/go-sql-driver/mysql)
    - [gomigrate](go get github.com/DavidHuie/gomigrate)
    - [uuid](go get github.com/nu7hatch/gouuid)
