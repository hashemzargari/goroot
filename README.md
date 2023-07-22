# GoRoot

GoRoot (I AM GOROOT) is a simple and mini web framework for Golang.

### How to use

#### 1- define your handler ( business logic)

```go
package main
   
import (
    "github.com/hashemzargari/goroot"
)

type SumNumbersRequest struct {
	Numbers []int `json:"numbers"`
}

type SumNumbersResponse struct {
	Result int `json:"result"`
}

type SumHandler struct{}

func (h SumHandler) Handle(ctx goroot.Context, request any) (any, error) {
	req := request.(SumNumbersRequest)

	result := 0
	for _, number := range req.Numbers {
		result += number
	}

	return SumNumbersResponse{Result: result}, nil
}

func (h SumHandler) GetRequestType() any {
	return SumNumbersRequest{}
}

func (h SumHandler) GetResponseType() any {
	return SumNumbersResponse{}
}

func (h SumHandler) GetPermissionClaims() []goroot.PermissionClaim {
	return []goroot.PermissionClaim{}
}

func (h SumHandler) GetAuthenticationClaims() goroot.AuthenticationClaim {
	return goroot.AllowAuthenticated | goroot.AllowAdmin
}

func (h SumHandler) GetApiRoute() string {
	return "/sum-numbers"
}

func (h SumHandler) GetMethod() goroot.ApiHttpMethod {
	return goroot.Post
}

```

#### 2- create app, register your handler and Run

```go
func main() {
    app := goroot.New("MyService")
    app.RegisterHandlers(SumHandler{})
    app.Run(":8080")
}
```

