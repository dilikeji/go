package App

import (
	"database/sql"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"math/rand"
)

// LoginRequest 入参结构
type LoginRequest struct {
	// https://goframe.org/docs/core/gvalid-rules 校验规则文档
	Username string `json:"username" v:"required#用户名必填"`
	Password string `json:"password" v:"required#密码必填"`
}

/*
示例sql
CREATE TABLE `user` (
  `user_id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

CREATE TABLE `log` (
  `log_id` int unsigned NOT NULL AUTO_INCREMENT,
  `log_info` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '日志内容',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间戳',
  PRIMARY KEY (`log_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='日志表';
*/

func Login(r *ghttp.Request) {
	var (
		// 用到的变量提前声明
		request  LoginRequest
		err      error
		record   gdb.Record
		result   sql.Result
		InsertId int64
		Rows     int64
	)
	if err = r.Parse(&request); err != nil {
		ReturnJson(r, gcode.New(CodeValidate, err.Error(), nil))
		return
	}
	// 获取上下文变量 *gvar.Var类型的数组或对象可以通过.MapStrVar()后获取其中的值
	g.Dump(r.GetCtxVar("Auth").MapStrVar()["timeNow"])
	// 数据库ORM操作 https://goframe.org/docs/core/gdb-chaining
	if record, err = g.Model("user").Fields("*").One(g.Map{
		"username": request.Username,
	}); err != nil {
		// 数据库连接或执行sql报错 例如username字段不存在 则报错如下:SELECT * FROM `user` WHERE `username`='111' LIMIT 1: Error 1054 (42S22): Unknown column 'username' in 'where clause'
		ReturnJson(r, gcode.New(CodeSql, err.Error(), nil))
		return
	}
	// 判断非空
	if record.IsEmpty() {
		ReturnJson(r, gcode.New(CodeError, "用户不存在", nil))
		return
	}
	// 字符串比较 https://goframe.org/docs/components/text-gstr#字符串比较
	if !gstr.Equal(record["password"].String(), request.Password) {
		ReturnJson(r, gcode.New(CodeError, "密码错误", nil))
		return
	}

	if result, err = g.Model("log").Insert(g.Map{
		// 类型转换 任意类型转换成字符串
		"log_info": g.NewVar(rand.Intn(101)).String(),
		// 时间对象 https://goframe.org/docs/components/os-gtime-time#时间对象
		"create_time": gtime.Now().Timestamp(),
	}); err != nil {
		ReturnJson(r, gcode.New(CodeSql, err.Error(), nil))
		return
	}
	InsertId, err = result.LastInsertId()
	g.Dump(InsertId)

	if result, err = g.Model("log").Data(g.Map{
		"create_time": gtime.Now().Timestamp(),
	}).Where(g.Map{
		"log_id": 1,
	}).Update(); err != nil {
		ReturnJson(r, gcode.New(CodeSql, err.Error(), nil))
		return
	}
	Rows, err = result.RowsAffected()
	g.Dump(Rows)

	ReturnJson(r, gcode.New(CodeSuccess, "SUCCESS", record))
}
