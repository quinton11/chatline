package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateRoomHash() string {
	uid := uuid.New()
	return uid.String()
}

func GetServerIp() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, inter := range interfaces {
		//ip := net.ParseIP(inter.)
		winEth := strings.Index(inter.Name, "Eth") // win ethernet interface
		winWi := strings.Index(inter.Name, "Wi")   // win wireless interface
		linEth := strings.Index(inter.Name, "eth") // unix ethernet interface
		linWL := strings.Index(inter.Name, "wlan") // unix wireless interface
		if winEth == 0 || winWi == 0 || linEth == 0 || linWL == 0 {

			if addrs, err := inter.Addrs(); err == nil {
				for _, addr := range addrs {

					if addr.Network() == "ip+net" {

						ip := strings.Split(addr.String(), "/")
						//check if address is ipv4
						if len(ip) == 2 && len(strings.Split(ip[0], ".")) == 4 {
							ipparse := net.ParseIP(ip[0])

							if ipparse.IsGlobalUnicast() {

								return ip[0], nil
							}

						}
					}

				}
			}
		}

	}
	return "", nil
}

func GenerateToken(room Room) (string, error) {
	jtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"Room": room, "exp": time.Now().Add(time.Hour * 1).Unix()})

	key := base64.StdEncoding.EncodeToString([]byte(room.Uuid))
	jtokenString, err := jtoken.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return jtokenString, nil
}

func ValidateToken(token string, key string) (Room, error) {
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, errors.New("error in parsing room token")
		}
		return []byte(key), nil
	})
	if err != nil {
		return Room{}, err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !(ok && tok.Valid) {
		return Room{}, errors.New("error in parsing token")

	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return Room{}, errors.New("token expired")
	}

	rBytes, err := json.Marshal(claims["Room"])
	if err != nil {
		return Room{}, err
	}

	var room Room
	err = json.Unmarshal(rBytes, &room)
	if err != nil {
		return Room{}, err
	}
	return room, nil
}

func GenerateHash(token string, room Room) (string, string) {
	tokenB := []byte(token)
	secretB := []byte(room.Uuid)

	tokenBase64 := base64.StdEncoding.EncodeToString(tokenB)
	secretBase64 := base64.StdEncoding.EncodeToString(secretB)

	return tokenBase64, secretBase64
}

// takes a base64 value and returns the decoded value
func Decode64(keybase64 string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(keybase64)
	if err != nil {
		return "", err
	}

	decodedKey := string(decodedBytes)

	return decodedKey, nil

}
