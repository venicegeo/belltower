package btorm

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert := assert.New(t)

	assert.True(isUser(UserRole))
	assert.True(isUser(CreatorRole))
	assert.True(isUser(AdminRole))

	assert.False(isCreator(UserRole))
	assert.True(isCreator(CreatorRole))
	assert.True(isCreator(AdminRole))

	assert.False(isAdmin(UserRole))
	assert.False(isAdmin(CreatorRole))
	assert.True(isAdmin(AdminRole))
}

func TestIsAuthorized(t *testing.T) {
	assert := assert.New(t)

	user := &User{
		Role: UserRole,
	}
	user.Id = "1"

	creator := &User{
		Role: CreatorRole,
	}
	creator.Id = "2"

	admin := &User{
		Role: AdminRole,
	}
	admin.Id = "3"

	myPublicRule := &Rule{}
	myPublicRule.Id = "10"
	myPublicRule.OwnerId = creator.Id
	myPublicRule.IsPublic = true

	myPrivateRule := &Rule{}
	myPrivateRule.Id = "11"
	myPrivateRule.OwnerId = creator.Id
	myPrivateRule.IsPublic = false

	yourPublicRule := &Rule{}
	yourPublicRule.Id = "11"
	yourPublicRule.Id = "12"
	yourPublicRule.OwnerId = "999"
	yourPublicRule.IsPublic = true

	yourPrivateRule := &Rule{}
	yourPrivateRule.Id = "13"
	yourPrivateRule.OwnerId = "99"
	yourPrivateRule.IsPublic = false

	_ = admin

	type X struct {
		user     *User
		obj      Authorizable
		op       Operation
		expected bool
	}
	x := []X{
		{user, myPublicRule, ReadOperation, true}, // 0
		{user, myPublicRule, CreateOperation, false},
		{user, myPublicRule, UpdateOperation, false},
		{user, myPublicRule, DeleteOperation, false},

		{user, myPrivateRule, ReadOperation, false}, // 4
		{user, myPrivateRule, CreateOperation, false},
		{user, myPrivateRule, UpdateOperation, false},
		{user, myPrivateRule, DeleteOperation, false},

		{creator, myPublicRule, ReadOperation, true}, // 8
		{creator, myPublicRule, CreateOperation, true},
		{creator, myPublicRule, UpdateOperation, true},
		{creator, myPublicRule, DeleteOperation, true},

		{creator, myPrivateRule, ReadOperation, true}, // 12
		{creator, myPrivateRule, CreateOperation, true},
		{creator, myPrivateRule, UpdateOperation, true},
		{creator, myPrivateRule, DeleteOperation, true},

		{creator, yourPublicRule, ReadOperation, true}, // 16
		{creator, yourPublicRule, CreateOperation, true},
		{creator, yourPublicRule, UpdateOperation, false},
		{creator, yourPublicRule, DeleteOperation, false},

		{creator, yourPrivateRule, ReadOperation, false}, // 20
		{creator, yourPrivateRule, CreateOperation, true},
		{creator, yourPrivateRule, UpdateOperation, false},
		{creator, yourPrivateRule, DeleteOperation, false},

		{admin, myPublicRule, ReadOperation, true}, // 24
		{admin, myPublicRule, CreateOperation, true},
		{admin, myPublicRule, UpdateOperation, true},
		{admin, myPublicRule, DeleteOperation, true},

		{admin, myPrivateRule, ReadOperation, true}, // 28
		{admin, myPrivateRule, CreateOperation, true},
		{admin, myPrivateRule, UpdateOperation, true},
		{admin, myPrivateRule, DeleteOperation, true},
	}

	for i, v := range x {
		ok := isAuthorized(v.user, v.obj, v.op)
		if ok != v.expected {
			log.Printf("** %d ** %#v", i, v)
		}
		assert.Equal(v.expected, ok)
	}
}
