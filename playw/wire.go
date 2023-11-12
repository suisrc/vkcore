package playw

import "github.com/google/wire"

// WireSet wire注入声明
var WireSet = wire.NewSet(
	NewBrowserHandler,

	wire.Struct(new(Injector), "*"), // 注册器
)

type Injector struct {
	PlayW *BrowserHandler
}

// Init 初始化
func (aa *Injector) PostInit() (func(), error) {
	return aa.PlayW.Close, nil
}
