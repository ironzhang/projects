# 数据过滤

## 7.1 组合WHERE子句

### 7.1.1 AND操作符

```
select prod_id, prod_price, prod_name from products where vend_id = 1003 and prod_price <= 10;
```

### 7.1.2 OR操作符

```
select prod_name, prod_price from products where vend_id = 1002 or vend_id = 1003;
```

### 7.1.3 计算次序

```
select prod_name, prod_price from products where vend_id = 1002 or vend_id = 1003 and prod_price >= 10;
select prod_name, prod_price from products where (vend_id = 1002 or vend_id = 1003) and prod_price >= 10;
```

## 7.2 IN操作符

```
select prod_name, prod_price from products where vend_id in (1002, 1003);
```

## 7.3 NOT操作符

```
select prod_name, prod_price from products where vend_id not in (1002, 1003);
```

