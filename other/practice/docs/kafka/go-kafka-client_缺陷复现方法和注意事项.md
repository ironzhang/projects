# go-kafka-client缺陷复现方法和注意事项

标签（空格分隔）： ablecloud

---

这个问题可能是客户端本身的缺陷，也可能是本身是不打算支持这种情况

---

## 1. 复现方法

---

**创建三个topic: test1, test2, test3**

```
./bin/kafka-topics.sh --create --zookeeper 127.0.0.1:2181 --topic test1 --replication-factor 1 --partitions 1
./bin/kafka-topics.sh --create --zookeeper 127.0.0.1:2181 --topic test2 --replication-factor 1 --partitions 1
./bin/kafka-topics.sh --create --zookeeper 127.0.0.1:2181 --topic test3 --replication-factor 1 --partitions 1
```

---

**实现三个producer分别往三个topic写数据**

* [simple_producer源码](img/simple_producer.go)
* [simple_producer配置文件](img/producer.conf)

```
./simple_producer --topic=test1 --config=producer.conf
./simple_producer --topic=test2 --config=producer.conf
./simple_producer --topic=test3 --config=producer.conf
```

---

**创建两个consumer，一个consumer消费test1,test2，一个consumer消费test1,test3**

* [simple_consumper源码](img/simple_consumer.go)
* [simple_consumer配置文件1](img/consumer1.conf)
* [simple_consumer配置文件2](img/consumer2.conf)

```
./simple_conumer --topic=test1,test2 --config=consumer1.conf

// 等待5秒

./simple_conumer --topic=test1,test3 --config=consumer2.conf
```

发现consumer1可以正常消费test1,test2，但是consumer2无法消息test3，并报错:

```
2018-03-01/17:25:55 [ERROR] [zk] Failed awaiting on state barrier 09d6ad5263dc29eed7efb047d877d032 [Timed out waiting for consensus on barrier path /consumers/ac/api/rebalance/09d6ad5263dc29eed7efb047d877d032
```

将consumer1停掉后，consumer2恢复正常 

---

## 2. 注意事项

* 当前go-kafka-client对同一个consumer-group中的不同consumer，消费相异的topic的处理有问题，所以避免这种使用方式，如果有这种需求，可以放到不周的consumer-group中
* 如果一个consumer同时消费多个topic，避免使用kafka自带的工具如kafka-console-consumer.sh使用同样的consumer-group，消息部分topic
