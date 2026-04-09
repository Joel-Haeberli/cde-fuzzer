# Derive Rules - Rule Derivation Tool

## Overview

The `derive-rules` tool automatically generates extraction rules from medical documents. It analyzes text files, identifies potential Common Data Elements (CDEs), and creates YAML rule configurations that can be used with the CDE-Fuzzer CLI.

## Installation

### From Release
Download the pre-built binary for your platform from the [GitHub Releases](https://github.com/Joel-Haeberli/cde-fuzzer/releases) page.

### From Source
```bash
git clone https://github.com/Joel-Haeberli/cde-fuzzer.git
cd cde-extractor
make build-linux  # or build-mac, build-windows
```

## Usage

### Basic Rule Derivation

```bash
./derive-rules-linux -data path/to/documents/ -output path/to/rules/
```

### Command Line Options

```
Usage: derive-rules -data <directory> -output <rules-dir> [-recursive]

Flags:
  -data string
        directory containing data sources (default "./data")
  -output string
        directory to save derived rules (default "./derived_rules")
  -recursive
        parse directories recursively
```

## How It Works

### 1. Document Parsing
The tool parses various document formats to extract text and metadata:
- **TXT files** - Plain text documents
- **CSV files** - Structured data with CDE identification
- **DOCX files** - Word documents (basic parsing)
- **PDF files** - PDF documents (basic parsing)
- **XLSX files** - Excel spreadsheets (basic parsing)

### 2. CDE Identification
Analyzes document content to identify potential Common Data Elements:
- **Patient identifiers** (age, ID, etc.)
- **Dates and times**
- **Clinical findings**
- **Medical procedures**
- **Test results**

### 3. Rule Generation
For each identified CDE, generates two types of rules:
- **Regex rule** - Pattern-based extraction
- **LLM rule** - AI-based extraction with prompts

### 4. Rule Saving
Saves rules as YAML files with automatic duplicate prevention and sequential numbering.

## Examples

### Example 1: Derive Rules from Patient Data

```bash
./derive-rules-linux -data ./patient_reports/ -output ./derived_rules/
```

### Example 2: Recursive Directory Processing

```bash
./derive-rules-linux -data ./all_medical_data/ -output ./rules/ -recursive
```

### Example 3: Process CSV Files

```bash
./derive-rules-linux -data ./cde_definitions/ -output ./cde_rules/
```

## Output Format

### Generated Rule Files

Each rule is saved as a YAML file with the following structure:

**Regex Rule Example:**
```yaml
name: "age_extractor"
type: "regex"
pattern: "\\b(\\d{1,3})\\s*(?:Jahre?|Jahr|yo|years?|months?|mos?|days?|dys?|Monate?|Tage?|Jahren?)\\b"
accuracy: 0.85
```

**LLM Rule Example:**
```yaml
name: "llm_age_extractor"
type: "llm"
prompt: "Extract the patient's age from the text. Look for numbers followed by age-related terms like 'years', 'Jahre', 'yo', etc. If no age information is present, respond with 'no data'."
accuracy: 0.90
```

**Note:** LLM rules use a mock client by default. To use real LLM APIs, set the `LLM_API_URL` and `LLM_API_KEY` environment variables. See the CDE Extractor Manual for details.

### File Naming Convention

- **Base name**: Derived from the CDE type (e.g., `age_extractor`)
- **Type prefix**: `llm_` for LLM rules
- **Number suffix**: Added if filename exists (e.g., `age_extractor_1.yaml`)
- **Extension**: `.yaml` for all rule files

## Advanced Features

### Duplicate Prevention

The tool automatically detects duplicate rules and skips them:

```
⚠️  Skipping duplicate rule: age_extractor
```

### Sequential Numbering

When rules would overwrite existing files, they are numbered sequentially:

```
age_extractor.yaml        # Original
age_extractor_1.yaml     # First duplicate
age_extractor_2.yaml     # Second duplicate
```

### CSV-Specific CDE Identification

For CSV files, the tool uses specialized pattern matching:
- Analyzes column headers for common CDE patterns
- Generates rules tailored to structured data
- Identifies patient IDs, dates, diagnoses, procedures, and results

## Rule Types Generated

### 1. Regex Rules

**Purpose**: Extract structured data using patterns

**Generated Patterns:**
- **Age**: Matches age expressions in multiple languages
- **Date**: Matches various date formats
- **Patient ID**: Matches ID patterns
- **Clinical Findings**: Matches diagnosis-related terms
- **Procedures**: Matches medical procedure names
- **Results**: Matches test result patterns

**Customization**: Edit the YAML files to refine patterns

### 2. LLM Rules

**Purpose**: Extract complex, unstructured data using AI

**Generated Prompts:**
- **Age**: "Extract the patient's age from the text..."
- **Date**: "Extract any dates mentioned in the text..."
- **Diagnosis**: "Extract the diagnosis or main findings..."
- **Procedures**: "Extract information about medical procedures..."
- **Results**: "Extract test results and findings..."

**Customization**: Edit prompts to improve extraction quality

## Best Practices

### Input Data Preparation

1. **Organize documents** by type (reports, CSV, etc.)
2. **Use consistent naming** for easy identification
3. **Clean text data** for better pattern matching
4. **Separate test data** from production data

### Rule Generation Strategy

1. **Start with CSV files** for structured CDE identification
2. **Process text reports** for unstructured data
3. **Review generated rules** before use
4. **Test rules** on sample data
5. **Refine rules** based on extraction results

### Rule Management

1. **Version your rules** for reproducibility
2. **Document rule purposes** in comments
3. **Group related rules** in subdirectories
4. **Archive old rules** instead of deleting
5. **Share rules** across projects

## Integration with CDE Extractor

### Using Derived Rules

```bash
# Derive rules
./derive-rules-linux -data ./reports/ -output ./rules/

# Use rules with CDE Extractor
./cde-extractor-linux -file new_report.txt -rules ./rules/
```

### Rule Refinement Workflow

```bash
# 1. Derive initial rules
./derive-rules-linux -data sample_reports/ -output draft_rules/

# 2. Test rules
./cde-extractor-linux -file test_report.txt -rules draft_rules/

# 3. Refine rules manually
# Edit draft_rules/*.yaml files

# 4. Test refined rules
./cde-extractor-linux -file test_report.txt -rules refined_rules/

# 5. Deploy to production
cp refined_rules/* production_rules/
```

## Performance Tips

### Large Datasets

- **Process in batches**: `find reports/ -name "*.txt" | xargs -n 10 ./derive-rules-linux -data {} -output batch_rules/`
- **Use recursive mode**: For deeply nested directory structures
- **Parallel processing**: Run multiple instances on different data subsets

### Rule Quality

- **Review generated rules**: Not all auto-generated rules may be useful
- **Test on sample data**: Validate rules before full deployment
- **Combine approaches**: Use both regex and LLM rules for comprehensive extraction

## Troubleshooting

### Common Issues

**No CDEs found:**
- Check document content and format
- Verify file types are supported
- Test with simpler documents first

**Duplicate rules:**
- This is expected behavior for similar documents
- Use `-recursive` to process all documents at once
- Manually curate final rule set

**File parsing errors:**
- Ensure files are not corrupted
- Check file permissions
- Verify file encoding (UTF-8 recommended)

### Debug Output

The tool provides detailed logging:
```
🔍 CDE Extractor - Rule Derivation Tool
📁 Scanning directory: data/
📄 Found 5 documents
Analyzing document: data/report1.txt
Found 3 potential CDEs
🔧 Derived 6 rules
⚠️  Skipping duplicate rule: age_extractor
Saved rule: rules/date_extractor.yaml
```

## Example Workflows

### Workflow 1: Medical Report Processing

```bash
# 1. Organize reports
mkdir -p organized_reports/{mammography,pathology,biopsy}
mv reports/*mammography*.txt organized_reports/mammography/
# ... organize other report types

# 2. Derive rules for each type
./derive-rules-linux -data organized_reports/mammography/ -output rules/mammography/
./derive-rules-linux -data organized_reports/pathology/ -output rules/pathology/
./derive-rules-linux -data organized_reports/biopsy/ -output rules/biopsy/

# 3. Combine and deduplicate rules
./derive-rules-linux -data organized_reports/ -output rules/combined/ -recursive
```

### Workflow 2: CDE Definition Processing

```bash
# Process CSV files with CDE definitions
./derive-rules-linux -data cde_definitions/ -output cde_rules/

# Review and refine CDE rules
# Edit cde_rules/*.yaml files

# Use refined CDE rules for extraction
./cde-extractor-linux -file patient_report.txt -rules cde_rules/
```

### Workflow 3: Research Data Preparation

```bash
# 1. Derive rules from sample data
./derive-rules-linux -data sample_reports/ -output research_rules/

# 2. Test rules on validation set
for report in validation_reports/*.txt; do
  ./cde-extractor-linux -file "$report" -rules research_rules/ > "validation_results/$(basename ${report%.txt}).json"
done

# 3. Analyze extraction quality
# Review validation_results/*.json

# 4. Refine rules based on analysis
# Edit research_rules/*.yaml

# 5. Apply to full dataset
for report in full_dataset/*.txt; do
  ./cde-extractor-linux -file "$report" -rules research_rules/ > "extracted_data/$(basename ${report%.txt}).json"
done
```

## Limitations

- **Basic document parsing**: Simplified parsers for DOCX/PDF/XLSX
- **Pattern-based CDE identification**: May miss complex CDEs
- **Rule generation quality**: Auto-generated rules may need refinement
- **Performance**: Large document collections may take time

## Future Enhancements

- **Advanced document parsing**: Full text extraction from DOCX/PDF
- **Machine learning CDE identification**: Better pattern discovery
- **Rule quality scoring**: Automatic rule evaluation
- **Interactive rule refinement**: GUI for rule editing
- **Batch processing optimization**: Parallel document processing

## Support

For issues, questions, or contributions:
- **GitHub Issues**: Report bugs or request features
- **Documentation**: Check the project wiki
- **Contributing**: Submit pull requests with improvements

## License

This tool is released under the MIT License. See LICENSE file for details.