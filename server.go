package main

import (
	"log"
	"strconv"

	"github.com/kataras/iris"
)

func getDocument(ctx *iris.Context) {
	storePath, err := findBivouacFile()
	if err != nil {
		log.Fatal(err)
	}

	store, err := LoadStore(storePath)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(iris.StatusOK, store)
}

func getIssues(ctx *iris.Context) {
	storePath, err := findBivouacFile()
	if err != nil {
		log.Fatal(err)
	}

	store, err := LoadStore(storePath)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(iris.StatusOK, store.Issues)
}

func getIssue(ctx *iris.Context) {
	storePath, err := findBivouacFile()
	if err != nil {
		log.Fatal(err)
	}

	store, err := LoadStore(storePath)
	if err != nil {
		log.Fatal(err)
	}

	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Fatal(err)
	}

	issue, err := store.GetIssue(uint(id))
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(iris.StatusOK, issue)
}
