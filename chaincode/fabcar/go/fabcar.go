package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/jsonq"

	utils "github.com/cd1/utils-golang"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the SmartContract structure
type SmartContract struct {
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

type Author struct {
	OwnerId      string
	Name         string
	Sign         Signature
	ownDocsId    []string
	editedDocsId []string
}
type Signature struct {
	authorName string
	DocName    string
	DocId      string
	timeStamp  string
	Sha256     string
}

type SegmentWrite struct {
	segmentId string
	AuthorId  string
	timeStamp string
	text      string
}
type User struct {
	Doctype      string
	Name         string
	Email        string
	PasswordHash string
	Token        string
	Key          string
}
type Docs struct {
	DocId    string
	DocName  string
	OwnerId  string
	EditorId string
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Car struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

var DocId, DocName, OwnerId, EditorId string

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "amarNaam" {
		return s.tellMyName(APIstub, args)
	} else if function == "queryCar" {
		return s.queryCar(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createCar" {
		return s.createCar(APIstub, args)
	} else if function == "queryAllCars" {
		return s.queryAllCars(APIstub)
	} else if function == "changeCarOwner" {
		return s.changeCarOwner(APIstub, args)
	} else if function == "getData" {
		return s.getData(APIstub, args)
	} else if function == "setData" {
		return s.setData(APIstub, args)
	} else if function == "putObject" {
		return s.putObject(APIstub, args)
	} else if function == "getObject" {
		return s.getObject(APIstub, args)
	} else if function == "createDoc" {
		return s.createDoc(APIstub, args)
	} else if function == "createSuccess" {
		return s.createSuccess(APIstub, args)
	} else if function == "accessDoc" {
		return s.accessDoc(APIstub, args)
	} else if function == "register" {
		return s.register(APIstub, args)
	} else if function == "login" {
		return s.login(APIstub, args)
	} else if function == "logout" {
		return s.logout(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*#In invoke.js
    fcn: 'createDoc',
    args: ['firstDoc', 'Shoumik38','Nammi31'],
#In Query.js
    fcn: 'createSuccess',
    args: ['firstDoc']*/
func (s *SmartContract) register(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments, required 3, given " + strconv.Itoa(len(args)))
	}

	name := args[0]
	email := args[1]
	password := args[2]

	h := sha256.New()
	h.Write([]byte(password))
	passwordHash := fmt.Sprintf("%x", h.Sum(nil))

	token := utils.RandomString()

	key := utils.RandomString()

	user := User{"user", name, email, passwordHash, token, key}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = APIstub.PutState(key, jsonUser)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(token))
}
func (s *SmartContract) logout(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, required 1, given " + strconv.Itoa(len(args)))
	}

	key := args[0]
	var user User

	jsonUser, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal(jsonUser, &user)
	if err != nil {
		return shim.Error(err.Error())
	}

	user.Token = utils.RandomString()

	jsonUser, err = json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = APIstub.PutState(key, jsonUser)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) login(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments, required 2, given " + strconv.Itoa(len(args)))
	}

	email := args[0]
	password := args[1]

	h := sha256.New()
	h.Write([]byte(password))
	passwordHash := fmt.Sprintf("%x", h.Sum(nil))

	queryString := fmt.Sprintf("{\"selector\":{\"Doctype\":\"user\",\"Email\":\"%s\",\"PasswordHash\":\"%s\"}}", email, passwordHash)

	jsonData, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	value := string(jsonData)

	// Take substring of first word with runes.
	// ... This handles any kind of rune in the string.
	runes := []rune(value)
	// ... Convert back into a string from rune slice.
	safeSubstring := string(runes[1 : len(runes)-1])

	fmt.Println(safeSubstring)

	jsonData = []byte(safeSubstring)

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(jsonData)))
	err = dec.Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
	}
	jq := jsonq.NewQuery(data)

	//[{"Key":"YvSgD5xAV0", "Record":{"Doctype":"user","Email":"tanmoykrishnadas@gmail.com","Key":"YvSgD5xAV0","Name":"Tanmoy Krishna Das","PasswordHash":"ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f","Token":"Bd56ti2SMt"}}]

	token, err := jq.String("Record", "Token")
	key, err := jq.String("Record", "Key")
	if err != nil {
		fmt.Println(err.Error())
	}

	return shim.Success([]byte(token + " " + key))
}
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}
func (s *SmartContract) createDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	DocId = "674567"
	DocName = args[0]
	OwnerId = args[1]
	EditorId = args[2]

	doc := Docs{DocId, DocName, OwnerId, EditorId}

	docJSON, err := json.Marshal(doc)
	err = APIstub.PutState(DocName, docJSON)
	if err != nil {
		return shim.Error("Document Creating failed: " + err.Error())
	}

	ret := "Document creating successful ;-) "

	return shim.Success([]byte(ret))
}
func (s *SmartContract) createSuccess(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docName := args[0]

	objectData, err := APIstub.GetState(docName)
	if err != nil {
		return shim.Error("getting create docData is failed: " + err.Error())
	}

	return shim.Success(objectData)
}

/*#In Query.js
    fcn: 'accessDoc',
	args: ['firstDoc', '674567','Nammi31'],
*/

func (s *SmartContract) accessDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	DocName1 := args[0]
	DocId1 := args[1]
	EditorId1 := args[2]

	var accessRight bool = true
	if DocName1 != DocName {
		accessRight = false
	}
	if DocId1 != DocId {
		accessRight = false
	}
	if EditorId1 != EditorId {
		accessRight = false
	}
	var ret string
	if accessRight == true {
		ret = "Document accessing is successful. "
	}

	return shim.Success([]byte(ret))
}
func (s *SmartContract) getObject(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	objectData, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error("getting object failed: " + err.Error())
	}

	var person1 Person
	err = json.Unmarshal(objectData, &person1)
	if err != nil {
		return shim.Error("UnMarshal failed: " + err.Error())
	}

	//print(person1.FirstName + " " + person1.LastName)

	return shim.Success(objectData)
}

func (s *SmartContract) putObject(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	key := args[0]
	firstName := args[1]
	lastName := args[2]
	age := args[3]

	ageInt, err := strconv.Atoi(age)
	if err != nil {
		return shim.Error("Interger conversion failed: " + err.Error())
	}

	p := Person{firstName, lastName, ageInt}

	pJson, err := json.Marshal(p)

	err = APIstub.PutState(key, pJson)
	if err != nil {
		return shim.Error("Putting object failed: " + err.Error())
	}

	ret := "putting object successful ;-) "

	return shim.Success([]byte(ret))
}

func (s *SmartContract) getData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	data, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error("There was an error")
	}

	return shim.Success(data)
}

func (s *SmartContract) setData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	key := args[0]
	val := args[1]

	err := APIstub.PutState(key, []byte(val))
	if err != nil {
		return shim.Error("There was an error")
	}

	str := "operation successful"

	return shim.Success([]byte(str))
}

func (s *SmartContract) tellMyName(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	name := "TKD"
	return shim.Success([]byte(name))
}

func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(carAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Car{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Car{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Car{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Car{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Car{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}

	i := 0
	for i < len(cars) {
		fmt.Println("i is ", i)
		carAsBytes, _ := json.Marshal(cars[i])
		APIstub.PutState("CAR"+strconv.Itoa(i), carAsBytes)
		fmt.Println("Added", cars[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var car = Car{Make: args[1], Model: args[2], Colour: args[3], Owner: args[4]}

	carAsBytes, _ := json.Marshal(car)
	APIstub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllCars(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CAR0"
	endKey := "CAR999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeCarOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	car := Car{}

	json.Unmarshal(carAsBytes, &car)
	car.Owner = args[1]

	carAsBytes, _ = json.Marshal(car)
	APIstub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
