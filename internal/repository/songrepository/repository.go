package songrepository

import (
	"context"
	"fmt"

	msong "online-song-library/internal/model/song"
	"online-song-library/pkg/dbstore"
)

type SongRepository struct {
	store dbstore.Store
}

func NewSongRepository(store dbstore.Store) *SongRepository {
	return &SongRepository{
		store: store,
	}
}

func (sr *SongRepository) GetPaginatedSongs(
	ctx context.Context,
	fields map[string]string,
	offset, limit int,
) (*[]msong.Song, error) {

	sql := fmt.Sprintf(`
	select
		id,
		group,
		song,
		release_date,
		link
	from songs
	where $1
	offset $2
	limit $3;
	`)

	where := fmt.Sprintf("where (%s)", func() string {
		if len(fields) == 0 {
			return "1=1"
		}

		filter := ""
		for k, v := range fields {
			if len(filter) > 0 {
				filter += " and "
			}
			filter += fmt.Sprintf("%s=%s", k, v)
		}
		return filter
	}())

	rows, err := sr.store.Query(
		ctx,
		sql,
		where,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	songs := []msong.Song{}
	for rows.Next() {
		song := msong.Song{}
		if err := rows.Scan(
			&song.ID,
			&song.Group,
			&song.Song,
			&song.ReleaseDate,
			&song.Link,
		); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return &songs, nil

}

func (sr *SongRepository) GetPaginatedText(
	ctx context.Context,
	song msong.Song,
	offset, limit int,
) ([]string, error) {
	const sql = `
	select
		verses[$1:$2]
	from songs
	where id = $3;
	`

	verses := []string{}

	if err := sr.store.QueryRow(
		ctx,
		sql,
		offset,
		limit,
		song.ID,
	).Scan(
		&verses,
	); err != nil {
		return nil, err
	}

	return verses, nil
}

func (sr *SongRepository) Delete(ctx context.Context, song msong.Song) error {
	const sql = `
	delete from songs
	where id = $1;
	`
	if _, err := sr.store.Exec(
		ctx,
		sql,
		song.ID,
	); err != nil {
		return err
	}

	return nil
}

func (sr *SongRepository) Update(ctx context.Context, song msong.Song) error {
	const sql = `
	update
		songs
	set
		group = $1,
		song = $2,
		release_date = $3,
		text = $4,
		link = $5
	where id = $6;
	`

	if _, err := sr.store.Exec(
		ctx,
		sql,
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Verses,
		song.Link,
	); err != nil {
		return err
	}

	return nil
}

func (sr *SongRepository) Create(ctx context.Context, song msong.Song) error {
	const sql = `
	insert into songs(
		group,
		song,
		release_date,
		text,
		link
	) values ($1, $2, $3, $4, $5);
	`

	if _, err := sr.store.Exec(
		ctx,
		sql,
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Verses,
		song.Link,
	); err != nil {
		return err
	}

	return nil
}
