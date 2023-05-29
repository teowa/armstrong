package coverage_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/ms-henglu/armstrong/coverage"
)

func TestIndex(t *testing.T) {
	apiPath, modelName, modelSwaggerPath, err := coverage.PathPatternFromIdFromIndex(
		"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		"2022-06-01",
		"/home/wangta/.cache/armstrong/azure-rest-api-specs/specification",
		false,
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(*apiPath, *modelName, *modelSwaggerPath)

	model, err := coverage.Expand(*modelName, *modelSwaggerPath)
	if err != nil {
		t.Error(err)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Error(err)
	}
	out = out
	//fmt.Println("model", string(out))

	rawRequestJson := `
	{
	 "location": "eastus",
	 "properties": {
		"dataSources": {
		  "performanceCounters": [
			{
			  "name": "cloudTeamCoreCounters",
			  "streams": [
				"Microsoft-Perf"
			  ],
			  "samplingFrequencyInSeconds": 15,
			  "counterSpecifiers": [
				"\\Processor(_Total)\\% Processor Time",
				"\\Memory\\Committed Bytes",
				"\\LogicalDisk(_Total)\\Free Megabytes",
				"\\PhysicalDisk(_Total)\\Avg. Disk Queue Length"
			  ]
			},
			{
			  "name": "appTeamExtraCounters",
			  "streams": [
				"Microsoft-Perf"
			  ],
			  "samplingFrequencyInSeconds": 30,
			  "counterSpecifiers": [
				"\\Process(_Total)\\Thread Count"
			  ]
			}
		  ],
		  "windowsEventLogs": [
			{
			  "name": "cloudSecurityTeamEvents",
			  "streams": [
				"Microsoft-WindowsEvent"
			  ],
			  "xPathQueries": [
				"Security!"
			  ]
			},
			{
			  "name": "appTeam1AppEvents",
			  "streams": [
				"Microsoft-WindowsEvent"
			  ],
			  "xPathQueries": [
				"System![System[(Level = 1 or Level = 2 or Level = 3)]]",
				"Application!*[System[(Level = 1 or Level = 2 or Level = 3)]]"
			  ]
			}
		  ],
		  "syslog": [
			{
			  "name": "cronSyslog",
			  "streams": [
				"Microsoft-Syslog"
			  ],
			  "facilityNames": [
				"cron"
			  ],
			  "logLevels": [
				"Debug",
				"Critical",
				"Emergency"
			  ]
			},
			{
			  "name": "syslogBase",
			  "streams": [
				"Microsoft-Syslog"
			  ],
			  "facilityNames": [
				"syslog"
			  ],
			  "logLevels": [
				"Alert",
				"Critical",
				"Emergency"
			  ]
			}
		  ]
		},
		"destinations": {
		  "logAnalytics": [
			{
			  "workspaceResourceId": "/subscriptions/703362b3-f278-4e4b-9179-c76eaf41ffc2/resourceGroups/myResourceGroup/providers/Microsoft.OperationalInsights/workspaces/centralTeamWorkspace",
			  "name": "centralWorkspace"
			}
		  ]
		},
		"dataFlows": [
		  {
			"streams": [
			  "Microsoft-Perf",
			  "Microsoft-Syslog",
			  "Microsoft-WindowsEvent"
			],
			"destinations": [
			  "centralWorkspace"
			]
		  }
		]
	 }
	}
	`

	body := map[string]interface{}{}
	err = json.Unmarshal([]byte(rawRequestJson), &body)
	if err != nil {
		t.Error(err)
	}

	coverage.MarkCovered(body, model)

	var covered, uncovered []string
	coverage.SplitCovered(model, &covered, &uncovered)

	sort.Strings(covered)
	out, err = json.MarshalIndent(covered, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("coveredList", string(out))

	sort.Strings(uncovered)
	out, err = json.MarshalIndent(uncovered, "", "\t")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("uncoveredList", string(out))

	fmt.Printf("covered:%v, uncoverd: %v, total: %v\n", len(covered), len(uncovered), len(covered)+len(uncovered))

}
