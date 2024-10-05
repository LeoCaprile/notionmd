package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/jomei/notionapi"
)

func main() {

	notionKey := os.Getenv("NOTION_KEY")

	fmt.Printf("%s \n", notionKey)

	client := notionapi.NewClient(notionapi.Token(notionKey))

	md := NewMD(client, "266919f2772245a6b691217651ae0a17", "./hey/")

	md.ConvertPagesFromDBToMd()

}

type markdown struct {
	client     *notionapi.Client
	db         string
	folderPath string
}

func NewMD(client *notionapi.Client, db, folderPath string) markdown {

	return markdown{
		client,
		db,
		folderPath,
	}
}

func (md *markdown) PageToMarkdown(page notionapi.Page, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Starting converting page: " + page.ID.String())

	var fileName string
	var title string
	var fileContent strings.Builder

	fileName += page.ID.String() + "_"

	for _, prop := range page.Properties {

		switch prp := prop.(type) {
		case *notionapi.TitleProperty:
			for _, tlt := range prp.Title {
				fileName += tlt.PlainText
				title += tlt.PlainText
			}
		}
	}

	fileName = strings.ToLower(fileName)
	fileName = strings.ReplaceAll(fileName, " ", "-")
	fileName += ".md"

	metadataMap := map[string]string{
		"pubDateTime": page.CreatedTime.Format("YYYY-MM-DD"),
		"title":       title,
		"postSlug":    fileName,
	}

	fileContent.WriteString("---\n")

	for k, v := range metadataMap {
		fileContent.WriteString(k + ": " + v + "\n")
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

	fmt.Println("Done converting page: " + page.ID.String())
}

func (md *markdown) ConvertPagesFromDBToMd() {
	fmt.Println("STARTING CONVERTING PAGES MD")

	dbquery := notionapi.DatabaseQueryRequest{
		Filter: notionapi.PropertyFilter{
			Property: "Publish",
			Checkbox: &notionapi.CheckboxFilterCondition{Equals: true}},
	}
	var wg sync.WaitGroup

	db, err := md.client.Database.Query(context.TODO(), notionapi.DatabaseID(md.db), &dbquery)

	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir(md.folderPath, 0777)

	pagesCount := len(db.Results)
	fmt.Printf("PAGESCOUNT: %d", pagesCount)
	wg.Add(pagesCount)

	for _, page := range db.Results {

		go md.PageToMarkdown(page, &wg)

	}

	fmt.Println("Number of Goroutines:", runtime.NumGoroutine())

	wg.Wait()

	fmt.Println("Finished")

}
