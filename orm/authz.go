package orm

// if bit is set, you are allowed to do that operation
const (
	CanCreateFeed uint = 1 << iota
	CanReadFeed
	CanUpdateFeed
	CanDeleteFeed
	CanCreateAction
	CanReadAction
	CanUpdateAction
	CanDeleteAction
	CanCreateRule
	CanReadRule
	CanUpdateRule
	CanDeleteRule

	CanReadAll     = CanReadFeed | CanReadAction | CanReadAction
	CanAdminFeed   = CanCreateFeed | CanReadFeed | CanUpdateFeed | CanDeleteFeed
	CanAdminAction = CanCreateAction | CanReadAction | CanUpdateAction | CanDeleteAction
	CanAdminRule   = CanCreateRule | CanReadRule | CanUpdateRule | CanDeleteRule
	CanAdminAll    = CanAdminFeed | CanAdminAction | CanAdminRule
)

type Authorizable interface {
	GetOwnerID() uint
	GetIsPublic() bool
}

func isAuthorizedInGeneral(requestor *User, requestedPermissions uint) bool {

	// admins can do anything
	if requestor.IsAdmin {
		return true
	}

	// finally, the user must have the specific rights
	if requestor.Rights&requestedPermissions != requestedPermissions {
		return false
	}

	// we made it
	return true
}

func isAuthorizedForObject(requestor *User, object Authorizable, requestedPermissions uint) bool {

	// admins can do anything
	if requestor.IsAdmin {
		return true
	}

	// owners can do anything
	if requestor.ID == object.GetOwnerID() {
		return true
	}

	// is the object itself public?
	if !object.GetIsPublic() {
		return false
	}

	// finally, the user must have the specific rights
	if requestor.Rights&requestedPermissions != requestedPermissions {
		return false
	}

	// we made it
	return true
}
