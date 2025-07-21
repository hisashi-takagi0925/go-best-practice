# Go 言語 ポインタの概念と使い道 学習メモ

## 目次

1. [ポインタの基本概念](#ポインタの基本概念)
2. [ポインタの宣言と初期化](#ポインタの宣言と初期化)
3. [ポインタの操作](#ポインタの操作)
4. [ポインタの使い道](#ポインタの使い道)
5. [ポインタとスライス](#ポインタとスライス)
6. [ポインタと構造体](#ポインタと構造体)
7. [ポインタの注意点](#ポインタの注意点)
8. [実用的な例](#実用的な例)

## ポインタの基本概念

### ポインタとは

**ポインタは、メモリ上の特定の場所（アドレス）を指し示す変数です。**

```
メモリ上のイメージ:
┌─────────┬─────────┬─────────┬─────────┐
│ アドレス │   0x1000  │   0x1001  │   0x1002  │   0x1003  │
├─────────┼─────────┼─────────┼─────────┤
│   値    │    42    │    0    │    0    │    0    │
└─────────┴─────────┴─────────┴─────────┘
                    ↑
                ポインタが指している場所
```

### 変数とポインタの関係

```go
var x int = 42        // 変数xに値42を格納
var p *int = &x       // ポインタpに変数xのアドレスを格納

fmt.Printf("xの値: %d\n", x)           // 42
fmt.Printf("xのアドレス: %p\n", &x)    // 0xc000018030
fmt.Printf("pの値: %p\n", p)           // 0xc000018030
fmt.Printf("pが指す値: %d\n", *p)      // 42
```

## ポインタの宣言と初期化

### 基本的な宣言

```go
// 1. ゼロ値で初期化（nil）
var p1 *int

// 2. 変数のアドレスで初期化
var x int = 42
var p2 *int = &x

// 3. 短縮記法
p3 := &x

// 4. new()関数で初期化
p4 := new(int)  // ゼロ値で初期化されたint型のポインタ
```

### ポインタの型

```go
// 異なる型のポインタ
var p1 *int        // int型のポインタ
var p2 *string     // string型のポインタ
var p3 *bool       // bool型のポインタ
var p4 *[]int      // intスライスのポインタ
var p5 *struct{}   // 構造体のポインタ
```

## ポインタの操作

### アドレス演算子（&）と間接参照演算子（\*）

```go
var x int = 42

// &演算子: 変数のアドレスを取得
p := &x
fmt.Printf("xのアドレス: %p\n", p)

// *演算子: ポインタが指す値を取得
value := *p
fmt.Printf("pが指す値: %d\n", value)

// *演算子: ポインタが指す値を変更
*p = 100
fmt.Printf("変更後のx: %d\n", x)  // 100
```

### ポインタの比較

```go
var x int = 42
p1 := &x
p2 := &x

// 同じアドレスを指しているかチェック
fmt.Printf("p1 == p2: %t\n", p1 == p2)  // true

// nilとの比較
var p3 *int
fmt.Printf("p3 == nil: %t\n", p3 == nil)  // true
```

## ポインタの使い道

### 1. 効率的な関数呼び出し

#### 値渡し vs ポインタ渡し

```go
// 値渡し（コピーが作成される）
func incrementByValue(x int) {
    x++  // コピーが変更されるだけ
}

// ポインタ渡し（元の値が変更される）
func incrementByPointer(x *int) {
    *x++  // 元の値が変更される
}

func main() {
    value := 10

    // 値渡し
    incrementByValue(value)
    fmt.Printf("値渡し後: %d\n", value)  // 10（変更されない）

    // ポインタ渡し
    incrementByPointer(&value)
    fmt.Printf("ポインタ渡し後: %d\n", value)  // 11（変更される）
}
```

#### 大きな構造体の効率性

```go
type LargeStruct struct {
    Data [1000]int
    Name string
}

// 非効率: 構造体全体がコピーされる
func processByValue(data LargeStruct) {
    // 処理...
}

// 効率的: アドレスのみが渡される
func processByPointer(data *LargeStruct) {
    // 処理...
}
```

### 2. 複数の値を返す

```go
// 従来の方法（複数の戻り値）
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("ゼロ除算")
    }
    return a / b, nil
}

// ポインタを使った方法
func divideWithPointer(a, b int, result *int, err *error) {
    if b == 0 {
        *err = errors.New("ゼロ除算")
        return
    }
    *result = a / b
    *err = nil
}
```

### 3. オプショナルな引数

```go
type Config struct {
    Host     string
    Port     int
    Timeout  time.Duration
}

// デフォルト値を使用
func NewServer(host string, port int) *Server {
    return &Server{
        Config: Config{
            Host:     host,
            Port:     port,
            Timeout:  30 * time.Second,  // デフォルト値
        },
    }
}

// カスタム設定を指定
func NewServerWithConfig(host string, port int, config *Config) *Server {
    defaultConfig := Config{
        Host:    host,
        Port:    port,
        Timeout: 30 * time.Second,
    }

    if config != nil {
        if config.Timeout != 0 {
            defaultConfig.Timeout = config.Timeout
        }
    }

    return &Server{Config: defaultConfig}
}
```

### 4. メソッドレシーバー

```go
type Counter struct {
    count int
}

// 値レシーバー（コピーが作成される）
func (c Counter) GetCount() int {
    return c.count
}

// ポインタレシーバー（元の構造体が変更される）
func (c *Counter) Increment() {
    c.count++
}

func main() {
    counter := Counter{count: 0}

    counter.Increment()  // ポインタレシーバー
    fmt.Printf("カウント: %d\n", counter.GetCount())  // 1
}
```

## ポインタとスライス

### スライスは既にポインタ的な動作

```go
// スライスは内部的にポインタを含んでいる
func modifySlice(s []int) {
    s[0] = 100  // 元のスライスが変更される
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    fmt.Printf("変更前: %v\n", numbers)  // [1 2 3 4 5]

    modifySlice(numbers)
    fmt.Printf("変更後: %v\n", numbers)  // [100 2 3 4 5]
}
```

### スライスのポインタ（まれな使用）

```go
// スライス自体のポインタ（通常は不要）
func modifySlicePointer(s *[]int) {
    (*s)[0] = 100
    *s = append(*s, 999)  // スライスの長さも変更
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    fmt.Printf("変更前: %v\n", numbers)  // [1 2 3 4 5]

    modifySlicePointer(&numbers)
    fmt.Printf("変更後: %v\n", numbers)  // [100 2 3 4 5 999]
}
```

## ポインタと構造体

### 構造体のポインタ

```go
type Person struct {
    Name string
    Age  int
}

func main() {
    // 値として作成
    person1 := Person{Name: "田中", Age: 25}

    // ポインタとして作成
    person2 := &Person{Name: "佐藤", Age: 30}

    // ポインタを取得
    person3 := &person1

    fmt.Printf("person1: %+v\n", person1)
    fmt.Printf("person2: %+v\n", person2)
    fmt.Printf("person3: %+v\n", person3)
}
```

### 構造体のメソッド

```go
type BankAccount struct {
    Balance int
    Owner   string
}

// 値レシーバー（読み取り専用）
func (b BankAccount) GetBalance() int {
    return b.Balance
}

// ポインタレシーバー（変更可能）
func (b *BankAccount) Deposit(amount int) {
    b.Balance += amount
}

func (b *BankAccount) Withdraw(amount int) bool {
    if b.Balance >= amount {
        b.Balance -= amount
        return true
    }
    return false
}

func main() {
    account := &BankAccount{Balance: 1000, Owner: "田中"}

    fmt.Printf("残高: %d\n", account.GetBalance())  // 1000

    account.Deposit(500)
    fmt.Printf("預け入れ後: %d\n", account.GetBalance())  // 1500

    success := account.Withdraw(200)
    fmt.Printf("引き出し成功: %t, 残高: %d\n", success, account.GetBalance())
}
```

## ポインタの注意点

### 1. nil ポインタの参照

```go
var p *int  // nilポインタ

// 危険: nilポインタの参照
// fmt.Printf("値: %d\n", *p)  // パニック！

// 安全: nilチェック
if p != nil {
    fmt.Printf("値: %d\n", *p)
} else {
    fmt.Println("ポインタはnilです")
}
```

### 2. ダングリングポインタ

```go
func createPointer() *int {
    x := 42
    return &x  // 関数終了後もxは有効（Goのガベージコレクション）
}

func main() {
    p := createPointer()
    fmt.Printf("値: %d\n", *p)  // 42（安全）
}
```

### 3. ポインタの比較

```go
func main() {
    x := 42
    p1 := &x
    p2 := &x
    p3 := &x

    fmt.Printf("p1 == p2: %t\n", p1 == p2)  // true
    fmt.Printf("p1 == p3: %t\n", p1 == p3)  // true

    y := 42
    p4 := &y
    fmt.Printf("p1 == p4: %t\n", p1 == p4)  // false（異なるアドレス）
}
```

## 実用的な例

### 1. 設定管理

```go
type Config struct {
    DatabaseURL string
    Port        int
    Debug       bool
}

func NewConfig() *Config {
    return &Config{
        DatabaseURL: "localhost:5432",
        Port:        8080,
        Debug:       false,
    }
}

func (c *Config) SetDebug(debug bool) {
    c.Debug = debug
}

func (c *Config) SetPort(port int) {
    c.Port = port
}
```

### 2. ビルダーパターン

```go
type QueryBuilder struct {
    table   string
    fields  []string
    where   []string
    orderBy string
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{}
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
    qb.table = table
    return qb
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
    qb.fields = fields
    return qb
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
    qb.where = append(qb.where, condition)
    return qb
}

func (qb *QueryBuilder) OrderBy(field string) *QueryBuilder {
    qb.orderBy = field
    return qb
}

func (qb *QueryBuilder) Build() string {
    // SQLクエリを構築
    query := "SELECT "
    if len(qb.fields) > 0 {
        query += strings.Join(qb.fields, ", ")
    } else {
        query += "*"
    }
    query += " FROM " + qb.table

    if len(qb.where) > 0 {
        query += " WHERE " + strings.Join(qb.where, " AND ")
    }

    if qb.orderBy != "" {
        query += " ORDER BY " + qb.orderBy
    }

    return query
}

func main() {
    query := NewQueryBuilder().
        From("users").
        Select("id", "name", "email").
        Where("age > 18").
        Where("active = true").
        OrderBy("name").
        Build()

    fmt.Println(query)
    // SELECT id, name, email FROM users WHERE age > 18 AND active = true ORDER BY name
}
```

### 3. オブザーバーパターン

```go
type Observer interface {
    Update(data string)
}

type Subject struct {
    observers []Observer
    data      string
}

func NewSubject() *Subject {
    return &Subject{
        observers: make([]Observer, 0),
    }
}

func (s *Subject) Attach(observer Observer) {
    s.observers = append(s.observers, observer)
}

func (s *Subject) SetData(data string) {
    s.data = data
    s.notify()
}

func (s *Subject) notify() {
    for _, observer := range s.observers {
        observer.Update(s.data)
    }
}

type Logger struct {
    name string
}

func NewLogger(name string) *Logger {
    return &Logger{name: name}
}

func (l *Logger) Update(data string) {
    fmt.Printf("[%s] データが更新されました: %s\n", l.name, data)
}

func main() {
    subject := NewSubject()

    logger1 := NewLogger("Logger1")
    logger2 := NewLogger("Logger2")

    subject.Attach(logger1)
    subject.Attach(logger2)

    subject.SetData("新しいデータ")
    // [Logger1] データが更新されました: 新しいデータ
    // [Logger2] データが更新されました: 新しいデータ
}
```

## 重要なポイント

### 1. ポインタの使い分け

- **値渡し**: 小さなデータ、変更不要な場合
- **ポインタ渡し**: 大きなデータ、変更が必要な場合
- **メソッドレシーバー**: 変更が必要な場合はポインタレシーバー

### 2. 安全性

- **nil チェック**: ポインタを使用する前に nil チェック
- **初期化**: 適切な初期化を忘れない
- **スコープ**: ポインタの有効期間を意識する

### 3. パフォーマンス

- **大きな構造体**: ポインタ渡しで効率化
- **小さなデータ**: 値渡しで十分
- **スライス**: 既にポインタ的な動作

### 4. Go 言語の特徴

- **ガベージコレクション**: メモリ管理は自動
- **値セマンティクス**: デフォルトは値渡し
- **明示的なポインタ**: 必要に応じて明示的に使用

## 参考資料

- [Go 言語公式ドキュメント - Pointers](https://golang.org/doc/effective_go.html#pointers)
- [Go 言語公式ドキュメント - Methods](https://golang.org/doc/effective_go.html#methods)
- [Go 言語公式ドキュメント - Structs](https://golang.org/doc/effective_go.html#structs)
