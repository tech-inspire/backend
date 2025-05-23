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
	ImageSize        int32

	SoundCloudSongURL        *string
	SoundCloudSongStartMilli *int
	Description              string
}
