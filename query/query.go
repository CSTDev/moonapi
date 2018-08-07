package query

import (
	"bytes"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// Order provides the different options for sorting
type Order string

// Filter specifies the type of problems to filter by
type Filter string

// Configuration defines the Board angle
type Configuration string

// HoldSet specifies the types of hold to include in the query
type HoldSet string

// Grade specifies the grade range available
type Grade int

const descString = "-desc"
const ascString = "-asc"

const (
	Newest     Order = "New"
	Difficulty Order = "Grade"
	Rating     Order = "Rating"
	Repeats    Order = "Repeats"
)

const (
	Forty  Configuration = "40° MoonBoard"
	Twenty Configuration = "25° MoonBoard"
)

const (
	OS   HoldSet = "original school holds"
	Wood HoldSet = "wooden holds"
	A    HoldSet = "hold set a"
	B    HoldSet = "hold set b"
	C    HoldSet = "hold set c"
)

const (
	Benchmarks Filter = "Benchmarks"
	SetByMe    Filter = "Setbyme"
	MyAscents  Filter = "Myascents"
)

const (
	FivePlus Grade = iota
	SixA
	SixAPlus
	SixB
	SixBPlus
	SixC
	SixCPlus
	SevenA
	SevenAPlus
	SevenB
	SevenBPlus
	SevenC
	SevenCPlus
	EightA
	EightAPlus
	EightB
	EightBPlus
)

var gradeStrings = [...]string{"5+", "6A", "6A+", "6B", "6B+", "6C", "6C+", "7A", "7A+", "7B", "7B+", "7C", "7C+", "8A", "8A+", "8B", "8B+"}

// Query contains the required strings to pass to the website in a request
type Query interface {
	Sort() string
	Filter() string
	Page() int
	PageSize() int
}

type query struct {
	term     string
	sort     string
	filter   string
	page     int
	pageSize int
}

func (q *query) Filter() string {
	return q.filter
}

func (q *query) Sort() string {
	return q.sort
}

func (q *query) Page() int {
	return q.page
}

func (q *query) PageSize() int {
	return q.pageSize
}

// QueryBuilder assists with the building of a Query
type QueryBuilder interface {
	Term(searchTerm string) QueryBuilder
	Sort(order Order, asc bool) QueryBuilder
	Configuration(filter Configuration) QueryBuilder
	HoldSet(filter HoldSet) QueryBuilder
	Filter(filter Filter) QueryBuilder
	MinGrade(min Grade) QueryBuilder
	MaxGrade(max Grade) QueryBuilder
	Page(page int) QueryBuilder
	PageSize(pageSize int) QueryBuilder
	Build() (Query, []error)
}

type queryBuilder struct {
	term          string
	order         string
	configuration string
	holdSet       string
	filter        string
	minGrade      Grade
	maxGrade      Grade
	page          int
	pageSize      int
	error         []error
}

// New creates a new QueryBuilder with a default min and max grade
func New() QueryBuilder {
	qb := queryBuilder{
		minGrade: FivePlus,
		maxGrade: EightBPlus,
		order:    "",
		pageSize: 15,
		page:     1,
	}
	return &qb
}

// Term adds a search term to the query, usually the name of
// the problem being searched for.
// Default: empty string
func (qb *queryBuilder) Term(searchTerm string) QueryBuilder {
	qb.term = "Name~contains~'" + searchTerm + "'"
	return qb
}

// Sort sets the sort order that will be used in the query.
// Only one sort order can be provided, if more than one is then
// an error is added to be returned and the last provided Order is used.
// Default: is Newest problems first
func (qb *queryBuilder) Sort(order Order, asc bool) QueryBuilder {
	if qb.order != "" {
		qb.error = append(qb.error, errors.New("can only sort by one parameter, defaulting to the last provided"))
	}

	switch order {
	case Newest:
		qb.order = string(order) + descString
	case Difficulty:
		qb.order = string(order) + ascending(asc)
	case Rating:
		qb.order = string(order) + descString
	case Repeats:
		qb.order = string(order) + ascending(asc)
	}
	return qb
}

func ascending(asc bool) string {
	if asc {
		return "Asc" + ascString
	}
	return "Desc" + descString
}

// Configuration sets the angle configuration of the board to use.
// Default: is all board configurations (40 and 20 degree)
func (qb *queryBuilder) Configuration(filter Configuration) QueryBuilder {
	unescaped, err := url.QueryUnescape(string(filter))
	if err != nil {
		panic(err)
	}
	if qb.configuration == "" {
		if filter == Forty {
			qb.minGrade = SixAPlus
		}
		qb.configuration = "Configuration~eq~'" + string(filter) + "'"
	} else {
		qb.minGrade = FivePlus
		qb.configuration = strings.TrimSuffix(qb.configuration, "'") + "," + unescaped + "'"
	}
	return qb
}

// HoldSet sets the hold sets that problems can include.
// Default: is all hold sets
func (qb *queryBuilder) HoldSet(filter HoldSet) QueryBuilder {
	var buffer bytes.Buffer
	if qb.holdSet == "" {
		buffer.WriteString("Holdsets~eq~'")
		buffer.WriteString(string(filter))
		buffer.WriteString("'")
		qb.holdSet = buffer.String()
	} else {
		qb.holdSet = strings.TrimSuffix(qb.holdSet, "'") + "," + string(filter) + "'"
	}
	return qb
}

// Filter specifies how to filter problems.
// Default: is not to filter
func (qb *queryBuilder) Filter(filter Filter) QueryBuilder {
	var buffer bytes.Buffer
	if qb.filter == "" {
		buffer.WriteString(string(filter))
		buffer.WriteString("~eq~''")
	} else {
		buffer.WriteString(qb.filter)
		addAnd(&buffer, string(filter))
		buffer.WriteString("~eq~''")
	}

	qb.filter = buffer.String()

	return qb
}

// MinGrade sets the mininum grade for the problems being searched.
// Default: FivePlus unless configuration is set to only Forty
func (qb *queryBuilder) MinGrade(min Grade) QueryBuilder {
	qb.minGrade = min
	return qb
}

// MaxGrade sets the maximum grade for the problems being searched.
// Default: EightBPlus
func (qb *queryBuilder) MaxGrade(max Grade) QueryBuilder {
	qb.maxGrade = max
	return qb
}

// Page specifies which page of results to return
func (qb *queryBuilder) Page(page int) QueryBuilder {
	if page < 1 {
		qb.error = append(qb.error, errors.New("page number cannot be below 1"))
	} else {
		qb.page = page
	}
	return qb
}

// PageSize specifies the number of results to return per page
func (qb *queryBuilder) PageSize(pageSize int) QueryBuilder {
	if pageSize > 100 || pageSize < 1 {
		qb.error = append(qb.error, errors.New("Page size must be between 1 and 100"))
	} else {
		qb.pageSize = pageSize
	}
	return qb
}

// Build takes the queryBuilder and constructs a query.
// All set values from the queryBuilder are converted to strings and added to either
// sort or fitler parameters of a query. This is the format required by the website.
// An array of errors are returned if any of the values set are invalid.
func (qb *queryBuilder) Build() (Query, []error) {

	if qb.minGrade > qb.maxGrade {
		qb.error = append(qb.error, errors.New("min grade cannot be higher than max grade"))
	}

	var buffer bytes.Buffer
	buffer.WriteString(qb.configuration)

	addAnd(&buffer, qb.term)
	addAnd(&buffer, qb.holdSet)
	addAnd(&buffer, qb.filter)
	addAnd(&buffer, strings.Replace("MinGrade~eq~'5+'", "5+", gradeStrings[qb.minGrade], -1))
	addAnd(&buffer, strings.Replace("MaxGrade~eq~'8B+'", "8B+", gradeStrings[qb.maxGrade], -1))

	filterString := buffer.String()

	query := &query{
		sort:     qb.order,
		filter:   filterString,
		page:     qb.page,
		pageSize: qb.pageSize,
	}

	if len(qb.error) > 0 {
		return query, qb.error
	}

	return query, nil
}

func addAnd(buffer *bytes.Buffer, nextString string) {
	if buffer.Len() != 0 && nextString != "" {
		buffer.WriteString("~and~")
	}
	buffer.WriteString(nextString)
}

func isValidGrade(grade string) bool {
	match, _ := regexp.MatchString("[6-8]([A|B|C])?(\\+)?", grade)
	return match
}

// ToOrder takes a string and returns its corresponding Order value
// errors if the string passed is not valid
func ToOrder(order string) (*Order, error) {
	var orderType Order
	switch strings.ToLower(order) {
	case "new":
		orderType = Newest
	case "grade":
		orderType = Difficulty
	case "rating":
		orderType = Rating
	case "repeats":
		orderType = Repeats
	default:
		return nil, errors.New("String passed to ToOrder was not a valid order value")
	}
	return &orderType, nil
}

// ToConfiguration takes a string and returns its corresponding Configuration value
// errors if the string passed is not valid
func ToConfiguration(config string) (*Configuration, error) {

	var configType Configuration
	switch strings.ToLower(config) {
	case "forty":
		configType = Forty
	case "twenty":
		configType = Twenty
	default:
		return nil, errors.New("String passed to ToConfiguration was not a valid configuration")
	}
	return &configType, nil
}

// ToHoldSet takes a string and returns its corresponding HoldSet value
// errors if the string passed is not valid
func ToHoldSet(holdSet string) (*HoldSet, error) {
	var holdSetType HoldSet
	switch strings.ToLower(holdSet) {
	case "os":
		holdSetType = OS
	case "wood":
		holdSetType = Wood
	case "a":
		holdSetType = A
	case "b":
		holdSetType = B
	case "c":
		holdSetType = C
	default:
		return nil, errors.New("String passed to ToHoldSet was not a valid Hold Set")
	}
	return &holdSetType, nil
}

// ToFilter takes a string and returns its corresponding Filter value
// errors if the string passed is not valid
func ToFilter(filter string) (*Filter, error) {
	var filterType Filter
	switch strings.ToLower(filter) {
	case "benchmarks":
		filterType = Benchmarks
	case "setbyme":
		filterType = SetByMe
	case "myascents":
		filterType = MyAscents
	default:
		return nil, errors.New("String passed to ToFilter was not a valid Filter")
	}
	return &filterType, nil
}

// ToGrade takes a string and returns its corresponding Grade value
// errors if the string passed is not valid
func ToGrade(grade string) (*Grade, error) {

	if !isValidGrade(grade) {
		return nil, errors.New("String passed to ToGrade was not a valid Grade")
	}
	var gradeType Grade
	switch strings.ToUpper(grade) {
	case "5+":
		gradeType = FivePlus
	case "6A":
		gradeType = SixA
	case "6A+":
		gradeType = SixAPlus
	case "6B":
		gradeType = SixB
	case "6B+":
		gradeType = SixBPlus
	case "6C":
		gradeType = SixC
	case "6C+":
		gradeType = SixCPlus
	case "7A":
		gradeType = SevenA
	case "7A+":
		gradeType = SevenAPlus
	case "7B":
		gradeType = SevenB
	case "7B+":
		gradeType = SevenBPlus
	case "7C":
		gradeType = SevenC
	case "7C+":
		gradeType = SevenCPlus
	case "8A":
		gradeType = EightA
	case "8A+":
		gradeType = EightAPlus
	case "8B":
		gradeType = EightB
	case "8B+":
		gradeType = EightBPlus
	}
	return &gradeType, nil
}
