package utils

type Room struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RoomCrypted struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}
