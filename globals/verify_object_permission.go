package globals

import (
	"github.com/PretendoNetwork/nex-go"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"golang.org/x/exp/slices"
)

func VerifyObjectPermission(ownerPID, accessorPID uint32, permission *datastore_types.DataStorePermission) uint32 {
	if permission.Permission > 3 {
		return nex.Errors.DataStore.InvalidArgument
	}

	// * Allow anyone
	if permission.Permission == 0 {
		return 0
	}

	// * Allow friends
	// TODO - Implement this
	if permission.Permission == 1 {
		return nex.Errors.DataStore.PermissionDenied
	}

	// * Allow people in permission.RecipientIDs
	if permission.Permission == 2 {
		if !slices.Contains(permission.RecipientIDs, accessorPID) {
			return nex.Errors.DataStore.PermissionDenied
		}
	}

	// * Allow only the owner
	if permission.Permission == 3 {
		if ownerPID != accessorPID {
			return nex.Errors.DataStore.PermissionDenied
		}
	}

	return 0
}
