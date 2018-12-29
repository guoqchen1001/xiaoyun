# mysql数据库表结构迁移

## 升级文件规则

1. 文件命名规则`{id}_{name}_{upOrdown}.sql`,其中`id`为版本号，`name`为文件名，`up`为升级版本，`down`为回滚该版本

2. 同一个文件的`down`和`up`必须成对出现，`id`和`name`相同时，认为该文件为一对，即提供了升级数据库结构必须提供提供降级语句

3. 同一个版本允许多个升级文件，只要版本号相同即可

4. 一个文件内语句必须以分号结尾，允许多个语句存在，多个语句以分号隔开

## 运行机制

1. 程序启动时，会检查配置文件中`mysql`连接，连接不成功则忽略步骤

2. 连接`mysql`成功后，会自动运行`mysql`数据库中`client`对象的`migrate`方法，进行数据库结构迁移

3. 程序会比对`migrations`目录中的版本在数据库中是否升级，若没升级，则自动按版本号执行语句。每个版本单独事务，其中一个失败，则后续不再执行

4. 数据库迁移日志记录在`app.log`中

5. 若需要回滚数据库，需要调用`migrate`对象的`down`方法，数据库会被回滚到最近一次迁移的前一个版本。

## 使用方法

1. 若开发过程中对数据库表结构有所修改，请在本地开发完毕后，将升级文件按照版本放置到此目录。然后提交版本控制。