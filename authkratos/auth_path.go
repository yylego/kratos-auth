package authkratos

import (
	"fmt"
	"maps"
)

// RouteScope defines the scope of routes to match
// RouteScope 定义需要匹配的路由范围
type RouteScope struct {
	Side         SelectSide         // INCLUDE or EXCLUDE mode // 包含或排除模式
	OperationSet map[Operation]bool // Set of operations to match // 需要匹配的操作集合
}

// NewInclude creates a RouteScope that matches specified operations
// NewInclude 创建仅匹配指定操作的 RouteScope
func NewInclude(operations ...Operation) *RouteScope {
	return &RouteScope{
		Side:         INCLUDE,
		OperationSet: newSet(operations),
	}
}

// NewExclude creates a RouteScope that matches all except specified operations
// NewExclude 创建排除指定操作后匹配所有其他操作的 RouteScope
func NewExclude(operations ...Operation) *RouteScope {
	return &RouteScope{
		Side:         EXCLUDE,
		OperationSet: newSet(operations),
	}
}

// Match checks if operation is within the scope
// Match 检查操作是否在范围内
func (c *RouteScope) Match(operation Operation) bool {
	switch c.Side {
	case INCLUDE:
		return c.OperationSet[operation]
	case EXCLUDE:
		return !c.OperationSet[operation]
	default:
		panic(fmt.Errorf("unknown select side %q", c.Side))
	}
}

// Opposite returns a RouteScope with inverted side and a clone of the same operations
// Opposite 返回反转 side 的 RouteScope，操作集合复制一份
func (c *RouteScope) Opposite() *RouteScope {
	return &RouteScope{
		Side:         c.Side.Opposite(),
		OperationSet: maps.Clone(c.OperationSet),
	}
}

func newSet[T comparable](slice []T) map[T]bool {
	set := make(map[T]bool, len(slice))
	for _, v := range slice {
		set[v] = true
	}
	return set
}
