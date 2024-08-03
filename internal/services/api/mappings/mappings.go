package mappings

import (
	"github.com/tuvuanh27/go-crawler/internal/pkg/mapper"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/user/v1/dtos"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*model.User, *dtos.RegisterUserResponseDto]()
	if err != nil {
		return err
	}
	return err
}
