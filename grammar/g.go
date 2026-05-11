package grammar

import (
	"fmt"
	"regexp"
	"time"
)

var (
	i int    = 10
	s string = "Go"
)

// 大文字はエクスポートされる
const Pi = 3.14

const (
	Username = "admin"
	Password = "password"
)

func main() {
	// LIFO: スタッキング
	defer fmt.Println("defer World1") // main関数の最後に実行される
	defer fmt.Println("defer World2") // main関数の最後に実行される
	defer fmt.Println("defer World3") // main関数の最後に実行される
	fmt.Println("Hello", time.Now())

	fmt.Printf("i: %d, s: %s\n", i, s)

	xi := 20
	xi = 30
	xt, xf := true, 3.14

	fmt.Printf("%T\n", xf) // 型定義を出力

	fmt.Printf("xi: %d, xt: %t, xf: %.2f\n", xi, xt, xf)

	fmt.Println("Hello" + "World")

	// キャスト
	var x int = 100
	xx := float64(x)
	fmt.Printf("%T %v", xx, xx) // 型定義Tと値v

	var a [2]int
	a[0] = 100
	b := [2]int{1000, 2000}
	fmt.Println(b)

	// スライス
	c := []int{1, 2, 3}
	c = append(c, 4)
	fmt.Printf("c: %v\n", c)
	fmt.Println(c[2:3])
	fmt.Println(c[1:])
	fmt.Println(c[:])
	d := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Printf("d: %v\n", d)

	// 初期値が決まっている:スライス
	// サイズだけ決まっている：makeで長さ指定
	// 追記量が予測できる：makeで容量指定
	// 両方決まっている：makeで長さと容量指定
	n := make([]int, 3, 5) // 長さ3、容量5のスライスを作成
	fmt.Printf("len=%d cap=%d value=%v", len(n), cap(n), n)
	n = append(n, 1, 2)
	fmt.Printf("len=%d cap=%d value=%v", len(n), cap(n), n)

	// c2 := make([]int, 5)
	c2 := make([]int, 0, 5)

	for i := 0; i < 5; i++ {
		c = append(c2, i)
		fmt.Println(c)
	}
	fmt.Println(c2)

	// マップ
	m := map[string]int{"apple": 5, "banana": 10}
	m2 := map[string]int{"Mike": 20, "Nancy": 24, "Messi": 30}
	fmt.Println(m2)
	fmt.Println(m["apple"])
	m["orange"] = 15
	fmt.Println(m)
	v, ok := m["banana"]
	fmt.Printf("Value: %d, Exists: %t\n", v, ok)
	delete(m, "banana")

	r1, r2 := add(1, 2)
	fmt.Printf("r1: %d, r2: %d\n", r1, r2) // 10進数

	s := []int{1, 2, 3}
	sum(s...) // 展開して渡す

	s2 := []int{1, 2, 5, 6, 2, 3, 1}
	fmt.Println(s2[2:4]) // 5,6

	num := sum(1, 2, 3, 4, 5)
	if num%2 == 0 {
		fmt.Printf("%d is even\n", num)
	} else {
		fmt.Printf("%d is odd\n", num)
	}

	// 定義とifを同時に定義
	if r := sum(1, 2, 3); r > 5 {
		fmt.Printf("Sum is greater than 5: %d\n", r)
	}

	for i := 0; i < 5; i++ {
		if i == 3 {
			continue // 3のときはスキップ
		}
		if i == 4 {
			break // 4のときはループ終了
		}
	}

	// レンジ
	l := []string{"Go", "Python", "Java"}
	for index, value := range l {
		fmt.Printf("Index: %d, Value: %s\n", index, value)
	}
	m3 := map[string]int{"Alice": 30, "Bob": 25, "Charlie": 35}
	for key, value := range m3 {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}

	// スイッチ
	switch num := sum(1, 2, 3); {
	case num < 10:
		fmt.Printf("Sum is less than 10: %d\n", num)
	case num == 10:
		fmt.Printf("Sum is equal to 10: %d\n", num)
	default:
		fmt.Printf("Sum is greater than 10: %d\n", num)
	}

	// パニックとリカバー（基本的にパニックをあえて書くことはない）
	defer func() {
		// recoverを使うことでパニックのハンドリングができる
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()
	// panic("Something went wrong!")

	ll2 := []int{54, 2, 23, 11, 4, 5, 65}
	var current int
	for index, value := range ll2 {
		if index == 0 {
			current = value
			continue
		}
		if current > value {
			current = value
		}
	}
	fmt.Println(current)

	mm2 := map[string]int{
		"apple":  200,
		"banana": 300,
		"grapes": 150,
		"orange": 80,
		"papaya": 500,
		"kiwi":   90,
	}

	var res int
	for _, v := range mm2 {
		res += v
	}
	fmt.Printf("Total price: %d\n", res)

	// ポインタ
	v2 := 100       // このときメモリ上のどこかにv2の値が格納され、さらにどこに格納されているかを示すアドレスが確保される（合計で2つの番地が確保される）
	p := &v2        // アンパサンドで`値の番地`が書き込まれているアドレスを取得（値そのもののアドレスではないので注意）
	fmt.Println(*p) // 値の番地から値を取得（デリファレンス）

	// makeとnewの違い
	st := &struct{}{} // newよりも&T{}の方が慣用的
	ch := make(chan int)
	fmt.Printf("st: %v\n", st) // st: &{}
	fmt.Printf("ch: %v\n", ch) // ch: 0x1400008e120

	// 構造体
	v3 := Vertex{X: 1, Y: 2}
	v3.X = 1
	v3.Y = 2
	fmt.Printf("v3: %v\n", v3) // v3: {1 2}
	v4 := Vertex{}             // フィールドを指定しないとゼロ値になる
	fmt.Printf("v4: %v\n", v4) // v4: {0 0}
	v5 := &Vertex{}            // ポインタであるとわかりやすいためnewよりこちら側が使われやすい
	fmt.Printf("v5: %v\n", v5) // v5: &{0 0}
	v3.Area()
	v3.Scale(10) // ポインタを使う理由：構造体の値を書き換える

	main2()

	v6 := Vertex{X: 3, Y: 4}
	fmt.Println(v6.Plus())

	fmt.Println(v6)

	// 日付
	now := time.Now()
	fmt.Println(now.Format(time.RFC3339)) // PostgreSQLのtimestamp型と互換性のあるフォーマット
	// 正規表現
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`, now.Format(time.RFC3339))
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)
	re.MatchString(now.Format(time.RFC3339))
	fmt.Println(match)
}

type Vertex struct {
	X, Y int
}

// 継承みたいなやつ
type Vertex3D struct {
	Vertex
	Z int
}

// コンストラクタみたいなやつ（goにはクラスがないからこうやって書く。ちなみに現状Vertexはパブリックな構造体なのでNew関数を使用する必要ない。プロパティが小文字の場合はNew関数を作るのが慣習的）
func NewVertex3D(x, y, z int) *Vertex3D {
	return &Vertex3D{Vertex{x, y}, z}
}

// 構造体にロジックをもたせる（メソッド）
func (v Vertex) Area() int {
	return v.X * v.Y
}

// ポインタレシーバー
func (v *Vertex) Scale(i int) {
	v.X = v.X * i
	v.Y = v.Y * i
}

func add(x, y int) (int, int) {
	return x + y, x - y
}

func calc(x, y int) (sum, diff int) {
	// 戻り値を事前に定義する。（なので変数定義ではなく代入になっている。）
	sum = x + y
	diff = x - y
	return
}

// 可変長
func sum(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

func test() {
	var i int = 100 // 300
	var j int = 200
	var p1 *int
	var p2 *int // p1 // 300
	p1 = &i
	p2 = &j
	i = *p1 + *p2
	p2 = p1
	j = *p2 + i // 600
	fmt.Println(j)
}

// インターフェース
type Human interface {
	Say() string
}

type Person struct {
	Name string
}

type Dog struct {
	Name string
}

func (p *Person) Say() string {
	p.Name = "Mr." + p.Name
	return p.Name
}

func DriveCar(human Human) {
	if human.Say() == "Mr.Mike" {
		fmt.Println("Mike is driving the car.")
	} else {
		fmt.Println("Unknown human is driving the car.")
	}
}

func main2() {
	var mike Human = &Person{Name: "Mike"} // レシーバー内で初期化したPerson構造体を変更するためポインタで渡す必要がある
	DriveCar(mike)
	// var dog Dog = Dog{Name: "Buddy"}
	// DriveCar(dog) // DogはHumanインターフェースを実装していないためエラー
}

// jsのany
func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %s\n", v)
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

func main3() {
	v := interface{}(42)                 // 空のインターフェースは任意の型を受け入れることができる
	ii := v.(int)                        // 型アサーションでint型に変換(キャストとほぼ同じでinterfaceの場合の書き方)
	fmt.Printf("v: %v, ii: %d\n", v, ii) // v: 42, ii: 42
	do(1)
	do("Hello")
	do(3.14)
}

func (v Vertex) Plus() int {
	return v.X + v.Y
}

func (v Vertex) String() string {
	return fmt.Sprintf("X is %d! Y is %d", v.X, v.Y)
}
