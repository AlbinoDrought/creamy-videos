package cmd

import (
	"log"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/videostore"
	"github.com/go-pg/pg"
)

type application struct {
	config appConfig
	fs     files.FileSystem
	repo   videostore.VideoRepo
}

func (instance application) makeDummyRepo() videostore.VideoRepo {
	return videostore.NewDummyVideoRepo(instance.fs)
}

func (instance application) makePostgresRepo() videostore.VideoRepo {
	db := pg.Connect(&pg.Options{
		User:     instance.config.PostgresUser,
		Password: instance.config.PostgresPassword,
		Addr:     instance.config.PostgresAddress,
		Database: instance.config.PostgresDatabase,
	})
	// db never closed

	return videostore.NewPostgresVideoRepo(*db, instance.fs)
}

func makeApp(cfg appConfig) (instance application) {
	instance.config = cfg

	instance.fs = files.TransformFileSystem(
		files.LocalFileSystem(instance.config.LocalVideoDirectory),
		func(p []byte) {
			for i := 0; i < len(p); i++ {
				p[i] = p[i] ^ instance.config.FilesystemKey
			}
		},
	)

	if instance.config.UsePostgres {
		log.Println("Video Repo: Postgres")
		instance.repo = instance.makePostgresRepo()
	} else {
		log.Println("Video Repo: JSON")
		instance.repo = instance.makeDummyRepo()
	}

	return instance
}
