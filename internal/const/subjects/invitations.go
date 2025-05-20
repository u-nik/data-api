package subjects

type InvitationsSubjects struct {
	Create string
	Update string
	Delete string
	Accept string
}

const InvitationsPrefix = "invitations."

var Invitations = InvitationsSubjects{
	Create: InvitationsPrefix + "create",
	Update: InvitationsPrefix + "update",
	Delete: InvitationsPrefix + "delete",
	Accept: InvitationsPrefix + "accept",
}
