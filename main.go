package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jomei/notionapi"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {

	notionKey := os.Getenv("NOTION_KEY")

	db := flag.String("db", "", "Id of the db were are the pages to create")

	folder := flag.String("folder", "./", "The folder were the pages will be created")

	flag.Parse()

	log.Debug("Notion key", "notionkey", notionKey)

	client := notionapi.NewClient(notionapi.Token(notionKey))

	md := NewMDConverter(client, *db, *folder)

	md.ConvertPagesFromDBToMd()

}

type markdown struct {
	client     *notionapi.Client
	db         string
	folderPath string
}

func NewMDConverter(client *notionapi.Client, db, folderPath string) markdown {

	return markdown{
		client,
		db,
		folderPath,
	}
}

func (md *markdown) PageToMarkdown(page notionapi.Page, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Info("Starting converting page: " + page.ID.String())

	var fileName string
	var title string
	var desc string
	var tags strings.Builder
	var fileContent strings.Builder

	fileName += page.ID.String() + "_"

	for _, prop := range page.Properties {

		switch prp := prop.(type) {

		case *notionapi.TitleProperty:
			for _, tlt := range prp.Title {
				fileName += tlt.PlainText
				title += tlt.PlainText
			}

		case *notionapi.RichTextProperty:
			for _, d := range prp.RichText {
				desc += d.PlainText
			}

		case *notionapi.MultiSelectProperty:
			for _, t := range prp.MultiSelect {
				tags.WriteString("  - " + t.Name + "\n")
			}
		}
	}

	fileName = strings.ToLower(fileName)
	fileName = strings.ReplaceAll(fileName, " ", "-")
	fileName += ".md"

	metadataMap := map[string]string{
		"pubDateTime": page.CreatedTime.Format(time.RFC3339),
		"title":       title,
		"postSlug":    fileName,
		"description": desc,
		"tags":        tags.String(),
	}

	fileContent.WriteString("---\n")

	for k, v := range metadataMap {

		if k == "tags" {
			fileContent.WriteString(k + ": \n" + v + "\n")
		} else {
			fileContent.WriteString(k + ": " + v + "\n")
		}

	}

	fileContent.WriteString("---\n")

	blocks, err := md.client.Block.GetChildren(context.TODO(), notionapi.BlockID(page.ID), &notionapi.Pagination{})

	if err != nil {
		log.Fatal(err)
	}

	for _, block := range blocks.Results {
		switch blk := block.(type) {

		case *notionapi.CodeBlock:
			fileContent.WriteString("```" + blk.Code.Language + "\n")
			for _, code := range blk.Code.RichText {
				fileContent.WriteString(code.Text.Content)
			}
			fileContent.WriteString("\n ``` \n")
			break
		case *notionapi.Heading1Block:
			for _, h2 := range blk.Heading1.RichText {
				fileContent.WriteString("# " + h2.PlainText + "\n")
			}
			break
		case *notionapi.Heading2Block:
			for _, h2 := range blk.Heading2.RichText {
				fileContent.WriteString("## " + h2.PlainText + "\n")
			}
			break
		case *notionapi.Heading3Block:
			for _, h2 := range blk.Heading3.RichText {
				fileContent.WriteString("### " + h2.PlainText + "\n")
			}
			break
		case *notionapi.ParagraphBlock:
			for _, b := range blk.Paragraph.RichText {
				if b.Annotations.Code {
					fileContent.WriteString("`" + b.Text.Content + "`")
				} else {
					fileContent.WriteString(b.Text.Content)
				}
			}

			fileContent.WriteString("\n")

			break
		default:
			fmt.Printf("%T\n", blk)
		}
	}

	errFile := os.WriteFile(md.folderPath+fileName, []byte(fileContent.String()), 0666)

	if errFile != nil {
		log.Fatal(errFile)
	}

	log.Info("Done converting page: " + page.ID.String())
}

func (md *markdown) ConvertPagesFromDBToMd() {
	log.Info("STARTING CONVERTING PAGES MD")

	var wg sync.WaitGroup

	dbquery := notionapi.DatabaseQueryRequest{
		Filter: notionapi.PropertyFilter{
			Property: "Publish",
			Checkbox: &notionapi.CheckboxFilterCondition{Equals: true}},
	}

	db, err := md.client.Database.Query(context.TODO(), notionapi.DatabaseID(md.db), &dbquery)

	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir(md.folderPath, 0777)

	pagesCount := len(db.Results)

	wg.Add(pagesCount)

	for _, page := range db.Results {

		go md.PageToMarkdown(page, &wg)

	}

	wg.Wait()

	log.Info("Finished")

}
