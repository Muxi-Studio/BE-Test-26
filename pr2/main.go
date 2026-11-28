package main

import (
	"fmt"
)

var fa []int

func find() {
	//找有几组
}

func main() {
	var n, m, op, u, v int
	var g int

	fmt.Scan(&m, &n)

	for i := 0; i < m; i++ {
		fmt.Scan(&op, &u, &v)
		if u == v && op == 2 {
			fmt.Println("0")
		}
		if op == 1 { // 相同

		} else { // 不同

		}
	}

	//最终结果为2的组数次方对998244353取余
	fmt.Println(2 ^ g%998244353)

}
