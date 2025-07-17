package users

import "context"

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.usersRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
