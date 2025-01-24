package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	emailRegex *regexp.Regexp
)

func init() {
	var err error
	emailRegex, err = regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if err != nil {
		panic("regex inválido")
	}
}

type SafeMap struct {
	mu    sync.Mutex
	files map[int]*os.File
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		files: make(map[int]*os.File),
	}
}

func (m *SafeMap) GetFile(index int, tempDir string, create bool) (*os.File, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	file, exists := m.files[index]
	if !exists && create {
		newFilename := fmt.Sprintf("%s%d.txt", tempDir, index)
		var err error
		file, err = os.Create(newFilename)
		if err != nil {
			return nil, err
		}
		m.files[index] = file
	}
	return file, nil
}

func (m *SafeMap) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, file := range m.files {
		file.Close()
	}
}

func main() {
	dir := "./"
	tempDir := "./temp/"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Println("Erro ao criar diretório temporário:", err)
		return
	}

	tempFiles := NewSafeMap()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Erro ao listar arquivos:", err)
		return
	}

	var wg sync.WaitGroup

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") || file.Name() == "limpos.txt" {
			continue
		}

		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			processFile(dir+filename, tempFiles, tempDir)
		}(file.Name())
	}

	wg.Wait()
	tempFiles.CloseAll()

	outputFile, err := os.Create(dir + "limpos.txt")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo limpos.txt:", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for i := range tempFiles.files {
		filename := fmt.Sprintf("%s%d.txt", tempDir, i)
		if err := finalizeTempFile(filename, writer); err != nil {
			fmt.Println("Erro ao processar arquivo temporário:", err)
		}
	}
	writer.Flush()

	os.RemoveAll(tempDir)
}

func processFile(filename string, tempFiles *SafeMap, tempDir string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo:", filename, err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if shouldBeIgnored(line) || len(line) == 0 {
			continue
		}
		index := getIndexForLine(line)
		tempFile, err := tempFiles.GetFile(index, tempDir, true)
		if err != nil {
			fmt.Println("Erro ao acessar arquivo temporário:", err)
			continue
		}
		if _, err := tempFile.WriteString(line + "\n"); err != nil {
			fmt.Println("Erro ao escrever no arquivo temporário:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler arquivo:", filename, err)
	}
}

func getIndexForLine(line string) int {
	if len(line) < 2 {
		return 0
	}
	chars := []rune(line[:2])
	index := 0
	for i, char := range chars {
		if char >= '0' && char <= '9' {
			index = index*36 + int(char-'0')
		} else if char >= 'a' && char <= 'z' {
			index = index*36 + int(char-'a') + 10
		} else {
			index = index*36 + 35
		}
		if i == 0 {
			index *= 36
		}
	}
	return index
}

func finalizeTempFile(tempFilename string, writer *bufio.Writer) error {
	f, err := os.Open(tempFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	lines := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if _, exists := lines[line]; !exists {
			lines[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func shouldBeIgnored(line string) bool {
	if strings.Contains(line, "gmail.com") || strings.Contains(line, "facebook.com") ||
		strings.Contains(line, "hotmail.com") || strings.Contains(line, "outlook.com") ||
		strings.Contains(line, "msn.com") {
		return true
	}
	if Armadilhas(line) || NomesRaros(line) || !emailRegex.MatchString(line) || !strings.HasSuffix(line, ".br") {
		return true
	}
	return false
}

func Armadilhas(texto string) bool {
	nomes := []string{"abuse", "abuso", "autoresponse", "autoresposta", "auto-resposta", "bounce", "hacker", "honeypot", "nao_responder", "naoresponder", "noreply", "postmaster", "spam", "spammer"}
	for _, nome := range nomes {
		if strings.Contains(texto, nome) {
			return true
		}
	}
	return false
}

func NomesRaros(texto string) bool {
	nomes := []string{"aadiv", "aahva", "aaradhya", "adhira", "akanksh", "anaisha", "anant", "andrew", "anushka", "ashley", "asmee", "ayaan", "beverly", "billy", "bradley", "brandon", "brittany", "bryan", "cheryl", "dasya", "debra", "dorothy", "drishti", "edward", "heather", "idhant", "ishank", "ishita", "jeffrey", "joseph", "kabir", "kahaan", "kashvi", "kathryn", "keith", "kenneth", "kimaya", "krisha", "laksh", "larry", "lawrence", "mahika", "marilyn", "matthew", "mehar", "mishka", "nehrika", "nimit", "pahal", "parv", "pranay", "prisha", "raunak", "raymond", "reyansh", "rishaan", "rishit", "rohan", "rushil", "saanvi", "sadhil", "sahana", "scott", "shanaya", "shrishti", "sneha", "stephen", "svenn", "taahira", "taarush", "taksh", "tanvi", "timothy", "tyler", "vihaan", "vivaan", "willie", "zachary"}
	for _, nome := range nomes {
		if strings.Contains(texto, nome) {
			return true
		}
	}
	return false
}
