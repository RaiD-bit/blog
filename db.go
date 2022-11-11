package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type DBconfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBNAME   string `yaml:"dbname"`
}

func ReadConfig(configPath string) (*DBconfig, error) {
	// create a config structure
	config := &DBconfig{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func dbCreateArticle(artice *Article) error {
	query, err := db.Prepare(`insert into articles(title,content) values ($1,$2);`)
	defer query.Close()

	if err != nil {
		return err
	}

	_, err = query.Exec(artice.Title, artice.Content)

	if err != nil {
		return err
	}
	return nil
}

func getAllArticles() ([]*Article, error) {
	query, err := db.Prepare("select id, title, content from articles;")
	defer query.Close()

	if err != nil {
		return nil, err
	}
	result, err := query.Query()
	articlesList := make([]*Article, 0)
	for result.Next() {
		data := new(Article)
		err := result.Scan(&data.ID, &data.Title, &data.Content)
		if err != nil {
			return nil, err
		}
		articlesList = append(articlesList, data)
	}
	return articlesList, nil
}

func getArticle(articleId string) (*Article, error) {
	query, err := db.Prepare("select id, title, content from articles where id = $1")
	defer query.Close()
	if err != nil {
		return nil, err
	}
	result := query.QueryRow(articleId)
	data := new(Article)
	err = result.Scan(&data.ID, &data.Title, &data.Content)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func updateArticle(articleId string, article *Article) error {
	query, err := db.Prepare("update articles set (title, content) = ($1,$2) where id=$3")
	defer query.Close()

	if err != nil {
		return err
	}
	_, err = query.Exec(article.Title, article.Content, article.ID)
	if err != nil {
		return err
	}
	return nil
}

func deleteArticle(articleId string) error {
	query, err := db.Prepare("delete from articles where id=$1")
	if err != nil {
		return err
	}
	_, err = query.Exec(articleId)

	if err != nil {
		return err
	}
	return nil
}

func connect() (*sql.DB, error) {
	var err error
	dbConf, err := ReadConfig("./dbconfig.yaml")
	if err != nil {
		log.Fatal(err)
	}
	// prepare db string
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBNAME)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		return nil, err
	}
	sqlStmt1 := `
    create table if not exists articles (id serial primary key, title text, content text);
  `
	_, err = db.Exec(sqlStmt1)
	if err != nil {
		return nil, err
	}
	return db, nil
}
