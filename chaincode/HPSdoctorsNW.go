
package main

import (
	"errors"
	"fmt"
	// "strconv"
	// "time"
	// "strings"
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
}

type SearchList struct {
    DocNPI_IDS   []string            `json:"DocNPI_IDS"`
}

func (self *DoctorsNWChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("In Init start ")

	var NPI_ID, DoctorName, MedicalCouncilName, MedicalCouncilRegNumber, LicenseID, ExpiryDate, LicenseStatus, Hospital, Speciality string

	DoctorName = `John Doe`
	MedicalCouncilName = `Indian Medial Council`
	MedicalCouncilRegNumber =  `007`
	LicenseID = `LICID_1234`
	ExpiryDate = `2017/05/05`
	LicenseStatus =`expired`
	Hospital = `Columbia Asia`
	Speciality = `Cardiologist`

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
	//bytes, err:= query.Query(stub, function,args)
	//if err != nil {
		//fmt.Println("Error retrieving function  ")
		//return nil, err
	//}

	bytes, err:= QueryDetails(stub, function,args)
	if err != nil {
		fmt.Println("Error retrieving function  ")
		return nil, err
	}
	return bytes,nil

}

func QueryDetails(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "GetDoctors" {
		fmt.Println("Invoking GetDoctors " + function)
		var doctors Doctor
		doctors,err := GetDoctors(args[0], stub)
		if err != nil {
			fmt.Println("Error receiving  the Doctor details")
			return nil, errors.New("Error receiving  Doctor details")
		}
		fmt.Println("All success, returning doctor details")
		return json.Marshal(doctors)
	}

	return nil, errors.New("Received unknown query function name")

}

func GetDoctors(NPI_ID string, stub shim.ChaincodeStubInterface)(Doctor, error) {
	fmt.Println("In query.GetDoctors start ")

	key := NPI_ID
	var doctors Doctor
	userBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving Doctors" , NPI_ID)
		return doctors, errors.New("Error retrieving Doctor Details" + NPI_ID)
	}
	err = json.Unmarshal(userBytes, &doctors)
	fmt.Println("Doctor   : " , doctors);
	fmt.Println("In query.GetDoctors end ")
	return doctors, nil
}


func AddDoctor(userJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.AddDoctor start ")

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

	fmt.Println("Created Docter with Key : "+ res.NPI_ID)
	
	fmt.Println("Trying to Add Doctor to a Speciality")
	err = stub.PutState(res.Speciality, []byte(res.NPI_ID))
	if err != nil {
		fmt.Println("Failed to add Doctor to Speciality ")
	}
	fmt.Println("Added Doc to Speciality")
	fmt.Println("In initialize.AddDoctor end ")
	
	
	
	return nil,nil

}



func main() {
	err := shim.Start(new(DoctorsNWChainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}


}
