package utils

import (
	mathRand "math/rand"
	"strings"
	"sync"

	"github.com/bxcodec/faker/v4"
)

type ConcurrentFaker struct {
	sync.Mutex
}

func (cf *ConcurrentFaker) Email() string {
	cf.Lock()
	defer cf.Unlock()
	return strings.ToLower(faker.Email())
}

func (cf *ConcurrentFaker) Password() string {
	cf.Lock()
	defer cf.Unlock()
	return faker.Password()
}

func (cf *ConcurrentFaker) Name() string {
	cf.Lock()
	defer cf.Unlock()
	return faker.Name()
}

func (cf *ConcurrentFaker) Username() string {
	cf.Lock()
	defer cf.Unlock()
	return faker.Username()
}

func (cf *ConcurrentFaker) LastName() string {
	cf.Lock()
	defer cf.Unlock()
	return faker.LastName()
}

func (cf *ConcurrentFaker) UUID() string {
	cf.Lock()
	defer cf.Unlock()
	return faker.UUIDHyphenated()
}

func (cf *ConcurrentFaker) Words(count int) string {
	words := cf.WordsList(count)

	cf.Lock()
	defer cf.Unlock()

	return strings.Join(words, " ")
}

func (cf *ConcurrentFaker) WordsList(count int) []string {
	cf.Lock()
	defer cf.Unlock()

	words := make([]string, 0, count)

	for i := 0; i < count; i++ {
		words = append(words, faker.Word())
	}

	return words
}

func (cf *ConcurrentFaker) TenantName() string {
	return "t" + cf.RandomLetterNumberBytes(10) + "z"
}

func (cf *ConcurrentFaker) SwarmName() string {
	return "s" + cf.RandomLetterNumberBytes(10) + "z"
}

func (cf *ConcurrentFaker) DistributorName() string {
	return "d" + cf.RandomLetterNumberBytes(10) + "z"
}

func (cf *ConcurrentFaker) UserName() string {
	return "u" + cf.RandomLetterNumberBytes(10) + "z"
}

func (cf *ConcurrentFaker) PolicyName() string {
	return "p" + cf.RandomLetterNumberBytes(10) + "z"
}

func (cf *ConcurrentFaker) ProjectName() string {
	return "pj" + cf.RandomLetterNumberBytes(9) + "z"
}

const letterNumberBytes = "abcdefghijklmnopqrstuvwxyz0123456789-"

func (cf *ConcurrentFaker) RandomLetterNumberBytes(n int) string {
	letterNumberBytes := "abcdefghijklmnopqrstuvwxyz0123456789-"

	return cf.randomString(letterNumberBytes, n)
}

func (cf *ConcurrentFaker) RandomBase64String(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return cf.randomString(letterBytes, n)
}

func (cf *ConcurrentFaker) RandomLowerCaseString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyz"
	return cf.randomString(letterBytes, n)
}

func (cf *ConcurrentFaker) randomString(source string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = source[mathRand.Intn(len(source))]
	}
	return string(b)
}

var Faker = &ConcurrentFaker{}
