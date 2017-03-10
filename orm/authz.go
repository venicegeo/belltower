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
)

type Authorizable interface {
	GetOwnerID() uint
	GetIsPublic() bool
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

	// the user must have the specific rights
	if requestor.Rights&requestedPermissions != requestedPermissions {
		return false
	}

	// finally, is the object itself public?
	if !object.GetIsPublic() {
		return false
	}

	// we made it
	return true
}
