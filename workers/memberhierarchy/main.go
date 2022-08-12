package main

import (
	"encoding/json"
	"myworkers/driver/database"
	"myworkers/driver/logger"
	"myworkers/workers/memberhierarchy/model"
)

func main() {
	// 初始化mysql
	database.MysqlInit()
	// 初始化Redis
	database.RedisInit()
	// 初始化日志
	logger.LoggerInit()

	hanlde("")
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
	res := model.UpdateHierarchyByMemberId(taskInfo.Lang, taskInfo.Data.MemberId, taskInfo.Data.HierarchyId, taskInfo.Data.Time)
	if !res {
		logger.Trace("error", taskInfo.Lang, "更新失败    "+taskJson)
		return false
	}
	return true
}
