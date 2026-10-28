package model

type ServiceImage struct {
	ServiceID   int64
	Key         string
	IsThumbnail bool
	SortOrder   uint32

	Service *Service
}
