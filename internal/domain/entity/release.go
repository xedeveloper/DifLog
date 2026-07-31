package entity

type ReleaseAsset struct{
	Name string
	DownloadURL string
}

type Release struct{
	TagName string
	Assets []ReleaseAsset
}
