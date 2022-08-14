package model

import (
	"myworkers/database"
	"myworkers/logger"
	"strconv"
)

var tableName string = "tbl_member_hierarchys"

type memberHierarchyInfo struct {
	Id          int    `json:"id" form:"id"`
	MemberId    int    `json:"member_id" form:"member_id"`
	HierarchyId string `json:"hierarchy_id" form:"hierarchy_id"`
	HandleTime  string `json:"handle_time" form:"handle_time"`
	CreatedAt   string `json:"created_at" form:"created_at"`
	UpdatedAt   string `json:"updated_at" form:"updated_at"`
}

func GetMemberHierarchyByMemberId(land string, memberId int) (memberHierarchyInfo, error) {
	var chapter = memberHierarchyInfo{}
	databaseName := database.MainDatabaseName[land]
	row := database.GetMainDb(land).QueryRow("select id,member_id,hierarchy_id,handle_time from "+databaseName+"."+tableName+" where member_id = ?;", memberId)
	err := row.Scan(&chapter.Id, &chapter.MemberId, &chapter.HierarchyId, &chapter.HandleTime)
	return chapter, err
}

func UpdateHierarchyByMemberId(land string, memberId int, hierarchyId int, handleTime string) bool {
	databaseName := database.MainDatabaseName[land]
	r, err := database.GetMainDb(land).Exec("update "+databaseName+"."+tableName+" set hierarchy_id = ? and handle_time = ? where member_id = ?;", hierarchyId, handleTime, memberId)

	if err != nil {
		logger.Trace("error", land, "修改用户失败    "+strconv.Itoa(memberId))
		return false
	}
	_, err = r.RowsAffected()
	if err != nil {
		logger.Trace("info", land, "未修改当前用户分组    "+strconv.Itoa(memberId))
		return false
	}
	return true
}
