# CDE-Fuzzer - Command Line Interface

## Overview

The CDE-Fuzzer CLI is a tool for extracting Common Data Elements (CDEs) from medical text using configurable rules. It supports multiple rule types including regex patterns, similarity matching, and LLM-based extraction.

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

### Basic Extraction

```bash
./cde-extractor-linux -file input.txt -rules path/to/rules/
```

### Command Line Options

```
Usage: cde-extractor -file <path> -rules <rules-dir>

Flags:
  -file string
        path to text file to extract from
  -rules string
        path to directory containing rule YAML files
```

## Rule Configuration

Rules are defined in YAML files and placed in a rules directory. The CLI automatically loads all `.yaml` files from the specified directory.

### Rule Types

#### 1. Regex Rules

Extract data using regular expressions:

```yaml
name: "age_extractor"
type: "regex"
pattern: "\\b(\\d{1,3})\\s*(?:years?|yo|Jahre?)\\b"
accuracy: 0.95
```

**Pattern Examples:**
- Age: `\b(\d{1,3})\s*(?:years?|yo|Jahre?)\b`
- Date: `\b(\d{2}\.\d{2}\.\d{4}|\d{4}-\d{2}-\d{2})\b`
- Patient ID: `\b(?:Patienten-ID|ID):?\s*([A-Z0-9\-]+)\b`

#### 2. LLM Rules

Use language models for complex extraction:

```yaml
name: "llm_diagnosis_extractor"
type: "llm"
prompt: "Extract the diagnosis or main findings from the text. Look for sections headed by 'Diagnose', 'Befund', 'Findings', or 'Diagnosis' and extract the relevant clinical information."
accuracy: 0.90
```

**Prompt Tips:**
- Be specific about what to extract
- Mention common section headers
- Provide examples if helpful
- Specify output format

#### 3. Similarity Rules

Match based on string similarity:

```yaml
name: "procedure_matcher"
type: "similarity"
target: "mammography"
threshold: 0.85
accuracy: 0.80
```

### Rule Accuracy

Each rule has an `accuracy` field (0.0-1.0) that represents the confidence in the extraction. This is used to:
- Rank multiple matching rules
- Calculate overall extraction confidence
- Filter low-confidence results

## Input Formats

### Supported Text Formats
- Plain text files (`.txt`)
- Medical reports
- Clinical notes
- Radiology reports
- Pathology reports

### Input Example

```
PATIENT INFORMATION:
Age: 65 years
Gender: Female

CLINICAL HISTORY:
Patient presents with right breast lump detected during self-examination.

PROCEDURES PERFORMED:
Bilateral mammography with tomosynthesis

FINDINGS:
1. Right breast: 1.2 cm irregular mass at 3 o'clock position
2. No suspicious microcalcifications

IMPRESSION:
BIRADS category 4 - Suspicious abnormality
```

## Output Format

The CLI outputs structured extraction results in JSON format:

```json
{
  "CDEID": "age_extractor",
  "Answer": "65 years",
  "Accuracy": 0.95,
  "Traces": [
    {
      "RuleName": "age_extractor",
      "Match": {
        "Value": "65 years",
        "Start": 42,
        "End": 50
      },
      "Accuracy": 0.95
    }
  ]
}
```

## Advanced Usage

### Multiple Rules

The CLI can apply multiple rules to the same text. Rules are loaded from all YAML files in the specified directory.

### Rule Chaining

Rules are applied in sequence, and the best match is selected based on accuracy scores.

### Custom Rules

Create custom rules for your specific extraction needs by adding YAML files to your rules directory.

## Examples

### Example 1: Extract Patient Age

**Rule (`age_rule.yaml`):**
```yaml
name: "age_extractor"
type: "regex"
pattern: "\\b(\\d{1,3})\\s*(?:years?|yo|Jahre?|years old)\\b"
accuracy: 0.95
```

**Command:**
```bash
./cde-extractor-linux -file patient_report.txt -rules ./rules/
```

### Example 2: Extract Diagnosis with LLM

**Rule (`diagnosis_rule.yaml`):**
```yaml
name: "diagnosis_extractor"
type: "llm"
prompt: "Extract the primary diagnosis from this radiology report. Return only the diagnosis text."
accuracy: 0.90
```

**Command:**
```bash
./cde-extractor-linux -file radiology_report.txt -rules ./rules/
```

### Example 3: Batch Processing

```bash
# Process multiple files
for file in reports/*.txt; do
  echo "Processing $file..."
  ./cde-extractor-linux -file "$file" -rules ./rules/ > "results/$(basename ${file%.txt}).json"
done
```

## Performance Tips

### Rule Optimization
- **Start with high-accuracy rules** for critical data elements
- **Use LLM rules** for complex, unstructured data
- **Combine rule types** for comprehensive extraction
- **Test rules** on sample data before production use

### Large Files
- Process large files in chunks if memory is limited
- Use streaming approaches for very large datasets
- Consider parallel processing for batch operations

## Troubleshooting

### Common Issues

**No matches found:**
- Check rule patterns match your text format
- Verify text encoding (UTF-8 recommended)
- Test with simpler patterns first

**Low accuracy:**
- Adjust rule accuracy scores
- Refine regex patterns
- Improve LLM prompts with more context

**Performance issues:**
- Reduce number of rules if response is slow
- Optimize complex regex patterns
- Process files in batches

### Debugging

Add debug output by modifying the CLI code to log:
- Rules being loaded
- Match attempts and results
- Accuracy calculations

## Integration

### With Other Tools

```bash
# Pipe output to other tools
./cde-extractor-linux -file report.txt -rules ./rules/ | jq '.Answer' > extracted_data.txt

# Use in scripts
RESULT=$(./cde-extractor-linux -file report.txt -rules ./rules/ | jq -r '.Answer')
echo "Extracted: $RESULT"
```

### API Server

For programmatic access, use the `cde-extractor-server` binary which provides an HTTP API.

## Best Practices

1. **Start simple** - Begin with basic regex rules
2. **Test thoroughly** - Validate rules on diverse sample data
3. **Document rules** - Keep track of rule purposes and performance
4. **Version rules** - Maintain rule history for reproducibility
5. **Monitor accuracy** - Track extraction quality over time

## Rule Examples

### Common Medical Data Elements

**Age:**
```yaml
name: "age_extractor"
type: "regex"
pattern: "\\b(\\d{1,3})\\s*(?:years?|yo|Jahre?|years old|years-old)\\b"
accuracy: 0.95
```

**Date:**
```yaml
name: "date_extractor"
type: "regex"
pattern: "\\b(\\d{2}\\.\\d{2}\\.\\d{4}|\\d{4}-\\d{2}-\\d{2}|[A-Z][a-z]{2}\\s+\\d{1,2},\\s+\\d{4})\\b"
accuracy: 0.90
```

**Patient ID:**
```yaml
name: "patient_id_extractor"
type: "regex"
pattern: "\\b(?:Patient\\s*ID|ID|MRN):?\\s*([A-Z0-9\\-]+)\\b"
accuracy: 0.85
```

**Diagnosis:**
```yaml
name: "diagnosis_extractor"
type: "llm"
prompt: "Extract the primary diagnosis from this medical report. Look for terms like 'diagnosis', 'impression', 'findings', or 'conclusion'. Return only the diagnosis text without any labels or headers."
accuracy: 0.90
```

## Limitations

- **Text-only input** - Does not process images or binary formats
- **Rule complexity** - Very complex rules may impact performance
- **Language support** - Primarily English and German (extendable)
- **Context awareness** - Rules operate independently without global context

## Future Enhancements

- **Rule chaining** with conditional logic
- **Context-aware extraction** using document structure
- **Multi-language support** with language detection
- **Performance profiling** for rule optimization
- **Interactive rule testing** interface

## Support

For issues, questions, or contributions:
- **GitHub Issues**: Report bugs or request features
- **Documentation**: Check the project wiki
- **Contributing**: Submit pull requests with improvements

## License

This tool is released under the MIT License. See LICENSE file for details.