# Copilot Instructions for Terraform AWS Provider

## Repository Overview

This is the **Terraform AWS Provider**, a Go-based Terraform plugin that manages AWS resources. It enables infrastructure-as-code workflows for AWS services through Terraform configuration.

- **Language**: Go 1.24+
- **Type**: Terraform Provider Plugin
- **Size**: Large codebase with 250+ AWS services in `internal/service/`
- **Main entry**: `main.go`

## Build and Development

### Prerequisites

- Go 1.24+ (version specified in `.go-version`)
- Make
- Docker (for some linting checks)
- Terraform CLI 0.12.26+ (for acceptance tests)

### Essential Commands

**Always run `make tools` first** to install required development tools:

```bash
make tools
```

**Build the provider:**

```bash
make build
```

**Run unit tests:**

```bash
make test                    # All unit tests
make test PKG=ec2           # Specific service package
```

**Run acceptance tests** (requires AWS credentials):

```bash
make testacc TESTS=TestAccIAMRole_basic PKG=iam
```

**Lint and format code:**

```bash
make fmt                     # Format Go code
make golangci-lint PKG=ec2  # Run linters on specific package
make provider-lint PKG=ec2  # Provider-specific lints
make semgrep PKG=ec2        # Semgrep checks
```

**Run generators:**

```bash
make gen                     # Run all code generators
```

### Quick Validation Workflow

For changes to a specific service package (e.g., `ec2`):

```bash
make fmt
make golangci-lint1 PKG=ec2
make provider-lint PKG=ec2
make test PKG=ec2
```

## Project Structure

```
internal/
├── service/          # AWS service implementations (one dir per service)
│   ├── ec2/         # Example: EC2 resources, data sources, tests
│   ├── iam/         # Example: IAM resources, data sources, tests
│   └── ...          # 250+ service packages
├── provider/         # Provider configuration
├── acctest/          # Acceptance test utilities
├── conns/            # AWS client connections
├── flex/             # Data flattening/expansion helpers
├── framework/        # Plugin Framework utilities
├── sdkv2/            # SDK V2 utilities
└── sweep/            # Resource sweepers
names/                # Service naming constants and caps
website/docs/         # User documentation (r/ for resources, d/ for data sources)
docs/                 # Contributor documentation
.ci/                  # CI configuration and linting rules
skaff/                # Scaffolding tool for new resources
.changelog/           # Changelog entries (one file per PR)
```

## Adding New Resources

1. **Use skaff scaffolding tool:**
   ```bash
   make skaff
   skaff resource -n ExampleThing -s exampleservice
   ```

2. **Resource files go in:** `internal/service/<service>/<resource>.go`
3. **Data source files go in:** `internal/service/<service>/<resource>_data_source.go`
4. **Test files go in:** `internal/service/<service>/<resource>_test.go`
5. **Documentation goes in:** `website/docs/r/<service>_<resource>.html.markdown`

## Code Conventions

### Resource Registration

Resources use annotations for self-registration:

```go
// @FrameworkResource("aws_service_thing", name="Thing")
func newResourceThing(context.Context) (resource.ResourceWithConfigure, error)

// @SDKResource("aws_service_thing", name="Thing")
func ResourceThing() *schema.Resource
```

Run `make gen` after adding or modifying annotations.

### Naming Rules

- Package names: lowercase, no underscores (e.g., `accessanalyzer`)
- Resource names: `aws_<service>_<thing>` in snake_case
- Go functions: Mixed caps (e.g., `ResourceVPCEndpoint`, not `ResourceVpcEndpoint`)
- Keep initialisms uppercase: `VPC`, `IAM`, `EC2`, `ARN`, `ID`

### Test Naming

- Basic test: `TestAcc{Service}{Resource}_basic`
- Disappears test: `TestAcc{Service}{Resource}_disappears`
- Per-attribute test: `TestAcc{Service}{Resource}_{Attribute}`

## Changelog Requirements

For user-impacting changes, create `.changelog/{PR_NUMBER}.txt`:

```
```release-note:enhancement
resource/aws_example_thing: Add `new_attribute` attribute
```
```

Headers: `new-resource`, `new-data-source`, `enhancement`, `bug`, `note`, `breaking-change`

## CI Checks

Key CI checks that run on PRs:

- `golangci-lint` - Go linting (5 stages)
- `provider-lint` - Provider-specific rules
- `semgrep` - Code quality and naming scans
- `go_generate` - Ensures generated code is current
- `testacc-lint` - Terraform config formatting in tests
- `copyright` - License headers
- `website` - Documentation checks

Run `make ci-quick` to run most CI checks locally.

## Common Issues and Solutions

1. **Generated code out of sync:** Run `make gen` and commit changes
2. **Import order issues:** Run `make fix-imports`
3. **Semgrep constant errors:** Run `make semgrep-fix`
4. **Terraform format in tests:** Run `make testacc-lint-fix`
5. **Website formatting:** Run `make website-terrafmt-fix`

## Key Documentation

- Contributor Guide: `docs/` directory
- Makefile reference: `docs/makefile-cheat-sheet.md`
- CI details: `docs/continuous-integration.md`
- Adding resources: `docs/add-a-new-resource.md`
- Writing tests: `docs/running-and-writing-acceptance-tests.md`
- Error handling: `docs/error-handling.md`
- Naming conventions: `docs/naming.md`

## Trust These Instructions

These instructions reflect tested, working commands. Only search for additional information if these instructions are incomplete or produce errors.
