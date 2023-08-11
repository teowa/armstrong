package report

import (
	_ "embed"
	"strings"

	"github.com/ms-henglu/armstrong/types"
)

//go:embed cleanup_error_report.md
var cleanupErrorReportTemplate string

func CleanupErrorMarkdownReport(report types.Error, logs []types.RequestTrace) string {
	parts := strings.Split(report.Type, "@")
	resourceType := ""
	apiVersion := ""
	if len(parts) == 2 {
		resourceType = parts[0]
		apiVersion = parts[1]
	}
	requestTraces := CleanupAllRequestTracesContent(report.Id, logs)
	content := cleanupErrorReportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", resourceType)
	content = strings.ReplaceAll(content, "${api_version}", apiVersion)
	content = strings.ReplaceAll(content, "${request_traces}", requestTraces)
	content = strings.ReplaceAll(content, "${error_message}", report.Message)
	return content
}

func CleanupAllRequestTracesContent(id string, logs []types.RequestTrace) string {
	content := ""
	for i := len(logs) - 1; i >= 0; i-- {
		if !strings.EqualFold(id, logs[i].ID) {
			continue
		}
		log := logs[i]
		if log.HttpMethod == "GET" && strings.Contains(log.Content, "REQUEST/RESPONSE") {
			st := strings.Index(log.Content, "GET https")
			ed := strings.Index(log.Content, ": timestamp=")
			trimContent := log.Content
			if st < ed {
				trimContent = log.Content[st:ed]
			}
			content = trimContent + "\n\n\n" + content
		} else if log.HttpMethod == "DELETE" {
			if strings.Contains(log.Content, "REQUEST/RESPONSE") {
				st := strings.Index(log.Content, "RESPONSE Status")
				ed := strings.Index(log.Content, ": timestamp=")
				trimContent := log.Content
				if st < ed {
					trimContent = log.Content[st:ed]
				}
				content = trimContent + "\n\n\n" + content
			} else if strings.Contains(log.Content, "OUTGOING REQUEST") {
				st := strings.Index(log.Content, "DELETE https")
				ed := strings.Index(log.Content, ": timestamp=")
				trimContent := log.Content
				if st < ed {
					trimContent = log.Content[st:ed]
				}
				content = trimContent + "\n\n" + content
			}
		}
	}

	return content
}
