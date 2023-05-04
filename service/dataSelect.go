package service

import (
	"sort"
	"time"
)

// 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDataList   []DataCell
	dataSelectorQuery *DataSelectorQuery
}

// DataCell 用于各种资源list的类型转换，转换后可以使用dataSelector的排序、过滤、分页
type DataCell interface {
	// GetCreation 用于排序
	GetCreation() time.Time
	// GetName 用于过滤
	GetName() string
}

// DataSelectorQuery 定义过滤和分页的属性，过滤：Name，分页：Limit 和 page
type DataSelectorQuery struct {
	// 用于过滤
	FilterQuery *FilterQuery
	// 用于分页
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

// 排序，实现自定义结构的排序，需要重写Len、Swap、Less方法

// Len 方法用于获取数组长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap 方法用于数组中的元素在比较大小后怎么交换位置，可定义升降序
func (d *dataSelector) Swap(i, j int) {
	// i,j 是切片的下标
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 方法用于定义数组中元素排序的“大小”的比较方式
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

// Sort 重写以上3个方法用使用 sort.Sort 进行排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}
