package service

import (
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
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

// Filter 过滤
// 用于过滤元素，比较元素的 Name 属性，若包含，则返回
func (d *dataSelector) Filter() *dataSelector {
	// 若 Name 的传参为空，则返回所有元素
	if d.dataSelectorQuery.FilterQuery.Name == "" {
		return d
	}
	// 若 Name 的传参不为空，则返回元素中包含Name的所有元素
	// 定义一个列表
	filteredList := []DataCell{}
	// 遍历所有数据
	for _, value := range d.GenericDataList {
		// 定义一个是否匹配的值
		matched := true
		// 获取对象的 Name
		objName := value.GetName()
		// 如果数据中不包含要过滤的 Name,则 matched 为 flase
		if !strings.Contains(objName, d.dataSelectorQuery.FilterQuery.Name) {
			matched = false
			continue
		}
		// 如果匹配到要过滤的 Name,则将数据添加到 filteredList 中
		if matched {
			filteredList = append(filteredList, value)
		}
	}
	// 将 filteredList 赋值给 d.GenericDataList 并返回
	d.GenericDataList = filteredList
	return d
}

// Paginate 用于数组分页，根据 Limit 和 Page 的传参，返回数据
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectorQuery.PaginateQuery.Limit
	page := d.dataSelectorQuery.PaginateQuery.Page
	// 验证参数合法，若参数不合法，则返回所有数据
	if limit <= 0 || page <= 0 {
		return d
	}
	// 定义offset
	// 举例：25个元素的切片 limit10
	// page1 start0 end 10
	// page2 start10 end 20
	// page3 start20 end 30
	// 含头不含尾
	startIndex := limit * (page - 1)
	endIndex := limit * page
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// 定义 podCell 类型，实现两个方法 GetCreation GetName，可进行类型转换
type podCell corev1.Pod

// GetCreation 获取 创建时间 这里 (p podCell) 没有使用指针
func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

// GetName 获取 名称
func (p podCell) GetName() string {
	return p.Name
}
