# Generate Synthetic - Synthetic Report Generator

## Overview

The `generate-synthetic` tool creates realistic synthetic medical reports for testing, training, and development purposes. It generates diverse radiology reports with configurable variability to simulate real-world medical documentation.

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

### Basic Report Generation

```bash
./generate-synthetic-linux -count 5 -variability 0.8
```

### Command Line Options

```
Usage: generate-synthetic -count <number> -variability <factor> [-output <directory>]

Flags:
  -count int
        number of synthetic reports to generate (default 3)
  -variability float
        variability factor (0.0-1.0) for report diversity (default 0.8)
  -output string
        directory to save synthetic reports (default "./synthetic_reports")
```

## How It Works

### 1. Data Generation
Creates realistic medical data including:
- **Patient information** (age, gender)
- **Clinical scenarios** (presenting issues)
- **Exam indications** (reason for examination)
- **Procedures performed** (imaging techniques)
- **Findings** (observations and results)
- **Impressions** (diagnoses and conclusions)
- **Recommendations** (follow-up actions)

### 2. Report Assembly
Combines generated data into structured report formats with:
- **Standardized sections** (Patient Info, Clinical Info, Procedures, Findings, etc.)
- **Medical terminology** (BIRADS categories, procedure names)
- **Realistic formatting** (headers, section labels, spacing)

### 3. Variability Control
Adjusts report diversity based on the variability factor:
- **Low variability (0.0-0.3)**: More uniform reports
- **Medium variability (0.4-0.7)**: Balanced diversity
- **High variability (0.8-1.0)**: Very diverse reports

### 4. File Output
Saves reports as text files with sequential numbering.

## Examples

### Example 1: Generate Standard Reports

```bash
./generate-synthetic-linux -count 10 -variability 0.6
```

### Example 2: Generate Highly Diverse Reports

```bash
./generate-synthetic-linux -count 20 -variability 0.9 -output training_data/
```

### Example 3: Generate Minimal Reports

```bash
./generate-synthetic-linux -count 3 -variability 0.2
```

## Output Format

### Report Structure

```
[HEADER]

PATIENT INFORMATION:
[Age, Gender]

CLINICAL INFORMATION:
[Clinical Scenario, Exam Indication]

PROCEDURES PERFORMED:
[Imaging Procedure, Optional: Contrast Usage]

FINDINGS:
[Observations, Optional: Additional Details]

IMPRESSION:
[Diagnosis, BIRADS Category]

RECOMMENDATIONS:
[Follow-up Actions]

ADDITIONAL INFORMATION:
[Optional: Radiation Dose, Comparisons]

[FOOTER]
```

### Sample Report

```
RADIOLOGY REPORT

PATIENT INFORMATION:
Patient: 45-year-old, female

CLINICAL INFORMATION:
Patient presents with right breast lump detected during self-examination. Indication: diagnostic workup.

PROCEDURES PERFORMED:
Bilateral mammography with tomosynthesis

FINDINGS:
The breasts demonstrate heterogeneous fibroglandular tissue. No suspicious microcalcifications are present.

IMPRESSION:
1. Probably benign finding. BIRADS category 3. Short-term follow-up recommended.

RECOMMENDATIONS:
Follow-up imaging in 6 months is suggested for the probably benign finding.

ADDITIONAL INFORMATION:
Radiation dose: 2.5 mGy

End of report
```

## Variability Control

### How Variability Affects Reports

**Variability: 0.2 (Low)**
- Minimal additional details
- Standard phrasing
- Basic findings
- Few optional sections

**Variability: 0.6 (Medium)**
- Some additional details
- Mixed phrasing
- Expanded findings
- Some optional sections

**Variability: 0.9 (High)**
- Maximum additional details
- Diverse phrasing
- Comprehensive findings
- Most optional sections
- Randomized structure

### Variability Examples

**Low Variability (0.2):**
```
FINDINGS:
The breasts demonstrate heterogeneous fibroglandular tissue.
```

**High Variability (0.9):**
```
FINDINGS:
The breasts demonstrate heterogeneous fibroglandular tissue. No suspicious microcalcifications are present. The skin and nipple-areolar complexes appear normal. No architectural distortion or focal asymmetry is seen. The findings are stable compared to prior examination.
```

## Advanced Usage

### Batch Generation for Training Data

```bash
# Generate 100 diverse reports for model training
./generate-synthetic-linux -count 100 -variability 0.85 -output training_dataset/

# Generate 50 uniform reports for baseline testing
./generate-synthetic-linux -count 50 -variability 0.3 -output baseline_dataset/
```

### Integration with Extraction Pipeline

```bash
# 1. Generate synthetic reports
./generate-synthetic-linux -count 10 -variability 0.7 -output test_reports/

# 2. Test extraction rules
for report in test_reports/*.txt; do
  ./cde-extractor-linux -file "$report" -rules ./rules/ > "extraction_results/$(basename ${report%.txt}).json"
done

# 3. Analyze extraction quality
# Review extraction_results/*.json
```

### Report Quality Assessment

```bash
# Generate reports with different variability
for variability in 0.2 0.5 0.8; do
  ./generate-synthetic-linux -count 5 -variability $variability -output "quality_test/v${variability}/"
done

# Test extraction on each variability level
for dir in quality_test/v*; do
  echo "Testing variability: $(basename $dir)"
  for report in "$dir"/*.txt; do
    ./cde-extractor-linux -file "$report" -rules ./rules/ >> "results_$(basename $dir).txt"
  done
done
```

## Report Content Details

### Patient Information

**Age:** Random selection from realistic distribution (25-81 years)
**Gender:** Male or female

**Examples:**
- "45-year-old, female"
- "67-year-old, male"
- "32-year-old, female"

### Clinical Scenarios

**Common Scenarios:**
- Patient presents with right breast lump
- Follow-up examination for known breast cancer
- Screening mammography for high-risk patient
- Evaluation of palpable mass in left breast
- Post-treatment surveillance
- Assessment of breast pain and nipple discharge
- Preoperative staging for newly diagnosed carcinoma
- Routine screening examination

### Exam Indications

**Indication Types:**
- screening
- diagnostic workup
- follow-up
- staging
- treatment response assessment
- preoperative planning
- surveillance

### Procedures Performed

**Imaging Procedures:**
- Bilateral mammography
- Digital breast tomosynthesis
- Breast MRI with contrast
- Targeted ultrasound examination
- Stereotactic core biopsy
- MRI-guided vacuum-assisted biopsy
- Diagnostic mammography with magnification views
- Whole breast ultrasound

**Contrast Usage (optional):**
- intravenous contrast administration
- gadolinium-based contrast agent
- contrast-enhanced imaging

### Findings

**Base Findings:**
- The breasts demonstrate heterogeneous fibroglandular tissue
- Scattered fibroglandular densities are present bilaterally
- The breast parenchyma shows age-appropriate involution
- No suspicious masses, calcifications, or architectural distortions
- Benign-appearing cysts in the upper outer quadrants
- Well-circumscribed oval mass in specific location
- Multiple bilateral simple cysts
- Implants appear intact with no evidence of rupture

**Additional Details (with variability):**
- No suspicious microcalcifications present
- No axillary adenopathy identified
- Skin and nipple-areolar complexes appear normal
- No architectural distortion or focal asymmetry
- Findings stable compared to prior examination
- No significant interval change noted

### Impressions

**BIRADS Categories:**
- **BIRADS 1**: Normal screening examination
- **BIRADS 2**: Benign findings
- **BIRADS 3**: Probably benign finding (short-term follow-up)
- **BIRADS 4**: Suspicious abnormality (biopsy consideration)
- **BIRADS 5**: Highly suggestive of malignancy

**Impression Examples:**
- "1. Normal screening examination. BIRADS category 1."
- "2. Benign findings. BIRADS category 2."
- "3. Probably benign finding. BIRADS category 3. Short-term follow-up recommended."
- "No evidence of malignancy or suspicious findings."
- "Findings consistent with benign breast disease."

### Recommendations

**Follow-up Actions:**
- Correlation with clinical examination recommended
- Follow-up imaging in 6 months suggested
- Biopsy recommended for suspicious mass
- Continue annual screening mammography
- No additional imaging necessary at this time
- Clinical correlation and biopsy consideration advised
- Multidisciplinary tumor board review recommended

### Additional Information

**Radiation Dose (optional):**
- 2.5 mGy, 3.1 mGy, 1.8 mGy, 2.2 mGy
- within standard limits
- as low as reasonably achievable

**Comparisons (optional):**
- Findings stable compared to prior examination from 6 months ago
- No significant interval change noted
- Current findings represent expected postoperative changes
- Mass has decreased in size since previous study
- New findings identified since last examination

## Performance Characteristics

### Generation Speed
- **5-10 reports/second** on typical hardware
- **Scalable** for batch generation
- **Low memory usage** per report

### Variability Impact
- **Low variability (0.2)**: Fastest generation
- **High variability (0.9)**: More processing for additional details
- **Optimal balance (0.6-0.8)**: Recommended for most use cases

### File Size
- **Average report**: 500-1500 bytes
- **1000 reports**: ~1-1.5 MB total
- **Compressible**: Text format compresses well

## Use Cases

### 1. Extraction Rule Testing
Generate synthetic reports to validate and optimize extraction rules before using on real data.

### 2. Model Training
Create diverse training datasets for machine learning models in medical NLP.

### 3. Software Development
Use synthetic data for testing and debugging medical text processing applications.

### 4. Quality Assurance
Test extraction pipelines with controlled, realistic test data.

### 5. Demonstration
Showcase extraction capabilities with sample synthetic reports.

### 6. Benchmarking
Compare different extraction approaches using standardized synthetic data.

## Best Practices

### Report Generation Strategy

1. **Start with medium variability** (0.6-0.8) for balanced diversity
2. **Generate test batches** before full datasets
3. **Review sample reports** to ensure quality
4. **Adjust variability** based on use case needs
5. **Document generation parameters** for reproducibility

### Data Organization

```bash
# Organize by use case
mkdir -p synthetic_data/{training,testing,validation,demo}

# Generate different sets
./generate-synthetic-linux -count 100 -variability 0.8 -output synthetic_data/training/
./generate-synthetic-linux -count 50 -variability 0.6 -output synthetic_data/testing/
./generate-synthetic-linux -count 20 -variability 0.9 -output synthetic_data/demo/
```

### Quality Control

```bash
# Generate sample reports
./generate-synthetic-linux -count 3 -variability 0.7 -output quality_check/

# Manually review
cat quality_check/*.txt

# Adjust parameters if needed
# Then generate full dataset
```

## Troubleshooting

### Common Issues

**Reports too similar:**
- Increase variability factor (try 0.8-0.9)
- Check that random seed is working

**Reports too diverse:**
- Decrease variability factor (try 0.3-0.5)
- Review base templates

**Performance issues:**
- Reduce count parameter for large batches
- Process in smaller batches
- Check system resources

**Output directory issues:**
- Verify write permissions
- Check directory path exists
- Review file naming conflicts

### Debugging

The tool provides progress logging:
```
🔮 CDE Extractor - Synthetic Report Generator
📝 Generating 5 synthetic reports with variability 0.8
✅ Generated: ./synthetic_reports/synthetic_report_1.txt
✅ Generated: ./synthetic_reports/synthetic_report_2.txt
✅ Generated: ./synthetic_reports/synthetic_report_3.txt
🎉 Successfully generated 5 synthetic reports
📁 Reports saved to: ./synthetic_reports/
```

## Integration Examples

### With CDE Extractor Pipeline

```bash
# Full pipeline: generate reports → extract data → analyze results

# 1. Generate synthetic reports
./generate-synthetic-linux -count 10 -variability 0.7 -output test_reports/

# 2. Derive rules from real data
./derive-rules-linux -data real_reports/ -output extracted_rules/

# 3. Test rules on synthetic data
for report in test_reports/*.txt; do
  ./cde-extractor-linux -file "$report" -rules extracted_rules/ > "test_results/$(basename ${report%.txt}).json"
done

# 4. Analyze test results
# Compare with expected outcomes
```

### With Machine Learning Pipeline

```bash
# 1. Generate training data
./generate-synthetic-linux -count 1000 -variability 0.85 -output ml_training/train/
./generate-synthetic-linux -count 200 -variability 0.85 -output ml_training/val/
./generate-synthetic-linux -count 100 -variability 0.85 -output ml_training/test/

# 2. Extract labels using CDE Extractor
for report in ml_training/train/*.txt; do
  ./cde-extractor-linux -file "$report" -rules gold_rules/ > "ml_labels/train/$(basename ${report%.txt}).json"
done

# 3. Train model
# Use reports and labels for supervised learning

# 4. Evaluate on test set
# Test trained model on synthetic test data
```

## Limitations

- **Synthetic data**: Not real patient data
- **Pattern-based**: Uses templates and randomization
- **Domain-specific**: Focused on radiology/breast imaging
- **Language**: Primarily English medical terminology

## Future Enhancements

- **Custom templates**: User-defined report structures
- **Domain expansion**: Support for other medical specialties
- **Multi-language**: International medical terminology
- **Real data mixing**: Combine synthetic with real data
- **Conditional generation**: Rule-based report variations
- **Interactive editing**: GUI for report customization

## Support

For issues, questions, or contributions:
- **GitHub Issues**: Report bugs or request features
- **Documentation**: Check the project wiki
- **Contributing**: Submit pull requests with improvements

## License

This tool is released under the MIT License. See LICENSE file for details.