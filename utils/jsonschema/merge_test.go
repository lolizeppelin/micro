package jsonschema

import (
	"encoding/json"
	"fmt"
	"testing"
)

const a = `{
    "switch": {
        "type": "object",
        "additionalProperties": false,
        "properties": {
            "create.enable": {
                "type": "boolean",
                "default": true,
                "description": "启用创角上报"
            },
            "levelup.enable": {
                "type": "boolean",
                "default": true,
                "description": "升级上报"
            },
            "login.enable": {
                "type": "boolean",
                "default": true,
                "description": "登录上报"
            },
            "create.retry": {
                "type": "boolean",
                "default": true,
                "description": "启用创角上报重试"
            },
            "sync.attribution": {
                "type": "boolean",
                "default": false,
                "description": "启用同步归因"
            }
        }
    },
    "uniques": {
        "type": "object",
        "required": [
            "token"
        ],
        "additionalProperties": false,
        "properties": {
            "token": {
                "type": "string",
                "minLength": 4,
                "maxLength": 64,
                "description": "引力accessToken,从引力获取"
            }
        }
    },
    "option": {
        "type": "object",
        "additionalProperties": false,
        "properties": {
            "engine.mode": {
                "type": "string",
                "default": "none",
                "enum": [
                    "debug",
                    "none"
                ],
                "description": "引力debugMode"
            },
            "retry.storage.key": {
                "type": "string",
                "minLength": 4,
                "maxLength": 24,
                "default": "yinli.user.retry",
                "description": "引力重上报换成用户id的本地存储key"
            },
            "engine.events": {
                "type": "array",
                "items": {
                    "type": "string",
                    "enum": [
                        "switch.server",
                        "click.login",
                        "change.name",
                        "leave.game",
                        "tutorial.finish",
                        "tutorial.step",
                        "view.mall",
                        "view.activity",
                        "mission.completed",
                        "stage.completed",
                        "stage.cleared",
                        "game.played",
                        "goods.delivered",
                        "reward.video.click",
                        "reward.video.close"
                    ],
                    "default": [
                        "tutorial.finish",
                        "view.mall",
                        "stage.completed",
                        "mission.completed",
                        "reward.video.click"
                    ]
                },
                "description": "额外上报事件"
            },
            "engine.ocean": {
                "type": "object",
                "additionalProperties": false,
                "description": "代理引力上报",
                "properties": {
                    "required": [
                        "id"
                    ],
                    "id": {
                        "type": "string",
                        "description": "占位key"
                    },
                    "login": {
                        "type": "string",
                        "default": "$MPRegister",
                        "description": "登录事件"
                    }
                }
            },
            "engine.gdt": {
                "type": "object",
                "required": [
                    "user_action_set_id",
                    "secret_key"
                ],
                "description": "代理广点通上报",
                "additionalProperties": false,
                "properties": {
                    "user_action_set_id": {
                        "type": "integer",
                        "minimum": 1,
                        "description": "广点通数据源id"
                    },
                    "quantity": {
                        "type": "integer",
                        "minimum": 1,
                        "default": 100,
                        "description": "广点通支付比率"
                    },
                    "metric": {
                        "type": "boolean",
                        "default": true,
                        "description": "广点通上报追踪"
                    },
                    "role": {
                        "type": "boolean",
                        "default": true,
                        "description": "上报角色创建"
                    },
                    "secret_key": {
                        "type": "string",
                        "minLength": 32,
                        "maxLength": 32,
                        "description": "广点通前端密钥"
                    }
                }
            }
        }
    }
}
`

const b = `{
    "switch": {
        "type": "object",
        "additionalProperties": true,
        "properties": {
            "engine.debug": {
                "type": "boolean",
                "default": false,
                "description": "调试模式启用广点通"
            },
            "expired.enabled": {
                "type": "boolean",
                "default": true,
                "description": "启用过期不报"
            },
            "disabled.metric": {
                "type": "boolean",
                "default": false,
                "description": "禁用metric"
            }
        }
    },
    "option": {
        "type": "object",
        "required": [
            "engine.user_action_set_id",
            "engine.secret_key"
        ],
        "additionalProperties": true,
        "properties": {
            "engine.quantity": {
                "type": "integer",
                "minimum": 1,
                "default": 100,
                "description": "广点通金额比例"
            },
            "engine.metric": {
                "type": "string",
                "format": "uri",
                "minLength": 8,
                "description": "打点追踪上报地址(广点通上报追踪)"
            },
            "engine.user_action_set_id": {
                "type": "integer",
                "minimum": 1,
                "description": "广点通数据源id"
            },
            "engine.secret_key": {
                "type": "string",
                "minLength": 32,
                "maxLength": 32,
                "description": "广点通前端密钥"
            },
            "expired.day": {
                "type": "integer",
                "minimum": 1,
                "default": 2,
                "description": "过期时间(单位天,超过指定天数的订单不上报)"
            }
        }
    }
}`

func TestMerge(t *testing.T) {

	ma := map[string]any{}
	mb := map[string]any{}

	if err := json.Unmarshal([]byte(a), &ma); err != nil {
		fmt.Println("failed decode a")
		return
	}
	if err := json.Unmarshal([]byte(b), &mb); err != nil {
		fmt.Println("failed decode b")
		return
	}

	_ma, _ := ma["option"].(map[string]any)
	_mb, _ := mb["option"].(map[string]any)

	_mc, err := MergeJSONSchema(_mb, _ma)
	if err != nil {
		fmt.Printf("merge failed: %s", err.Error())
		return
	}

	bf, err := json.MarshalIndent(_mc, "", "  ")
	if err != nil {
		fmt.Printf("failed decode c: %s", err.Error())
		return
	}

	fmt.Println(string(bf))

}
