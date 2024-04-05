// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/waf/api/v1/auth/dologin": {
            "post": {
                "description": "dologin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "处理登录",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "api_helper.LoginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api_helper.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            }
        },
        "/waf/api/v1/auth/getCaptcha": {
            "get": {
                "description": "获取图片验证码",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "获取图片验证码",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            }
        },
        "/waf/api/v1/config": {
            "get": {
                "description": "获取配置信息",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Config"
                ],
                "summary": "获取配置信息",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "修改配置文件并重启服务",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Config"
                ],
                "summary": "修改配置文件",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "config.Config",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/config.Config"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            }
        },
        "/waf/api/v1/ip": {
            "get": {
                "description": "分页查询IP管理记录",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "IP"
                ],
                "summary": "分页查询IP管理记录",
                "parameters": [
                    {
                        "type": "string",
                        "description": "关键词",
                        "name": "keywords",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "IP类型",
                        "name": "type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "页面大小",
                        "name": "perPage",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "新增IP记录",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "IP"
                ],
                "summary": "新增IP管理记录",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "types.AddIPRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.AddIPRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "删除IP记录",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "IP"
                ],
                "summary": "删除IP记录",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "types.DeleteIPRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.DeleteIPRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            }
        },
        "/waf/api/v1/log": {
            "get": {
                "description": "查询指定天数和小时数的日志记录",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Log（日志模块）"
                ],
                "summary": "查询指定天数和小时数的日志记录",
                "parameters": [
                    {
                        "type": "string",
                        "description": "查询范围天数",
                        "name": "days",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            }
        },
        "/waf/api/v1/log/getBlockLog": {
            "get": {
                "description": "分页查询被拦截的流量日志",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Log（日志模块）"
                ],
                "summary": "分页查询被拦截的流量日志",
                "parameters": [
                    {
                        "type": "string",
                        "description": "查询IP",
                        "name": "keywords",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "页面大小",
                        "name": "perPage",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            }
        },
        "/waf/api/v1/routing": {
            "get": {
                "description": "分页查询路由管理记录",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Routing"
                ],
                "summary": "分页查询路由管理记录",
                "parameters": [
                    {
                        "type": "string",
                        "description": "关键词",
                        "name": "keywords",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "路由类型",
                        "name": "type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "页面大小",
                        "name": "perPage",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "新增路由管理记录",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Routing"
                ],
                "summary": "新增路由管理记录",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "types.AddRoutingRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.AddRoutingRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "删除路由管理记录",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Routing"
                ],
                "summary": "删除路由管理记录",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "types.DeleteRoutingRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.DeleteRoutingRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            }
        },
        "/waf/api/v1/sqlInject": {
            "get": {
                "description": "查询全部sql注入防护规则",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "SqlInject"
                ],
                "summary": "查询全部sql注入防护规则",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "新增sql注入防护规则",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SqlInject"
                ],
                "summary": "新增sql注入防护规则",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "types.AddSqlInjectRuleRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.AddSqlInjectRuleRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "删除sql注入防护规则",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "SqlInject"
                ],
                "summary": "删除sql注入防护规则",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "types.DeleteSqlInjectRuleRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.DeleteSqlInjectRuleRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            }
        },
        "/waf/api/v1/xssDetect": {
            "get": {
                "description": "查询全部xss攻击防护规则",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "XssDetect"
                ],
                "summary": "查询全部xss攻击防护规则",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "新增sql注入防护规则",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "XssDetect"
                ],
                "summary": "新增xss攻击防护规则",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "types.XssDetectRule",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.XssDetectRule"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "删除xss攻击防护规则",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "XssDetect"
                ],
                "summary": "删除xss攻击防护规则",
                "parameters": [
                    {
                        "description": "请求体",
                        "name": "types.DeleteXssDetectRuleRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.DeleteXssDetectRuleRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_helper.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api_helper.LoginRequest": {
            "type": "object",
            "required": [
                "captcha",
                "captchaId",
                "password",
                "username"
            ],
            "properties": {
                "captcha": {
                    "type": "string"
                },
                "captchaId": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "api_helper.Response": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "msg": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "config.Config": {
            "type": "object",
            "properties": {
                "rateLimiter": {
                    "description": "限速器配置",
                    "type": "object",
                    "properties": {
                        "banCounter": {
                            "description": "计数器常规值，超过触发常规封禁时间",
                            "type": "integer"
                        },
                        "banDuration": {
                            "description": "常规封禁时间，单位秒",
                            "type": "integer"
                        },
                        "fixedWindow": {
                            "description": "固定窗口模式配置",
                            "type": "object",
                            "properties": {
                                "maxRequest": {
                                    "description": "窗口内最大请求数",
                                    "type": "integer"
                                },
                                "windowSize": {
                                    "description": "窗口大小,单位秒",
                                    "type": "integer"
                                }
                            }
                        },
                        "leakyBucket": {
                            "description": "漏桶模式配置",
                            "type": "object",
                            "properties": {
                                "capacity": {
                                    "description": "桶容量",
                                    "type": "integer"
                                },
                                "leakyPerSecond": {
                                    "description": "漏水速度，单位秒",
                                    "type": "integer"
                                }
                            }
                        },
                        "maxCounter": {
                            "description": "计数器最大值，超过该值拉入永久IP黑名单",
                            "type": "integer"
                        },
                        "mode": {
                            "description": "限速算法，1:令牌桶2:漏桶3:固定窗口4:滑动窗口",
                            "type": "integer"
                        },
                        "slideWindow": {
                            "description": "滑动窗口模式配置",
                            "type": "object",
                            "properties": {
                                "maxRequest": {
                                    "description": "窗口内最大队列数",
                                    "type": "integer"
                                },
                                "windowSize": {
                                    "description": "窗口大小，单位秒",
                                    "type": "integer"
                                }
                            }
                        },
                        "tokenBucket": {
                            "description": "令牌桶模式配置",
                            "type": "object",
                            "properties": {
                                "maxToken": {
                                    "description": "最大令牌数",
                                    "type": "integer"
                                },
                                "tokenPerSecond": {
                                    "description": "令牌生成速度，单位秒",
                                    "type": "integer"
                                }
                            }
                        }
                    }
                },
                "secret": {
                    "type": "object",
                    "properties": {
                        "jwtSecretKey": {
                            "type": "string"
                        },
                        "sessionSecretKey": {
                            "type": "string"
                        }
                    }
                },
                "server": {
                    "type": "object",
                    "properties": {
                        "targetAddress": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "types.AddIPRequest": {
            "type": "object",
            "required": [
                "ip",
                "type"
            ],
            "properties": {
                "expiration": {
                    "type": "integer"
                },
                "ip": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "types.AddRoutingRequest": {
            "type": "object",
            "required": [
                "method",
                "routing",
                "type"
            ],
            "properties": {
                "method": {
                    "type": "string"
                },
                "routing": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "types.AddSqlInjectRuleRequest": {
            "type": "object",
            "required": [
                "rule"
            ],
            "properties": {
                "rule": {
                    "type": "string"
                }
            }
        },
        "types.DeleteIPRequest": {
            "type": "object",
            "required": [
                "ip",
                "type"
            ],
            "properties": {
                "ip": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "types.DeleteRoutingRequest": {
            "type": "object",
            "required": [
                "method",
                "routing",
                "type"
            ],
            "properties": {
                "method": {
                    "type": "string"
                },
                "routing": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                }
            }
        },
        "types.DeleteSqlInjectRuleRequest": {
            "type": "object",
            "required": [
                "rule"
            ],
            "properties": {
                "rule": {
                    "type": "string"
                }
            }
        },
        "types.DeleteXssDetectRuleRequest": {
            "type": "object",
            "required": [
                "rule"
            ],
            "properties": {
                "rule": {
                    "type": "string"
                }
            }
        },
        "types.XssDetectRule": {
            "type": "object",
            "properties": {
                "rule": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "v0.1",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "GoWAFer",
	Description: "Golang编写的一款基于反向代理模式的web防火墙应用 By supercat0867",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
