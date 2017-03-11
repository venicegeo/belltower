package orm

// These are the security policies:
//
// admin: can do CRUD
// owner: can do CRUD
// others: can do R, if the object is public
//
// additional restriction: if the object is an association object, you must have R to the two associated objects
// if bit is set, you are allowed to do that operation

const (
	// UserRole: can read (public) things
	// CreatorRole: UserRole powers, plus can create things (except Users), and can do anything to things it owns
	// AdminRole: can do any operation for any thing (including User objects)
	UserRole uint = iota
	CreatorRole
	AdminRole

	// ReadOperation is for Read, WriteOperation is for Create/Update/Delete
	ReadOperation = iota
	WriteOperation
)

type Authorizable interface {
	GetOwnerID() uint
	GetIsPublic() bool
}

func testRole(requestorRole uint, requestedRole uint) bool {
	return requestorRole >= requestedRole
}

func isAuthorizedForCreate(requestor *User) bool {

	// admins can do anything
	if testRole(requestor.Role, AdminRole) {
		return true
	}

	// creators can create things
	if testRole(requestor.Role, CreatorRole) {
		return true
	}

	// No one else can do anything: "You shall not pass!"
	return false
}

func isAuthorized(requestor *User, object Authorizable, operation uint) bool {

	// admins can do anything
	if testRole(requestor.Role, AdminRole) {
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
	if object.GetIsPublic() && operation == ReadOperation {
		return true
	}

	// "You shall not pass!"
	return false
}
