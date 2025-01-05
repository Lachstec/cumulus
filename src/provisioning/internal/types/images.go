package types

// Image represents all available VM images in the openstack cluster.
type Image string

const (
	Alpine3_20_3 Image = "29a24dc0-b24b-4cc8-b43b-a8a4c6916d0f" //nolint:all
	Ubuntu14_04  Image = "f211d66d-9167-4133-abd6-c40b1586394e" //nolint:all
	Ubuntu16_04  Image = "cc21acd3-8cf1-40be-8bb1-6fea9453b0bb" //nolint:all
	Ubuntu18_04  Image = "f56402aa-369e-42b7-a64d-2db29bebfebd" //nolint:all
	Ubuntu20_04  Image = "de55d4f6-5905-4e4e-b42b-db7ebabbcda4" //nolint:all
	Ubuntu22_04  Image = "1404d277-1fd2-4682-9fbd-80c7d05b80e1" //nolint:all
	Ubuntu24_04  Image = "d6d1835c-7180-4ca9-b4a1-470afbd8b398" //nolint:all
	Fedora41     Image = "715efd2c-b224-49d6-bbc7-00c204b1f04c" //nolint:all
	Cirros063    Image = "b3f88062-56d5-43b1-9d1e-96980ea0e16b" //nolint:all
)

func (i Image) Value() string {
	return string(i)
}
