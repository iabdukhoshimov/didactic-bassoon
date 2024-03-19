// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package sqlc

import (
	"database/sql/driver"
	"fmt"
	"time"

	zero "gopkg.in/guregu/null.v4/zero"
)

type OrganizationStatus string

const (
	OrganizationStatusACTIVE    OrganizationStatus = "ACTIVE"
	OrganizationStatusPENDING   OrganizationStatus = "PENDING"
	OrganizationStatusFREETRIAL OrganizationStatus = "FREE_TRIAL"
	OrganizationStatusINACTIVE  OrganizationStatus = "INACTIVE"
	OrganizationStatusDELETED   OrganizationStatus = "DELETED"
	OrganizationStatusBLOCKED   OrganizationStatus = "BLOCKED"
)

func (e *OrganizationStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrganizationStatus(s)
	case string:
		*e = OrganizationStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrganizationStatus: %T", src)
	}
	return nil
}

type NullOrganizationStatus struct {
	OrganizationStatus OrganizationStatus `json:"organization_status"`
	Valid              bool               `json:"valid"` // Valid is true if OrganizationStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrganizationStatus) Scan(value interface{}) error {
	if value == nil {
		ns.OrganizationStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrganizationStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrganizationStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrganizationStatus), nil
}

type Roles string

const (
	RolesOWNER    Roles = "OWNER"
	RolesADMIN    Roles = "ADMIN"
	RolesUSER     Roles = "USER"
	RolesEMPLOYEE Roles = "EMPLOYEE"
)

func (e *Roles) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Roles(s)
	case string:
		*e = Roles(s)
	default:
		return fmt.Errorf("unsupported scan type for Roles: %T", src)
	}
	return nil
}

type NullRoles struct {
	Roles Roles `json:"roles"`
	Valid bool  `json:"valid"` // Valid is true if Roles is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRoles) Scan(value interface{}) error {
	if value == nil {
		ns.Roles, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Roles.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRoles) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Roles), nil
}

type UserStatus string

const (
	UserStatusACTIVE    UserStatus = "ACTIVE"
	UserStatusPENDING   UserStatus = "PENDING"
	UserStatusFREETRIAL UserStatus = "FREE_TRIAL"
	UserStatusINACTIVE  UserStatus = "INACTIVE"
	UserStatusDELETED   UserStatus = "DELETED"
	UserStatusBLOCKED   UserStatus = "BLOCKED"
)

func (e *UserStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserStatus(s)
	case string:
		*e = UserStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UserStatus: %T", src)
	}
	return nil
}

type NullUserStatus struct {
	UserStatus UserStatus `json:"user_status"`
	Valid      bool       `json:"valid"` // Valid is true if UserStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserStatus) Scan(value interface{}) error {
	if value == nil {
		ns.UserStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserStatus), nil
}

type Organization struct {
	ID          string                 `json:"id"`
	Name        zero.String            `json:"name"`
	OwnerID     zero.String            `json:"owner_id"`
	Status      NullOrganizationStatus `json:"status"`
	PhoneNumber zero.String            `json:"phone_number"`
	CreatedAt   zero.Time              `json:"created_at"`
	UpdatedAt   zero.Time              `json:"updated_at"`
	DeletedAt   zero.Time              `json:"deleted_at"`
}

type OrganizationEmployee struct {
	ID        string      `json:"id"`
	OrgID     zero.String `json:"org_id"`
	UserID    string      `json:"user_id"`
	UpdatedAt zero.Time   `json:"updated_at"`
	UpdatedBy zero.String `json:"updated_by"`
	DeletedAt zero.Time   `json:"deleted_at"`
	DeletedBy zero.String `json:"deleted_by"`
}

type RefreshToken struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type User struct {
	ID         string      `json:"id"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Email      string      `json:"email"`
	Password   string      `json:"password"`
	Role       Roles       `json:"role"`
	UserStatus zero.String `json:"user_status"`
	CreatedAt  zero.Time   `json:"created_at"`
	UpdatedAt  zero.Time   `json:"updated_at"`
	DeletedAt  zero.Time   `json:"deleted_at"`
}
