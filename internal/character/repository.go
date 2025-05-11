package character

import (
	"context"

	"github.com/th1enq/go_coffee/db"
	"github.com/th1enq/go_coffee/internal/model"
)

type CharacterRepository struct {
	db *db.DB
}

func NewCharacterRepository(db *db.DB) *CharacterRepository {
	return &CharacterRepository{
		db: db,
	}
}

func (r *CharacterRepository) CreateCharacter(ctx context.Context, c *model.Character) error {
	result := r.db.Create(c)
	return result.Error
}

func (r *CharacterRepository) UpdateCharacter(ctx context.Context, id uint64, c *model.Character) error {
	result := r.db.Model(&model.Character{}).Where("id = ?", id).Updates(c)
	return result.Error
}

func (r *CharacterRepository) DeleteCharacter(ctx context.Context, id uint64) error {
	result := r.db.Delete(&model.Character{}, id)
	return result.Error
}

func (r *CharacterRepository) GetCharacterByName(ctx context.Context, name string) (*model.Character, error) {
	var character model.Character
	result := r.db.Where("name = ?", name).First(&character)
	if result.Error != nil {
		return nil, result.Error
	}
	return &character, nil
}

func (r *CharacterRepository) SearchCharacters(ctx context.Context, name, region, vision, weaponType string, rarity int) ([]*model.Character, error) {
	var characters []*model.Character
	query := r.db.Model(&model.Character{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if region != "" {
		query = query.Where("region = ?", region)
	}
	if vision != "" {
		query = query.Where("vision = ?", vision)
	}
	if weaponType != "" {
		query = query.Where("weapon_type = ?", weaponType)
	}
	if rarity > 0 {
		query = query.Where("rarity = ?", rarity)
	}

	result := query.Find(&characters)
	if result.Error != nil {
		return nil, result.Error
	}

	return characters, nil
}
