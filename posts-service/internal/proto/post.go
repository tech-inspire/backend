package proto

import (
	"log/slog"

	postsv1 "github.com/tech-inspire/api-contracts/api/gen/go/posts/v1"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/pkg/generics"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Image(variant models.ImageVariant) *postsv1.ImageVariant {
	variantType := postsv1.VariantType_VARIANT_TYPE_UNSPECIFIED
	switch variant.VariantType {
	case models.Thumbnail:
		variantType = postsv1.VariantType_THUMBNAIL
	case models.Original:
		variantType = postsv1.VariantType_ORIGINAL
	default:
		slog.Error("variant type is not supported", slog.String("variant", string(variant.VariantType)))
	}

	return &postsv1.ImageVariant{
		Url:         variant.URL,
		Width:       int32(variant.Width),
		Height:      int32(variant.Height),
		Size:        variant.Size,
		VariantType: variantType,
	}
}

func Post(post *models.Post) *postsv1.Post {
	var soundcloudSongStart *int32
	if post.SoundCloudSongStartMilli != nil {
		s := int32(*post.SoundCloudSongStartMilli)
		soundcloudSongStart = &s
	}

	return &postsv1.Post{
		PostId:              post.PostID.String(),
		AuthorId:            post.AuthorID.String(),
		Images:              generics.Convert(post.Images, Image),
		SoundcloudSong:      post.SoundCloudSongURL,
		SoundcloudSongStart: soundcloudSongStart,
		Description:         post.Description,
		CreatedAt:           timestamppb.New(post.CreatedAt),
	}
}
