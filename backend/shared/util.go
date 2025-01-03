package shared

import gonanoid "github.com/matoous/go-nanoid/v2"

func UniqueId() string {
	return gonanoid.MustGenerate("abcdefghijklmnopqrstuvwxyz012345679_-", 10)
}
