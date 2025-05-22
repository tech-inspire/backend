package dto

type UpdatePostParams struct {
	SoundCloudSongURL        *string
	SoundCloudSongStartMilli *int
	Description              *string
}

type CreatePostParams struct {
	UploadSessionKey string
	ImageWidth       int
	ImageHeight      int
	ImageSize        int64

	SoundCloudSongURL        *string
	SoundCloudSongStartMilli *int
	Description              string
}
