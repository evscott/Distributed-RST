package models

/* Information/Metadata about node */
type NodeInfo struct {
	NodeId  int  `json:"nodeId"`
	NodeIpAddr  string  `json:"nodeIpAddr"`
	Port  string  `json:"port"`
}