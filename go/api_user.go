/*
 * Swagger Blog
 *
 * Simple Blog
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"encoding/json"
	"net/http"

	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sysu327/Server/dal/db"
	"github.com/sysu327/Server/dal/model"
	"github.com/dgrijalva/jwt-go"
)

func ArticleIdCommentPost(w http.ResponseWriter, r *http.Request) {
	db.Init()
	token, isValid := ValidateToken(w, r)
	if isValid == false {
		Response(MyResponse{
			nil,
			"authentication fail",
		}, w, http.StatusBadRequest)
		return
	}

	var comment model.Comment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		Response(MyResponse{
			nil,
			err.Error(),
		}, w, http.StatusBadRequest)
		return
	}

	if v, ok := token.Claims.(jwt.MapClaims); ok {
		name, _ := v["name"].(string)
		comment.User = name
	}

	articleId := strings.Split(r.URL.Path, "/")[2]
	comment.ArticleId, err = strconv.ParseInt(articleId, 10, 64)
	if err != nil {
		Response(MyResponse{
			nil,
			err.Error(),
		}, w, http.StatusBadRequest)
		return
	}

	comment.Date = fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day())

	articles := db.GetArticles(comment.ArticleId, 0)

	if len(articles) == 0 {
		Response(MyResponse{
			nil,
			"articles not found",
		}, w, http.StatusBadRequest)
		return
	}

	for i := 0; i < len(articles); i++ {
		articles[i].Comments = append(articles[i].Comments, comment)
	}

	err = db.PutArticles(articles)

	if err != nil {
		Response(MyResponse{
			nil,
			err.Error(),
		}, w, http.StatusBadRequest)
		return
	}

	Response(MyResponse{
		comment,
		nil,
	}, w, http.StatusOK)
}

func UserLoginPost(w http.ResponseWriter, r *http.Request) {

	db.Init()

	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		Response(MyResponse{
			nil,
			"parameter error",
		}, w, http.StatusBadRequest)
		return
	}

	check := db.GetUser(user.Username)

	if check.Username != user.Username || check.Password != user.Password {
		Response(MyResponse{
			nil,
			"username or password error",
		}, w, http.StatusBadRequest)
		return
	}

	// token := jwt.New(jwt.SigningMethodHS256)
	// claims := make(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	// claims["iat"] = time.Now().Unix()
	// claims["name"] = user.Username
	// token.Claims = claims
	// tokenString, err := token.SignedString([]byte(SecretKey))

	tokenString, err := SignToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		tokenString = "signing token error"
	}

	Response(MyResponse{
		tokenString,
		nil,
	}, w, http.StatusOK)
}

func UserRegisterPost(w http.ResponseWriter, r *http.Request) {

	db.Init()

	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		Response(MyResponse{
			nil,
			"parameter error",
		}, w, http.StatusBadRequest)
		return
	}

	check := db.GetUser(user.Username)

	if check.Username != "" {
		Response(MyResponse{
			nil,
			"username existed",
		}, w, http.StatusBadRequest)
		return
	}

	err = db.PutUsers([]model.User{user})

	if err != nil {
		Response(MyResponse{
			nil,
			err.Error(),
		}, w, http.StatusBadRequest)
		return
	}
	tokenString, err := SignToken(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		tokenString = "signing token error"
	}
	Response(MyResponse{
		tokenString,
		nil,
	}, w, http.StatusOK)

}