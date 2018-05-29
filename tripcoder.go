package main

import "fmt"
import "crypto/sha256"
import "encoding/base64"
import "sync"
import "strings"
import "bufio"
import "os"

func md(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	md  := hash.Sum(nil)
	b64 := base64.StdEncoding.EncodeToString(md)
	
	return b64
}

//we can check if another routine has discovered the tripcode let's say every 1000 iterations
func main() {
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter search string: ")
  searchString, _ := reader.ReadString('\n')
  searchString = strings.Replace(searchString, "\n", "", -1)
  
  findMatch(searchString)
}

func findMatchRoutine(searchString string, seed string, matches chan string, wg *sync.WaitGroup) {
  defer wg.Done()
  
  for {
    hashText := md(seed)
    seed = hashText
    
    if strings.Contains(hashText,searchString) {
      matches <- hashText
    }
  }
}

func printMatchRoutine(matches chan string, wg *sync.WaitGroup) {
  defer wg.Done()
  
  for {
    hashText := <-matches
    fmt.Println(hashText)
  }
}

func findMatch(searchString string) {
  var wg sync.WaitGroup
  
  matches := make(chan string, 100)
  
  numRoutines := 4
  wg.Add(numRoutines+1)
  for i:=0; i<numRoutines; i++ {
    go findMatchRoutine(searchString,string(i),matches,&wg)
  }
  go printMatchRoutine(matches,&wg)
  wg.Wait()
}