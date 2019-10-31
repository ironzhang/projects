# 用正则表达式进行搜索

## 9.2 使用MySQL正则表达式

### 9.2.1 基本字符匹配

```
select prod_name from products where prod_name regexp '1000';
select prod_name from products where prod_name regexp '.000' order by prod_name;
select prod_name from products where prod_name regexp binary 'JetPack .000' order by prod_name;
```

### 9.2.2 进行OR匹配

```
select prod_name from products where prod_name regexp '1000|2000' order by prod_name;
```

### 9.2.3 匹配几个字符之一

```
select prod_name from products where prod_name regexp '[123] Ton' order by prod_name;
select prod_name from products where prod_name regexp '[1|2|3] Ton' order by prod_name;
select prod_name from products where prod_name regexp '[^123] Ton' order by prod_name;
```

### 9.2.4 匹配范围

```
select prod_name from products where prod_name regexp '[1-5] Ton' order by prod_name;
```

### 9.2.5 匹配特殊字符

```
select vend_name from vendors where vend_name regexp '.' order by vend_name;
select vend_name from vendors where vend_name regexp '\\.' order by vend_name;
```

### 9.2.7 匹配多个实例

```
select prod_name from products where prod_name regexp '\\([0-9] sticks?\\)' order by prod_name;
select prod_name from products where prod_name regexp '[[:digit:]]{4}' order by prod_name;
select prod_name from products where prod_name regexp '[:digit:]{4}' order by prod_name;
```

### 9.2.8 定位符

```
select prod_name from products where prod_name regexp '^[0-9\\.]' order by prod_name;
```
```
select 'hello' regexp '[0-9]';
```

