package main

import (
	"encoding/json"
	"myworkers/config"
	"myworkers/database"
	"myworkers/logger"
	"myworkers/sqs"
	"myworkers/workers/memberhierarchy/model"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// 接收终止信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	isStop := false // 是否接收到结束信号
	// 读取配置初始化
	config.Init()
	// 初始化mysql
	database.MysqlInit()
	// 初始化Redis
	database.RedisInit()
	// 初始化日志
	logger.LoggerInit()

	ch := make(chan struct{}, 3) // 限制协程数量
	for {
		select {
		case <-c: // 接收到终止信号后不再创建任务，等待当前任务执行完毕后退出程序
			isStop = true
			// return
		default:
			if isStop {
				if len(ch) == 0 {
					return
				} else {
					continue
				}
			}

			ch <- struct{}{}
			go func() {
				defer func() {
					<-ch
				}()
				message, err := sqs.ReceiveMessage()
				if err == nil {
					res := hanlde(*message.Body)
					logger.Trace("info", "", "队列消费"+strconv.FormatBool(res)+"    "+*message.Body)
				} else {
					logger.Trace("error", "", "队列服务报错    "+err.Error())
				}
			}()
		}
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
