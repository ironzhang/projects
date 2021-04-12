# API

## gc-token

```go

type Token struct {
	Value string
	Version uint64
	ExpireTS uint64
}

type Request struct {
	UID uint64
}

type Result struct {
	UID uint64
	AccessToken string
	AccessTokenExpireTS uint64
	RefreshToken string
	RefreshTokenExpireTS uint64
}

func UpdateToken() error
func UpdateAccessToken() error
func GetAccessToken() error
func GetRefreshToken() error
```

---

## gc-account

```

func Login(*LoginRequest, *LoginResult) error {
	//1. 验证密码

	//2. 更新AK,RK
	gctoken.UpdateToken(uid)

	//3. 返回AK,RK
}

func UpdateAccessToken(*UpdateAccessTokenRequest, *UpdateAccessTokenResult) error {
	//1. 获取RK以验证RK加密的签名
	gctoken.GetRefreshToken()

	//2. 更新AK
	gctoken.UpdateAccessToken()

	//3. 返回更新的AK
}

func GetAccessToken() error {
}

```

```

func VerifySignature() error {
	//1. 获取AK
	gcaccount.GetAccessToken(uid)

	//2. 判断签名版本

	//3. 判断签名过期时间

	//4. 验证签名
}

```

