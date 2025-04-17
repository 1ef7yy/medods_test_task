package domain

import "github.com/1ef7yy/medods_test_task/models"

func (d domain) Login(guid string) (models.Token, error)

func (d domain) Refresh(tokens models.Token) (models.Token, error)
