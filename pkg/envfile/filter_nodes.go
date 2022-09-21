package envfile

import (
	"strings"

	"go.firedancer.io/radiance/proto/env/v1"
)

func ParseOnlyFlag(only string) []string {
	if only == "" {
		return nil
	}
	return strings.Split(only, ",")
}

func FilterNodes(nodes []*envv1.RPCNode, only []string) []*envv1.RPCNode {
	if len(only) == 0 {
		return nodes
	}
	var filtered []*envv1.RPCNode
	for _, node := range nodes {
		for _, o := range only {
			if node.Name == o {
				filtered = append(filtered, node)
			}
		}
	}
	return filtered
}
