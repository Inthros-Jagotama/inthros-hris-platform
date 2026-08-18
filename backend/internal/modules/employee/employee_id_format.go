package employee

import (
	"context"
	"errors"
)

// Employee ID generation modes (mirrors setting.EmployeeIDGenerationMode*).
const (
	EmployeeIDGenerationModeManual = "MANUAL"
	EmployeeIDGenerationModeHybrid = "HYBRID"
	EmployeeIDGenerationModeAuto   = "AUTO"
)

// ErrEmployeeIDRequired is returned when generation_mode is MANUAL and the
// caller did not supply an employee_id.
var ErrEmployeeIDRequired = errors.New("employee_id is required")

// EmployeeIDFormatProvider is a narrow, package-local interface for the
// cross-module dependency on setting.EmployeeIDFormatService — kept minimal
// (rather than importing the setting package's concrete type) so employee
// only depends on the two operations it actually needs. This mirrors the
// existing cross-module adapter pattern used elsewhere in main.go (e.g.
// employeeHireAdapter, workforceGapAdapter): the concrete setting service is
// wired in from main.go, employee only sees this interface.
type EmployeeIDFormatProvider interface {
	// GenerationMode returns the tenant's current employee_id generation
	// mode (MANUAL/HYBRID/AUTO).
	GenerationMode(ctx context.Context) (string, error)
	// Generate atomically produces the next employee_id per the tenant's
	// configured format/sequence.
	Generate(ctx context.Context) (string, error)
}

// resolveEmployeeID decides the EmployeeID value for a new employee per the
// tenant's configured generation mode:
//
//	MANUAL — req.EmployeeID is required; used as-is.
//	AUTO   — always generated; req.EmployeeID (even if supplied) is ignored.
//	HYBRID — req.EmployeeID used if non-empty, otherwise generated.
//
// When no provider is wired (e.g. older deployments before this feature, or
// unit tests that don't set one up), behavior falls back to the historical
// fully-manual behavior: req.EmployeeID is required and used as-is.
func (s *Service) resolveEmployeeID(ctx context.Context, reqEmployeeID string) (string, error) {
	if s.employeeIDFormat == nil {
		if reqEmployeeID == "" {
			return "", ErrEmployeeIDRequired
		}
		return reqEmployeeID, nil
	}

	mode, err := s.employeeIDFormat.GenerationMode(ctx)
	if err != nil {
		return "", err
	}

	switch mode {
	case EmployeeIDGenerationModeAuto:
		return s.employeeIDFormat.Generate(ctx)
	case EmployeeIDGenerationModeHybrid:
		if reqEmployeeID != "" {
			return reqEmployeeID, nil
		}
		return s.employeeIDFormat.Generate(ctx)
	default: // MANUAL (and any unrecognized value — fail safe to manual)
		if reqEmployeeID == "" {
			return "", ErrEmployeeIDRequired
		}
		return reqEmployeeID, nil
	}
}
