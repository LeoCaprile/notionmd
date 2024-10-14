package main

import (
	"testing"

	"github.com/jomei/notionapi"
)

func TestHeader1BlockMdParsing(t *testing.T) {

	header := &notionapi.Heading1Block{
		Heading1: notionapi.Heading{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "Hey"}, PlainText: "Hey"}},
		},
	}

	result := BlockToMarkdown(header)

	expected := "# Hey\n"

	if result != expected {

		t.Fatalf("expected '%v', got %v", expected, result)

	}

}

func TestHeader2BlockMdParsing(t *testing.T) {

	header := &notionapi.Heading2Block{
		Heading2: notionapi.Heading{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "Hey"}, PlainText: "Hey"}},
		},
	}

	result := BlockToMarkdown(header)

	expected := "## Hey\n"

	if result != expected {

		t.Fatalf("expected '%v', got %v", expected, result)

	}
}

func TestHeader3BlockMdParsing(t *testing.T) {

	header := &notionapi.Heading3Block{
		Heading3: notionapi.Heading{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "Hey"}, PlainText: "Hey"}},
		},
	}

	result := BlockToMarkdown(header)

	expected := "### Hey\n"

	if result != expected {

		t.Fatalf("expected '%v', got %v", expected, result)

	}
}

func TestParagraphBlockMdParsing(t *testing.T) {

	paragraph := &notionapi.ParagraphBlock{

		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "Hey"}, Annotations: &notionapi.Annotations{Code: false}, PlainText: "Hey"}},
		},
	}

	result := BlockToMarkdown(paragraph)

	expected := "Hey\n"

	if result != expected {

		t.Fatalf("expected '%v', got %v", expected, result)

	}
}

func TestParagraphBlockWithCodeWithinMdParsing(t *testing.T) {

	paragraph := &notionapi.ParagraphBlock{

		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "Hey "}, Annotations: &notionapi.Annotations{Code: false}, PlainText: "Hey "},
				{Text: &notionapi.Text{Content: "Hey"}, Annotations: &notionapi.Annotations{Code: true}, PlainText: "Hey"}},
		},
	}

	result := BlockToMarkdown(paragraph)

	expected := "Hey `Hey`\n"

	if result != expected {

		t.Fatalf("expected '%v', got %v", expected, result)

	}
}
