package scylla

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/table"
)

// ImageVariant corresponds to the image_variant UDT in ScyllaDB.
type ImageVariant struct {
	gocqlx.UDT
	VariantType string `cql:"variant_type" db:"variant_type"`
	URL         string `cql:"url"         db:"url"`
	Width       int    `cql:"width"       db:"width"`
	Height      int    `cql:"height"     db:"height"`
	Size        int32  `cql:"size"        db:"size"`
}

// Post maps to the posts_by_id table.
type Post struct {
	PostID              gocql.UUID     `db:"post_id"`
	AuthorID            gocql.UUID     `db:"author_id"`
	Images              []ImageVariant `db:"images"`
	SoundCloudSong      *string        `db:"soundcloud_song"`
	SoundCloudSongStart *int           `db:"soundcloud_song_start"`
	Description         string         `db:"description"`
	CreatedAt           time.Time      `db:"created_at"`
}

var (
	postMetadata = table.Metadata{
		Name:    "posts.posts_by_id",
		Columns: []string{"post_id", "author_id", "images", "soundcloud_song", "soundcloud_song_start", "description", "created_at"},
		PartKey: []string{"post_id"},
	}
	postTable = table.New(postMetadata)
)
