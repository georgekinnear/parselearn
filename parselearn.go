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

func parseLearnReceipt(inputPath string) (Submission, error) {

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
			processSubmission(scanner.Text(), &sub)
		case strings.HasPrefix(line, "Comments:"):
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

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case strings.HasPrefix(line, "Original filename:"):
			processOriginalFilename(line, &sub)
			sub.NumberOfFiles++
		case strings.HasPrefix(line, "Filename:"):
			processFilename(line, &sub)
		default:
			continue
		}

	}

	return sub, scanner.Err()
}

//Name: First Last (sxxxxxxx)
func processName(line string, sub *Submission) {
	tokens := strings.Split(line, " ")
	sub.FirstName = strings.TrimSpace(tokens[1])
	sub.LastName = strings.TrimSpace(tokens[2])
	matric := strings.TrimSpace(tokens[3])
	matric = strings.TrimPrefix(matric, "(")
	matric = strings.TrimSuffix(matric, ")")
	sub.Matriculation = strings.TrimSpace(matric)
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

}

//Comments:
//There are no student comments for this assignment
func processComments(line string, sub *Submission) {

}

//Files:
//	Original filename: OnlineExam-Bxxxxxx.pdf
//	Filename: Practice Exam Drop Box_sxxxxxxx_attempt_yyyy-mm-dd-hh-mm-ss_OnlineExam-Bxxxxxx.pdf
func processOriginalFilename(line string, sub *Submission) {

}
func processFilename(line string, sub *Submission) {

}

func writeSubmissionsToCSV(subs []Submission, outputPath string) error {
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
