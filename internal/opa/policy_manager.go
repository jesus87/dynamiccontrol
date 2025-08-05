package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"dynamiccontrol/internal/types"

	"github.com/open-policy-agent/opa/rego"
)

// PolicyManager handles OPA policy loading and evaluation
type PolicyManager struct {
	policies map[string]*rego.PreparedEvalQuery
}

// NewPolicyManager creates a new policy manager
func NewPolicyManager() *PolicyManager {
	return &PolicyManager{
		policies: make(map[string]*rego.PreparedEvalQuery),
	}
}

// LoadPolicies loads all Rego policies from the policies directory
func (pm *PolicyManager) LoadPolicies(policiesDir string) error {
	files, err := ioutil.ReadDir(policiesDir)
	if err != nil {
		return fmt.Errorf("failed to read policies directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".rego") {
			continue
		}

		policyName := strings.TrimSuffix(file.Name(), ".rego")
		policyPath := filepath.Join(policiesDir, file.Name())

		if err := pm.loadPolicy(policyName, policyPath); err != nil {
			log.Printf("Failed to load policy %s: %v", policyName, err)
			continue
		}

		log.Printf("Loaded policy: %s", policyName)
	}

	return nil
}

// loadPolicy loads a single Rego policy file
func (pm *PolicyManager) loadPolicy(policyName, policyPath string) error {
	policyBytes, err := ioutil.ReadFile(policyPath)
	if err != nil {
		return fmt.Errorf("failed to read policy file %s: %w", policyPath, err)
	}

	query := rego.New(
		rego.Query("data."+policyName+".allow"),
		rego.Module(policyName+".rego", string(policyBytes)),
	)

	preparedQuery, err := query.PrepareForEval(context.Background())
	if err != nil {
		return fmt.Errorf("failed to prepare policy %s: %w", policyName, err)
	}

	pm.policies[policyName] = &preparedQuery
	return nil
}

// EvaluatePolicy evaluates a policy with the given input
func (pm *PolicyManager) EvaluatePolicy(policyName string, input map[string]interface{}) (*types.PolicyResult, error) {
	preparedQuery, exists := pm.policies[policyName]
	if !exists {
		return &types.PolicyResult{
			Allowed: false,
			Error:   fmt.Sprintf("Policy %s not found", policyName),
		}, nil
	}

	ctx := context.Background()
	results, err := preparedQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return &types.PolicyResult{
			Allowed: false,
			Error:   fmt.Sprintf("Policy evaluation error: %v", err),
		}, nil
	}

	if len(results) == 0 || len(results[0].Expressions) == 0 {
		return &types.PolicyResult{
			Allowed: false,
			Error:   "No policy result found",
		}, nil
	}

	allowed, ok := results[0].Expressions[0].Value.(bool)
	if !ok {
		return &types.PolicyResult{
			Allowed: false,
			Error:   "Policy result is not a boolean",
		}, nil
	}

	return &types.PolicyResult{
		Allowed: allowed,
	}, nil
}

// EvaluatePolicies evaluates multiple policies and returns combined result
func (pm *PolicyManager) EvaluatePolicies(policyNames []string, input map[string]interface{}) (*types.PolicyResult, error) {
	if len(policyNames) == 0 {
		return &types.PolicyResult{
			Allowed: true,
		}, nil
	}

	for _, policyName := range policyNames {
		result, err := pm.EvaluatePolicy(policyName, input)
		if err != nil {
			return &types.PolicyResult{
				Allowed: false,
				Error:   fmt.Sprintf("Policy evaluation error: %v", err),
			}, nil
		}

		if !result.Allowed {
			return &types.PolicyResult{
				Allowed: false,
				Error:   fmt.Sprintf("Policy %s denied the request", policyName),
			}, nil
		}
	}

	return &types.PolicyResult{
		Allowed: true,
	}, nil
}

// CreatePolicyInput creates the input for policy evaluation
func CreatePolicyInput(method, path string, headers map[string]string, body interface{}) map[string]interface{} {
	input := map[string]interface{}{
		"method":  method,
		"path":    path,
		"headers": headers,
	}

	if body != nil {
		// Convert body to map[string]interface{} for Rego evaluation
		if bodyBytes, err := json.Marshal(body); err == nil {
			var bodyMap map[string]interface{}
			if json.Unmarshal(bodyBytes, &bodyMap) == nil {
				input["body"] = bodyMap
			}
		}
	}

	return input
}

// ListLoadedPolicies returns a list of loaded policy names
func (pm *PolicyManager) ListLoadedPolicies() []string {
	policies := make([]string, 0, len(pm.policies))
	for policyName := range pm.policies {
		policies = append(policies, policyName)
	}
	return policies
}
