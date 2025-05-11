package character

import (
	"context"
	"time"

	"github.com/th1enq/go_coffee/db"
	"github.com/th1enq/go_coffee/internal/model"
	"github.com/th1enq/go_coffee/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CharacterService struct {
	proto.UnimplementedCharacterServiceServer
	repo *CharacterRepository
}

func NewCharacterService(db *db.DB) *CharacterService {
	return &CharacterService{
		repo: NewCharacterRepository(db),
	}
}

func (c *CharacterService) CreateCharacter(ctx context.Context, req *proto.CreateCharacterRequest) (*proto.CreateCharacterResponse, error) {
	birthDay, err := time.Parse("2-Jan", req.Character.Birthday)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid birthday date: %v", err)
	}

	realeaseDate, err := time.Parse("01/02/06", req.Character.ReleaseDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid realase data: %v", err)
	}

	character := &model.Character{
		Name:          req.Character.Name,
		Rarity:        int(req.Character.Rarity),
		Region:        req.Character.Region,
		Vision:        req.Character.Vision,
		WeaponType:    req.Character.WeaponType,
		Constellation: req.Character.Constellation,
		Birthday:      birthDay,
		Affilliation:  req.Character.Affilliation,
		ReleaseDate:   realeaseDate,
	}

	if err := c.repo.CreateCharacter(ctx, character); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create character: %v", err)
	}
	return &proto.CreateCharacterResponse{
		Character: req.Character,
	}, nil
}

func (c *CharacterService) UpdateCharacter(ctx context.Context, req *proto.UpdateCharacterRequest) (*proto.UpdateCharacterResponse, error) {
	birthDay, err := time.Parse("2-Jan", req.Character.Birthday)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid birthday date: %v", err)
	}

	releaseDate, err := time.Parse("01/02/06", req.Character.ReleaseDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid release date: %v", err)
	}

	character := &model.Character{
		Name:          req.Character.Name,
		Rarity:        int(req.Character.Rarity),
		Region:        req.Character.Region,
		Vision:        req.Character.Vision,
		WeaponType:    req.Character.WeaponType,
		Constellation: req.Character.Constellation,
		Birthday:      birthDay,
		Affilliation:  req.Character.Affilliation,
		ReleaseDate:   releaseDate,
	}

	if err := c.repo.UpdateCharacter(ctx, req.Id, character); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update character: %v", err)
	}

	return &proto.UpdateCharacterResponse{
		Character: req.Character,
	}, nil
}

func (c *CharacterService) DeleteCharacter(ctx context.Context, req *proto.DeleteCharacterRequest) (*proto.DeleteCharacterResponse, error) {
	if err := c.repo.DeleteCharacter(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete character: %v", err)
	}

	return &proto.DeleteCharacterResponse{}, nil
}

func (c *CharacterService) GetCharacterByName(ctx context.Context, req *proto.GetCharacterByNameRequest) (*proto.GetCharacterResponse, error) {
	character, err := c.repo.GetCharacterByName(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "character not found: %v", err)
	}

	characterMsg := &proto.CharacterMessage{
		Id:            uint64(character.ID),
		Name:          character.Name,
		Rarity:        int32(character.Rarity),
		Region:        character.Region,
		Vision:        character.Vision,
		WeaponType:    character.WeaponType,
		Constellation: character.Constellation,
		Birthday:      character.Birthday.Format("2-Jan"),
		Affilliation:  character.Affilliation,
		ReleaseDate:   character.ReleaseDate.Format("01/02/06"),
	}

	return &proto.GetCharacterResponse{
		Character: characterMsg,
	}, nil
}

func (c *CharacterService) SearchCharacters(ctx context.Context, req *proto.SearchCharactersRequest) (*proto.GetCharactersResponse, error) {
	characters, err := c.repo.SearchCharacters(ctx, req.Name, req.Region, req.Vision, req.WeaponType, int(req.Rarity))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search characters: %v", err)
	}

	var characterMsgs []*proto.CharacterMessage
	for _, character := range characters {
		characterMsg := &proto.CharacterMessage{
			Id:            uint64(character.ID),
			Name:          character.Name,
			Rarity:        int32(character.Rarity),
			Region:        character.Region,
			Vision:        character.Vision,
			WeaponType:    character.WeaponType,
			Constellation: character.Constellation,
			Birthday:      character.Birthday.Format("2-Jan"),
			Affilliation:  character.Affilliation,
			ReleaseDate:   character.ReleaseDate.Format("01/02/06"),
		}
		characterMsgs = append(characterMsgs, characterMsg)
	}

	return &proto.GetCharactersResponse{
		Characters: characterMsgs,
	}, nil
}
