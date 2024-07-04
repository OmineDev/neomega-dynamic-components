package main

import (
	"encoding/json"
	"fmt"
	"strings"

	neomega_backbone "github.com/OmineDev/neomega-backbone"
	"github.com/OmineDev/neomega-core/neomega"
)

// 对于 直接运行(go run .)框架，使用 go build -buildmode=plugin -o dynamic_component1.so main.go 编译
// 对于 vscode，使用 go build -buildmode=plugin -o dynamic_component1.so -gcflags "all=-N -l" 编译
// 对于 release 框架(此处存疑)，使用 go build -buildmode=plugin -trimpath -o dynamic_component1.so main.go
// Exported 供 omega 调用以创建一个插件实例
var Exported neomega_backbone.DynamicComponentFactory

func init() {
	// 因为一份插件可能对应多个配置文件，所以每个配置都需要创建一个实例
	Exported=func(name string, challengeFn neomega_backbone.ChallengeFn) neomega_backbone.DynamicComponent {
		// challengeFn 允许插件检查宿主程序是否满s足条件，若未得到预期回答则可以终止插件
		// if challengaFn("1+1")!=2{panic("stop")}
		// 当然，这里一般用非对称加密就是了
		// 当编译出来的插件文件包含多个插件的时候，可以用name选择实现
		// if name=="a" return &A{}, if name=="b" return &B{}
		// 默认 name=""
		return &MyDynamicComponent{}
	}
	func(neomega_backbone.DynamicComponentFactory){}(Exported)
	fmt.Println("dynamic component 1 程序成功被读取")
}

type MyDynamicComponentConfig struct {
	Vesion              string   `json:"Version"`
	Suffix              string   `json:"口癖"`
	Repeator            string   `json:"复读机"`
	Hint                string   `json:"请求输入提示"`
	TerminalTrigger     []string `json:"终端触发词"`
	GameTrigger         []string `json:"游戏内触发词"`
	DescripetionOnStart string   `json:"启动时提示"`
}

type MyDynamicComponent struct {
	neomega_backbone.BasicDynamicComponent
	config MyDynamicComponentConfig
}

func (c *MyDynamicComponent) Init(
	cfg neomega_backbone.DynamicComponentConfig,
	storage neomega_backbone.StorageAndLogAccess,
) {
	bs, err := json.Marshal(cfg.Configs())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Dynamic Component 1 原始配置: %v\n", string(bs))
	json.Unmarshal(bs, &c.config)
	if c.config.Vesion == "" { // 假如这个配置还没初始化
		c.config = MyDynamicComponentConfig{
			Vesion:          "0.0.1",
			Suffix:          "喵喵~ ",
			Repeator:        "2401PT",
			Hint:            "喵? (随便输入点什么)",
			TerminalTrigger: []string{"r", "echo", "repeat"},
			GameTrigger:     []string{"r", "echo", "repeat"},
		}
		cfg.Upgrade(c.config) // 更新配置文件
	}
	if c.config.Vesion == "0.0.1" { // 假如这个配置是老版本的
		c.config.Vesion = "0.0.2"
		c.config.DescripetionOnStart = "扣%v和复读机交互"
		cfg.Upgrade(c.config) // 更新配置文件
	}
	fmt.Println("Dynamic Component 1 Init!")
}

func (c *MyDynamicComponent) doGameEcho(pk neomega.PlayerKit, msg string) {
	pk.Say(fmt.Sprintf("%v: %v%v", c.config.Repeator, msg, c.config.Suffix))
}

func (c *MyDynamicComponent) onGameMenu(chat *neomega.GameChat) {
	pk, found := c.Frame.GetPlayerInteract().GetPlayerKit(chat.Name)
	if !found {
		return
	}
	if len(chat.Msg) > 0 {
		c.doGameEcho(pk, chat.Msg[0])
	} else {
		pk.Say(c.config.Hint)
		pk.InterceptJustNextInput(func(chat *neomega.GameChat) {
			c.doGameEcho(pk, chat.Msg[0])
		})
	}
}

func (c *MyDynamicComponent) doTerminalEcho(msg string) {
	c.Frame.Logger().Write(fmt.Sprintf("%v: %v%v", c.config.Repeator, msg, c.config.Suffix))
}

func (c *MyDynamicComponent) onTerminalMenu(chat []string) {
	if len(chat) > 0 {
		c.doTerminalEcho(chat[0])
	} else {
		fmt.Println(c.config.Hint)
		c.Frame.GetTerminalInput().AsyncCallBackInGoRoutine(func(nextInput string) {
			c.doTerminalEcho(strings.TrimSpace(nextInput))
		})
	}
}

func (c *MyDynamicComponent) Inject(frame neomega_backbone.ExtendOmega) {
	c.Frame = frame
	fmt.Println("Dynamic Component 1 Injected!")
	frame.AddGameMenuEntry(&neomega_backbone.GameMenuEntry{
		MenuEntry: neomega_backbone.MenuEntry{
			Triggers:     c.config.GameTrigger,
			ArgumentHint: "[要复读的内容]",
			Usage:        "和复读机互动",
		},
		OnTrigCallBack: c.onGameMenu,
	})
	frame.AddBackendMenuEntry(&neomega_backbone.BackendMenuEntry{
		MenuEntry: neomega_backbone.MenuEntry{
			Triggers:     c.config.TerminalTrigger,
			ArgumentHint: "[要复读的内容]",
			Usage:        "和复读机互动",
		},
		OnTrigCallBack: c.onTerminalMenu,
	})
}

func (c *MyDynamicComponent) Activate() {
	fmt.Println("Dynamic Component 1 Activate!")
	fmt.Println(fmt.Printf(c.config.DescripetionOnStart, c.config.TerminalTrigger[0]))
}
