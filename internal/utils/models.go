package utils

type Room struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RoomCrypted struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}
