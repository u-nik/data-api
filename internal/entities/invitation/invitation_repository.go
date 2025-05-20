package invitation

import (
	"context"

	"github.com/uptrace/bun"
)

type InvitationRepository struct {
	db *bun.DB
}

func NewInvitationRepository(db *bun.DB) *InvitationRepository {
	return &InvitationRepository{db: db}
}

func (r *InvitationRepository) Create(ctx context.Context, inv *Invitation) error {
	_, err := r.db.NewInsert().Model(inv).Exec(ctx)
	return err
}

func (r *InvitationRepository) FindByID(ctx context.Context, id string) (*Invitation, error) {
	inv := new(Invitation)
	err := r.db.NewSelect().Model(inv).Where("id = ?", id).Scan(ctx)
	return inv, err
}

func (r *InvitationRepository) Update(ctx context.Context, inv *Invitation) error {
	_, err := r.db.NewUpdate().Model(inv).WherePK().Exec(ctx)
	return err
}

func (r *InvitationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&Invitation{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *InvitationRepository) List(ctx context.Context) ([]Invitation, error) {
	var invitations []Invitation
	err := r.db.NewSelect().Model(&invitations).Scan(ctx)
	return invitations, err
}

func (r *InvitationRepository) FindByToken(ctx context.Context, token string) (*Invitation, error) {
	inv := new(Invitation)
	err := r.db.NewSelect().Model(inv).Where("token = ?", token).Scan(ctx)
	return inv, err
}

func (r *InvitationRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	_, err := r.db.NewDelete().Model(&Invitation{}).Where("user_id = ?", userID).Exec(ctx)
	return err
}
