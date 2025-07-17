package main

import "fmt"

/*
mapはkey-valueのペアを管理するデータ構造
keyは重複しない
*/

func mapLesson() {
	prices := map[string]int{
		"apple": 100,
		"banana": 200,
		"cherry": 300,
	}

	fmt.Println(prices)

	// mapはtrue/falseで要素の存在を確認することができる
	_, ok := prices["apple"]
	fmt.Println(ok)
}

/*
事前に配列の要素数がわかっているのであれば、 make関数を使用して指定の要素数の配列を作成する
そうすることで、事前にメモリを確保しているため配列の作成が高速になる
*/
func arrayMakeLesson() {
	arr := make([]int, 10)
	fmt.Println(arr)
}

/*
pointerはメモリアドレスを指す
*/
func pointerLesson() {
	num := 10
	pointer := &num
	fmt.Println(pointer)
	fmt.Println(*pointer)
}


/*
structは複数のデータ型をまとめて管理するデータ構造
*/

type User struct {
	Name string
	Age int
	showName func() string
}

func structLesson() {
	user := User{
		Name: "John",
		Age: 20,
		showName: func() string {
			return "John"
		},
	}
	fmt.Println(user)
	fmt.Println(user.showName())
}

func main() {
	arrayMakeLesson()
	mapLesson()
	pointerLesson()
	structLesson()
}

