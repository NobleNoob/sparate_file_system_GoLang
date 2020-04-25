package config

const (
	// AsyncTransferEnable : 是否开启文件异步转移(默认同步)
	AsyncTransferEnable = true
	// RabbitURL : rabbitmq服务的入口url
	RabbitURL = "amqp://guest:guest@127.0.0.1:5672/"
	TransExchangeName = "uploadserver.trans"
	TransS3QueueName = "uploadserver.trans.s3"
	TransS3ErrQueueName = "uploadserver.trans.s3.err"
	TransS3RoutingKey = "s3"
)
