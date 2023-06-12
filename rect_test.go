/**
 * Title: RectTest
 * Author: Anyazhou
 * Date: 6/9/23 2:59 PM
 * Description: This is a RectTest file.
 */
package main

import (
	"fmt"
	"testing"
)

func TestRectGet(t *testing.T) {
	var root RectNode
	root.root = true
	root.Init(100, 100)
	flag := root.Get(10, 10)
	if flag == nil {
		t.Error("Get error")
	} else {
		fmt.Printf("%v\n", flag)
	}
	flag = root.Get(50, 50)
	if flag == nil {
		t.Error("Get error")
	} else {
		fmt.Printf("%v\n", flag)
	}
	flag = root.Get(50, 50)
	if flag == nil {
		t.Error("Get error")
	} else {
		fmt.Printf("%v\n", flag)
	}
	flag = root.Get(50, 50)
	if flag == nil {
		t.Error("Get error")
	} else {
		fmt.Printf("%v\n", flag)
	}
	flag = root.Get(50, 50)
	if flag != nil {
		t.Error("Get error")
		fmt.Printf("%v\n", flag)
	}
}

func TestRectInsert(t *testing.T) {
	var root RectNode
	root.root = true
	root.Init(100, 100)
	flag := root.InsertPoint(Point{10, 10})
	if flag == false {
		t.Error("Insert (10, 10) error")
	}
	flag = root.InsertPoint(Point{0, 0})
	if flag == false {
		t.Error("Insert (0, 0) error")
	}

	flag = root.InsertPoint(Point{100, 100})
	if flag == true {
		t.Error("Insert (100, 100) error")
	}
}

func TestRectFree(t *testing.T) {
	var root RectNode
	root.root = true
	root.Init(100, 100)
	flag := root.InsertPoint(Point{10, 10})
	if flag == false {
		t.Error("Insert (10, 10) error")
	}
	flag = root.InsertPoint(Point{0, 0})
	if flag == false {
		t.Error("Insert (0, 0) error")
	}

	flag = root.FreeRect(Rect{Point{0, 0}, Point{1, 1}})
	if flag == false {
		t.Error("FreeRect (0,0)-(1,1) error")
	}
	flag = root.FreeRect(Rect{Point{10, 10}, Point{11, 11}})
	if flag == false {
		t.Error("FreeRect (10,10)-(11,11) error")
	}
	flag = root.FreeRect(Rect{Point{100, 100}, Point{101, 101}})
	if flag == true {
		t.Error("FreeRect (100,100)-(101,101) error")
	}
}
