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

	//字符串分割符号
	ConditionSplitChar = "-"

)

type RedisTaskStatus int

//任务的状态
const (
	_ RedisTaskStatus = iota

	RedisStatusNotStart

	RedisStautsRunning

	RedisStatusFailure

	RedisStatusOK
)

const (
	RedisStatusNotStartStr = "任务未开始"

	RedisStautsRunningStr = "任务运行中"

	RedisStatusFailureStr = "任务失败"

	RedisStatusOKStr = "任务运行成功"
)

func RedisStatus(s RedisTaskStatus) string {

	result := ""

	switch s {
	case RedisStatusNotStart:
		result = RedisStatusNotStartStr
	case RedisStautsRunning:
		result = RedisStautsRunningStr

	case RedisStatusFailure:
		result = RedisStatusFailureStr
	case RedisStatusOK:
		result = RedisStatusOKStr
	}

	return result
}
