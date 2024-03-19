package serializer

import (
	"encoding/json"
	"fmt"
	"log"

	"gitlab.com/tour/internal/core/repository/psql/sqlc"
)

func MarshalUnMarshal(input any, output any) error {
	rawBytes, err := json.Marshal(input)
	if err != nil {
		log.Println("error while marshaling", fmt.Errorf("%w", err))
		return err
	}

	return json.Unmarshal(rawBytes, output)
}

func SerializeEnumToString(r sqlc.Roles) (string, error) {
	switch r {
	case sqlc.RolesOWNER:
		return "OWNER", nil
	case sqlc.RolesADMIN:
		return "ADMIN", nil
	case sqlc.RolesEMPLOYEE:
		return "EMPLOYEE", nil
	case sqlc.RolesUSER:
		return "USER", nil
	default:
		return "", fmt.Errorf("unknown Enum type")
	}
}
