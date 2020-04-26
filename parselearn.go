package parselearn

import (
	"bufio"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

type Submission struct {
	FirstName          string  `csv:"FirstName"`
	LastName           string  `csv:"LastName"`
	Matriculation      string  `csv:"Matriculation"`
	Assignment         string  `csv:"Assignment"`
	DateSubmitted      string  `csv:"DateSubmitted"`
	SubmissionField    string  `csv:"SubmissionField"`
	Comments           string  `csv:"Comments"`
	OriginalFilename   string  `csv:"OriginalFilename"`
	Filename           string  `csv:"Filename"`
	ExamNumber         string  `csv:"ExamNumber"`
	MatriculationError string  `csv:"MatriculationError"`
	ExamNumberError    string  `csv:"ExamNumberError"`
	FiletypeError      string  `csv:"FiletypeError"`
	FilenameError      string  `csv:"FilenameError"`
	NumberOfPages      string  `csv:"NumberOfPages"`
	FilesizeMB         float64 `csv:"FilesizeMB"`
	NumberOfFiles      int     `csv:"NumberOfFiles"`
}

func ParseLearnReceipt(inputPath string) (Submission, error) {

	sub := Submission{}

	file, err := os.Open(inputPath)

	if err != nil {
		return sub, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

SCAN:
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "Name:"):
			processName(line, &sub)
		case strings.HasPrefix(line, "Assignment:"):
			processAssignment(line, &sub)
		case strings.HasPrefix(line, "Date Submitted:"):
			processDateSubmitted(line, &sub)
		case strings.HasPrefix(line, "Submission Field:"):
			scanner.Scan()
			processSubmission(scanner.Text(), &sub)
		case strings.HasPrefix(line, "Comments:"):
			scanner.Scan()
			processComments(scanner.Text(), &sub)
		case strings.HasPrefix(line, "Files:"):
			break SCAN
		default:
			continue
		}
	}

	// now read in the files ....
	// TODO figure out nested csv so we can record multiple files
	// meanwhile for safety, count the number of original files

	sub.NumberOfFiles = 0

	// we will take the first file as the one we expect to see renamed
	gotOriginal := false
	gotNew := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case strings.HasPrefix(line, "Original filename:"):
			if !gotOriginal {
				processOriginalFilename(line, &sub)
				gotOriginal = true
			}
			sub.NumberOfFiles++
		case strings.HasPrefix(line, "Filename:"):
			if !gotNew {
				processFilename(line, &sub)
				gotNew = true
			}
		default:
			continue
		}

	}

	return sub, scanner.Err()
}

//Name: First Last (sxxxxxxx)
func processName(line string, sub *Submission) {

	m := strings.Index(line, ":")
	n := strings.Index(line, "(")
	p := strings.Index(line, ")")

	name := strings.TrimSpace(line[m+1 : n])
	matric := strings.TrimSpace(line[n+1 : p])

	sub.FirstName = "-"
	sub.LastName = name
	sub.Matriculation = matric
}

//Assignment: Practice Exam Drop Box
func processAssignment(line string, sub *Submission) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "Assignment:")
	sub.Assignment = strings.TrimSpace(line)
}

//Date Submitted: Monday, dd April yyyy hh:mm:ss o'clock BST
func processDateSubmitted(line string, sub *Submission) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "Date Submitted:")
	sub.DateSubmitted = strings.TrimSpace(line)
}

//Submission Field:
//There is no student submission text data for this assignment.
func processSubmission(line string, sub *Submission) {
	sub.SubmissionField = strings.TrimSpace(line)

}

//Comments:
//There are no student comments for this assignment
func processComments(line string, sub *Submission) {
	sub.Comments = strings.TrimSpace(line)
}

//Files:
//	Original filename: OnlineExam-Bxxxxxx.pdf
//	Filename: Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_OnlineExam-Bxxxxxx.pdf
func processOriginalFilename(line string, sub *Submission) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "Original filename:")
	sub.OriginalFilename = strings.TrimSpace(line)
}
func processFilename(line string, sub *Submission) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "Filename:")
	sub.Filename = strings.TrimSpace(line)
}

func WriteSubmissionsToCSV(subs []Submission, outputPath string) error {
	// wrap the marshalling library in case we need converters etc later
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	return gocsv.MarshalFile(&subs, file)
}

//Name: First Last (sxxxxxxx)
//Assignment: Practice Exam Drop Box
//Date Submitted: Monday, dd April yyyy hh:mm:ss o'clock BST
//Current Mark: Needs Marking
//
//Submission Field:
//There is no student submission text data for this assignment.
//
//Comments:
//There are no student comments for this assignment.
//
//Files:
//	Original filename: OnlineExam-Bxxxxxx.pdf
//	Filename: Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_OnlineExam-Bxxxxxx.pdf
