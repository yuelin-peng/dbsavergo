# dbsavergo
	本项目用于在db层发生容灾切换时，且不保证数据不丢失的情况下，如何保证业务连续运转，且保证数据最终一致

# 原理
	db saver通过业务写入订单数据和版本号，同时订阅备机的binlog来消除存储的订单数据，最终维护一份未同步到备机的数据列表。
	当主db发生异常时，业务将访问db saver判断该订单是否已同步。如果未同步，业务可以拒绝为订单的服务。
	使用db saver，业务可以将db异常时，对业务的影响面控制在同步延迟范围内的订单上，不会大面积影响用户。当原主db恢复时，可以
	恢复之前未同步的数据。

# 如何保证db saver的可用性：
1、通过降级db saver server的存储来实现，比如：当db saver的第一级存储失败，可以降级到二级存储，然后db saver server再旁路从二级存储中获取数据写入一级存储
2、当db saver client远程调用服务失败时，可以存储到本地，然后再异步同步到db saver server，不过在服务器异常时，会有丢失数据风险

# 如何提高实时业务的可用性
	Db saver和db同时出现异常的可能性时两两者异常可能性的乘积，其远小于单一组件的异常可能性。
	为了达到增加一个组件，不降低可用性，当db saver异常时，直接降级，其降级逻辑如下：
	(1) 写入降级：同步写入数据到db saver，当db saver发生异常时，直接降级到本地存储，然后异步同步，此时可能会造成数据丢失
	(2) 读取：当db saver查询异常时，直接认为无需拦截

# db saver丢数据影响分析
	为了保证db saver不降低实时业务的可用性，db saver降级逻辑会存在丢失数据的情况，其影响可分下面三种情况分析：
	(1) 当业务db正常时，随着时间推移，比如：5分钟，数据已同步到备机，此时丢失的数据已无关紧要，是否同步到db saver都不影响数据的可用性。
	(2) 当业务db发生异常时，如果db saver丢失了最近的部分数据，只是造成部分订单没有被控制，远比完全不控制要好。
	(3) 如果在db saver存储异常，且业务db异常时，服务回到原始情况，并不会造成情况恶化。
