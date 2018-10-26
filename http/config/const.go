package config

const (

	//压缩文件生成的结果文件夹
	ServerLogDir = "server_log"


	//压缩文件生成的结果文件夹
	ZipResultDir = "download"


	//临时中转的文件夹
	TmpTransferDir = "tmpTransfer"

	// 队列名称
	RedisTaskName = "ai-dong-task"
)


type RedisTaskStatus int
//任务的状态
const (

	_  RedisTaskStatus = iota

	RedisStatusNotStart

	RedisStautsRunning

	RedisStatusFailure

	RedisStatusOK

)

