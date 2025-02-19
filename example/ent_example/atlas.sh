#!/bin/sh

## Schema 命令

# 查看 Ent Schema 对应的 SQL 结构，输出到控制台
#atlas schema inspect \
#  -u "ent://ent/schema" \
#  --format '{{ sql . "  " }}' \
#  --dev-url "mysql://root:123456@:3306/ent_test"

# 将声明式 Schema 直接应用到数据库（一般开发环境使用）
# --dev-url 是临时数据库，用于执行计算迁移计划，需要使用一个空的数据库
atlas schema apply \
 -u "mysql://root:123456@:3306/ent_test" \
 --to "ent://ent/schema" \
 --dev-url "mysql://root:123456@:3306/temp_db"


# 清空数据库
#atlas schema clean \
#  -u "mysql://root:123456@:3306/ent_test"


# 比较 Ent Schema 与数据库的差异，生成迁移文件
#atlas migrate diff init \
#  --dir "file://migrations" \
#  --to "ent://ent/schema" \
#  --dev-url "mysql://root:123456@:3306/ent_test"


# 执行迁移
#atlas migrate apply \
# --dir "file://migrations" \
# --url "mysql://root:123456@:3306/ent_test"