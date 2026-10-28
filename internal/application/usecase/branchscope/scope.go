package branchscope

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type Scope struct {
	BranchIDs []int64
}

func Load(ctx context.Context, schema string, userID int64, userRepo repository.UserRepository) (Scope, error) {
	user, err := userRepo.FindByID(ctx, schema, userID)
	if err != nil {
		return Scope{}, err
	}
	if user == nil {
		return Scope{}, exceptions.ErrUserNotFound
	}

	branchIDs := make([]int64, 0, len(user.Branches))
	for _, branch := range user.Branches {
		if branch == nil {
			continue
		}
		branchIDs = append(branchIDs, branch.ID)
	}

	return Scope{BranchIDs: branchIDs}, nil
}

func (s Scope) IsGlobal() bool {
	return len(s.BranchIDs) == 0
}

func (s Scope) Contains(branchID int64) bool {
	if s.IsGlobal() {
		return true
	}

	for _, allowedBranchID := range s.BranchIDs {
		if allowedBranchID == branchID {
			return true
		}
	}

	return false
}

func (s Scope) Apply(requestedBranchID int64) (int64, []int64, error) {
	if s.IsGlobal() {
		return requestedBranchID, nil, nil
	}

	if requestedBranchID != 0 {
		if !s.Contains(requestedBranchID) {
			return 0, nil, exceptions.ErrForbidden
		}
		return requestedBranchID, nil, nil
	}

	return 0, s.BranchIDs, nil
}

func (s Scope) Ensure(branchID int64) error {
	if s.IsGlobal() {
		return nil
	}
	if branchID != 0 && s.Contains(branchID) {
		return nil
	}

	return exceptions.ErrForbidden
}

func (s Scope) EnsureAll(branchIDs []int64) error {
	if s.IsGlobal() {
		return nil
	}

	if len(branchIDs) == 0 {
		return exceptions.ErrForbidden
	}

	for _, branchID := range branchIDs {
		if !s.Contains(branchID) {
			return exceptions.ErrForbidden
		}
	}

	return nil
}

func (s Scope) EnsureAny(branchIDs []int64) error {
	if s.IsGlobal() {
		return nil
	}

	for _, branchID := range branchIDs {
		if s.Contains(branchID) {
			return nil
		}
	}

	return exceptions.ErrForbidden
}
