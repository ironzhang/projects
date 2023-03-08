# Design

---

# Books

## heinemann GK

## At the market

apple  [ˈæpl]
I Can see the apples.

banana  [bəˈnɑ:nə]
I Can see the bananas.

## Rex

swing  [swɪŋ]
I like to swing

jump   [dʒʌmp]
I like to jump

## Peppa Pig
### Muddy puddles
### Mr Dinosaur is lost

---


## Words

|word|phonetic_alphabet|pic|sound|
|----|-----------------|---|-----|
|apple|[ˈæpl]|pic url|sound url|

## Tags

|word|tag|
|----|--|
|apple|food|
|apple|fruit|
|banana|food|
|banana|fruit|
|swing|play|

## Books

|id|tittle|class|
|--|------|-----|
|1|At the market|heinemann GK|
|2|Rex|heinemann GK|

## Contents

|book_id|word|sentence|
|-------|----|--------|
|1|apple|I can see the apples.|
|1|banana|I can see the bananas.|
|2|swing|I like to swing.|
|2|jump|I like to jump.|
|3|jump|George likes to jump in muddy puddles too.|

---

## Interface

### list groups

```
GET /api/groups
```

### add books

```
POST /api/books
```

### list books

```
GET /api/books?class=heinemann GK
```

3. get book content

```
GET /api/books/1
```

4. get word

```
GET /api/words/apple
```

5. set word tags

```
SET /api/words/apple
```

6. list tags

```
GET /api/tags
```

7. list tag words

```
GET /api/tags/fruit
```

---

## Support

```
http://dict-co.iciba.com/api/dictionary.php?key=9703E5C2F64060501DB26311D45CF2CA&type=json&w=apple
```

