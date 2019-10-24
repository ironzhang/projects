# 用通配符进行过滤

## 8.1 LIKE操作符

### 8.1.1 百分号(%)通配符

```
select prod_id, prod_name from products where prod_name like 'jet%';
select prod_id, prod_name from products where prod_name like '%anvil%';
select prod_id, prod_name from products where prod_name like 's%e';
```

### 8.1.2 下划线(_)通配符

```
select prod_id, prod_name from products where prod_name like '_ ton anvil';
select prod_id, prod_name from products where prod_name like '% ton anvil';
```

