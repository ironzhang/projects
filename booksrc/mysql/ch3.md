# 使用MySQL

## 3.2 选择数据库

```
use example;
```
或者
```
mysql -h127.0.0.1 -uroot -p123456 -D example;
```

## 3.3 了解数据库和表

```
show databases;
show tables;
```
```
show columns from customers;
describe customers;
desc customers;
```
```
show create database example;
show create table customers;
```
```
show status;
show grants;
show grants for root;
show errors;
show warnings;
```
```
help show;
```

