package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"rest-api-sqlBoiler/models"
)

func main() {
	ctx := context.Background()
	db := connectDB()
	boil.SetDB(db) //Global variant
	author := createAuthor(ctx)
	createArticle(ctx, author)
	createArticle(ctx, author)
	selectAuthorWithArticle(ctx, author.ID)
	selectAuthorWithArticleJoin(ctx, author.ID)
}

//function to connect db
func connectDB() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost user=root password=123 dbname=sqlBoiler_db port=5432 sslmode=disable ")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

//function that creates a new author
func createAuthor(ctx context.Context) models.Author {
	author := models.Author{
		Name:  "John Doe",
		Email: "johndoe@gmail.com",
	}
	err := author.InsertG(ctx, boil.Infer()) //InsertG : global variant which is used to insert data into db
	if err != nil {
		log.Fatal(err)
	}
	return author
}

//function that creates a new article
func createArticle(ctx context.Context, author models.Author) models.Article {
	article := models.Article{
		Title:    "Hello World",
		Body:     null.StringFrom("Hello world, this is an article."),
		AuthorID: author.ID,
	}

	err := article.InsertG(ctx, boil.Infer()) //  insert article
	if err != nil {
		log.Fatal(err)
	}

	return article
}

//  function that selects author and articles.
func selectAuthorWithArticle(ctx context.Context, authorID int) {
	author, err := models.Authors(models.AuthorWhere.ID.EQ(authorID)).OneG(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Author: \n\tID:%d \n\tName:%s \n\tEmail:%s\n", author.ID, author.Name, author.Email)

	articles, err := author.Articles().AllG(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range articles {
		fmt.Printf("Article: \n\tID:%d \n\tTitle:%s \n\tBody:%s \n\tCreatedAt:%v\n", a.ID, a.Title, a.Body.String, a.CreatedAt.Time)
	}
}

// since each database row will contain info for author and article.

func selectAuthorWithArticleJoin(ctx context.Context, authorID int) {
	type AuthorAndArticle struct {
		Article models.Article `boil:"article,bind"`
		Author  models.Author  `boil:"author,bind"`
	}

	authorAndArticles := make([]AuthorAndArticle, 0)

	err := models.NewQuery(
		qm.Select("*"),
		qm.From(models.TableNames.Author),
		qm.InnerJoin("article on article.author_id = author.id"),
		models.AuthorWhere.ID.EQ(authorID),
	).BindG(ctx, &authorAndArticles)
	if err != nil {
		log.Fatal(err)
	}

	for _, authorAndArticle := range authorAndArticles {
		author := authorAndArticle.Author
		a := authorAndArticle.Article

		fmt.Printf("Author: \n\tID:%d \n\tName:%s \n\tEmail:%s\n", author.ID, author.Name, author.Email)
		fmt.Printf("Article: \n\tID:%d \n\tTitle:%s \n\tBody:%s \n\tCreatedAt:%v\n", a.ID, a.Title, a.Body.String, a.CreatedAt.Time)
	}
}
