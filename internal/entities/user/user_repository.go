package user

import (
	"context"

	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *User) error {
	_, err := r.db.NewInsert().Model(u).Exec(ctx)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*User, error) {
	user := new(User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	return user, err
}

func (r *UserRepository) Update(ctx context.Context, u *User) error {
	_, err := r.db.NewUpdate().Model(u).WherePK().Exec(ctx)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model(&User{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *UserRepository) List(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.NewSelect().Model(&users).Scan(ctx)
	return users, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	user := new(User)
	err := r.db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	return user, err
}
