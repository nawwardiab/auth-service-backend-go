package service

import (
	"errors"
	"fmt"
	"server/internal/model"
	"server/internal/repo"
)

var ErrForbidden = errors.New("not allowed to access this resource")
var ErrCannotDeleteDefault = errors.New("cannot delete default address")

type AddressService struct {
	addrRepo *repo.AddressRepo
}


func NewAddressService(addrRepo *repo.AddressRepo) *AddressService {
	return &AddressService{addrRepo: addrRepo}
}

func (s *AddressService) CreateAddress(userID int, a *model.Address) error {
  a.UId = userID
	if a.IsDefault {
		clearAccErr := s.addrRepo.ClearDefaultForUser(a.UId)
    if clearAccErr != nil {
      return fmt.Errorf("service: clearing previous defaults: %w", clearAccErr)
    }
  }

	createAccErr := s.addrRepo.CreateAddress(a)
  if createAccErr != nil {
    return fmt.Errorf("service: CreateAddress failed: %w", createAccErr)
  }
  return nil
}

// GetAddress retrieves a single address by its ID
func (s *AddressService) GetAddress(userID, id int) (*model.Address, error) {
  addr, fetchErr := s.addrRepo.GetByID(id)
  if fetchErr != nil {
    return nil, fmt.Errorf("service: GetAddress failed: %w", fetchErr)
  }
  if addr.UId != userID {
    return nil, ErrForbidden
  }
  return addr, nil
}


// DeleteAddress removes an address record
func (s *AddressService) DeleteAddress(userID, id int) error {
  addr, err := s.addrRepo.GetByID(id)
  if err != nil {
    return err
  }
  if addr.UId != userID {
    return ErrForbidden
  }
  if addr.IsDefault {
    return ErrCannotDeleteDefault
  }
   deleteErr := s.addrRepo.Delete(id)
	if deleteErr != nil {
    return fmt.Errorf("service: DeleteAddress failed: %w", deleteErr)
  }
  return nil
}

// UpdateAddress applies updates, enforcing ownership and single-default rules.
func (s *AddressService) UpdateAddress(userID int, a *model.Address) error {
	// fetch existing to check ownership
	existing, err := s.addrRepo.GetByID(a.ID)
	if err != nil {
		return fmt.Errorf("service: fetch existing: %w", err)
	}
	if existing.UId != userID {
		return ErrForbidden
	}

	// if setting new default, clear old ones
	if a.IsDefault {
		if err := s.addrRepo.ClearDefaultForUser(userID); err != nil {
			return fmt.Errorf("service: clearing previous defaults: %w", err)
		}
	}

	// perform update
	if err := s.addrRepo.Update(a); err != nil {
		return fmt.Errorf("service: update address: %w", err)
	}

	return nil
}