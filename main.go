package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lingochamp/core"
	"github.com/lingochamp/xorm"
	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        int64     `xorm:"'id' BIGINT AUTOINCR PK"`
	Name      string    `xorm:"'name' VARCHAR(64)"`
	CreatedAt time.Time `xorm:"'created_at' BIGINT CREATED"`
	UpdatedAt time.Time `xorm:"'updated_at' BIGINT UPDATED"`
}

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int64("id", u.ID),
		slog.String("name", u.Name),
		slog.Time("created_at", u.CreatedAt),
		slog.Time("updated_at", u.UpdatedAt),
	)
}

var (
	engine *xorm.Engine
)

func Insert(ctx context.Context, session *xorm.Session, user *User) (err error) {
	if session == nil {
		session = engine.NewSession()
		defer session.Close()
	}

	_, err = session.InsertOne(ctx, user)
	return errors.WithStack(err)
}

func Find(ctx context.Context, session *xorm.Session, option func(*xorm.Session), filter User) (users []*User, err error) {
	if session == nil {
		session = engine.NewSession()
		defer session.Close()
	}

	if option != nil {
		option(session)
	}

	err = session.Find(ctx, &users, filter)
	return users, errors.WithStack(err)
}

func Init() {
	// Logger
	{
		var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
		slog.SetDefault(logger)
	}

	// DB
	{
		var err error
		engine, err = xorm.NewEngine("sqlite3", "file::memory:", xorm.LoggerOption(func(context.Context) core.ILogger {
			return xorm.NewSimpleLogger2(os.Stdout, "[XORM]", 0)
		}))
		if err != nil {
			slog.Error(`xorm.NewEngine("sqlite3", "file::memory:") failed`, slog.Any("err", err))
			return
		}

		// Note：因为程序结束数据库就销毁了，所以没有调用 engine.Close()

		engine.ShowSQL(true)

		err = engine.Sync2(context.Background(), new(User))
		if err != nil {
			slog.Error("engine.Sync2() failed", slog.Any("err", err))
			return
		}
	}
}

func main() {
	Init()

	var ctx = context.Background()
	var logger = slog.Default()
	var err error

	// 创建测试数据
	{
		var user = User{
			Name: "Alice",
		}
		err = Insert(ctx, nil, &user)
		if err != nil {
			logger.Error("Insert() failed", slog.Any("err", err))
			return
		}
		logger.Info("Insert() succeeded", slog.Any("user", user))
	}

	// 验证写入
	{
		users, err := Find(ctx, nil, nil, User{})
		if err != nil {
			logger.Error("Find(ctx, nil, nil, User{}) failed", slog.Any("err", err))
			return
		}
		if len(users) != 1 {
			logger.Error("len(users) != 1", slog.Any("users", users))
			return
		}
		if users[0].Name != "Alice" {
			logger.Error("users[0].Name != 'Alice'", slog.Any("users", users))
			return
		}
		logger.Info("Find() succeeded", slog.Any("user", users[0]))
	}

	// 验证 BUG
	{
		var selectStr = "MAX(id) AS id"
		option := func(session *xorm.Session) {
			session.Select(selectStr).Limit(1)
		}

		{
			_, err := Find(ctx, nil, option, User{})
			if err == nil {
				logger.Info("Find() should fail")
				return
			}
			slog.Info("Find() failed as expected", slog.Any("err", err))
		}

		selectStr = strings.ReplaceAll(selectStr, `(id)`, `( id)`)
		{
			users, err := Find(ctx, nil, option, User{})
			if err != nil {
				logger.Error("Find(ctx, nil, option, User{}) failed", slog.Any("err", err))
				return
			}
			logger.Info("Find() succeeded", slog.Any("users", users))
		}
	}
}
