/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-20 00:00:00
 * @FilePath: \go-swagger\constants\http.go
 * @Description: HTTP 常量 - 请求头、CORS、MIME 类型、HTTP 方法
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package constants

// ==================== HTTP 头与 CORS ====================

const (
	HeaderContentType               = "Content-Type"
	HeaderAccessControlAllowOrigin  = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"
	CORSAllowAll                    = "*"
	CORSDefaultMethods              = "GET, POST, PUT, DELETE, OPTIONS"
	CORSDefaultHeaders              = "Content-Type, Authorization"
	HTTPMethodOptions               = "OPTIONS"
)

// ==================== MIME 类型 ====================

const (
	MimeApplicationJSON        = "application/json"
	MimeApplicationJSONCharset = "application/json; charset=utf-8"
	MimeTextHTMLCharset        = "text/html; charset=utf-8"
	MimeYAML                   = "yaml"
	MimeYML                    = "yml"
)

// ==================== 文件扩展名 ====================

const (
	FileExtYAML = ".yaml"
	FileExtYML  = ".yml"
	FileExtJSON = ".json"
)
