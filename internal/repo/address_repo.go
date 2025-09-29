package repo

import (
	"fmt"
	"server/internal/model"

	"github.com/jackc/pgx"
)

type AddressRepo struct {
	db *pgx.Conn
}

func NewAddressRepo(db *pgx.Conn) *AddressRepo {
	return  &AddressRepo{db: db}
}

// CreateAddress inserts a new address and populates a.ID, CreatedAt, UpdatedAt.
func (r *AddressRepo) CreateAddress(a *model.Address) error{
	query := `INSERT INTO addresses
      (u_id, addr_1, addr_2, zip, city, country, is_default)
    VALUES
      ($1,$2,$3,$4,$5,$6,$7)
    RETURNING id, created_at, updated_at;
	`
	row := r.db.QueryRow(query, a.UId, a.Addr_1, a.Addr_2, a.Zip, a.City, a.Country, a.IsDefault,)
  scanErr := row.Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	if scanErr != nil {
		return fmt.Errorf("Create Address: %w", scanErr)
	} else {
		return nil
	}
}

// ClearDefaultForUser sets defaut=false on all addresses for the given user.
func (r *AddressRepo) ClearDefaultForUser(userID int) error {
	query := `
		UPDATE addresses
  	  SET is_default = FALSE
 		WHERE u_id = $1;
	`
	_, execErr := r.db.Exec(query, userID)
  if execErr != nil {
    return fmt.Errorf("ClearDefaultForUser: %w", execErr)
  }
  return nil
}

// GetByID fetches a single address by its primary key.
func (r *AddressRepo) GetByID(id int) (*model.Address, error){
	query := `
    SELECT id, u_id, addr_1, addr_2, zip, city, country,  is_default, created_at, updated_at
      FROM addresses
    WHERE id = $1;
  `

  a := new(model.Address)
  row := r.db.QueryRow(query, id)
	scanErr := row.Scan(
    &a.ID, &a.UId,
    &a.Addr_1, &a.Addr_2,
    &a.Zip, &a.City, &a.Country,
    &a.IsDefault, &a.CreatedAt, &a.UpdatedAt,
  )
  if scanErr != nil {
    if scanErr == pgx.ErrNoRows {
        return nil, fmt.Errorf("GetByID: no address with id %d", id)
    }
    return nil, fmt.Errorf("GetByID: %w", scanErr)
  }

  return a, nil
}

// Delete removes an address by its ID.
func (r *AddressRepo) Delete(id int) error {
	query := `DELETE FROM addresses
	 WHERE id = $1;
	`
	_, execErr := r.db.Exec(query, id)
  if execErr != nil {
    return fmt.Errorf("Delete: %w", execErr)
  }
  return nil
}

// Update modifies an existing address, flipping the default flag if requested.
func (r *AddressRepo) Update(a *model.Address) error {
	const query = `
		UPDATE addresses
		   SET addr_1     = $1,
		       addr_2     = $2,
		       zip        = $3,
		       city       = $4,
		       country    = $5,
		       is_default = $6,
		       updated_at = now()
		 WHERE id = $7
	`

  _, execErr := r.db.Exec(query,
		a.Addr_1, a.Addr_2,
		a.Zip, a.City, a.Country,
		a.IsDefault, a.ID,
	)
	if execErr != nil {
		return fmt.Errorf("AddressRepo.Update: %w", execErr)
	}
	return nil
}