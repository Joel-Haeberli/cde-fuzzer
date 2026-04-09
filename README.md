# CDE-Fuzzer

**Fuzzy algorithms for medical data extraction and synthetic report generation**

The tool is still under development. Verify your results and don't use the tool blindly. Do not hesitate filing an Issue if you find bugs or encounter problems.

## Quick Start

```bash
# Clone and build
git clone https://github.com/Joel-Haeberli/cde-fuzzer.git
cd cde-fuzzer
make build-linux  # or build-mac, build-windows

# Extract CDEs from medical text
./bin/cde-fuzzer-linux -file patient_report.txt -rules example_rules/

# Generate synthetic reports
./bin/generate-synthetic-linux -count 5 -variability 0.8

# Derive rules from medical data
./bin/derive-rules-linux -data your_medical_reports/ -output derived_rules/
```

## Overview

CDE-Fuzzer is a toolkit for extracting structured clinical data from unstructured medical text using fuzzy matching algorithms. The system processes Common Data Elements (CDEs) through pattern recognition and text analysis.

### Core Functionality: Fuzzy Data Extraction

The toolkit uses fuzzy matching algorithms to:
- Recognize patterns in medical text
- Extract [common data elements](https://cde.nlm.nih.gov) with context awareness
- Adapt to varied text formats and terminology
- Provide confidence scores for extraction results

### 1. CDE-Fuzzer - Extraction engine
- Extracts structured data from unstructured medical text
- Supports regex patterns, similarity matching, and LLM-based extraction
- Provides extraction traces with confidence estimates
- Available as command-line tool and HTTP service

### 2. Derive Rules - Rule generator
- Analyzes medical documents to identify potential CDEs
- Generates extraction rules from document patterns
- Prevents duplicate rules and manages file conflicts
- Supports TXT, CSV, DOCX, PDF, XLSX formats

### 3. Generate Synthetic - Report generator
- Creates synthetic medical reports from extracted data
- Configurable variability for content diversity
- Generates test datasets and training data
- Produces reports with medical terminology and structure

## Installation

### From Release (Recommended)
Download pre-built binaries from [GitHub Releases](https://github.com/Joel-Haeberli/cde-fuzzer/releases)

### From Source
```bash
git clone https://github.com/Joel-Haeberli/cde-fuzzer.git
cd cde-extractor
make build  # Builds all platforms
```

Binaries Available:
- cde-extractor-* - Main extraction tool
- cde-extractor-server-* - HTTP API server
- derive-rules-* - Rule derivation tool
- generate-synthetic-* - Synthetic report generator
- generate-report-* - Basic report generator
- generate-diverse-report-* - Advanced report generator

## Usage

### CDE Extraction
```bash
# Basic extraction
./cde-extractor-linux -file input.txt -rules example_rules/

# With custom rules
./cde-extractor-linux -file report.txt -rules your_rules/
```

[Full CDE Extractor Manual](manuals/CDE_EXTRACTOR_MANUAL.md)

### Rule Derivation
```bash
# Derive rules from documents
./derive-rules-linux -data medical_reports/ -output derived_rules/

# Recursive processing
./derive-rules-linux -data all_data/ -output rules/ -recursive
```

[Full Derive Rules Manual](manuals/DERIVE_RULES_MANUAL.md)

### Synthetic Report Generation
```bash
# Generate test reports
./generate-synthetic-linux -count 10 -variability 0.7

# Custom output directory
./generate-synthetic-linux -count 5 -output test_reports/
```

[Full Generate Synthetic Manual](manuals/GENERATE_SYNTHETIC_MANUAL.md)

## Key Features

### Fuzzy Extraction Engine
- Fuzzy matching algorithms for data extraction
- Multiple rule types: regex patterns, similarity matching, LLM-based extraction
- Pattern recognition for varied medical text formats
- Confidence scoring with extraction traces
- Cross-platform: Linux, macOS, Windows

### Rule Derivation
- Automatic CDE identification from document patterns
- Duplicate prevention for rule management
- Sequential numbering for rule versioning
- Multi-format document support: TXT, CSV, DOCX, PDF, XLSX
- Pattern-based rule generation

### Synthetic Reports
- Medical report generation from extracted data
- Configurable variability for content diversity
- Batch generation for dataset creation
- Medical terminology and structure
- Context-aware report assembly

## Complete Workflow

```
Medical Documents → Rule Derivation → CDE Extraction → Synthetic Reports → Testing/Training
```

Example Pipeline:
```bash
# 1. Derive rules from your data
./derive-rules-linux -data patient_reports/ -output my_rules/ -recursive

# 2. Extract CDEs using derived rules
./cde-extractor-linux -file new_report.txt -rules my_rules/

# 3. Generate synthetic reports for testing
./generate-synthetic-linux -count 20 -variability 0.8 -output test_data/

# 4. Validate extraction pipeline
for report in test_data/*.txt; do
  ./cde-extractor-linux -file "$report" -rules my_rules/ > "results/$(basename ${report%.txt}).json"
done
```

## Project Structure

```
cde-extractor/
├── bin/                  # Built binaries
├── cmd/                  # Source code for each tool
├── internal/             # Core libraries
├── example_rules/        # Sample extraction rules
├── manuals/              # Detailed documentation
├── templates/            # Report templates
└── test_data/            # Test files
```

## Example Rules

Rules are defined in YAML format:

Regex Rule:
```yaml
name: "age_extractor"
type: "regex"
pattern: "\b(\d{1,3})\s*(?:years?|Jahre?)\b"
accuracy: 0.95
```

LLM Rule:
```yaml
name: "diagnosis_extractor"
type: "llm"
prompt: "Extract the primary diagnosis from this radiology report."
accuracy: 0.90
```

Similarity Rule:
```yaml
name: "procedure_matcher"
type: "similarity"
target: "mammography"
threshold: 0.85
```

## Advanced Features

### Fuzzy Pattern Recognition
The core of CDE Extractor uses fuzzy matching algorithms for extraction from medical text:
- Pattern matching for variations in terminology and formatting
- Context-aware extraction of medical data elements
- Confidence-based results with accuracy estimation
- Processing for real-world medical documents

### Rule Generation
- CDE identification from document patterns
- Automatic rule generation from analyzed documents
- Pattern-based rule creation
- Rule refinement capabilities

### Data Management
- Duplicate prevention using comparison algorithms
- Sequential versioning for rule organization
- Conflict resolution for file management
- Storage with naming conventions

### CSV Parsing
- Structured data extraction with column matching
- CDE identification from headers and content
- Rule generation tailored to tabular data
- Parsing for CSV files

## Documentation

- [CDE Extractor Manual](manuals/CDE_EXTRACTOR_MANUAL.md) - Complete CLI guide
- [Derive Rules Manual](manuals/DERIVE_RULES_MANUAL.md) - Rule generation guide
- [Generate Synthetic Manual](manuals/GENERATE_SYNTHETIC_MANUAL.md) - Report generation guide
- [Release Process](manuals/RELEASE_PROCESS.md) - Version tagging and releases

## Related Tools

- ACR Appropriateness Criteria: Pre-built rules for ACR data
- Custom Templates: Flexible report formatting
- Rule Validation: Quality assurance for extraction rules

## License

GPLv3 License - See [LICENSE](LICENSE) for details.

## Support

For issues, questions, or contributions:
- GitHub Issues: Report bugs or request features
- Documentation: Check the manuals for detailed guides
- Contributing: Submit pull requests with improvements

---

Built for medical data processing using LLM: Mistral's Vibe (and some Claude).

credits: ideas, concepts, architecture and prompting by [häbu.ch](https://häbu.ch)
