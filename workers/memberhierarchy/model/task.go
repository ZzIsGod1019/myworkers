package model

// {"lang":"ko", "action":"set_hierarchy", "data":{"member_id":12345, "hierarchy_id":"600", "time":"1660216823"}}
type TaskInfo struct {
	Lang   string `json:"lang"`
	Action string `json:"action"`
	Data   Data   `json:"data"`
}

type Data struct {
	MemberId    int    `json:"member_id"`
	HierarchyId int    `json:"hierarchy_id"`
	Time        string `json:"time"`
}
