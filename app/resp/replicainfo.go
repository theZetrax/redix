package resp

import (
	"fmt"
	"reflect"
	"strings"
)

// NodeRole represents the role of the node in the Redis cluster.
type NodeRole string

const (
	RoleMaster  NodeRole = "master"
	RoleReplica NodeRole = "slave"
)

// It contains the host and port of the current node.
type NodeInfo struct {
	Role                NodeRole `resp:"role"`
	MasterReplicaId     string   `resp:"master_replid"`
	MasterReplicaOffset string   `resp:"master_repl_offset"`
	Host                string   `resp:"-"`
	Port                string   `resp:"-"`
	MasterPort          string   `resp:"-"`
	MasterHost          string   `resp:"-"`
	OffsetCount         int      `resp:"-"`
}

func NewNodeInfo(host string, port string, replicaid string, replicaoffset string, master_addr string) *NodeInfo {
	info := &NodeInfo{
		Host:                host,
		Port:                port,
		Role:                RoleMaster,
		MasterReplicaId:     replicaid,
		MasterReplicaOffset: replicaoffset,
		OffsetCount:         -1,
	}

	if master_addr != "" {
		info.Role = RoleReplica
		master_host, master_port, found := strings.Cut(master_addr, ":")
		if !found {
			panic("Invalid master node address: " + master_addr)
		}
		info.MasterHost = master_host
		info.MasterPort = master_port
	}

	return info
}

func EncodeNodeInfo(node_info interface{}) []byte {
	response := []byte("# Replication" + CRLF)
	size := len(response)
	lnCount := 1

	entries := reflect.ValueOf(node_info)
	types := entries.Type()

	for i := 0; i < entries.NumField(); i++ {
		entry := entries.Field(i).String()
		if tag := types.Field(i).Tag.Get("resp"); tag != "-" {
			response = append(response, []byte(tag+":"+entry+CRLF)...)
			size += len(tag) + len(entry) + 1
			lnCount += 1
		}
	}

	// Append the size of the response
	size_of_response := []byte("$" + fmt.Sprint(size+lnCount) + CRLF)
	response = append(size_of_response, response...)

	return response
}
