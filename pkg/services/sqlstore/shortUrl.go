package sqlstore

import (
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetFullUrlByUid)
	bus.AddHandler("sql", UpdateShortUrlLastSeenAt)
	bus.AddHandler("sql", CreateShortUrl)
}

func GetFullUrlByUid(query *models.GetFullUrlQuery) error {
	var shortUrl models.ShortUrl
	exists, err := x.Where("uid=?", query.Uid).Get(&shortUrl)
	if err != nil {
		return err
	}

	if !exists {
		return models.ErrShortUrlNotFound
	}

	query.Result = &shortUrl
	return nil
}

func UpdateShortUrlLastSeenAt(cmd *models.UpdateShortUrlLastSeenAtCommand) error {
	now := getTime().Unix()
	var shortUrl = models.ShortUrl{
		Uid:        cmd.Uid,
		LastSeenAt: now,
	}

	_, err := x.ID(cmd.Uid).Update(&shortUrl)
	return err
}

func CreateShortUrl(command *models.CreateShortUrlCommand) error {
	now := getTime().Unix()
	shortUrl := models.ShortUrl{
		Uid:       command.Uid,
		Path:      command.Path,
		CreatedBy: command.CreatedBy,
		CreatedAt: now,
	}

	_, err := x.Insert(&shortUrl)
	if err != nil {
		return err
	}

	command.Result = &shortUrl
	return err
}