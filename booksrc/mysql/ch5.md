# 排序检索数据

## 5.1 排序数据

```
select prod_name from products order by prod_name;
```

## 5.2 按多个列排序

```
select prod_id, prod_price, prod_name from products order by prod_price, prod_name;
```

## 5.3 指定排序方向

```
select prod_id, prod_price, prod_name from products order by prod_price desc;
select prod_id, prod_price, prod_name from products order by prod_price desc, prod_name;
select prod_id, prod_price, prod_name from products order by prod_price desc, prod_name asc;
```
```
select prod_price from products order by prod_price desc limit 1;
```

