package main

import (
	"encoding/json"
	"myworkers/config"
	"myworkers/database"
	"myworkers/logger"
	"myworkers/sqs"
	"myworkers/workers/memberhierarchy/model"
	"strconv"
)

func main() {
	// 读取配置初始化
	config.Init()
	// 初始化mysql
	database.MysqlInit()
	// 初始化Redis
	database.RedisInit()
	// 初始化日志
	logger.LoggerInit()

	ch := make(chan struct{}, 3)
	for {
		ch <- struct{}{}
		go func() {
			defer func() {
				<-ch
			}()
			message, err := sqs.ReceiveMessage()
			if err == nil {
				res := hanlde(*message.Body)
				logger.Trace("info", "", "队列消费"+strconv.FormatBool(res)+"    "+*message.Body)
			}
		}()
	}
}

func hanlde(taskJson string) bool {
	taskInfo := model.TaskInfo{}
	json.Unmarshal([]byte(taskJson), &taskInfo)
	info, err := model.GetMemberHierarchyByMemberId(taskInfo.Lang, taskInfo.Data.MemberId)
	if err != nil {
		logger.Trace("error", taskInfo.Lang, "未找到用户数据    "+taskJson)
		return false
	}
	if info.HandleTime > taskInfo.Data.Time {
		logger.Trace("error", taskInfo.Lang, "旧数据不做更新    "+taskJson)
		return false
	}
	res, err := model.UpdateHierarchyByMemberId(taskInfo.Lang, taskInfo.Data.MemberId, taskInfo.Data.HierarchyId, taskInfo.Data.Time)
	if !res && err != nil {
		logger.Trace("error", taskInfo.Lang, err.Error()+"    "+taskJson)
		return false
	}
	return true
}
