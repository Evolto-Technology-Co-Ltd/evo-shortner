package endpoint

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ariebrainware/rama-shortner-backend/external"
	"github.com/ariebrainware/rama-shortner-backend/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/itchyny/base58-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type shortURLRequest struct {
	URL string `json:"url"`
}

func ShortURL(c *gin.Context) {
	request := &shortURLRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, &model.Response{
			Success: true,
			Error:   fmt.Errorf("fail to short the url"),
			Msg:     "",
			Data:    nil,
		})
	}

	// Prepare Insert to MongoDB
	collection := external.GetMongoConn(os.Getenv("MONGO_COLLECTION"))
	shortLink := generateShortLink(request.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.D{
		{"url", request.URL},
		{"short_url", shortLink},
	})
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, &model.Response{
			Success: false,
			Error:   fmt.Errorf("fail to short the url"),
			Msg:     "",
			Data:    nil,
		})
	}
	data := map[string]interface{}{
		"id":  res.InsertedID,
		"key": fmt.Sprintf("%s/%s", os.Getenv("ROOT_HOST"), shortLink),
	}
	c.JSON(http.StatusOK, &model.Response{
		Success: true,
		Error:   nil,
		Msg:     "success",
		Data:    data,
	})
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

func generateShortLink(initialLink string) string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	urlHashBytes := sha256Of(initialLink + u.String())
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	_keyLength := os.Getenv("KEY_LENGTH")
	keyLength, err := strconv.Atoi(_keyLength)
	if err != nil {
		log.Error(err)
	}
	return finalString[:keyLength]
}

type result struct {
	URL string
}

func GetURL(c *gin.Context) {
	key := c.Param("key")
	filter := bson.D{{"short_url", key}}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := external.GetMongoConn(os.Getenv("MONGO_COLLECTION"))
	res := &result{}
	err := collection.FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		// Do something when no record was found
		fmt.Println("record does not exist")
	} else if err != nil {
		log.Fatal(err)
	}
	// Do something with result...
	c.Redirect(http.StatusTemporaryRedirect, res.URL)
	return
}
