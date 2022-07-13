package main

import (
	"go.uber.org/zap"
	"net/http"
)

var sugarLogger *zap.SugaredLogger

func main()  {
	INitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.google.com")
	simpleHttpGet("http://www.google.com")
}

func INitLogger()  {
	logger, _:= zap.NewProduction()
	sugarLogger =logger.Sugar()
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		sugarLogger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}