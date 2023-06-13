/**
 * Title: quad_tree
 * Author: Anyazhou
 * Date: 6/13/23 2:54 PM
 * Description: This is a quad_tree file.
 */
package main

import (
	"image"
	"math/rand"
	"time"
)

type Point struct {
	x int
	y int
}

// IsLine 判断两个点是否是在一条水平线或者是垂直线上面
func IsLine(point1, point2 Point) bool {
	if point1.x == point2.x || point1.y == point2.y {
		return true
	}
	return false
}

type Rect struct {
	min Point
	max Point
}

func (rect Rect) IsPointInRect(point Point) bool {
	return (point.x >= rect.min.x && point.y >= rect.min.y) && (point.x < rect.max.x && point.y < rect.max.y)
}

func (rect Rect) IsRectInRect(inRect Rect) bool {
	return (inRect.min.x >= rect.min.x && inRect.min.y >= rect.min.y) && (inRect.max.x <= rect.max.x && inRect.max.y <= rect.max.y)
}
func (rect Rect) ToImageRect() image.Rectangle {
	return image.Rect(rect.min.x, rect.min.y, rect.max.x, rect.max.y)
}

// MergeRect 判断两个Rect是否可以合并
func (rect *Rect) MergeRect(rect2 Rect) bool {
	//判断是否有一条边事完全重合的
	if rect.min.y == rect2.min.y && rect.max.x == rect2.min.x && rect.max.y == rect2.max.y {
		rect.max = rect2.max
		return true
		//return &Rect{rect.min, rect2.max}
	} else if rect2.min.y == rect.min.y && rect2.max.x == rect.min.x && rect2.max.y == rect.max.y {
		rect.min = rect2.min
		return true
		//return &Rect{rect2.min, rect.max}
	} else if rect.max.y == rect2.min.y && rect.min.x == rect2.min.x && rect.max.x == rect2.max.x {
		rect.max = rect2.max
		return true
		//return &Rect{rect.min, rect2.max}
	} else if rect2.max.y == rect.min.y && rect2.min.x == rect.min.x && rect2.max.x == rect.max.x {
		rect.min = rect2.min
		return true
		//return &Rect{rect2.min, rect.max}
	}

	return false
}

const MAXLAYER = 1

// RectNode 默认左上角角是0，0坐标点。如果项目组不是的话，那么拿到坐标以后再进行转换就好了。
// X轴代表的是length,Y轴代表的width
// 后面添加一个可重入的锁
type RectNode struct {
	//树的信息
	parent           *RectNode
	leftTopChild     *RectNode
	rightTopChild    *RectNode
	leftButtomChild  *RectNode
	RightButtomChild *RectNode

	//数据信息
	leftTop     Point
	rightButtom Point
	rect        []Rect //当前片中可以分配出去的快。

	leaf bool //只有叶子节点才会有，叶子节点中被分配出去的空间
	root bool //标记是否是根节点
}

func (rectNode *RectNode) Init(length, width int) {
	if !rectNode.root {
		return
	}

	rectNode.leftTop = Point{0, 0}
	rectNode.rightButtom = Point{length, width}

	rectNode.CreateChild(1)
}

func (rectNode *RectNode) CreateChild(layer int) {
	//当前方块已经被处理了所以就不在进行分配了。
	if rectNode.leftTopChild != nil || rectNode.leftButtomChild != nil || rectNode.rightTopChild != nil || rectNode.RightButtomChild != nil {
		return
	}
	if layer > MAXLAYER {
		return
	}

	var mid = Point{(rectNode.rightButtom.x + rectNode.leftTop.x) / 2, (rectNode.rightButtom.y + rectNode.leftTop.y) / 2}

	leftTop := new(RectNode)
	leftTop.leftTop = rectNode.leftTop
	leftTop.rightButtom = mid
	leftTop.parent = rectNode

	leftButtom := new(RectNode)
	leftButtom.leftTop = Point{rectNode.leftTop.x, mid.y}
	leftButtom.rightButtom = Point{mid.x, rectNode.rightButtom.y}
	leftButtom.parent = rectNode

	rightTop := new(RectNode)
	rightTop.leftTop = Point{mid.x, rectNode.leftTop.y}
	rightTop.rightButtom = Point{rectNode.rightButtom.x, mid.y}
	rightTop.parent = rectNode

	rightButtom := new(RectNode)
	rightButtom.leftTop = mid
	rightButtom.rightButtom = rectNode.rightButtom
	rightButtom.parent = rectNode

	rectNode.leftTopChild = leftTop
	rectNode.leftButtomChild = leftButtom
	rectNode.rightTopChild = rightTop
	rectNode.RightButtomChild = rightButtom

	if layer == MAXLAYER {
		leftTop.leaf = true
		leftButtom.leaf = true
		rightTop.leaf = true
		rightButtom.leaf = true

		leftTop.rect = make([]Rect, 0, 5)
		leftTop.rect = append(leftTop.rect, Rect{leftTop.leftTop, leftTop.rightButtom})

		leftButtom.rect = make([]Rect, 0, 5)
		leftButtom.rect = append(leftButtom.rect, Rect{leftButtom.leftTop, leftButtom.rightButtom})

		rightTop.rect = make([]Rect, 0, 5)
		rightTop.rect = append(rightTop.rect, Rect{rightTop.leftTop, rightTop.rightButtom})

		rightButtom.rect = make([]Rect, 0, 5)
		rightButtom.rect = append(rightButtom.rect, Rect{rightButtom.leftTop, rightButtom.rightButtom})
	} else {
		leftTop.CreateChild(layer + 1)
		leftButtom.CreateChild(layer + 1)
		rightTop.CreateChild(layer + 1)
		rightButtom.CreateChild(layer + 1)
	}
}

// InsertPoint Insert 这个是插入地图中不可达点信息
func (rectNode *RectNode) InsertPoint(point Point) bool {
	var rect = Rect{min: point, max: Point{point.x + 1, point.y + 1}}
	return rectNode.InsertRect(rect)
}

// InsertRect Insert 这个是插入地图中不可达点信息
func (rectNode *RectNode) InsertRect(inRect Rect) bool {
	if rectNode.leaf {
		// 如果是叶子节点的话，那么这个Rect肯定是被包含在这个里面的，直接进行判断就可以了。
		for index, rect := range rectNode.rect {
			if rect.IsRectInRect(inRect) {
				//对这个rect进行切割，把剩下来的还放入到对应的里面去
				left := Rect{rect.min, inRect.min}
				right := Rect{Point{inRect.max.x, rect.min.y}, rect.max}
				top := Rect{Point{inRect.min.x, rect.min.y}, Point{inRect.max.x, inRect.min.y}}
				buttom := Rect{Point{inRect.min.x, inRect.max.y}, Point{inRect.max.x, rect.max.y}}
				results := append(rectNode.rect[:index], rectNode.rect[index+1:]...)
				rectNode.rect = results
				//查看是否可以merge到当前的格子里面，可以的话进行merge，不可以的话直接进行插入操作
				var leftFlag = false
				var rightFlag = false
				var topFlag = false
				var buttomFlag = false
				//判断切割的四个rect是否合法，如果不合法的话直接把标记位设置为false
				if IsLine(left.min, left.max) {
					leftFlag = true
				}
				if IsLine(right.min, right.max) {
					rightFlag = true
				}
				if IsLine(top.min, top.max) {
					topFlag = true
				}
				if IsLine(buttom.min, buttom.max) {
					buttomFlag = true
				}
				for inner, result := range results {
					if leftFlag && rightFlag && topFlag && buttomFlag {
						break
					}
					if !leftFlag && result.MergeRect(left) {
						rectNode.rect[inner] = result
						leftFlag = true
					} else if !rightFlag && result.MergeRect(right) {
						rectNode.rect[inner] = result
						rightFlag = true
					} else if !topFlag && result.MergeRect(top) {
						rectNode.rect[inner] = result
						topFlag = true
					} else if !buttomFlag && result.MergeRect(buttom) {
						rectNode.rect[inner] = result
						buttomFlag = true
					}
				}
				if !leftFlag {
					rectNode.rect = append(rectNode.rect, left)
				}
				if !rightFlag {
					rectNode.rect = append(rectNode.rect, right)
				}
				if !topFlag {
					rectNode.rect = append(rectNode.rect, top)
				}
				if !buttomFlag {
					rectNode.rect = append(rectNode.rect, buttom)
				}
				return true
			}
		}
		return false
	} else {
		//查看这个rect是否在一个区域里面是的话直接传输进去，不是的话，那么需要进行分割处理
		innerFunc := func(min, max Point) bool {
			var rect Rect
			rect.min = min
			rect.max = max
			if rect.IsRectInRect(inRect) {
				return true
			}
			return false
		}

		//首先进行查看插入的rect是否在当前区域内，如果不在的话，直接返回去false
		if !innerFunc(rectNode.leftTop, rectNode.rightButtom) {
			return false
		}

		if innerFunc(rectNode.leftTopChild.leftTop, rectNode.leftTopChild.rightButtom) {
			return rectNode.leftTopChild.InsertRect(inRect)
		} else if innerFunc(rectNode.rightTopChild.leftTop, rectNode.rightTopChild.rightButtom) {
			return rectNode.rightTopChild.InsertRect(inRect)
		} else if innerFunc(rectNode.RightButtomChild.leftTop, rectNode.RightButtomChild.rightButtom) {
			return rectNode.RightButtomChild.InsertRect(inRect)
		} else if innerFunc(rectNode.leftButtomChild.leftTop, rectNode.leftButtomChild.rightButtom) {
			return rectNode.leftButtomChild.InsertRect(inRect)
		} else {
			//查看重合进行插入
			var leftTopRectangle = image.Rectangle{Min: image.Point{X: rectNode.leftTopChild.leftTop.x, Y: rectNode.leftTopChild.leftTop.y}, Max: image.Point{X: rectNode.leftTopChild.rightButtom.x, Y: rectNode.leftTopChild.rightButtom.y}}
			leftTopResult := leftTopRectangle.Intersect(inRect.ToImageRect())
			leftTopflag := false

			var rightTopRectangle = image.Rectangle{Min: image.Point{X: rectNode.rightTopChild.leftTop.x, Y: rectNode.rightTopChild.leftTop.y}, Max: image.Point{X: rectNode.rightTopChild.rightButtom.x, Y: rectNode.rightTopChild.rightButtom.y}}
			rightTopResult := rightTopRectangle.Intersect(inRect.ToImageRect())
			rightTopflag := false

			var rightButtomRectangle = image.Rectangle{Min: image.Point{X: rectNode.RightButtomChild.leftTop.x, Y: rectNode.RightButtomChild.leftTop.y}, Max: image.Point{X: rectNode.RightButtomChild.rightButtom.x, Y: rectNode.RightButtomChild.rightButtom.y}}
			rightButtomResult := rightButtomRectangle.Intersect(inRect.ToImageRect())
			rightButtomFlag := false

			var leftButtomRectangle = image.Rectangle{Min: image.Point{X: rectNode.leftButtomChild.leftTop.x, Y: rectNode.leftButtomChild.leftTop.y}, Max: image.Point{X: rectNode.leftButtomChild.rightButtom.x, Y: rectNode.leftButtomChild.rightButtom.y}}
			leftButtomResult := leftButtomRectangle.Intersect(inRect.ToImageRect())
			leftButtomFlag := false

			if !leftTopResult.Empty() {
				leftTopflag = rectNode.leftTopChild.InsertRect(Rect{Point{leftTopResult.Min.X, leftTopResult.Min.Y}, Point{leftTopResult.Max.X, leftTopResult.Max.Y}})
			} else {
				leftTopflag = true
			}
			if !leftTopflag {
				goto resultLabel
			}

			if !rightTopResult.Empty() {
				rightTopflag = rectNode.rightTopChild.InsertRect(Rect{Point{rightTopResult.Min.X, rightTopResult.Min.Y}, Point{rightTopResult.Max.X, rightTopResult.Max.Y}})
			} else {
				rightTopflag = true
			}
			if !rightTopflag {
				//需要把前面已经插入的进行回退
				if !leftTopResult.Empty() {
					rectNode.leftTopChild.FreeRect(Rect{Point{leftTopResult.Min.X, leftTopResult.Min.Y}, Point{leftTopResult.Max.X, leftTopResult.Max.Y}})
				}
				goto resultLabel
			}

			if !rightButtomResult.Empty() {
				rightButtomFlag = rectNode.RightButtomChild.InsertRect(Rect{Point{rightButtomResult.Min.X, rightButtomResult.Min.Y}, Point{rightButtomResult.Max.X, rightButtomResult.Max.Y}})
			} else {
				rightButtomFlag = true
			}
			if !rightButtomFlag {
				//需要把前面已经插入的进行回退
				if !leftTopResult.Empty() {
					rectNode.leftTopChild.FreeRect(Rect{Point{leftTopResult.Min.X, leftTopResult.Min.Y}, Point{leftTopResult.Max.X, leftTopResult.Max.Y}})
				}
				if !rightTopResult.Empty() {
					rectNode.rightTopChild.FreeRect(Rect{Point{rightTopResult.Min.X, rightTopResult.Min.Y}, Point{rightTopResult.Max.X, rightTopResult.Max.Y}})
				}
				goto resultLabel
			}

			if !leftButtomResult.Empty() {
				leftButtomFlag = rectNode.leftButtomChild.InsertRect(Rect{Point{leftButtomResult.Min.X, leftButtomResult.Min.Y}, Point{leftButtomResult.Max.X, leftButtomResult.Max.Y}})
			} else {
				leftButtomFlag = true
			}
			if !leftButtomFlag {
				if !leftTopResult.Empty() {
					rectNode.leftTopChild.FreeRect(Rect{Point{leftTopResult.Min.X, leftTopResult.Min.Y}, Point{leftTopResult.Max.X, leftTopResult.Max.Y}})
				}
				if !rightTopResult.Empty() {
					rectNode.rightTopChild.FreeRect(Rect{Point{rightTopResult.Min.X, rightTopResult.Min.Y}, Point{rightTopResult.Max.X, rightTopResult.Max.Y}})
				}
				if !rightButtomResult.Empty() {
					rectNode.RightButtomChild.FreeRect(Rect{Point{rightButtomResult.Min.X, rightButtomResult.Min.Y}, Point{rightButtomResult.Max.X, rightButtomResult.Max.Y}})
				}
				goto resultLabel
			}
		resultLabel:
			return leftTopflag && rightTopflag && rightButtomFlag && leftButtomFlag
		}
	}
	return false
}

func (rectNode *RectNode) FreeRect(inRect Rect) bool {
	//插入到当前的格子里面，这里是否需要校验一下插入的是否合法呢？比如本身就没有分配出去，但是现在又还了回来，我这里感觉是需要操作的。
	if rectNode.leaf {
		// 查看两个矩形是否有重合，如果存在重合的话，那么就会直接返回错误
		for _, rect := range rectNode.rect {
			if !rect.ToImageRect().Intersect(inRect.ToImageRect()).Empty() {
				return false
			}
		}
		//进行合并
		result := rectNode.rect
		for index, rect := range rectNode.rect {
			if rect.MergeRect(inRect) {
				result[index] = rect
				return true
			}
		}
		//无法合并，直接插入
		rectNode.rect = append(rectNode.rect, inRect)
		return true
	} else {
		//查看这个rect是否在一个区域里面是的话直接传输进去，不是的话，那么需要进行分割处理
		innerFunc := func(min, max Point) bool {
			var rect Rect
			rect.min = min
			rect.max = max
			if rect.IsRectInRect(inRect) {
				return true
			}
			return false
		}

		//首先进行查看插入的rect是否在当前区域内，如果不在的话，直接返回去false
		if !innerFunc(rectNode.leftTop, rectNode.rightButtom) {
			return false
		}

		if innerFunc(rectNode.leftTopChild.leftTop, rectNode.leftTopChild.rightButtom) {
			return rectNode.leftTopChild.FreeRect(inRect)
		} else if innerFunc(rectNode.rightTopChild.leftTop, rectNode.rightTopChild.rightButtom) {
			return rectNode.rightTopChild.FreeRect(inRect)
		} else if innerFunc(rectNode.RightButtomChild.leftTop, rectNode.RightButtomChild.rightButtom) {
			return rectNode.RightButtomChild.FreeRect(inRect)
		} else if innerFunc(rectNode.leftButtomChild.leftTop, rectNode.leftButtomChild.rightButtom) {
			return rectNode.leftButtomChild.FreeRect(inRect)
		} else {
			var leftTopRectangle = image.Rectangle{Min: image.Point{X: rectNode.leftTopChild.leftTop.x, Y: rectNode.leftTopChild.leftTop.y}, Max: image.Point{X: rectNode.leftTopChild.rightButtom.x, Y: rectNode.leftTopChild.rightButtom.y}}
			leftTopResult := leftTopRectangle.Intersect(inRect.ToImageRect())
			leftTopflag := false

			var rightTopRectangle = image.Rectangle{Min: image.Point{X: rectNode.rightTopChild.leftTop.x, Y: rectNode.rightTopChild.leftTop.y}, Max: image.Point{X: rectNode.rightTopChild.rightButtom.x, Y: rectNode.rightTopChild.rightButtom.y}}
			rightTopResult := rightTopRectangle.Intersect(inRect.ToImageRect())
			rightTopflag := false

			var rightButtomRectangle = image.Rectangle{Min: image.Point{X: rectNode.RightButtomChild.leftTop.x, Y: rectNode.RightButtomChild.leftTop.y}, Max: image.Point{X: rectNode.RightButtomChild.rightButtom.x, Y: rectNode.RightButtomChild.rightButtom.y}}
			rightButtomResult := rightButtomRectangle.Intersect(inRect.ToImageRect())
			rightButtomFlag := false

			var leftButtomRectangle = image.Rectangle{Min: image.Point{X: rectNode.leftButtomChild.leftTop.x, Y: rectNode.leftButtomChild.leftTop.y}, Max: image.Point{X: rectNode.leftButtomChild.rightButtom.x, Y: rectNode.leftButtomChild.rightButtom.y}}
			leftButtomResult := leftButtomRectangle.Intersect(inRect.ToImageRect())
			leftButtomFlag := false

			if !leftTopResult.Empty() {
				leftTopflag = rectNode.leftTopChild.FreeRect(Rect{Point{leftTopResult.Min.X, leftTopResult.Min.Y}, Point{leftTopResult.Max.X, leftTopResult.Max.Y}})
			} else {
				leftTopflag = true
			}
			if !leftTopflag {
				goto resultLabel
			}

			if !rightTopResult.Empty() {
				rightTopflag = rectNode.rightTopChild.FreeRect(Rect{Point{rightTopResult.Min.X, rightTopResult.Min.Y}, Point{rightTopResult.Max.X, rightTopResult.Max.Y}})
			} else {
				rightTopflag = true
			}
			if !rightTopflag {
				//需要把前面已经插入的进行回退
				if !leftTopResult.Empty() {
					rectNode.leftTopChild.InsertRect(Rect{Point{leftTopResult.Min.X, leftTopResult.Min.Y}, Point{leftTopResult.Max.X, leftTopResult.Max.Y}})
				}
				goto resultLabel
			}

			if !rightButtomResult.Empty() {
				rightButtomFlag = rectNode.RightButtomChild.FreeRect(Rect{Point{rightButtomResult.Min.X, rightButtomResult.Min.Y}, Point{rightButtomResult.Max.X, rightButtomResult.Max.Y}})
			} else {
				rightButtomFlag = true
			}
			if !rightButtomFlag {
				//需要把前面已经插入的进行回退
				if !leftTopResult.Empty() {
					rectNode.leftTopChild.InsertRect(Rect{Point{leftTopResult.Min.X, leftTopResult.Min.Y}, Point{leftTopResult.Max.X, leftTopResult.Max.Y}})
				}
				if !rightTopResult.Empty() {
					rectNode.rightTopChild.InsertRect(Rect{Point{rightTopResult.Min.X, rightTopResult.Min.Y}, Point{rightTopResult.Max.X, rightTopResult.Max.Y}})
				}
				goto resultLabel
			}

			if !leftButtomResult.Empty() {
				leftButtomFlag = rectNode.leftButtomChild.FreeRect(Rect{Point{leftButtomResult.Min.X, leftButtomResult.Min.Y}, Point{leftButtomResult.Max.X, leftButtomResult.Max.Y}})
			} else {
				leftButtomFlag = true
			}
			if !leftButtomFlag {
				if !leftTopResult.Empty() {
					rectNode.leftTopChild.InsertRect(Rect{Point{leftTopResult.Min.X, leftTopResult.Min.Y}, Point{leftTopResult.Max.X, leftTopResult.Max.Y}})
				}
				if !rightTopResult.Empty() {
					rectNode.rightTopChild.InsertRect(Rect{Point{rightTopResult.Min.X, rightTopResult.Min.Y}, Point{rightTopResult.Max.X, rightTopResult.Max.Y}})
				}
				if !rightButtomResult.Empty() {
					rectNode.RightButtomChild.InsertRect(Rect{Point{rightButtomResult.Min.X, rightButtomResult.Min.Y}, Point{rightButtomResult.Max.X, rightButtomResult.Max.Y}})
				}
				goto resultLabel
			}
		resultLabel:
			return leftTopflag && rightTopflag && rightButtomFlag && leftButtomFlag
		}
		//查看重合进行插入
	}
	return false
}

// Get 获取指定长宽的矩形块
func (rectNode *RectNode) Get(length, width int) *Rect {
	if rectNode.leaf {
		for index, rect := range rectNode.rect {
			if rect.max.x-rect.min.x >= length && rect.max.y-rect.min.y >= width {
				//那么就从左上角拿出来一块区域放置
				//对这个rect进行切割，剩下来的区域分割成两个部分分别是右边的一个AABB，以及下面的一个AABB
				inRect := &Rect{Point{rect.min.x, rect.min.y}, Point{rect.min.x + length, rect.min.y + width}}
				right := Rect{Point{rect.min.x + length, rect.min.y}, rect.max}
				buttom := Rect{Point{rect.min.x, rect.min.y + width}, Point{rect.min.x + length, rect.max.y}}

				results := append(rectNode.rect[:index], rectNode.rect[index+1:]...)
				rectNode.rect = results
				//查看是否可以merge到当前的格子里面，可以的话进行merge，不可以的话直接进行插入操作
				var rightFlag = false
				var buttomFlag = false
				//判断切割的四个rect是否合法，如果不合法的话直接把标记位设置为false
				if IsLine(right.min, right.max) {
					rightFlag = true
				}
				if IsLine(buttom.min, buttom.max) {
					buttomFlag = true
				}
				for inner, result := range results {
					if rightFlag && buttomFlag {
						break
					}
					if !rightFlag && result.MergeRect(right) {
						rectNode.rect[inner] = result
						rightFlag = true
					} else if !buttomFlag && result.MergeRect(buttom) {
						rectNode.rect[inner] = result
						buttomFlag = true
					}
				}
				if !rightFlag {
					rectNode.rect = append(rectNode.rect, right)
				}
				if !buttomFlag {
					rectNode.rect = append(rectNode.rect, buttom)
				}
				return inRect
			}
		}
	} else {
		//在这里随机找一个自己下面的子树去进行查找。
		rand.Seed(time.Now().UnixNano())
		randomFlag := rand.Intn(4)

		for index := 0; index < 4; index++ {
			var result *Rect
			switch randomFlag {
			case 0:
				result = rectNode.leftTopChild.Get(length, width)
			case 1:
				result = rectNode.rightTopChild.Get(length, width)
			case 2:
				result = rectNode.RightButtomChild.Get(length, width)
			case 3:
				result = rectNode.leftButtomChild.Get(length, width)
			}
			if result != nil {
				return result
			}
			randomFlag = (randomFlag + 1) % 4
		}
	}
	return nil
}
