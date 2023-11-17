package procv

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/suisrc/vkcore/httpv"
	"github.com/suisrc/vkcore/mgo"
	"github.com/suisrc/vkcore/solver"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// ==============================================================
// 初始化 WAF 监控
func InitWAF(domain, challengeJs string) func() {
	// 初始化 AWS WAF 令牌
	ccc := make(chan int)
	go solver.ListenToAwsWAF(domain, challengeJs, ccc)
	// ccc <- 1 // 也可以关闭监听 WAF 浏览器
	// 脚本执行过程，不要关闭，用于实时更新WAF
	for solver.AwsWaf == "" {
		logrus.Info("main process wait for waf token init...")
		time.Sleep(time.Second) // 等待 WAF 初始化令牌
	}
	return func() {
		close(ccc)
	}
}

// ==============================================================
// 返回值为 false 时，终止执行
type Execute func(context.Context, *mongo.Database, int)
type Process func(context.Context, *mongo.Collection, chan int, UserUpdate)

// ==============================================================
// 返回是一个mongo迭代器
type Find func(ctx context.Context, cll *mongo.Collection) (*mongo.Cursor, error)
type Exec func(*UserData, int, UserUpdate)

// ==============================================================
// 多协程执行
func RunWithEnv(execute Process) {
	cpath, count := ParseEnv()
	RunWithProcess(execute, cpath, count)
}

func RunWithExe(execute Execute) {
	cpath, count := ParseEnv()
	RunWithExecute(execute, cpath, count)
}

func ParseEnv() (cpath string, count int) {
	if len(os.Args) > 1 {
		cpath = os.Args[1]
	}
	if cpath == "" {
		cpath = "mongo.json" // 默认配置文件
	} else if cpath == "dev" {
		cpath = "data/conf/mongo.json"
	}
	logrus.Infof("conf path: %s", cpath)
	if len(os.Args) > 2 {
		count, _ = strconv.Atoi(os.Args[2])
	}
	if count == 0 {
		count = 1 // 默认并发数
	}
	logrus.Infof("parallel : %d", count)
	return
}

// ==============================================================
// 多协程执行
func RunWithExecute(execute Execute, cpath string, count int) {
	// 数据库连接
	cli, clx, err := mgo.NewDatabaseByFile(cpath)
	if err != nil {
		logrus.Panic("init mongo db err: ", err) // 直接终止程序
	}
	defer clx()
	ctx := context.TODO()

	//============================================================
	// 创建一个通道来接收中断信号，启动一个 goroutine 来等待中断信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		// 阻塞等待中断信号
		<-interrupt
		logrus.Info("receive interrupt signal")
		// 退出应用程序
		os.Exit(0)
	}()
	//============================================================
	// 执行业务相关操作
	execute(ctx, cli, count)
}

// ==============================================================
// 多协程执行
func RunWithProcess(process Process, cpath string, count int) {
	RunWithExecute(func(ctx context.Context, cli *mongo.Database, count int) {
		cll := cli.Collection(CollName)
		upt := UpdateUserByMgo(ctx, cll)

		logrus.Info("start process")
		wc := make(chan int, count) // 限制并发
		defer func() { close(wc) }()
		process(ctx, cll, wc, upt)
		// 等待并发完成
		for len(wc) > 0 {
			logrus.Info("wait for process done, left: ", len(wc))
			time.Sleep(time.Second) // 等待并发完成
		}
		logrus.Info("process done")

	}, cpath, count)
}

func RunWithFind(find Find, exec Exec, cpath string, count int) {
	// 执行业务相关操作
	RunWithProcess(func(ctx context.Context, cll *mongo.Collection, wc chan int, upt UserUpdate) {
		cur, err := find(ctx, cll)
		if err != nil {
			logrus.Error("find user error: ", err.Error())
			return
		}
		defer cur.Close(ctx)
		hdl, _ := httpv.NewPlayFC("")
		defer hdl.Close()

		usrs := []*UserData{}
		err = cur.All(ctx, &usrs)
		if err != nil {
			logrus.Error("decode user error: ", err.Error())
			return
		}
		for cc := 0; cc < len(usrs); cc++ {
			user := usrs[cc]
			// if err := republik.SyncAccessToken(hdl, user, upt, false); err != nil {
			// 	logrus.Error("sync user token error: ", err.Error(), " <- ", user.Username)
			// 	continue
			// }
			// else if err := republik.CheckUserToken(user); err != nil {
			// 	logrus.Error("check user token error: ", err.Error(), " <- ", user.Username)
			// 	continue
			// }
			wc <- 1                           // 并发锁
			go func(user *UserData, cc int) { // 异步并发处理
				exec(user, cc, upt)
				<-wc // 释放锁
			}(user, cc)

			// if cc > 10 {
			// 	break // 调试用
			// }
		}
	}, cpath, count)
}

// ==============================================================
// 多协程执行

type Exec2 func(*httpv.PlayFC, *UserData, int, UserUpdate, chan *UserData, chan int) // 需要管控wc锁

func RunWithExec2(find Find, exec Exec2, cpath string, tcnt, ucnt int) {
	// 执行业务相关操作
	RunWithProcess(func(ctx context.Context, cll *mongo.Collection, wc chan int, upt UserUpdate) {
		cur, err := find(ctx, cll)
		if err != nil {
			logrus.Error("find user error: ", err.Error())
			return
		}
		defer cur.Close(ctx)
		hdl, _ := httpv.NewPlayFC("")
		defer hdl.Close()
		uc := make(chan *UserData, tcnt*2)
		go func() {
			for {
				usr := <-uc
				if usr == nil {
					break
				}
				wc <- 1 // 并发锁
				go exec(hdl, usr, usr.OperateTemp["cc"].(int), upt, uc, wc)
			}
		}()
		usrs := []*UserData{}
		err = cur.All(ctx, &usrs)
		if err != nil {
			logrus.Error("decode user error: ", err.Error())
			return
		}
		for cc := 0; cc < len(usrs) && ucnt > cc; cc++ {
			for len(wc) == tcnt {
				time.Sleep(time.Second) // 防止竞争，等待完成
			}
			user := usrs[cc]
			user.OperateTemp = make(map[string]interface{})
			user.OperateTemp["cc"] = cc
			uc <- user
		}
		time.Sleep(time.Second) // 子弹飞一会
		for len(uc) > 0 || len(wc) > 0 {
			logrus.Info("wait for process done, left: ", len(wc), " -> ", len(uc))
			time.Sleep(time.Second) // 等待并发完成
		}
	}, cpath, tcnt)
}
