package root

import (
	"bytes"
	"fmt"
)

// Error 错误处理结构
type Error struct {
	// 错误码
	Code string

	// 错误信息，显示给终端客户用
	Message string

	// 错误详情，用于追溯
	Op  string
	Err error
}

// 通用错误码
const (
	ECONFLICT = "conflict"  // 冲突
	EINTERNAL = "internal"  // 内部错误
	EINVALID  = "invalid"   // 未通过验证
	ENOFOUND  = "not_found" // 未找到
)

// 配置信息错误
const (
	ECONFIGNOTFOUND      = "config_not_found"       // 配置文件不存在
	ECONFIGNOTINVALID    = "config_not_invalid"     // 配置文件格式不合法
	ECONFIGMYSQLNOTFOUND = "config_mysql_not_found" // mysql配置信息不存在
	ECONFIGMSSQLNOTFOUND = "config_mssql_not_found" // mssql配置信息不存在
	ECONFIGHTTPNOTFOUND  = "config_http_not_found"  // http配置不存在
	ECONFIGAUTHNOTFOUND  = "config_auth_not_found"  // auth配置不存在
)

// 数据库错误
const (
	ESERVICEWITHNILSESSION = "service_with_nil_session" // session对象为空
	ESERVICEWITHNILDB      = "service_with_nil_db"      // 数据库对象为空
	EDBQUERYERROR          = "db_query_error"           // 数据库查询错误
	EDBBEGINERROR          = "db_begin_error"           // 数据库开始事务错误
	EDBPREPAREERROR        = "db_prepare_error"         // 数据库准备stmt失败
	EDBEXECERROR           = "db_exec_error"            // 数据库语执行错误
	EDBOPENERROR           = "db_open_error"            // 数据库打开失败
)

//EAUTHERROR token解析错误
const EAUTHERROR = "auth_error" // 身份验证错误

// Error 实现错误接口
func (e *Error) Error() string {
	var buf bytes.Buffer

	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}

		buf.WriteString(e.Message)
	}

	return buf.String()
}

// ErrorCode 获取错误码
func ErrorCode(err error) string {

	if err == nil {
		return ""
	}

	if e, ok := err.(*Error); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}
	return EINTERNAL
}

// ErrorMessage 获取错误信息
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}

	return "系统内部错误，请联系系统管理员"
}
