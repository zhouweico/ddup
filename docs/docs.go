// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/auth/login": {
            "post": {
                "description": "用户登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "登录信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "登录成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.LoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "无效的请求参数",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/logout": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "用户退出登录",
                "produces": [
                    "application/json"
                ],
                "summary": "退出登录",
                "responses": {
                    "200": {
                        "description": "退出成功",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/register": {
            "post": {
                "description": "用户注册",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "注册信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "注册成功",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "400": {
                        "description": "无效的请求参数",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/organizations": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "获取当前用户所属的所有组织",
                "produces": [
                    "application/json"
                ],
                "summary": "获取用户组织列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/dto.OrganizationResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "创建新的组织，创建者默认成为管理员",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "创建组织",
                "parameters": [
                    {
                        "description": "组织信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateOrganizationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.OrganizationResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/organizations/{org_name}": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "更新组织基本信息（仅管理员可操作）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "更新组织信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "组织名称",
                        "name": "org_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "更新信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UpdateOrganizationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "删除组织（仅管理员可操作）",
                "produces": [
                    "application/json"
                ],
                "summary": "删除组织",
                "parameters": [
                    {
                        "type": "string",
                        "description": "组织名称",
                        "name": "org_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/organizations/{org_name}/join": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "加入已存在的组织",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "加入组织",
                "parameters": [
                    {
                        "type": "string",
                        "description": "组织名称",
                        "name": "org_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/organizations/{org_name}/members": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "获取组织的所有成员",
                "produces": [
                    "application/json"
                ],
                "summary": "获取组织成员列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "组织名称",
                        "name": "org_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/dto.MemberResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "添加新成员到组织（仅管理员可操作）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "添加组织成员",
                "parameters": [
                    {
                        "type": "string",
                        "description": "组织名称",
                        "name": "org_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "成员信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.AddOrganizationMemberRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/organizations/{org_name}/members/{username}": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "更新组织成员信息（仅管理员可操作）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "更新组织成员",
                "parameters": [
                    {
                        "type": "string",
                        "description": "组织名称",
                        "name": "org_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "更新信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UpdateMemberRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "从组织中移除成员（仅管理员可操作）",
                "produces": [
                    "application/json"
                ],
                "summary": "移除组织成员",
                "parameters": [
                    {
                        "type": "string",
                        "description": "组织名称",
                        "name": "org_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "username",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/users": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "获取当前登录用户信息",
                "produces": [
                    "application/json"
                ],
                "summary": "获取用户信息",
                "responses": {
                    "200": {
                        "description": "获取成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.UserResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "更新用户基本信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "更新用户信息",
                "parameters": [
                    {
                        "description": "用户信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.UpdateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "更新成功",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "400": {
                        "description": "无效的请求参数",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "删除用户账号",
                "produces": [
                    "application/json"
                ],
                "summary": "删除用户",
                "responses": {
                    "200": {
                        "description": "删除成功",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/users/password": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "修改用户密码",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "修改密码",
                "parameters": [
                    {
                        "description": "密码信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ChangePasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "修改成功",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "400": {
                        "description": "无效的请求参数",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/users/socials": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "获取当前用户的所有社交媒体账号",
                "produces": [
                    "application/json"
                ],
                "summary": "获取社交媒体账号列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/handler.SocialResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "为当前用户创建社交媒体账号",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "创建社交媒体账号",
                "parameters": [
                    {
                        "description": "社交媒体账号信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SocialRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/users/socials/{id}": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "更新指定的社交媒体账号",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "更新社交媒体账号",
                "parameters": [
                    {
                        "type": "string",
                        "description": "社交媒体账号ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "社交媒体账号信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SocialRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "删除指定的社交媒体账号",
                "produces": [
                    "application/json"
                ],
                "summary": "删除社交媒体账号",
                "parameters": [
                    {
                        "type": "string",
                        "description": "社交媒体账号ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "检查服务和数据库连接状态",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "系统"
                ],
                "summary": "健康检查",
                "responses": {
                    "200": {
                        "description": "服务正常",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "503": {
                        "description": "服务异常",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.AddOrganizationMemberRequest": {
            "type": "object",
            "required": [
                "role",
                "username"
            ],
            "properties": {
                "role": {
                    "type": "string",
                    "enum": [
                        "admin",
                        "member"
                    ]
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dto.ChangePasswordRequest": {
            "type": "object",
            "required": [
                "newPassword",
                "oldPassword"
            ],
            "properties": {
                "newPassword": {
                    "type": "string",
                    "minLength": 6
                },
                "oldPassword": {
                    "type": "string",
                    "minLength": 6
                }
            }
        },
        "dto.CreateOrganizationRequest": {
            "type": "object",
            "required": [
                "display_name",
                "name"
            ],
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "display_name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 2
                },
                "email": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 2
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "dto.LoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dto.LoginResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "expiredAt": {
                    "type": "string"
                },
                "expiresIn": {
                    "type": "integer"
                },
                "token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/dto.UserResponse"
                }
            }
        },
        "dto.MemberResponse": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "joined_at": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dto.OrganizationResponse": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "display_name": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "role": {
                    "description": "当前用户在组织中的角色",
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "dto.RegisterRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 8
                },
                "username": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 3
                }
            }
        },
        "dto.UpdateMemberRequest": {
            "type": "object",
            "required": [
                "role"
            ],
            "properties": {
                "role": {
                    "type": "string",
                    "enum": [
                        "admin",
                        "member"
                    ]
                }
            }
        },
        "dto.UpdateOrganizationRequest": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "display_name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 2
                },
                "email": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 2
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "dto.UpdateUserRequest": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "bio": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "mobile": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                }
            }
        },
        "dto.UserResponse": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "bio": {
                    "type": "string"
                },
                "birthday": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastLogin": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "mobile": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "handler.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "响应码",
                    "type": "integer"
                },
                "data": {
                    "description": "响应数据"
                },
                "message": {
                    "description": "响应信息",
                    "type": "string"
                }
            }
        },
        "handler.SocialRequest": {
            "type": "object",
            "required": [
                "platform",
                "username"
            ],
            "properties": {
                "description": {
                    "description": "描述",
                    "type": "string"
                },
                "platform": {
                    "description": "平台名称",
                    "type": "string"
                },
                "url": {
                    "description": "个人主页链接",
                    "type": "string"
                },
                "username": {
                    "description": "平台用户名",
                    "type": "string"
                }
            }
        },
        "handler.SocialResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "platform": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                },
                "verified": {
                    "type": "boolean"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "DDUP API",
	Description:      "DDUP 服务 API 文档",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
