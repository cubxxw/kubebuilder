```markdown
# kubebuilder Development Patterns

> Auto-generated skill from repository analysis

## Overview

This skill provides guidance for contributing to the `kubebuilder` project, a toolkit for building Kubernetes APIs using Go. It covers coding conventions, dependency upgrade workflows, and testing patterns found in the repository. Whether you're fixing bugs, adding features, or upgrading dependencies, following these patterns ensures consistency and maintainability.

## Coding Conventions

### File Naming

- Use **snake_case** for file names.
  - Example: `groupversion_info.go`, `zz_generated.deepcopy.go`

### Import Style

- Use **relative imports** within the module.
  - Example:
    ```go
    import (
        "github.com/example/kubebuilder/internal/controller"
        "github.com/example/kubebuilder/api/v1"
    )
    ```

### Export Style

- Use **named exports** for functions, types, and variables that should be accessible outside their package.
  - Example:
    ```go
    // Exported type
    type MyController struct {}

    // Exported function
    func NewMyController() *MyController {
        return &MyController{}
    }
    ```

### Commit Patterns

- Commit messages are **freeform** (no strict prefixes), averaging ~99 characters.

## Workflows

### Upgrade Kubernetes or Controller Runtime Dependencies

**Trigger:** When updating the supported Kubernetes or controller-runtime version and ensuring all sample/test projects reflect the change.  
**Command:** `/upgrade-k8s-deps`

**Step-by-step instructions:**

1. **Update Dependencies**
   - Edit `go.mod` and `go.sum` in each `testdata` and tutorial project to the new dependency versions.
     ```bash
     go get sigs.k8s.io/controller-runtime@vX.Y.Z
     go get k8s.io/api@vX.Y.Z
     go mod tidy
     ```
2. **Update CRD YAMLs and Generated Files**
   - Regenerate or update CRD YAMLs and code-generated files (e.g., `deepcopy.go`) to match new API versions or codegen output.
     ```bash
     make generate
     ```
   - Example:
     - `api/v1/zz_generated.deepcopy.go`
     - `config/crd/bases/mygroup.example.com_myresources.yaml`
3. **Update Main and Controller Files**
   - Modify `main.go` and controller files if there are breaking API or dependency changes.
     ```go
     // main.go
     import (
         "sigs.k8s.io/controller-runtime/pkg/manager"
     )
     ```
4. **Regenerate Deployment Artifacts**
   - Update Helm charts, install manifests, or other deployment files.
     - Example: `dist/chart/templates/crd/*.yaml`, `dist/install.yaml`
5. **Update Scaffolding Templates**
   - If necessary, update scaffolding templates to support new versions.
     - Example: `pkg/plugins/golang/vX/scaffolds/internal/templates/api/*.go`

**Files involved:**
- `**/go.mod`, `**/go.sum`
- `**/api/**/groupversion_info.go`, `**/api/**/types.go`, `**/api/**/zz_generated.deepcopy.go`
- `**/cmd/main.go`
- `**/internal/controller/**/*.go`
- `**/config/crd/bases/*.yaml`
- `**/dist/chart/templates/crd/*.yaml`
- `**/dist/install.yaml`
- `pkg/plugins/golang/**/scaffolds/internal/templates/api/*.go`

**Frequency:** ~1-2 times per major release

## Testing Patterns

- **Test File Pattern:** Files matching `*.test.*`
- **Testing Framework:** Not explicitly detected; likely uses Go's standard `testing` package.
- **Example:**
  ```go
  // api/v1/myresource_test.go
  import "testing"

  func TestMyResourceDefaults(t *testing.T) {
      // test implementation
  }
  ```

## Commands

| Command           | Purpose                                                                 |
|-------------------|-------------------------------------------------------------------------|
| /upgrade-k8s-deps | Upgrade controller-runtime or Kubernetes dependencies in all subprojects |

```