/**
 * Title: main
 * Author: Anyazhou
 * Date: 5/23/23 3:31 PM
 * Description: This is a main file.
 */

package main

type Nums struct {
	num int
}

func (num *Nums) Add(x int) {
	num.num += x
	println(num.num)
}

func main() {

}
