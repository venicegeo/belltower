package orm

// These are the security policies:
//
// admin: can do CRUD
// owner: can do CRUD
// others: can do R, if the object is public
//
// additional restriction: if the object is an association object, you must have R to the two associated objects
// if bit is set, you are allowed to do that operation

type Role uint
type Operation uint

const (
	// UserRole: can read (public) things
	// CreatorRole: UserRole powers, plus can create things (except Users), and can do anything to things it owns
	// AdminRole: can do any operation for any thing (including User objects)
	UserRole Role = iota
	CreatorRole
	AdminRole

	// ReadOperation is for Read, WriteOperation is for Create/Update/Delete
	ReadOperation Operation = iota
	CreateOperation
	UpdateOperation
	DeleteOperation
)

type Authorizable interface {
	GetOwnerID() uint
	GetIsPublic() bool
}

func isUser(requestorRole Role) bool {
	return isRole(requestorRole, UserRole)
}

func isCreator(requestorRole Role) bool {
	return isRole(requestorRole, CreatorRole)
}

func isAdmin(requestorRole Role) bool {
	return isRole(requestorRole, AdminRole)
}

func isRole(requestorRole Role, requestedRole Role) bool {
	return uint(requestorRole) >= uint(requestedRole)
}

func isAuthorizedForCreate(requestor *User) bool {

	// admins can do anything
	if isAdmin(requestor.Role) {
		return true
	}

	// creators can create things
	if isCreator(requestor.Role) {
		return true
	}

	// No one else can do anything: "You shall not pass!"
	return false
}

func isAuthorizedForRead(requestor *User, object Authorizable) bool {

	// admins can do anything
	if isAdmin(requestor.Role) {
		return true
	}

	// owners can do anything to objects they own
	if requestor.ID == object.GetOwnerID() {
		return true
	}

	// We're either a user or a creator. Creators can't do operations on
	// objects it doesn't own, and it doesn't own this one,  so if we
	// are a creator, we're effectively just a user.
	//
	// And users can only read public things.
	if object.GetIsPublic() {
		return true
	}

	// "You shall not pass!"
	return false
}

func isAuthorizedForUpdate(requestor *User, object Authorizable) bool {

	// admins can do anything
	if isAdmin(requestor.Role) {
		return true
	}

	// owners can do anything to objects they own
	if requestor.ID == object.GetOwnerID() {
		return true
	}

	// "You shall not pass!"
	return false
}

func isAuthorizedForDelete(requestor *User, object Authorizable) bool {
	return isAuthorizedForUpdate(requestor, object)
}

func isAuthorized(requestor *User, object Authorizable, operation Operation) bool {
	switch operation {
	case ReadOperation:
		return isAuthorizedForRead(requestor, object)
	case CreateOperation:
		return isAuthorizedForCreate(requestor)
	case UpdateOperation:
		return isAuthorizedForUpdate(requestor, object)
	case DeleteOperation:
		return isAuthorizedForDelete(requestor, object)
	}
	panic("internal error")
}
