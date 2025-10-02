package scaffolding

import (
	"strings"
	"testing"
)

func TestBasicWorkflowTemplate(t *testing.T) {
	if basicWorkflowTemplate == "" {
		t.Error("basicWorkflowTemplate should not be empty")
	}
	
	// Check for essential keywords
	keywords := []string{"task(", "workflow.define", "description", "command"}
	for _, keyword := range keywords {
		if !strings.Contains(basicWorkflowTemplate, keyword) {
			t.Errorf("basicWorkflowTemplate should contain '%s'", keyword)
		}
	}
}

func TestCICDWorkflowTemplate(t *testing.T) {
	if cicdWorkflowTemplate == "" {
		t.Error("cicdWorkflowTemplate should not be empty")
	}
	
	// Check for CI/CD specific keywords
	keywords := []string{"build", "test", "deploy"}
	for _, keyword := range keywords {
		if !strings.Contains(cicdWorkflowTemplate, keyword) {
			t.Errorf("cicdWorkflowTemplate should contain '%s'", keyword)
		}
	}
}

func TestInfrastructureWorkflowTemplate(t *testing.T) {
	if infrastructureWorkflowTemplate == "" {
		t.Error("infrastructureWorkflowTemplate should not be empty")
	}
	
	// Check for infra specific keywords
	keywords := []string{"task("}
	for _, keyword := range keywords {
		if !strings.Contains(infrastructureWorkflowTemplate, keyword) {
			t.Errorf("infrastructureWorkflowTemplate should contain '%s'", keyword)
		}
	}
}

func TestDataPipelineWorkflowTemplate(t *testing.T) {
	if dataPipelineWorkflowTemplate == "" {
		t.Error("dataPipelineWorkflowTemplate should not be empty")
	}
	
	// Check for data pipeline specific keywords
	keywords := []string{"task("}
	for _, keyword := range keywords {
		if !strings.Contains(dataPipelineWorkflowTemplate, keyword) {
			t.Errorf("dataPipelineWorkflowTemplate should contain '%s'", keyword)
		}
	}
}

func TestMicroservicesWorkflowTemplate(t *testing.T) {
	if microservicesWorkflowTemplate == "" {
		t.Error("microservicesWorkflowTemplate should not be empty")
	}
	
	// Check for microservices specific keywords
	keywords := []string{"task("}
	for _, keyword := range keywords {
		if !strings.Contains(microservicesWorkflowTemplate, keyword) {
			t.Errorf("microservicesWorkflowTemplate should contain '%s'", keyword)
		}
	}
}

func TestConfigurationTemplates(t *testing.T) {
	if readmeTemplate == "" {
		t.Error("readmeTemplate should not be empty")
	}
	
	if gitignoreTemplate == "" {
		t.Error("gitignoreTemplate should not be empty")
	}
	
	if configTemplate == "" {
		t.Error("configTemplate should not be empty")
	}
}

func TestAllTemplatesHaveModernDSL(t *testing.T) {
	templates := []struct {
		name     string
		template string
	}{
		{"basic", basicWorkflowTemplate},
		{"cicd", cicdWorkflowTemplate},
		{"infrastructure", infrastructureWorkflowTemplate},
		{"microservices", microservicesWorkflowTemplate},
		{"dataPipeline", dataPipelineWorkflowTemplate},
	}
	
	for _, tt := range templates {
		t.Run(tt.name, func(t *testing.T) {
			// Check for Modern DSL pattern
			if !strings.Contains(tt.template, "task(") {
				t.Errorf("%s template should use task() function", tt.name)
			}
			if !strings.Contains(tt.template, ":build()") {
				t.Errorf("%s template should use :build() method", tt.name)
			}
		})
	}
}
