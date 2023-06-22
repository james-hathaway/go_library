package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var libraryTree = &BinaryTree{} // create our main tree

func main() {
	// Load the library from the file at the start
	LoadLibraryFromFile()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Add a book")
		fmt.Println("2. Update a book")
		fmt.Println("3. Delete a book")
		fmt.Println("4. Get details of a book")
		fmt.Println("5. List all books")
		fmt.Println("6. Save to file and exit")
		fmt.Print("Enter option number: ")
		text, _ := reader.ReadString('\n')
		option, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil {
			fmt.Println("Invalid input. Try again.")
			continue // jump back to the top of the loop
		}
		switch option {
		case 1:
			book := readBookDetails(reader)
			AddBook(book)
		case 2:
			fmt.Print("Enter the title of the book to update: ")
			title, _ := reader.ReadString('\n')
			book := readBookDetails(reader)
			UpdateBook(strings.TrimSpace(title), book)
		case 3:
			fmt.Print("Enter the title of the book to delete: ")
			title, _ := reader.ReadString('\n')
			DeleteBook(strings.TrimSpace(title))
		case 4:
			fmt.Print("Enter the title of the book: ")
			title, _ := reader.ReadString('\n')
			book := GetBook(strings.TrimSpace(title))
			if book != nil {
				printBookDetails(*book)
			} else {
				fmt.Println("Book not found.")
			}
		case 5:
			printAllBooks()
		case 6:
			// Save the library to the file before exiting
			SaveLibraryToFile()
			return
		default:
			fmt.Println("Invalid option. Try again.")
		}
	}
}

func readBookDetails(reader *bufio.Reader) Book {
	fmt.Print("Enter title: ")
	title, _ := reader.ReadString('\n')
	fmt.Print("Enter author: ")
	author, _ := reader.ReadString('\n')
	fmt.Print("Enter publication year: ")
	yearStr, _ := reader.ReadString('\n')
	year, _ := strconv.Atoi(strings.TrimSpace(yearStr))
	fmt.Print("Enter genre: ")
	genre, _ := reader.ReadString('\n')
	return Book{
		Title:           strings.TrimSpace(title),
		Author:          strings.TrimSpace(author),
		PublicationYear: year,
		Genre:           strings.TrimSpace(genre),
	}
}

func printBookDetails(book Book) {
	fmt.Println("\nTitle:", book.Title)
	fmt.Println("Author:", book.Author)
	fmt.Println("Publication Year:", book.PublicationYear)
	fmt.Println("Genre:", book.Genre)
}

func printAllBooks() {
	if len(BooksMap) == 0 {
		fmt.Println("No books in the library.")
		return
	}
	for _, book := range BooksMap {
		printBookDetails(*book)
		fmt.Println()
	}
}

type Book struct { // Make a struct for the book's traits
	Title           string
	Author          string
	PublicationYear int
	Genre           string
}

type Node struct { //create our tree nodes
	Book  Book
	Left  *Node
	Right *Node
}

type BinaryTree struct { // create binary tree
	Root *Node
}

var BooksMap = make(map[string]*Book) //create hashmap to store the books efficiently

func (t *BinaryTree) Insert(book Book) *BinaryTree {
	if t.Root == nil {
		t.Root = &Node{Book: book}
	} else {
		t.Root.Insert(book)
	}
	return t
}

func (n *Node) Insert(book Book) {
	if n == nil {
		return
	} else if book.Title < n.Book.Title {
		if n.Left == nil {
			n.Left = &Node{Book: book}
		} else {
			n.Left.Insert(book)
		}
	} else {
		if n.Right == nil {
			n.Right = &Node{Book: book}
		} else {
			n.Right.Insert(book)
		}
	}
}

func AddBook(book Book) {
	BooksMap[book.Title] = &book
	// Also add the book to the binary tree
	libraryTree.Insert(book)
}

func UpdateBook(title string, book Book) {
	// Update the hashmap
	BooksMap[title] = &book
	// Rebuild the tree as it's not easy to update one node in the tree
	rebuildTree()
}

func DeleteBook(title string) {
	// Remove from the hashmap
	delete(BooksMap, title)
	// Rebuild the tree
	rebuildTree()
}

func GetBook(title string) *Book {
	return BooksMap[title]
}

func rebuildTree() {
	libraryTree = &BinaryTree{}
	for _, book := range BooksMap {
		libraryTree.Insert(*book)
	}
}

func SaveLibraryToFile() {
	// Convert the BooksMap to JSON
	data, err := json.Marshal(BooksMap)
	if err != nil {
		fmt.Println("Error while saving library:", err)
		return
	}

	// Write the JSON to the file
	err = ioutil.WriteFile("library.txt", data, 0644)
	if err != nil {
		fmt.Println("Error while saving library:", err)
	}
}

func LoadLibraryFromFile() {
	// Read the file
	data, err := ioutil.ReadFile("library.txt")
	if err != nil {
		fmt.Println("Error while loading library:", err)
		return
	}

	// Convert the JSON data to BooksMap
	err = json.Unmarshal(data, &BooksMap)
	if err != nil {
		fmt.Println("Error while loading library:", err)
		return
	}

	// Rebuild the tree from the loaded BooksMap
	rebuildTree()
}
