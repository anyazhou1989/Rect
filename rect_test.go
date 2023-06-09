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
