package service

import (
	"context"
	"strings"

	"online-song-library/internal/clients/infoservice"
	msong "online-song-library/internal/model/song"
	"online-song-library/internal/repository/songrepository"

	"github.com/sirupsen/logrus"
)

type Service struct {
	songRepository *songrepository.SongRepository
	client         *infoservice.Client
}

func NewService(
	songRepository *songrepository.SongRepository,
	client *infoservice.Client,
) *Service {
	return &Service{
		songRepository: songRepository,
		client:         client,
	}
}

func (service *Service) GetPaginatedSongs(
	ctx context.Context,
	fields map[string]string,
	offset, limit int,
) (*[]msong.Song, error) {
	return service.songRepository.GetPaginatedSongs(
		ctx,
		fields,
		offset,
		limit,
	)
}

func (service *Service) GetPaginatedText(
	ctx context.Context,
	song msong.Song,
	offset, limit int,
) (string, error) {
	verses, err := service.songRepository.GetPaginatedText(
		ctx,
		song,
		offset,
		limit,
	)
	if err != nil {
		return "", err
	}

	return strings.Join(verses, "\n\n"), nil
}

func (service *Service) CreateSong(ctx context.Context, song msong.Song) error {
	songDetail, err := service.client.GetSongInfo(song)
	if err != nil {
		logrus.Error("unable to get SongDetail: ", err)
	}

	return service.songRepository.Create(ctx, msong.Song{
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: songDetail["releaseDate"],
		Verses:      msong.SplitIntoVerses(songDetail["text"]),
		Link:        songDetail["link"],
	})
}

func (service *Service) UpdateSong(ctx context.Context, song msong.Song) error {
	return service.songRepository.Update(ctx, song)
}

func (service *Service) DeleteSong(ctx context.Context, song msong.Song) error {
	return service.songRepository.Delete(ctx, song)
}
