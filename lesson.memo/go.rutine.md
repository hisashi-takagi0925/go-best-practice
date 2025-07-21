# Go 言語 同期処理・ゴルーチン・チャネル・ワーカープール 学習メモ

## 目次

1. [同期処理の基本](#同期処理の基本)
2. [ゴルーチンの基本](#ゴルーチンの基本)
3. [チャネルの概念](#チャネルの概念)
4. [双方向通信](#双方向通信)
5. [ブロッキング（待機）](#ブロッキング待機)
6. [ワーカープール](#ワーカープール)
7. [実用的な例](#実用的な例)

## 同期処理の基本

### 概要

複数のゴルーチンで共有する変数は、適切な同期処理が必要です。

### 問題のあるコード例

```go
var counter int = 0

func increment(wg *sync.WaitGroup) {
    for i := 0; i < 100; i++ {
        counter++ // 競合状態！
    }
    wg.Done()
}
```

**問題点：**

- `counter++`は実際には 3 つのステップに分解される
  1. 値を読み取り
  2. 加算
  3. 結果を書き込み
- 複数のゴルーチンが同時実行すると、期待値より小さい結果になる

## 同期処理の選択肢

### 1. sync/atomic（アトミック操作）

**適している場合：**

- 単純な数値の操作（加算、減算、読み取り、書き込み）
- パフォーマンスが重要な場合
- 複雑な操作が不要な場合

```go
var counter int64 = 0

func increment(wg *sync.WaitGroup) {
    for i := 0; i < 100; i++ {
        atomic.AddInt64(&counter, 1) // アトミックな加算
    }
    wg.Done()
}
```

**主な関数：**

- `atomic.AddInt64(&var, delta)` - アトミックな加算
- `atomic.LoadInt64(&var)` - アトミックな読み取り
- `atomic.StoreInt64(&var, value)` - アトミックな書き込み
- `atomic.CompareAndSwapInt64(&var, old, new)` - 条件付き交換
- `atomic.SwapInt64(&var, new)` - 値の交換

### 2. sync.Mutex（ミューテックス）

**適している場合：**

- 複雑な操作が必要な場合
- 複数の変数を同時に更新する場合
- 条件分岐を含む操作

```go
var counter int = 0
var mu sync.Mutex

func increment(wg *sync.WaitGroup) {
    for i := 0; i < 100; i++ {
        mu.Lock()
        counter++
        mu.Unlock()
    }
    wg.Done()
}
```

### 3. sync.RWMutex（読み書き分離ミューテックス）

**適している場合：**

- 読み取りが多く、書き込みが少ない場合
- データ構造へのアクセス

```go
var data map[string]int
var rwmu sync.RWMutex

// 読み取り
rwmu.RLock()
value := data["key"]
rwmu.RUnlock()

// 書き込み
rwmu.Lock()
data["key"] = 42
rwmu.Unlock()
```

## 複雑な操作の例

```go
var balance int = 1000
var mu sync.Mutex

// 引き出し処理（条件分岐が必要）
func withdraw(amount int) {
    mu.Lock()
    if balance >= amount {
        balance -= amount
        fmt.Printf("引き出し: 残高 %d\n", balance)
    }
    mu.Unlock()
}

// 預け入れ処理
func deposit(amount int) {
    mu.Lock()
    balance += amount
    fmt.Printf("預け入れ: 残高 %d\n", balance)
    mu.Unlock()
}
```

## 選択基準

| 用途             | 推奨方法  | 理由                  |
| ---------------- | --------- | --------------------- |
| 単純な数値操作   | `atomic`  | パフォーマンスが良い  |
| 複雑な操作       | `mutex`   | 柔軟性が高い          |
| 読み取りが多い   | `RWMutex` | 並行性が高い          |
| ゴルーチン間通信 | `channel` | Go 言語の推奨パターン |

## ゴルーチンの基本

### ゴルーチンとは

- **軽量スレッド**: メモリ使用量が少ない
- **並行実行**: 複数の処理を同時に実行
- **`go`キーワード**: ゴルーチンを起動

### 基本的な使い方

```go
// 順次実行（通常の処理）
for i := 1; i <= 3; i++ {
    fmt.Printf("処理 %d\n", i)
    time.Sleep(100 * time.Millisecond)
}

// 並行実行（ゴルーチン）
for i := 1; i <= 3; i++ {
    go func(id int) {
        fmt.Printf("ゴルーチン %d 開始\n", id)
        time.Sleep(100 * time.Millisecond)
        fmt.Printf("ゴルーチン %d 終了\n", id)
    }(i)
}
```

## チャネルの概念

### チャネルとは

- **ゴルーチン間の通信パイプ**
- **値の受け渡しの中継機**
- **自然な同期機能**

### 基本的な使い方

```go
// チャネルの作成
ch := make(chan string)

// ゴルーチンで送信
go func() {
    ch <- "Hello from goroutine!"
}()

// メインゴルーチンで受信
message := <-ch
fmt.Printf("受信した値 = %s\n", message)
```

### バッファ付きチャネル

```go
// バッファサイズ2のチャネル
ch := make(chan string, 2)

// バッファに値を送信（ブロックしない）
ch <- "First message"
ch <- "Second message"

// バッファから値を取り出す
fmt.Printf("受信1: %s\n", <-ch)
fmt.Printf("受信2: %s\n", <-ch)
```

## 双方向通信

### 双方向通信の基本

```go
ch := make(chan string)

// ゴルーチンA: 送信と受信の両方を行う
go func() {
    ch <- "Hello from A"  // 送信
    response := <-ch       // 受信
    fmt.Printf("応答を受信 = %s\n", response)
}()

// メインゴルーチン: 送信と受信の両方を行う
message := <-ch            // 受信
ch <- "Hello from Main"    // 送信
```

### 会話形式の双方向通信

```go
ch := make(chan string)

// 質問する側
go func() {
    questions := []string{"お名前は？", "年齢は？", "趣味は？"}
    for _, question := range questions {
        ch <- question
        answer := <-ch
        fmt.Printf("回答 = %s\n", answer)
    }
}()

// 回答する側
answers := []string{"田中太郎です", "25歳です", "プログラミングです"}
for _, answer := range answers {
    question := <-ch
    ch <- answer
}
```

## ブロッキング（待機）

### ブロッキングの動作

**ゴルーチン内で受信が来るまで処理は待機します**

```go
ch := make(chan string)

// ゴルーチンA: 受信待機
go func() {
    fmt.Println("受信待機開始")
    message := <-ch  // ここでブロック（待機）
    fmt.Printf("受信完了 = %s\n", message)
}()

// メインゴルーチン: 3秒後に送信
time.Sleep(3 * time.Second)
ch <- "Hello from Main"
```

### 送信側のブロッキング

```go
ch := make(chan string)

// ゴルーチンA: 送信（受信側が準備できるまで待機）
go func() {
    ch <- "Hello from Sender"  // ここでブロック
    fmt.Println("送信完了")
}()

// メインゴルーチン: 3秒後に受信
time.Sleep(3 * time.Second)
message := <-ch
```

### 非ブロッキング受信（select 文）

```go
ch := make(chan string)

select {
case message := <-ch:
    fmt.Printf("受信成功 = %s\n", message)
default:
    fmt.Println("受信待機中（ブロックしない）")
}
```

## ワーカープール

### ワーカープールとは

**工場の作業員システム**を想像してください：

1. **作業員（ワーカー）**: 複数の作業員が待機している
2. **作業台（チャネル）**: 仕事が置かれる場所
3. **仕事（ジョブ）**: 処理すべきタスク
4. **結果**: 処理が完了した成果物

### 基本構造

```go
// 1. 仕事を入れる箱
jobs := make(chan int, 5)

// 2. 結果を入れる箱
results := make(chan int, 5)

// 3. 作業員（ワーカー）
worker := func(id int) {
    for job := range jobs {
        fmt.Printf("ワーカー %d: ジョブ %d を処理中...\n", id, job)
        time.Sleep(500 * time.Millisecond)
        result := job * 2
        results <- result
    }
}

// ワーカーを3人起動
for i := 1; i <= 3; i++ {
    go worker(i)
}

// 仕事を追加
for i := 1; i <= 5; i++ {
    jobs <- i
}
close(jobs)

// 結果を収集
for i := 1; i <= 5; i++ {
    result := <-results
    fmt.Printf("結果 %d: %d\n", i, result)
}
```

### ワーカープールの利点

#### 順次処理 vs ワーカープール

```
順次処理: 仕事1 → 仕事2 → 仕事3 → 仕事4 → 仕事5 → 仕事6
         (500ms)  (500ms)  (500ms)  (500ms)  (500ms)  (500ms)
         合計: 3秒

ワーカープール: 仕事1,2,3,4,5,6 を同時に処理
                (500ms) で全て完了
                合計: 1秒
```

**速度向上: 3 倍速くなりました！**

### 実用的なワーカープール

```go
type Task struct {
    ID   int
    Name string
    Data int
}

type Result struct {
    TaskID int
    Result int
    Worker int
}

// ワーカー関数
worker := func(id int) {
    for task := range tasks {
        var result int
        switch task.Name {
        case "計算":
            result = task.Data * 2
        case "変換":
            result = task.Data + 100
        case "検証":
            result = task.Data * task.Data
        }
        results <- Result{TaskID: task.ID, Result: result, Worker: id}
    }
}
```

## 実用的な例

### 1. 基本的なゴルーチン

```go
func basicGoroutine() {
    // ゴルーチンなし（順次実行）
    for i := 1; i <= 3; i++ {
        fmt.Printf("処理 %d\n", i)
        time.Sleep(100 * time.Millisecond)
    }

    // ゴルーチンあり（並行実行）
    for i := 1; i <= 3; i++ {
        go func(id int) {
            fmt.Printf("ゴルーチン %d 開始\n", id)
            time.Sleep(100 * time.Millisecond)
            fmt.Printf("ゴルーチン %d 終了\n", id)
        }(i)
    }
}
```

### 2. チャネルを使った通信

```go
func channelCommunication() {
    ch := make(chan int)

    // 送信側ゴルーチン
    go func() {
        for i := 1; i <= 5; i++ {
            ch <- i
            time.Sleep(100 * time.Millisecond)
        }
        close(ch)
    }()

    // 受信側（メインゴルーチン）
    for value := range ch {
        fmt.Printf("受信: %d\n", value)
    }
}
```

### 3. select 文の例

```go
func selectExample() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    // ゴルーチン1: 1秒後にch1に送信
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "from channel 1"
    }()

    // ゴルーチン2: 2秒後にch2に送信
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "from channel 2"
    }()

    // select文で最初に準備できたチャネルを処理
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Printf("%s\n", msg1)
        case msg2 := <-ch2:
            fmt.Printf("%s\n", msg2)
        case <-time.After(3 * time.Second):
            fmt.Println("タイムアウト")
        }
    }
}
```

## 重要なポイント

### 1. 同期処理の選択

- **単純な数値操作**: `atomic`
- **複雑な操作**: `mutex`
- **読み取りが多い**: `RWMutex`
- **ゴルーチン間通信**: `channel`

### 2. チャネルの特徴

- **ブロッキング**: 送信/受信が準備できるまで待機
- **自然な同期**: 明示的な同期処理が不要
- **双方向通信**: 1 つのチャネルで送受信可能

### 3. ワーカープールの利点

- **並行処理**: 複数のワーカーが同時に処理
- **効率性**: リソースの有効活用
- **スケーラビリティ**: ワーカー数を調整可能

### 4. デッドロックの回避

- **チャネルの閉じ忘れ**: `close(ch)`を忘れない
- **送信側と受信側の不整合**: 適切な同期を確保
- **select 文の活用**: 非ブロッキング操作

### 5. 共有変数の注意点

- **共有変数は必ず同期処理が必要**
- **atomic は単純な操作に最適**
- **mutex は複雑な操作に必要**
- **パフォーマンスと安全性のバランスを考慮**
- **Go 言語では「共有による通信」より「通信による共有」が推奨**

## 参考資料

- [Go 言語公式ドキュメント - Goroutines](https://golang.org/doc/effective_go.html#goroutines)
- [Go 言語公式ドキュメント - Channels](https://golang.org/doc/effective_go.html#channels)
- [Go 言語公式ドキュメント - sync/atomic](https://golang.org/pkg/sync/atomic/)
- [Go 言語公式ドキュメント - sync](https://golang.org/pkg/sync/)
