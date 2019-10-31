# 检索数据

## 4.2 检索单个列

```
select prod_name from products;
```

## 4.3 检索多个列

```
select prod_id, prod_name, prod_price from products;
```

## 4.4 检索所有列

```
select * from products;
```

## 4.5 检索不同的行

```
select vend_id from products;
select distinct(vend_id) from products;
select distinct vend_id from products;
```

## 4.6 限制结果

```
select prod_name from products limit 5;
select prod_name from products limit 4, 5;
select prod_name from products limit 5 offset 4;
```

## 4.7 使用完全限定的表名

```
select products.prod_name from products;
select products.prod_name from example.products;
```

