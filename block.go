package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/jomei/notionapi"
)

func BlockToMarkdown(block notionapi.Block) string {

	var blockContent strings.Builder

	switch blk := block.(type) {

	case *notionapi.CodeBlock:
		blockContent.WriteString("```" + blk.Code.Language + "\n")
		for _, code := range blk.Code.RichText {
			blockContent.WriteString(code.Text.Content)
		}
		blockContent.WriteString("\n ``` \n")
		break
	case *notionapi.Heading1Block:
		for _, h2 := range blk.Heading1.RichText {
			blockContent.WriteString("# " + h2.PlainText + "\n")
		}
		break
	case *notionapi.Heading2Block:
		for _, h2 := range blk.Heading2.RichText {
			blockContent.WriteString("## " + h2.PlainText + "\n")
		}
		break
	case *notionapi.Heading3Block:
		for _, h2 := range blk.Heading3.RichText {
			blockContent.WriteString("### " + h2.PlainText + "\n")
		}
		break
	case *notionapi.ParagraphBlock:
		for _, b := range blk.Paragraph.RichText {
			if b.Annotations.Code {
				blockContent.WriteString("`" + b.Text.Content + "`")
			} else {
				blockContent.WriteString(b.Text.Content)
			}
		}

		blockContent.WriteString("\n")

		break
	case *notionapi.BulletedListItemBlock:
		for _, b := range blk.BulletedListItem.RichText {
			blockContent.WriteString("- " + b.PlainText + "\n")
		}
		break
	default:
		log.Warn("Unkown block type", "block", fmt.Sprintf("%T\n", blk))
	}

	return blockContent.String()

}
