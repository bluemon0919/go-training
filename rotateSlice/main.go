/*
right rotate sample
*/

package main

import (
	"fmt"
)

func main() {
	s := []int{1, 2, 3, 4, 5}
	fmt.Println(s)

	rotate3(s, 2)
	fmt.Println(s)
}

func rotate1(s []int, r int) {
	r = r % len(s)
	tmp := make([]int, len(s)-r)
	copy(tmp, s[:len(s)-r])
	copy(s, s[len(s)-r:])
	copy(s[r:], tmp)
}

func rotate2(s []int, r int) {
	r = r % len(s)
	tmp := []int{}
	tmp = append(tmp, s[len(s)-r:]...)
	tmp = append(tmp, s[:len(s)-r]...)
	copy(s, tmp)
}

// プログラミング言語Goの4.2章のサンプルを右シフトに改造したもの
func rotate3(s []int, r int) {
	reverse(s[:len(s)-r])
	reverse(s[len(s)-r:])
	reverse(s)
}

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
