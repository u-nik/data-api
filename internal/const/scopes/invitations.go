package scopes

type InvitationsScopes struct {
	Create string
	Update string
	Delete string
	Read   string
	List   string
	Accept string
}

const InvitationPrefix = "invitations."

var Invitations = InvitationsScopes{
	Create: InvitationPrefix + "create",
	Update: InvitationPrefix + "update",
	Delete: InvitationPrefix + "delete",
	Read:   InvitationPrefix + "read",
	List:   InvitationPrefix + "list",
	Accept: InvitationPrefix + "accept",
}
