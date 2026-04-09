package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

// RuleConfig represents the configuration for a rule in YAML.
type RuleConfig struct {
	Name     string  `yaml:"name"`
	Type     string  `yaml:"type"`
	Pattern  string  `yaml:"pattern,omitempty"`
	Accuracy float64 `yaml:"accuracy,omitempty"`
	Target   string  `yaml:"target,omitempty"`
	Threshold float64 `yaml:"threshold,omitempty"`
	Prompt   string  `yaml:"prompt,omitempty"`
}

// LoadRulesFromDirectory loads all rule configurations from YAML files in the specified directory.
func LoadRulesFromDirectory(dirPath string, enableLLM bool) ([]Rule, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var rules []Rule
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".yaml" && filepath.Ext(file.Name()) != ".yml" {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
		}

		var ruleConfig RuleConfig
		if err := yaml.Unmarshal(data, &ruleConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML in %s: %v", filePath, err)
		}

		// Skip LLM rules if LLM is not enabled
		if ruleConfig.Type == "llm" && !enableLLM {
			fmt.Printf("⚠️  Skipping LLM rule %s (LLM not enabled, use --llm flag)\n", ruleConfig.Name)
			continue
		}

		rule, err := createRuleFromConfig(ruleConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule from %s: %v", filePath, err)
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func createRuleFromConfig(config RuleConfig) (Rule, error) {
	switch config.Type {
	case "regex":
		if config.Pattern == "" {
			return nil, fmt.Errorf("pattern is required for regex rule")
		}
		regex, err := regexp.Compile(config.Pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %v", err)
		}
		return NewRegexRule(config.Name, regex, config.Accuracy), nil

	case "similarity":
		if config.Target == "" {
			return nil, fmt.Errorf("target is required for similarity rule")
		}
		return NewSimilarityRule(config.Name, config.Target, config.Threshold, Levenshtein), nil

	case "llm":
		if config.Prompt == "" {
			return nil, fmt.Errorf("prompt is required for llm rule")
		}
		// Try to use RemoteLLMClient if environment variables are set, otherwise fall back to MockLLMClient
		var llmClient LLMClient
		if os.Getenv("LLM_API_URL") != "" && os.Getenv("LLM_API_KEY") != "" {
			remoteClient, err := NewRemoteLLMClient()
			if err != nil {
				fmt.Printf("Warning: Failed to create RemoteLLMClient: %v, falling back to MockLLMClient\n", err)
				llmClient = &MockLLMClient{}
			} else {
				llmClient = remoteClient
			}
		} else {
			llmClient = &MockLLMClient{}
		}
		return NewLLMRule(config.Name, config.Prompt, config.Accuracy, llmClient), nil

	default:
		return nil, fmt.Errorf("unknown rule type: %s", config.Type)
	}
}