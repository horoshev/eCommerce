// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/": {
            "get": {
                "description": "Returns name, version and link to the swagger.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "general"
                ],
                "summary": "Welcome message from microservice.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.ServerResponse"
                        }
                    }
                }
            }
        },
        "/order": {
            "post": {
                "description": "Processing user order request. When ` + "`" + `uid` + "`" + ` param missed - creates new user and wallet as well.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Creates new order.",
                "parameters": [
                    {
                        "description": "User ID",
                        "name": "uid",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Order"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/orders": {
            "get": {
                "description": "Find and return created orders of the user using paging.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "orders"
                ],
                "summary": "Returns list of created orders for user.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Page size",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Order"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/requests": {
            "get": {
                "description": "Find and return created requests of the users using paging.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "requests"
                ],
                "summary": "Returns list of created requests for users.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Page size",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.UserRequest"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/requests/{uid}": {
            "get": {
                "description": "Find and return created requests of the user using paging.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "requests"
                ],
                "summary": "Returns list of created requests for user.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Page size",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.UserRequest"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.Response": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "x-order": "0"
                }
            }
        },
        "api.ServerResponse": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "x-order": "0"
                },
                "version": {
                    "type": "string",
                    "x-order": "1"
                },
                "swagger": {
                    "type": "string",
                    "x-order": "2"
                }
            }
        },
        "models.Order": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "x-order": "0"
                },
                "user_id": {
                    "type": "string",
                    "x-order": "1"
                },
                "status": {
                    "type": "string",
                    "x-order": "2"
                },
                "amount": {
                    "type": "number",
                    "x-order": "3"
                },
                "timestamp": {
                    "type": "string",
                    "x-order": "4"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.OrderProduct"
                    },
                    "x-order": "5"
                },
                "updates": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.OrderUpdate"
                    },
                    "x-order": "6"
                }
            }
        },
        "models.OrderProduct": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "x-order": "0"
                },
                "quantity": {
                    "type": "integer",
                    "x-order": "1"
                }
            }
        },
        "models.OrderUpdate": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "x-order": "0"
                },
                "status": {
                    "type": "string",
                    "x-order": "1"
                },
                "timestamp": {
                    "type": "string",
                    "x-order": "2"
                }
            }
        },
        "models.UserRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "x-order": "0"
                },
                "user_id": {
                    "type": "string",
                    "x-order": "1"
                },
                "type": {
                    "type": "string",
                    "x-order": "2"
                },
                "path": {
                    "type": "string",
                    "x-order": "3"
                },
                "timestamp": {
                    "type": "string",
                    "x-order": "4"
                }
            }
        }
    },
    "x-extension-openapi": {
        "example": "value on a json format"
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
	Version:     "1.0",
	Host:        "localhost:80",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Registry",
	Description: "This API is entrypoint for user requests.\nUser has option to order some products and list created orders.\nAlso, user is able to list his requests to this API.",
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
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
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
	swag.Register("swagger", &s{})
}
