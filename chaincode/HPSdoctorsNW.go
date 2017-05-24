
package main

import (
	"errors"
	"fmt"
	"bytes"
	"encoding/gob"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

type DoctorsNWChainCode struct {
}

type Doctor struct{

	NPI_ID string `json:"NPI_ID"`
	DoctorName string `json:"DoctorName"`
	MedicalCouncilName string `json:"MedicalCouncilName"`
	MedicalCouncilRegNumber string `json:"MedicalCouncilRegNumber"`
	LicenseID string `json:"LicenseID"`
	ExpiryDate string `json:"ExpiryDate"`
	LicenseStatus string `json:"LicenseStatus"`
	Hospital string `json:"Hospital"`
	Speciality string `json:"Speciality"`
	Area string `json:"Area"`
}

type SearchList struct{
	NPI_ID string `json:"NPI_ID"`
	SearchKeyWord string `json:"SearchKeyWord"`
}

type DocSearchList struct{
	DocList []Doctor `json:"DocList"`
}

func (self *DoctorsNWChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("In Init start ")

	var NPI_ID, DoctorName, MedicalCouncilName, MedicalCouncilRegNumber, LicenseID, ExpiryDate, LicenseStatus, Hospital, Speciality, Area string

	DoctorName = `John Doe`
	MedicalCouncilName = `Indian Medial Council`
	MedicalCouncilRegNumber =  `007`
	LicenseID = `LICID_1234`
	ExpiryDate = `2017/05/05`
	LicenseStatus =`expired`
	Hospital = `Columbia Asia`
	Speciality = `Cardiologist`
	Area = `SanFranscisco`

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting NPI_ID")
	}

	NPI_ID = args[0]

	res := &Doctor{}
	res.NPI_ID = NPI_ID
	res.DoctorName = DoctorName
	res.MedicalCouncilName = MedicalCouncilName
	res.MedicalCouncilRegNumber = MedicalCouncilRegNumber
	res.LicenseID = LicenseID
	res.ExpiryDate = ExpiryDate
	res.LicenseStatus = LicenseStatus
	res.Hospital = Hospital
	res.Speciality = Speciality
	res.Area = Area

	body, err := json.Marshal(res)
	if err != nil {
        panic(err)
    }
    fmt.Println(string(body))
	
	
	
	if function == "InitializeUser" {
		userBytes, err := AddDoctor(string(body),stub)
		if err != nil {
			fmt.Println("Error receiving  the User Details")
			return nil, err
		}
		fmt.Println("Initialization of User complete")
		
		return userBytes, nil
	}
	fmt.Println("Initialization No functions found ")
	return nil, nil
}


func (self *DoctorsNWChainCode) Invoke(stub shim.ChaincodeStubInterface,
	function string, args []string) ([]byte, error) {
	fmt.Println("In Invoke with function  " + function)

	if function == "AddDoctor" {
		fmt.Println("invoking AddDoctor " + function)
		testBytes,err := AddDoctor(args[0],stub)
		if err != nil {
			fmt.Println("Error performing AddDoctor ")
			return nil, err
		}
		fmt.Println("Processed AddDoctor successfully. ")
		return testBytes, nil
	}

	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation: " + function)
}

func (self *DoctorsNWChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
	fmt.Println("In Query with function " + function)

	bytes, err:= QueryDetails(stub, function,args)
	if err != nil {
		fmt.Println("Error retrieving function  ")
		return nil, err
	}
	return bytes,nil

}

func QueryDetails(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "GetDoctorDetails" {
		fmt.Println("Invoking GetDoctorDetails " + function)
		var doctors Doctor
		doctors,err := GetDoctorDetails(args[0], stub)
		if err != nil {
			fmt.Println("Error receiving  the Doctor details")
			return nil, errors.New("Error receiving  Doctor details")
		}
		fmt.Println("All success, returning doctor details")
		return json.Marshal(doctors)
	}

	if function == "GetDocDetailsBySearchKey" {
		fmt.Println("Invoking GetDocDetailsBySearchKey " + function)
		
		sList,err := GetDocDetailsBySearchKey(args[0], stub)
		if err != nil {
			fmt.Println("Error receiving  the Speciality details")
			return nil, errors.New("Error receiving  Speciality details")
		}
		fmt.Println("All success, returning Speciality details")
		return json.Marshal(sList)
	}
    return nil, errors.New("Received unknown query function name")

}

func GetDoctorDetails(NPI_ID string, stub shim.ChaincodeStubInterface)(Doctor, error) {
	fmt.Println("In query.GetDoctorDetails start ")

	key := NPI_ID
	var doctors Doctor
	userBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving Doctors" , NPI_ID)
		return doctors, errors.New("Error retrieving Doctor Details" + NPI_ID)
	}
	err = json.Unmarshal(userBytes, &doctors)
	fmt.Printf("%q",userBytes)
	fmt.Println("\nDoctor   : " , doctors);
	fmt.Println("In query.GetDoctorDetails end ")
	return doctors, nil
}

func GetDocDetailsBySearchKey(DocSpec string, stub shim.ChaincodeStubInterface)([]Doctor, error) {
	fmt.Println("In query.GetDocDetailsBySearchKey start ")
	
	key := DocSpec
	var doctors[] Doctor;

	sList, _ := GetDocList(key, stub)
	length := len(sList)
	
	for i := 0; i < length; i++ {
    doctorVal,err := GetDoctorDetails(sList[i], stub)
	if err != nil {
			fmt.Println("Error receiving  the Doctor details")
			return doctors, errors.New("Error receiving  Doctor details")
		}
			
			//Creating a list of doctors.
			doctors = append(doctors, doctorVal)
			fmt.Println(doctorVal)
			fmt.Println(doctors)
	
	}
	

	fmt.Println("In query.GetDocDetailsBySearchKey end ")
	fmt.Println("Printing doctors outside loop")
	fmt.Println(doctors)
	return doctors, nil
}

func GetDocList(DocSpec string, stub shim.ChaincodeStubInterface)([]string, error) {
	fmt.Println("In query.GetDocList start ")
	key := DocSpec
	userBytes, err := stub.GetState(key)
	sList := []string{}
	sListBytes := bytes.NewBuffer(userBytes)
	gob.NewDecoder(sListBytes).Decode(&sList)
	fmt.Println("sList after conversion to String")
	fmt.Println(sList)
	if err != nil {
		fmt.Println("Error retrieving Speciality" , DocSpec)
		return sList, errors.New("Error retrieving speciality Details" + DocSpec)
	}
    
	fmt.Println("Speciality   : " , userBytes);
	fmt.Println("In query.GetDocList end ")
	return sList, nil
}

func AddDocToSearchList(DocSpec SearchList, stub shim.ChaincodeStubInterface)([]byte, error) {
	// Checking if NPI_ID already present.
	var flag int = 0
	s, err := GetDocList(DocSpec.SearchKeyWord, stub)
	length :=len(s)
	
	for i := 0; i < length; i++ {
		if s[i] == DocSpec.NPI_ID {
			flag = 1
		}
	}

	// Adding Doctor's NPI_ID to SearchKeyWord List
	if flag != 1 {
		s = append(s, DocSpec.NPI_ID)   // appending NPI_ID to existing list only if not present.
		
		// convert from []string to []byte to put into ledger
		buf := &bytes.Buffer{}
		gob.NewEncoder(buf).Encode(s)
		bs := buf.Bytes()
		
		
		fmt.Println("Here is the String array in Byte format-->")
		fmt.Printf("%q", bs)
		fmt.Println("Adding ", DocSpec.NPI_ID," to SearchKeyWord ", DocSpec.SearchKeyWord); 
		err = stub.PutState(DocSpec.SearchKeyWord, bs)
		
		if err != nil {
			fmt.Println("Failed to add Doctor to SearchKeyWord ")
		}
	}
	return nil, nil
}

func AddDoctor(userJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.AddDoctor start ")
	//var s []string
	var SL SearchList
	res := &Doctor{}

	 
	err := json.Unmarshal([]byte(userJSON), res)
	if err != nil {
		fmt.Println("Failed to unmarshal user ")
	}
	fmt.Println("NPI_ID : ",res.NPI_ID)


	body, err := json.Marshal(res)
	if err != nil {
        panic(err)
    }
    fmt.Println(string(body))
	err = stub.PutState(res.NPI_ID, []byte(string(body)))
	if err != nil {
		fmt.Println("Failed to create Doctor ")
	}
	
	if err != nil {
		fmt.Println("Failed to create Doctor ")
	}

	fmt.Println("Created Doctor with Key : "+ res.NPI_ID)
	
	// Setting SearchList for Speciality search
	SL.NPI_ID = res.NPI_ID
	SL.SearchKeyWord = res.Speciality
	
	testBytes,err1 := AddDocToSearchList(SL, stub)
	if err1 != nil {
		fmt.Println("Failed to add Doctor to Speciality search ")
	}	
	
	// Setting SearchList for Area search
	SL.NPI_ID = res.NPI_ID
	SL.SearchKeyWord = res.Area
	
	testBytes,err = AddDocToSearchList(SL, stub)
	
	fmt.Printf("%q", testBytes)
	
	if err != nil {
		fmt.Println("Failed to add Doctor to Area search ")
	}
	

	fmt.Println("In initialize.AddDoctor end ")	
	return nil,nil

}


func main() {
	err := shim.Start(new(DoctorsNWChainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}


}
