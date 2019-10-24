# 过滤数据

## 6.1 使用WHERE子句

```
select prod_name, prod_price from products where prod_price=2.50;
```

## 6.2 WHERE子句操作符

```
select prod_name, prod_price from products where prod_price between 2.5 and 5;
select prod_name, prod_price from products where prod_price >= 2.5 and prod_price <= 5;
```

### 6.2.1 检查单个值

mysql在执行匹配时默认不区分大小写。
```
select prod_name, prod_price from products where prod_name='fuses';
```
```
select prod_name, prod_price from products where prod_price < 10;
select prod_name, prod_price from products where prod_price <= 10;
```

### 6.2.2 不匹配检查

```
select vend_id, prod_name from products where vend_id <> 1003;
select vend_id, prod_name from products where vend_id != 1003;
```

### 6.2.3 范围值检查

```
select prod_name, prod_price from products where prod_price between 5 and 10;
```

### 6.2.4 空值检查

```
select prod_name from products where prod_price is null;
select cust_id from customers where cust_email is null;
```

