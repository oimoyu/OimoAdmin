package front

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/oimoyu/OimoAdmin/src/utils/_type"
	"github.com/oimoyu/OimoAdmin/src/utils/functions"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
	"path"
	"strings"
)

func GeneratePages() error {
	//// skip if exist
	//if functions.IsPathExists(pagesPath) {
	//	return nil
	//}

	// login page
	loginPageResult := `
	{
	  "type": "page",
	  "name": "Login",
	  "css": {
		".login-content": {
		  "max-width": "500px"
		}
	  },
	  "className": "login-content",
	  "body": [
		{
		  "type": "form",
		  "api": "post:./login",
		  "title": "Login",
		  "body": [
			{
			  "name": "username",
			  "type": "input-text",
			  "label": "Username"
			},
			{
			  "name": "password",
			  "type": "input-password",
			  "label": "Password"
			}
		  ],
		  "actions": [
			  {
				"type": "submit",
				"label": "Login",
				"level": "primary"
			  }
		  ],
		  "onEvent": {
			"submitSucc": {
			  "actions": [
				{
				  "actionType": "custom",
				  "script": "window.location.hash = '/dashboard/';"
				}
			  ]
			}
		  }

		}
	  ],
	  "onEvent": {
		"init": {
		  "actions": [
			{
			  "actionType": "custom",
			  "script": ""
			}
		  ]
		}

	  }
	}
`

	// dashboard page
	dashboardPageResult := `
	{
	  "type": "page",
	  "name": "Dashboard",
	  "body": [
		{
			"type": "service",
			"api": "post:./sys_info",
			"interval": 100,
			"silentPolling": true,
			"body":[
				{
				  "type": "grid",
				"className": "text-center",
				  "columns": [
					{
						"body": [
							{
								"className": "m-0 text-center",
								"type": "static",
								"value": "CPU"
							},
							{
								"className": "place-content-center",
								"type": "progress",
								"value": "${cpu_percentage_usage}",
								"mode": "circle"
							},
							{
								"className": "m-0 text-center",
								"type": "static",
								"value": "${cpu_cores} Cores"
							}
						]
					},
					{
						"body": [
							{
								"className": "m-0 text-center",
								"type": "static",
								"value": "Memory"
							},
							{
								"className": "place-content-center",
								"type": "progress",
								"value": "${memory_percentage_usage}",
								"mode": "circle"
							},
							{
								"className": "m-0 text-center",
								"type": "static",
								"value": "${memory_usage} / ${memory_total} MB"
							}
						]
					},
					{
						"body": [
							{
								"className": "m-0 text-center",
								"type": "static",
								"value": "Disk"
							},
							{
								"className": "place-content-center",
								"type": "progress",
								"value": "${disk_percentage_usage}",
								"mode": "circle"
							},
							{
								"className": "m-0 text-center",
								"type": "static",
								"value": "${disk_usage} / ${disk_total} MB"
							}
						]
					},
					{
						"body": [
							{
								"className": "m-0 text-center",
								"type": "static",
								"value": "Network"
							},
							{
								"type": "tpl",
								"tpl": "<div class='mb-1'><i class='fa-solid fa-arrow-up'></i> ${upload_speed} MB/s</div>"
							},
							{
								"type": "tpl",
								"tpl": "<div class='mb-1'><i class='fa-solid fa-arrow-down'></i> ${download_speed} MB/s</div>"
							}
						]
					}
				  ]
				}
			]
		},
		{
			"type": "divider"
		},
		{
			"type": "service",
			"api": "post:./dashboard",
			"body":[
				{
					"type": "property",
					"title": "Database",
					"items": [
					  {
						"label": "DB Name",
						"content": "${db_name}",
						"span": 1.5
					  },
					  {
						"label": "Driver",
						"content": "${driver_name}",
						"span": 1.5
					  }
					]
				  }
			]
		},
		{
			"type": "divider"
		},
		{
			"type": "service",
			"id": "raw_sql_service",
			"api": "post:./raw_sql_records",
			"body":[
				{
				  "type": "form",
				  "id": "raw_sql_form",
				  "api": {
					"method": "post",
					"url": "./execute_raw_sql",
					"data": {
						"dry_run": true,
						"raw_sql": "$raw_sql"
					}
				  },
				  "title": "",
				  "body": [
					{
						"type": "each",
						"className": "pb-4",
						"name": "$raw_sql_records",
						"items": {
							"type": "tpl",
							"tpl": "<span class='label label-success m-l-sm'>${item}</span> ",
							"onEvent": {
								"click": {
								  "actions": [
									{
									  "actionType": "setValue",
									  "componentId": "raw_sql_form",
									  "args": {
										"value": {
											"raw_sql": "${item}"
										}
									  }
									}
								  ]
								}
							}
						}
					},
					{
						"type": "textarea",
						"name": "raw_sql",
						"label": "Execute Raw Sql",
						"placeholder": "UPDATE users set username='test' where id=123456"
					}
				  ],
				  "actions": [
					  {
						"type": "submit",
						"label": "Execute",
						"className": "float-left",
						"level": "success"
					  }
				  ],
				  "onEvent": {
					  "submitSucc": {
						"actions": [
							{
							  "actionType": "custom",
							  "script": ""
							},
							{
							  "actionType": "custom",
							  "script": "if(event.data.result.data.affect_rows===0){doAction({actionType: 'toast', args: {'msg': 'This Sql Affect 0 Rows.'}})};if(!event.data.result.data.affect_rows){event.stopPropagation()}"
							},
							{
							  "actionType": "confirmDialog",
							  "dialog": {
								"title": "Execute Confirm",
								"msg": "Are you sure to execute <br /><br /> <span style='color:red'>${raw_sql}</span> <br /><br />Affect Rows: <span style='color:red'>${result.data.affect_rows}</span>"
							  }
							},
							{
							  "actionType": "ajax",
							  "api": {
								"url": "./execute_raw_sql",
								"method": "post",
								"data": {
									"dry_run": false,
									"raw_sql": "$raw_sql"
								}
							  }
							},
							{
							  "actionType": "reload",
							  "componentId": "raw_sql_service"
							}
						]
					  }
				  }
				}
			]
		}
	  ]
	}
`
	logPageResult := `
  {
	"label": "Log",
	"url": "/log/",
	"icon": "fa-solid fa-clock",
	"schema": {
	  "type": "page",
	  "title": "Log",
	  "body": [
	  	{
			"type": "crud",
			"api": "post:./log",
			"columns": [
			  {
				"name": "create_time",
				"label": "Time"
			  },
			  {
				"name": "msg",
				"label": "Message"
			  }
			]
		}
	  ]
	}
  }
`

	tables := g.OimoAdmin.Tables
	var pagesTemplateResult string

	for _, table := range tables {
		pagesTemplateContent := `
  {
	"label": "{{ .Table.Name }}",
	"url": "/{{ .Table.Name }}/",
	"icon": "fa-solid fa-list",
	"schemaApi": {
		"method": "post",
		"url": "./page",
		"data": {
			"table_name": "{{ .Table.Name }}"
		}
	}


  },
`
		pageTemplateData := map[string]interface{}{
			"Table": table,
		}

		pageTemplateResult, err := functions.ParseTemplate(pagesTemplateContent, pageTemplateData)
		if err != nil {
			err := fmt.Errorf("parse pages template error: %w", err)
			g.OimoAdmin.Logger.Error(err.Error())
			panic(err)
		}

		pagesTemplateResult = pagesTemplateResult + pageTemplateResult
	}

	pagesTemplateResult = strings.TrimSpace(pagesTemplateResult)
	pagesTemplateResult = strings.TrimSuffix(pagesTemplateResult, ",")

	siteTemplateContent := `
{
  "status": 0,
  "msg": "",
  "data": {
    "pages": [
      {
        "label": "Home",
        "url": "/",
        "redirect": "/dashboard/"
      },
      {
        "label": "Login",
        "url": "/login/",
		"schema": {{ .LoginPage }}
      },
      {
        "children": [
		  {
			"label": "Dashboard",
			"url": "/dashboard/",
			"icon": "fa fa-tachometer",
			"schema": {{ .DashboardPage }}
		  },
          {{ .Pages }},
		  {{ .LogPage }}
        ]
      }
    ]
  }
}
`

	siteTemplateData := map[string]interface{}{
		"Pages":         pagesTemplateResult,
		"LoginPage":     loginPageResult,
		"DashboardPage": dashboardPageResult,
		"LogPage":       logPageResult,
	}
	siteTemplateResult, err := functions.ParseTemplate(siteTemplateContent, siteTemplateData)
	if err != nil {
		err := fmt.Errorf("parse site template error: %w", err)
		g.OimoAdmin.Logger.Error(err.Error())
		panic(err)
	}

	siteTemplateResult, err = functions.FormatJsonString(siteTemplateResult)
	if err != nil {
		err := fmt.Errorf("format json fail: %w", err)
		g.OimoAdmin.Logger.Error(err.Error())
		panic(err)
	}

	wwwRootDir := functions.GetWWWRootDir()
	pagesPath := path.Join(wwwRootDir, "site.json")

	err = functions.CreateFileWithContent(pagesPath, siteTemplateResult)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	return nil
}

func PageHandler(c *gin.Context) {
	var requestData struct {
		TableName string `json:"table_name" binding:"required"`
	}
	if err := c.ShouldBindBodyWith(&requestData, binding.JSON); err != nil {
		restful.ParamErr(c, fmt.Sprintf("parse input params failed: %v", err))
		return
	}

	var table *_type.TableStruct
	for _, currentTable := range g.OimoAdmin.Tables {
		if requestData.TableName == currentTable.Name {
			table = &currentTable
			break
		}
	}
	if table == nil {
		restful.ParamErr(c, "Unknown table")
		return
	}

	pageTemplateContent := `
	{
	  "type": "page",
	  "id": "fetch_list_page",
	  "title": "{{ .Table.Name }}",
	  "body": {
		"type": "crud",
		"id": "fetch_list_crud",
		"syncLocation": false,
		"api": {
			"method": "post",
			"url": "./fetch_list",
			"data": {
				"&": "$$",
				"page": "${page}",
				"perPage": "${perPage}",
				"table_name": "{{ .Table.Name }}"
			}
		},
		"quickSaveApi": {
			"method": "post",
			"url": "./update_rows",
			"data": {
				"&": "$$",
				"table_name": "{{ .Table.Name }}"
			}
		},
		"footerToolbar": [
		  "switch-per-page",
		  "pagination"
		],
		"defaultParams": {
		  "perPage": 50
		},

		"filter": {
			"title": null,
			"body": [
			  {
				"type": "input-text",
				"name": "keyword",
				"clearable": true,
				"placeholder": "search by keyword",
				"size": "sm"
			  },
			  {
				"type": "submit",
				"level": "primary",
				"label": "query"
			  },
			  {
				"type": "button",
				"label": "Create Row",
				"level": "success",
				"icon": "fa-solid fa-plus",
				"actionType": "dialog",
				"dialog": {
				  "title": "Create Row",
				  "body": {
					"type": "form",
					"api": {
						"method": "post",
						"url": "./create_row",
						"data": {
							"&": "$$",
							"table_name": "{{ .Table.Name }}"
						}
					},
					"reload": "fetch_list_crud",
					"body": [
						{{- range $index, $Column := .Table.Columns -}}
						{
						  "type": "input-text",
						  "name": "{{ $Column.Name }}",
						  "label": "{{ $Column.Name }}"
						},
						{{- end -}}
						{}
					]
				  }
				}
			  }
			]
		},
		"bulkActions": [
			{
			  "label": "Batch delete",
			  "actionType": "ajax",
			  "level": "danger",
			  "api": {
				"method": "post",
				"url": "./delete_rows",
				"data": {
					"&": "$$",
					"table_name": "{{ .Table.Name }}"
				}
			  },
			  "confirmText": "Are you sure to batch delete?"
			}
		],
		"columns": [
			{{- range $index, $Column := .Table.Columns -}}
			{
			  "name": "{{ $Column.Name }}",
			  "label": "{{ $Column.Name }}",
			  "sortable": true,
			  "quickEdit": 
				{{ if eq $Column.Name "id" }}
					false
				{{ else }}
					{
					  "type": "input-text"
					}
				{{ end }}
			},
			{{- end -}}
			{
			  "type": "operation",
			  "label": "Action",
			  "buttons": [
				{
				  "type": "button-group",
				  "buttons": [
					{
					  "type": "button",
					  "label": "Delete",
					  "level": "danger",
					  "actionType": "ajax",
					  "confirmText": "Are you sure to delete this row?",
					  "api": {
						"method": "post",
						"url": "./delete_rows",
						"data": {
							"ids": "${id},",
							"table_name": "{{ .Table.Name }}"
						}
					  }

					}

				  ]
				}
			  ],
			  "placeholder": "-",
			  "fixed": "right"
			}
		],
		"onEvent": {
			"quickSaveSucc": {
			  "actions": [
			  ]
			}
		}
	  }
	}
`
	pageTemplateData := map[string]interface{}{
		"Table": table,
	}
	pageTemplateResult, err := functions.ParseTemplate(pageTemplateContent, pageTemplateData)
	if err != nil {
		err := fmt.Errorf("parse page template error: %w", err)
		g.OimoAdmin.Logger.Error(err.Error())
		panic(err)
	}

	c.String(200, pageTemplateResult)
	return
}
