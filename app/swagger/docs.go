// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package swagger

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
        "contact": {
            "name": "suisrc",
            "email": "susirc@outlook.com"
        },
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/demo/get": {
            "get": {
                "description": "Get",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Demo id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.ErrorInfo"
                        }
                    }
                }
            }
        },
        "/demo/get1": {
            "get": {
                "description": "Get",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Demo id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.ErrorInfo"
                        }
                    }
                }
            }
        },
        "/demo/hello": {
            "get": {
                "description": "Hello world",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Hello",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/demo/set": {
            "post": {
                "description": "Set",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Set",
                "parameters": [
                    {
                        "description": "Demo Info",
                        "name": "item",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/schema.DemoSet"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.ErrorInfo"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "helper.ErrorInfo": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "响应数据",
                    "type": "object"
                },
                "errorCode": {
                    "description": "错误代码",
                    "type": "string"
                },
                "errorMessage": {
                    "description": "向用户显示消息",
                    "type": "string"
                },
                "showType": {
                    "description": "错误显示类型：0静音； 1条消息警告； 2消息错误； 4通知； 9页",
                    "type": "integer"
                },
                "success": {
                    "description": "请求成功, false",
                    "type": "boolean"
                },
                "traceId": {
                    "description": "方便进行后端故障排除：唯一的请求ID",
                    "type": "string"
                }
            }
        },
        "schema.DemoSet": {
            "type": "object",
            "required": [
                "code",
                "name"
            ],
            "properties": {
                "code": {
                    "description": "编号",
                    "type": "string"
                },
                "memo": {
                    "description": "备注",
                    "type": "string"
                },
                "name": {
                    "description": "名称",
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
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
	Version:     "0.0.1",
	Host:        "",
	BasePath:    "/api",
	Schemes:     []string{"https", "http"},
	Title:       "zgo",
	Description: "GIN + ENT/SQLX + CASBIN + WIRE",
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
