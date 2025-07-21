# クリーンアーキテクチャ設計書

## 概要

このプロジェクトは、Robert C. Martin (Uncle Bob) が提唱するクリーンアーキテクチャの原則に従って設計されています。ビジネスロジックを中心に据え、外部の技術的詳細から独立した構造を実現しています。

## アーキテクチャの原則

### 1. 依存性逆転の原則 (DIP)
- 内側の層は外側の層に依存しない
- 外側の層は内側の層に依存する
- インターフェースを通じて依存関係を逆転させる

### 2. 関心の分離
- 各層は明確な責任を持つ
- ビジネスロジックはドメイン層に集約
- 技術的な詳細は外側の層に配置

## レイヤー構造

```
┌─────────────────────────────────────────────────┐
│                Infrastructure                    │
│  (HTTPサーバー、外部API、データベース等)           │
├─────────────────────────────────────────────────┤
│                  Interface                       │
│  (ハンドラー、DTO、ゲートウェイ実装)              │
├─────────────────────────────────────────────────┤
│                  Use Case                        │
│  (アプリケーション固有のビジネスルール)            │
├─────────────────────────────────────────────────┤
│                   Domain                         │
│  (エンティティ、値オブジェクト、リポジトリIF)      │
└─────────────────────────────────────────────────┘
```

## ディレクトリ構造

```
web-api/
├── cmd/
│   └── api/
│       └── main.go              # アプリケーションのエントリーポイント
├── internal/
│   ├── domain/                  # ドメイン層
│   │   ├── entity/              # エンティティ
│   │   ├── valueobject/         # 値オブジェクト
│   │   └── repository/          # リポジトリインターフェース
│   ├── usecase/                 # ユースケース層
│   │   ├── post/                # Post関連のビジネスロジック
│   │   └── user/                # User関連のビジネスロジック
│   ├── interface/               # インターフェース層
│   │   ├── api/                 # Web API関連
│   │   │   ├── handler/         # HTTPハンドラー
│   │   │   ├── dto/             # データ転送オブジェクト
│   │   │   └── router/          # ルーティング
│   │   └── gateway/             # 外部サービスゲートウェイ
│   └── infrastructure/          # インフラストラクチャ層
│       ├── http/                # HTTPクライアント
│       └── server/              # HTTPサーバー
└── config/                      # 設定管理
```

## 各層の詳細

### 1. Domain層（ドメイン層）

最も内側の層で、ビジネスの核となる概念を表現します。

#### Entity（エンティティ）
```go
type Post struct {
    id     valueobject.PostID
    userID valueobject.UserID
    title  string
    body   string
}
```
- ビジネスの中核となるオブジェクト
- ビジネスルールをカプセル化
- 外部の技術的詳細に依存しない

#### Value Object（値オブジェクト）
```go
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    // バリデーションロジック
}
```
- 不変のオブジェクト
- 独自のバリデーションルールを持つ
- ドメインの表現力を高める

#### Repository Interface（リポジトリインターフェース）
```go
type UserRepository interface {
    FindAll(ctx context.Context) ([]*entity.User, error)
    FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error)
}
```
- データアクセスの抽象化
- ドメイン層が具体的な実装に依存しないようにする

### 2. UseCase層（ユースケース層）

アプリケーション固有のビジネスルールを実装します。

```go
type Service struct {
    userRepo repository.UserRepository
}

func (s *Service) GetUserByID(ctx context.Context, idStr string) (*entity.User, error) {
    // ビジネスロジックの実装
}
```
- ドメインオブジェクトを組み合わせてユースケースを実現
- 外部との入出力を調整
- トランザクション境界の管理

### 3. Interface層（インターフェース層）

外部との接点を提供し、データ形式の変換を行います。

#### Handler（ハンドラー）
```go
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    // HTTPリクエストの処理
    // UseCaseの呼び出し
    // レスポンスの返却
}
```

#### DTO（Data Transfer Object）
```go
type UserResponse struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
```
- 外部とのデータ交換用オブジェクト
- ドメインモデルとは独立

#### Gateway（ゲートウェイ）
- 外部APIとの通信を実装
- リポジトリインターフェースの実装

### 4. Infrastructure層（インフラストラクチャ層）

技術的な詳細を実装する最外層です。

- HTTPサーバーの設定
- HTTPクライアントの設定
- データベース接続（将来的に追加可能）
- 外部サービスとの統合

## データフローの例

### ユーザー情報取得のフロー

1. **HTTPリクエスト受信**
   - `GET /users/123`

2. **Handler層**
   - リクエストからIDを抽出
   - UseCase層のサービスを呼び出し

3. **UseCase層**
   - 文字列IDをValue Objectに変換
   - リポジトリを通じてデータを取得
   - ビジネスルールの適用

4. **Gateway層**
   - 外部API（JSONPlaceholder）にリクエスト
   - レスポンスをドメインオブジェクトに変換

5. **Handler層**
   - ドメインオブジェクトをDTOに変換
   - JSONレスポンスとして返却

### データ更新のフロー（ユーザー登録の例）

1. **HTTPリクエスト受信**
   ```json
   POST /users
   {
     "name": "John Doe",
     "username": "johndoe",
     "email": "john@example.com"
   }
   ```

2. **Handler層（Presentation）**
   ```go
   func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
       // リクエストボディをDTOにデコード
       var dto dto.CreateUserRequest
       json.NewDecoder(r.Body).Decode(&dto)
       
       // 基本的な入力チェック
       if err := dto.Validate(); err != nil {
           http.Error(w, err.Error(), http.StatusBadRequest)
           return
       }
       
       // UseCase層を呼び出し
       user, err := h.userService.CreateUser(r.Context(), dto)
   }
   ```

3. **DTO層でのバリデーション**
   ```go
   type CreateUserRequest struct {
       Name     string `json:"name"`
       Username string `json:"username"`
       Email    string `json:"email"`
   }
   
   func (d *CreateUserRequest) Validate() error {
       // 必須フィールドチェック
       if d.Name == "" || d.Username == "" {
           return errors.New("required fields missing")
       }
       return nil
   }
   ```

4. **UseCase層（Application Business Rules）**
   ```go
   func (s *UserService) CreateUser(ctx context.Context, dto dto.CreateUserRequest) (*entity.User, error) {
       // Value Objectの生成（ドメインルールのバリデーション）
       email, err := valueobject.NewEmail(dto.Email)
       if err != nil {
           return nil, fmt.Errorf("invalid email: %w", err)
       }
       
       // ビジネスルールチェック（重複確認）
       existingUser, _ := s.userRepo.FindByEmail(ctx, email)
       if existingUser != nil {
           return nil, errors.New("email already registered")
       }
       
       // エンティティの生成
       userID := valueobject.GenerateUserID() // ID生成
       user := entity.NewUser(userID, dto.Name, dto.Username, email)
       
       // リポジトリ経由で永続化
       if err := s.userRepo.Save(ctx, user); err != nil {
           return nil, fmt.Errorf("failed to save user: %w", err)
       }
       
       return user, nil
   }
   ```

5. **Domain層でのバリデーション**
   ```go
   // Value Objectでのドメインルール適用
   func NewEmail(value string) (Email, error) {
       value = strings.TrimSpace(value)
       
       // フォーマットチェック
       if !emailRegex.MatchString(value) {
           return Email{}, ErrInvalidEmailFormat
       }
       
       // ドメイン固有のルール（例：ブラックリスト）
       if isBlacklistedDomain(extractDomain(value)) {
           return Email{}, ErrBlacklistedDomain
       }
       
       return Email{value: value}, nil
   }
   ```

6. **Repository実装（Gateway/Infrastructure）**
   ```go
   func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
       // 実際のデータストアへの保存
       // 例：データベース
       query := `INSERT INTO users (id, name, username, email) VALUES (?, ?, ?, ?)`
       _, err := r.db.ExecContext(ctx, query,
           user.ID().String(),
           user.Name(),
           user.Username(),
           user.Email().String(),
       )
       
       // 例：外部API
       payload := map[string]interface{}{
           "id":       user.ID().Value(),
           "name":     user.Name(),
           "username": user.Username(),
           "email":    user.Email().String(),
       }
       _, err := r.httpClient.Post(r.apiURL+"/users", "application/json", payload)
       
       return err
   }
   ```

7. **レスポンス返却**
   ```go
   // Handler層でエンティティをDTOに変換
   response := dto.UserResponse{
       ID:       user.ID().Value(),
       Name:     user.Name(),
       Username: user.Username(),
       Email:    user.Email().String(),
   }
   
   w.Header().Set("Content-Type", "application/json")
   w.WriteHeader(http.StatusCreated)
   json.NewEncoder(w).Encode(response)
   ```

### バリデーションの階層と責任

1. **Presentation層（Handler/DTO）**
   - 必須フィールドの存在確認
   - 基本的な型チェック
   - リクエスト形式の検証

2. **Domain層（Value Object）**
   - ビジネスルールに基づくフォーマット検証
   - ドメイン固有の制約チェック
   - 値の正規化

3. **UseCase層**
   - エンティティ間の整合性チェック
   - ビジネスプロセスの検証
   - 外部システムとの整合性確認

4. **Infrastructure層**
   - データベース制約（ユニーク制約など）
   - トランザクション管理
   - 外部サービスのエラーハンドリング

### エラーハンドリングの流れ

```go
// Domain層のエラー
var (
    ErrInvalidEmailFormat = errors.New("invalid email format")
    ErrBlacklistedDomain  = errors.New("email domain is blacklisted")
)

// UseCase層でのエラーハンドリング
if err != nil {
    switch {
    case errors.Is(err, valueobject.ErrInvalidEmailFormat):
        return nil, &ValidationError{Field: "email", Message: "Invalid format"}
    case errors.Is(err, valueobject.ErrBlacklistedDomain):
        return nil, &BusinessRuleError{Message: "Domain not allowed"}
    default:
        return nil, fmt.Errorf("unexpected error: %w", err)
    }
}

// Handler層でのHTTPステータスマッピング
switch err.(type) {
case *ValidationError:
    http.Error(w, err.Error(), http.StatusBadRequest)
case *BusinessRuleError:
    http.Error(w, err.Error(), http.StatusUnprocessableEntity)
case *NotFoundError:
    http.Error(w, err.Error(), http.StatusNotFound)
default:
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
```

## 利点

1. **テスタビリティ**
   - 各層が独立しているため単体テストが容易
   - モックを使用した統合テストが可能

2. **保守性**
   - 責任が明確に分離されている
   - 変更の影響範囲が限定的

3. **拡張性**
   - 新機能の追加が既存コードに影響しにくい
   - 外部サービスの切り替えが容易

4. **ビジネスロジックの保護**
   - ドメイン層が技術的詳細から独立
   - ビジネスルールの一貫性を保ちやすい

## 今後の拡張可能性

1. **データベースの追加**
   - Infrastructure層にDB接続を追加
   - Gateway層にリポジトリ実装を追加

2. **認証・認可**
   - Interface層にミドルウェアを追加
   - UseCase層に認可ロジックを実装

3. **キャッシュ**
   - Infrastructure層にキャッシュストアを追加
   - Gateway層でキャッシュ制御

4. **イベント駆動**
   - Domain層にドメインイベントを追加
   - Infrastructure層にメッセージキューを統合