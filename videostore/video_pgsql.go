package videostore

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// postgresVideoRepo stores models to a Postgres DB
type postgresVideoRepo struct {
	db pg.DB
	fs files.FileSystem
}

func NewPostgresVideoRepo(db pg.DB, fs files.FileSystem) *postgresVideoRepo {
	err := db.CreateTable((*Video)(nil), &orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		log.Fatalf("failed to create table: %+v", err)
	}

	return &postgresVideoRepo{
		db,
		fs,
	}
}

func (repo *postgresVideoRepo) FindById(id uint) (*Video, error) {
	video := &Video{
		ID: id,
	}

	err := repo.db.Select(video)

	if err == pg.ErrNoRows {
		return video, ErrorVideoNotFound
	}

	return video, err
}

func (repo *postgresVideoRepo) All(filter *VideoFilter, limit uint, offset uint) ([]*Video, error) {
	var videos []*Video

	query := repo.db.Model(&videos)

	if !filter.Empty() {
		query = query.Apply(func(q *orm.Query) (*orm.Query, error) {
			if len(filter.Title) > 0 {
				q = q.Where("LOWER(title) LIKE LOWER(?)", "%"+filter.Title+"%")
			}

			if len(filter.Tags) > 0 {
				q = q.Where("tags \\?& ?", pg.Array(filter.Tags))
			}

			if len(filter.Any) > 0 {
				q = q.WhereOr("LOWER(title) LIKE LOWER(?)", "%"+filter.Any+"%")
				q = q.WhereOr("tags \\?& ?", pg.Array([]string{filter.Any}))
			}

			return q, nil
		})
	}

	if filter.Sort() {
		query = query.Apply(func(q *orm.Query) (*orm.Query, error) {
			if !filter.ValidSortField() {
				return nil, fmt.Errorf("invalid sort field %v", filter.SortField)
			}

			if !filter.ValidSortDirection() {
				return nil, fmt.Errorf("invalid sort direction %v", filter.SortDirection)
			}

			q = q.Order(fmt.Sprintf("%v %v", filter.SortField, filter.SortDirection))

			return q, nil
		})
	}

	query = query.Apply(func(q *orm.Query) (*orm.Query, error) {
		return q.Limit(int(limit)).Offset(int(offset)), nil
	})

	err := query.Select()

	return videos, err
}

func (repo *postgresVideoRepo) Save(video *Video) (*Video, error) {
	var err error

	if video.Exists() {
		video.TimeUpdated = time.Now().Format(time.RFC3339)
		err = repo.db.Update(&video)
	} else {
		video.TimeCreated = time.Now().Format(time.RFC3339)
		video.TimeUpdated = time.Now().Format(time.RFC3339)
		err = repo.db.Insert(&video)
	}

	return video, err
}

func (repo *postgresVideoRepo) Delete(video *Video) error {
	err := repo.db.Delete(video)
	if err != nil {
		return errors.Wrap(err, "failed to delete from db")
	}

	_, err = repo.fs.Stat(video.Source)
	if !repo.fs.IsNotExist(err) {
		// video exists, attempt to delete
		err = repo.fs.Remove(video.Source)
		if err != nil {
			log.Print(errors.Wrap(err, "failed to remove video from disk"))
		}
	}

	_, err = repo.fs.Stat(video.Thumbnail)
	if !repo.fs.IsNotExist(err) {
		// thumbnail exists, attempt to delete
		err = repo.fs.Remove(video.Thumbnail)
		if err != nil {
			log.Print(errors.Wrap(err, "failed to remove thumbnail from disk"))
		}
	}

	return nil
}

func (repo *postgresVideoRepo) Upload(video *Video, reader io.Reader) (*Video, error) {
	video.Thumbnail = ""
	video.Source = ""
	video, err := repo.Save(video)

	if err != nil {
		return video, err
	}

	rootDir := strconv.Itoa(int(video.ID))
	if _, err := repo.fs.Stat(rootDir); repo.fs.IsNotExist(err) {
		repo.fs.MkdirAll(rootDir, os.ModePerm)
	}

	videoPath := path.Join(rootDir, "video"+path.Ext(video.OriginalFileName))

	files.PipeTo(repo.fs, videoPath, reader)

	video.Source = videoPath
	go eventuallyMakeThumbnail(video, repo, repo.fs)

	return repo.Save(video)
}
