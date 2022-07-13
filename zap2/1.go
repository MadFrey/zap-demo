package main

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var sugarLogger *zap.SugaredLogger

func main()  {
	INitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.google.com") //刷新缓存区
	simpleHttpGet("http://www.google.com")
}

func INitLogger()  {
	// 进行日志分级
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error  以下分均为lowPriority debug等级为-1 info为0
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool{ //error及以上分为highPriority 等级为1及以上
		return lev >= zap.ErrorLevel
	})
	LowLogWriter:=getLowLogWriter()
	HighLogWriter:= getHighLogWriter()
	encoder := GetEncoder()
	LowCore := zapcore.NewCore(encoder, LowLogWriter,lowPriority) // DebugLevel 配置日志级别 表示那种级别的日志将被写入
	HighCore := zapcore.NewCore(encoder, HighLogWriter, highPriority)
	logger:=zap.New(zapcore.NewTee(HighCore,LowCore),zap.AddCaller()) //加上zap.Addcaller（）后，记录器将会以程序的文件名，行号和函数名注释每条消息
	sugarLogger=logger.Sugar()

}


func GetEncoder()zapcore.Encoder { //编码器 如何写入日志
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime=zapcore.ISO8601TimeEncoder //修改时间编码器
	encoderConfig.EncodeLevel=zapcore.CapitalLevelEncoder //在日志中使用大写字母记录日志
	return zapcore.NewConsoleEncoder(encoderConfig)  //使用普通Encode，非json
}

//使用Lumberjack进行日志切割归档
func getLowLogWriter() zapcore.WriteSyncer { //指定日志的写入位置
	LowJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",  //日志文件的位置
		MaxSize:    1, // 在进行切割前，日志文件的最大大小
		MaxBackups: 5,  //保留旧文件的最大个数
		MaxAge:     30,  //保留旧文件的最大天数
		Compress:   false, //是否压缩/归档旧文件
	}


	return zapcore.AddSync(LowJackLogger)  //AddSync 自动匹配合适的Sync函数处理
}

func getHighLogWriter() zapcore.WriteSyncer { //指定日志的写入位置
	HighJackLogger := &lumberjack.Logger{
		Filename:   "./test2.log",  //日志文件的位置
		MaxSize:    1, // 在进行切割前，日志文件的最大大小
		MaxBackups: 5,  //保留旧文件的最大个数
		MaxAge:     30,  //保留旧文件的最大天数
		Compress:   false, //是否压缩/归档旧文件
	}
	return zapcore.AddSync(HighJackLogger)
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)  //对应debug 等级 ，其他以此类推
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}










