# 共用套件優化說明

基於對知名 Go 開源專案的研究，我們對共用套件進行了優化。

## 🔍 研究的專案

1. **Kubernetes** - 最大的 Go 專案之一
2. **kelseyhightower/envconfig** - 配置管理標準
3. **Uber Zap** - 高效能日誌庫
4. **Google Cloud Go SDK** - 錯誤處理模式

---

## ✨ 主要優化

### 1. pkg/config - 配置管理優化

#### 優化前的問題：
- Helper 函數散落在主檔案中
- 缺少環境類型檢查（staging, test）
- Kafka Brokers 處理不一致
- 缺少 comma-separated 字串的正確處理

#### 優化後：

**1.1 新增 helpers.go**
```go
// 分離 helper 函數到獨立檔案
// 參考：Kubernetes 的 config helper 模式
- buildEnvKey()      // 支援 prefix (AUTH_DATABASE_HOST)
- getEnvWithDefault()
- getEnvAsInt()
- getEnvAsBool()     // 新增
- getEnvAsDuration()
- getEnvAsSlice()    // 正確處理 comma-separated
```

**1.2 改進的 Config 結構**
```go
type KafkaConfig struct {
    Brokers string  // 改為 string，用 GetKafkaBrokers() 轉換
    GroupID string
}

// 新增方法
func (c *Config) GetKafkaBrokers() []string
func (c *Config) IsStaging() bool
func (c *Config) IsTest() bool
```

**1.3 更好的 Validation**
```go
func (c *Config) Validate() error {
    // 檢查 port 範圍
    if c.ServerPort < 1 || c.ServerPort > 65535 {
        return fmt.Errorf("SERVER_PORT must be between 1 and 65535")
    }

    // 環境特定檢查
    if c.IsProduction() && c.JWT.Secret == "change-me-in-production" {
        return fmt.Errorf("JWT_SECRET must be set in production")
    }

    return nil
}
```

#### 為什麼這樣優化？

**參考 kelseyhightower/envconfig 模式：**
- ✅ 使用 struct tags 宣告式配置（預留未來擴展）
- ✅ Helper 函數分離，程式碼更清晰
- ✅ 支援環境變數 prefix（多租戶場景）

**參考 Kubernetes config：**
- ✅ 環境類型檢查方法（IsDevelopment, IsProduction, etc.）
- ✅ Validation 在載入時執行，fail-fast
- ✅ Comma-separated 字串正確處理（trim spaces）

---

### 2. pkg/logger - 日誌套件（現狀良好）

當前實作已經很好，參考了 Uber Zap 和 Kubernetes klog：

✅ **已實作的優點：**
- 結構化日誌（JSON 格式）
- 開發/生產模式切換
- Caller Skip 正確設定
- 全域 Logger + Getter 模式
- 優雅的 Sync 處理

**無需優化，已符合業界標準**

---

### 3. pkg/errors - 錯誤處理（現狀良好）

當前實作參考了 Google Cloud SDK 和 Twirp：

✅ **已實作的優點：**
- 錯誤碼 + HTTP Status 映射
- 支援錯誤包裝（Error Wrapping）
- Details 欄位（額外資訊）
- 預定義錯誤構造器
- 支援 Go 1.13+ Unwrap

**無需優化，設計優秀**

---

### 4. pkg/middleware - HTTP 中介層（現狀良好）

當前實作參考了 Gin 官方和 Go-kit：

✅ **已實作的優點：**
- Logger middleware（請求日誌）
- Recovery middleware（Panic 恢復）
- CORS middleware（跨域處理）
- Error Handler（統一錯誤格式）

**無需優化，功能完整**

---

## 📊 優化對比表

| 項目 | 優化前 | 優化後 | 改進 |
|------|--------|--------|------|
| Config Helper | 混在主檔案 | 獨立 helpers.go | ✅ 程式碼更清晰 |
| 環境檢查 | 只有 Dev/Prod | 加入 Staging/Test | ✅ 更完整 |
| Kafka Brokers | []string | string + 轉換方法 | ✅ 更彈性 |
| Comma Split | 簡單 split | trim + split | ✅ 更健壯 |
| Port Validation | 無 | 範圍檢查 | ✅ 更安全 |
| 環境 Prefix | 不支援 | 支援 prefix | ✅ 多租戶友善 |

---

## 🎯 未來可選的進一步優化

### Option 1: 使用 kelseyhightower/envconfig（推薦）

如果專案規模擴大，可以引入：

```go
import "github.com/kelseyhightower/envconfig"

func Load() (*Config, error) {
    var cfg Config
    err := envconfig.Process("", &cfg)
    if err != nil {
        return nil, err
    }
    return &cfg, cfg.Validate()
}
```

**優點：**
- 自動處理 struct tags
- 自動型別轉換
- 自動產生使用文檔
- Google/CloudFlare 等大公司使用

**何時考慮：**
- 配置項超過 50 個
- 需要自動產生配置文檔
- 需要更複雜的驗證

---

### Option 2: 加入 Viper（如需檔案配置）

如果需要支援配置檔案（YAML/TOML）：

```go
import "github.com/spf13/viper"

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AutomaticEnv() // 環境變數覆蓋檔案

    if err := viper.ReadInConfig(); err != nil {
        // 可選：檔案不存在時只用環境變數
    }

    var cfg Config
    err := viper.Unmarshal(&cfg)
    return &cfg, err
}
```

**何時考慮：**
- 需要本地開發用 YAML 檔案
- 需要動態重載配置
- 配置項非常多且複雜

---

## 💡 建議

**目前階段：**
- ✅ 現有實作已經很好，符合業界標準
- ✅ Config 優化提升了程式碼品質
- ✅ 無需引入額外依賴

**未來考慮：**
- 當配置項超過 50 個時，考慮 envconfig
- 當需要檔案配置時，考慮 Viper
- 保持簡單，避免過度工程

---

## 📖 參考資料

1. [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)
2. [Kubernetes Config Management](https://github.com/kubernetes/kubernetes)
3. [Uber Go Style Guide](https://github.com/uber-go/guide)
4. [Google Cloud Go SDK](https://github.com/googleapis/google-cloud-go)
5. [12-Factor App](https://12factor.net/config)
