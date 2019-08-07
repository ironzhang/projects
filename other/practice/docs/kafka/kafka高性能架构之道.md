# kafka 高性能架构之道

---

## 1. 利用Partition实现并行处理

**Partition提供并行处理的能力**

Kafka是一个Pub-Sub的消息系统，无论是发布还是订阅，都须指定Topic。Topic只是一个逻辑的概念。每个Topic都包含一个或多个Partition，不同Partition可位于不同节点。同时Partition在物理上对应一个本地文件夹，每个Partition包含一个或多个Segment，每个Segment包含一个数据文件和一个与之对应的索引文件。在逻辑上，可以把一个Partition当作一个非常长的数组，可通过这个“数组”的索引（offset）去访问其数据。

一方面，由于不同Partition可位于不同机器，因此可以充分利用集群优势，实现机器间的并行处理。另一方面，由于Partition在物理上对应一个文件夹，即使多个Partition位于同一个节点，也可通过配置让同一节点上的不同Partition置于不同的disk drive上，从而实现磁盘间的并行处理，充分发挥多磁盘的优势。

利用多磁盘的具体方法是，将不同磁盘mount到不同目录，然后在server.properties中，将log.dirs设置为多目录（用逗号分隔）。Kafka会自动将所有Partition尽可能均匀分配到不同目录也即不同目录（也即不同disk）上。

注：虽然物理上最小单位是Segment，但Kafka并不提供同一Partition内不同Segment间的并行处理。因为对于写而言，每次只会写Partition内的一个Segment，而对于读而言，也只会顺序读取同一Partition内的不同Segment。

**Partition是最小并发粒度**

多Consumer消费同一个Topic时，同一条消息只会被同一Consumer Group内的一个Consumer所消费。而数据并非按消息为单位分配，而是以Partition为单位分配，也即同一个Partition的数据只会被一个Consumer所消费（在不考虑Rebalance的前提下）。

如果Consumer的个数多于Partition的个数，那么会有部分Consumer无法消费该Topic的任何数据，也即当Consumer个数超过Partition后，增加Consumer并不能增加并行度。

简而言之，Partition个数决定了可能的最大并行度。如下图所示，由于Topic 2只包含3个Partition，故group2中的Consumer 3、Consumer 4、Consumer 5 可分别消费1个Partition的数据，而Consumer 6消费不到Topic 2的任何数据。

![](img/kg_1.png)

以Spark消费Kafka数据为例，如果所消费的Topic的Partition数为N，则有效的Spark最大并行度也为N。即使将Spark的Executor数设置为N+M，最多也只有N个Executor可同时处理该Topic的数据

## 2. ISR实现可用性与数据一致性的动态平衡

**CAP理论**

CAP理论是指，分布式系统中，一致性、可用性和分区容忍性最多只能同时满足两个。

* 一致性
	* 通过某个节点的写操作结果对后面通过其它节点的读操作可见
	* 如果更新数据后，并发访问情况下后续读操作可立即感知该更新，称为强一致性
	* 如果允许之后部分或者全部感知不到该更新，称为弱一致性
	* 若在之后的一段时间（通常该时间不固定）后，一定可以感知到该更新，称为最终一致性
* 可用性
	* 任何一个没有发生故障的节点必须在有限的时间内返回合理的结果
* 分区容忍性
	* 部分节点宕机或者无法与其它节点通信时，各分区间还可保持分布式系统的功能

一般而言，都要求保证分区容忍性。所以在CAP理论下，更多的是需要在可用性和一致性之间做权衡。

**常用数据复制及一致性方案**

* Master-Slave
	* RDBMS的读写分离即为典型的Master-Slave方案
	* 同步复制可保证强一致性但会影响可用性
	* 异步复制可提供高可用性但会降低一致性
* WNR
	* 主要用于去中心化的分布式系统中。DynamoDB与Cassandra即采用此方案或其变种
	* N代表总副本数，W代表每次写操作要保证的最少写成功的副本数，R代表每次读至少要读取的副本数
	* 当W+R>N时，可保证每次读取的数据至少有一个副本拥有最新的数据
	* 多个写操作的顺序难以保证，可能导致多副本间的写操作顺序不一致。Dynamo通过向量时钟保证最终一致性
* Paxos及其变种
	* Google的Chubby，Zookeeper的原子广播协议（Zab），RAFT等

## 基于ISR的数据复制方案

Kafka的数据复制是以Partition为单位的。而多个备份间的数据复制，通过Follower向Leader拉取数据完成。从一这点来讲，Kafka的数据复制方案接近于上文所讲的Master-Slave方案。不同的是，Kafka既不是完全的同步复制，也不是完全的异步复制，而是基于ISR的动态复制方案。

ISR，也即In-sync Replica。每个Partition的Leader都会维护这样一个列表，该列表中，包含了所有与之同步的Replica（包含Leader自己）。每次数据写入时，只有ISR中的所有Replica都复制完，Leader才会将其置为Commit，它才能被Consumer所消费。

这种方案，与同步复制非常接近。但不同的是，这个ISR是由Leader动态维护的。如果Follower不能紧“跟上”Leader，它将被Leader从ISR中移除，待它又重新“跟上”Leader后，会被Leader再次加加ISR中。每次改变ISR后，Leader都会将最新的ISR持久化到Zookeeper中。

至于如何判断某个Follower是否“跟上”Leader，不同版本的Kafka的策略稍微有些区别。

* 对于0.8.*版本，如果Follower在replica.lag.time.max.ms时间内未向Leader发送Fetch请求（也即数据复制请求），则Leader会将其从ISR中移除。如果某Follower持续向Leader发送Fetch请求，但是它与Leader的数据差距在replica.lag.max.messages以上，也会被Leader从ISR中移除。
* 从0.9.0.0版本开始，replica.lag.max.messages被移除，故Leader不再考虑Follower落后的消息条数。另外，Leader不仅会判断Follower是否在replica.lag.time.max.ms时间内向其发送Fetch请求，同时还会考虑Follower是否在该时间内与之保持同步。
0* .10.* 版本的策略与0.9.*版一致
* 对于0.8.*版本的replica.lag.max.messages参数，很多读者曾留言提问，既然只有ISR中的所有Replica复制完后的消息才被认为Commit，那为何会出现Follower与Leader差距过大的情况。原因在于，Leader并不需要等到前一条消息被Commit才接收后一条消息。事实上，Leader可以按顺序接收大量消息，最新的一条消息的Offset被记为LEO（Log end offset）。而只有被ISR中所有Follower都复制过去的消息才会被Commit，Consumer只能消费被Commit的消息，最新被Commit的Offset被记为High watermark。换句话说，LEO 标记的是Leader所保存的最新消息的offset，而High watermark标记的是最新的可被消费的（已同步到ISR中的Follower）消息。而Leader对数据的接收与Follower对数据的复制是异步进行的，因此会出现Hight watermark与LEO存在一定差距的情况。0.8.*版本中replica.lag.max.messages限定了Leader允许的该差距的最大值。

Kafka基于ISR的数据复制方案原理如下图所示。

![](img/kg_2.png)

如上图所示，在第一步中，Leader A总共收到3条消息，故其high watermark为3，但由于ISR中的Follower只同步了第1条消息（m1），故只有m1被Commit，也即只有m1可被Consumer消费。此时Follower B与Leader A的差距是1，而Follower C与Leader A的差距是2，均未超过默认的replica.lag.max.messages，故得以保留在ISR中。在第二步中，由于旧的Leader A宕机，新的Leader B在replica.lag.time.max.ms时间内未收到来自A的Fetch请求，故将A从ISR中移除，此时ISR={B，C}。同时，由于此时新的Leader B中只有2条消息，并未包含m3（m3从未被任何Leader所Commit），所以m3无法被Consumer消费。第四步中，Follower A恢复正常，它先将宕机前未Commit的所有消息全部删除，然后从最后Commit过的消息的下一条消息开始追赶新的Leader B，直到它“赶上”新的Leader，才被重新加入新的ISR中。

**使用ISR方案的原因**

由于Leader可移除不能及时与之同步的Follower，故与同步复制相比可避免最慢的Follower拖慢整体速度，也即ISR提高了系统可用性。
ISR中的所有Follower都包含了所有Commit过的消息，而只有Commit过的消息才会被Consumer消费，故从Consumer的角度而言，ISR中的所有Replica都始终处于同步状态，从而与异步复制方案相比提高了数据一致性。
ISR可动态调整，极限情况下，可以只包含Leader，极大提高了可容忍的宕机的Follower的数量。与Majority Quorum方案相比，容忍相同个数的节点失败，所要求的总节点数少了近一半。

**ISR相关配置说明**

Broker的min.insync.replicas参数指定了Broker所要求的ISR最小长度，默认值为1。也即极限情况下ISR可以只包含Leader。但此时如果Leader宕机，则该Partition不可用，可用性得不到保证。
只有被ISR中所有Replica同步的消息才被Commit，但Producer发布数据时，Leader并不需要ISR中的所有Replica同步该数据才确认收到数据。Producer可以通过acks参数指定最少需要多少个Replica确认收到该消息才视为该消息发送成功。acks的默认值是1，即Leader收到该消息后立即告诉Producer收到该消息，此时如果在ISR中的消息复制完该消息前Leader宕机，那该条消息会丢失。而如果将该值设置为0，则Producer发送完数据后，立即认为该数据发送成功，不作任何等待，而实际上该数据可能发送失败，并且Producer的Retry机制将不生效。更推荐的做法是，将acks设置为all或者-1，此时只有ISR中的所有Replica都收到该数据（也即该消息被Commit），Leader才会告诉Producer该消息发送成功，从而保证不会有未知的数据丢失。

e
