package main

import (
	"fmt"
	"strconv"

	"github.com/kataras/iris"
)

func getDocument(ctx *iris.Context) {
	storePath, err := findBivouacFile()
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	store, err := LoadStore(storePath)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(iris.StatusOK, store)
}

func getIssues(ctx *iris.Context) {
	storePath, err := findBivouacFile()
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	store, err := LoadStore(storePath)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(iris.StatusOK, store.Issues)
}

func getIssue(ctx *iris.Context) {
	storePath, err := findBivouacFile()
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	store, err := LoadStore(storePath)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	issueIDParam := ctx.Param("issue_id")

	issueID, err := strconv.Atoi(issueIDParam)
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	issue, err := store.GetIssue(uint(issueID))
	if err != nil {
		ctx.EmitError(iris.StatusNotFound)
		return
	}

	ctx.JSON(iris.StatusOK, issue)
}

func addIssue(ctx *iris.Context) {
	issue := new(Issue)
	if err := ctx.ReadJSON(issue); err != nil {
		fmt.Println(string(ctx.PostBody()))
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	storePath, err := findBivouacFile()
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	store, err := LoadStore(storePath)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	issue.ID = store.getNextID()
	store.AddIssue(*issue)
	store.Write()

	ctx.JSON(iris.StatusCreated, *issue)
}

func getComments(ctx *iris.Context) {
	storePath, err := findBivouacFile()
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	store, err := LoadStore(storePath)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	issueIDParam := ctx.Param("issue_id")

	issueID, err := strconv.Atoi(issueIDParam)
	if err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	issue, err := store.GetIssue(uint(issueID))
	if err != nil {
		ctx.EmitError(iris.StatusNotFound)
		return
	}

	ctx.JSON(iris.StatusOK, issue.Comments)
}
