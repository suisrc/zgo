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
        "/authz": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "授权接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Authorize",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.Success"
                        }
                    }
                }
            }
        },
        "/signin": {
            "post": {
                "description": "登陆",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "Signin",
                "parameters": [
                    {
                        "description": "SigninBody Info",
                        "name": "item",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schema.SigninBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.Success"
                        }
                    }
                }
            }
        },
        "/signin/refresh": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "刷新令牌",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "Refresh",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.Success"
                        }
                    }
                }
            }
        },
        "/signout": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "登出",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "Signin",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.Success"
                        }
                    }
                }
            }
        },
        "/signup": {
            "post": {
                "description": "登陆",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "Signup",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.Success"
                        }
                    }
                }
            }
        },
        "/user/current": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "查询当前用户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "user current",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.Success"
                        }
                    }
                }
            }
        },
        "/user/current/access": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "查询当前用户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "user current access",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.Success"
                        }
                    }
                }
            }
        },
        "/user/current/notices": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "查询当前用户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "user current notices",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/helper.Success"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "helper.Success": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "响应数据",
                    "type": "object"
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
        "schema.SigninBody": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "attach": {
                    "description": "reset:重置登入 refresh:刷新令牌",
                    "type": "string"
                },
                "captcha": {
                    "description": "验证码",
                    "type": "string"
                },
                "client": {
                    "description": "应用, 默认为主应用, 为空",
                    "type": "string"
                },
                "code": {
                    "description": "标识码",
                    "type": "string"
                },
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "role": {
                    "description": "角色",
                    "type": "string"
                },
                "username": {
                    "description": "账户",
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
