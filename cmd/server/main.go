package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/Joel-Haeberli/cde-fuzzer/internal/cde"
	"github.com/Joel-Haeberli/cde-fuzzer/internal/core"
)

type extractRequest struct {
	Text string `json:"text"`
}

var rules []core.Rule

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	// Load rules from directory specified by environment variable
	rulesDir := os.Getenv("RULES_DIR")
	if rulesDir != "" {
		var err error
		// Enable LLM rules by default in server mode
		rules, err = core.LoadRulesFromDirectory(rulesDir, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading rules: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Fallback to demo rule if no rules directory is provided
		fmt.Fprintln(os.Stderr, "No rules directory provided, using demo rule")
		rules = append(rules, core.NewRegexRule(
			"age-regex",
			regexp.MustCompile(`\b\d{1,3}\s*(?:years?\s*old|yo)\b`),
			0.85,
		))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /extract", handleExtract)

	fmt.Printf("cde-extractor server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func handleExtract(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req extractRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Demo extraction process (same as CLI).
	demoCDE := cde.CDE{
		ID:       "demo-1",
		Question: "What is the patient's age?",
	}

	chain := core.NewRuleChain("extraction-chain", rules...)
	estimator := &core.DefaultAccuracyEstimator{}
	process := core.NewExtractionProcess(demoCDE, chain, estimator)

	result, err := process.Run(req.Text)
	if err != nil {
		http.Error(w, "extraction failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
