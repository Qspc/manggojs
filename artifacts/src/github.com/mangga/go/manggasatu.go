package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"log"
	"math"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ManggaContract struct {
	contractapi.Contract
}

type UserContract struct {
	contractapi.Contract
}

// User struct
type User struct {
	ID   string `json:"id"`
	NoHP string `json:"noHP"`
	Email string `json:"email"`
	NamaLengkap 	string `json:"namaLengkap"`

	Username	string `json:"username"`
	Password 	string `json:"password"`

	TanggalLahir  	string	`json:"tanggalLahir"`
	NIK 	int64 	`json:"nik"`
	Role 	int64 	`json:"role"`
	Alamat 	int64	`json:"alamat"`

}

// Mangga struct
type Mangga struct {
	ID 				string `json:"id"` // for query
	BenihID 	string `json:"benihID"`
	ManggaID 	string `json:"manggaID"`

	NamaPengirim	string `json:"namaPengirim"`
	NamaPenerima	string `json:"namaPenerima"`

	KuantitasBenihKg 	float64 `json:"kuantitasBenihKg"`
	HargaBenihPerKg  	float64 `json:"hargaBenihPerKg"`
	HargaBenihTotal		float64	`json:"hargaBenihTotal"`

	KuantitasManggaKg 	float64 `json:"kuantitasManggaKg"`
	HargaManggaPerKg  	float64 `json:"hargaManggaPerKg"`
	HargaManggaTotal	float64	`json:"hargaManggaTotal"`

	TanggalTransaksi 	int64	`json:"tanggaltransaksi"` // createdat with unix

	// Unique Value
	// From Penangkar
	VarietasBenih   string `json:"varietasBenih"`
	UmurBenih       string `json:"umurBenih"`
	// UmurPanen       string `json:"umurPanen"`

	// From Petani
	Pupuk         	string 	`json:"pupuk"`
	TanggalTanam	string `json:"tanggalTanam"`
	LokasiLahan		string `json:"lokasiLahan"`
	
	Ukuran    	string 	`json:"ukuran"`
	Pestisida     	string 	`json:"pestisida"`
	KadarAir 	float64 `json:"kadarAir"`
	Perlakuan     	string 	`json:"perlakuan"`
	Produktivitas 	string 	`json:"produktivitas"`
	TanggalPanen 	int64 	`json:"tanggalPanen"`

	// From Pengumpul
	TanggalMasuk     int64 `json:"tanggalMasuk"`
	TeknikSorting    string `json:"teknikSorting"`
	MetodePengemasan string `json:"metodePengemasan"`
	Pengangkutan string `json:"pengangkutan"`

	//pedagang
	Pembeli		string `json:"pembeli"`

	CaraPembayaran string `json:"caraPembayaran"`

	TxID1 string `json:"txID1"` // penangkar - petani
	TxID2 string `json:"txID2"` // petani - pengumpul
	TxID3 string `json:"txID3"` // pengumpul - pedagang
	TxID4 string `json:"txID4"` // pedagang besar - konsumen

	IsAsset 	bool `json:"isAsset"`
	IsConfirmed bool `json:"isConfirmed"`
	IsEmpty		bool `json:"isEmpty"`
	IsRejected 	bool `json:"isRejected"`

	RejectReason	string	`json:"rejectReason"`
}

// tahap 1 = proses input/registrasi benih (penangkar) 
// fungsi invoke
func (s *ManggaContract) RegistrasiBenih(ctx contractapi.TransactionContextInterface, manggaData string) (string, error) {
	if len(manggaData) == 0 {
		return "", fmt.Errorf("Please pass the correct benih data")
	}

	// data yg dibawa : manggaData (data dari fe)
	// data yg dikirim : varietas, kuantitas, umur benih, ;createdat, ;ID, ;benihID
	mangga := new(Mangga)

	// set createdAt
	mangga.CreatedAt = time.Now().Unix()
	// Set ID (key)
	mangga.ID = ctx.GetStub().GetTxID()
	// Set BenihID
	mangga.BenihID = ctx.GetStub().GetTxID()

	// set isAsset
	mangga.IsAsset = true
	mangga.IsConfirmed = false
	mangga.IsRejected = false

	// insert varietas, kuantitas, umur benih dari manggaData ke mangga
	err := json.Unmarshal([]byte(manggaData), &mangga)

	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling benih. %s", err.Error())
	}

	manggaAsBytes, err := json.Marshal(mangga)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling benih. %s", err.Error())
	}

	// kirim pesan "createdasset" ke peer selanjutnya
	ctx.GetStub().SetEvent("CreateAsset", manggaAsBytes)

	// Put state using key and data
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(mangga.ID, manggaAsBytes)
}

// tahap 2 = proses update/menambah kuantitas dari benih (penangkat)
// fungsi invoke
func (s *ManggaContract) AddKuantitasBenihByID(ctx contractapi.TransactionContextInterface, quantity float64, manggaID string) (*Mangga, error) {
	if len(manggaID) == 0 {
		return nil, fmt.Errorf("Please pass the correct mangga id")
	}

	// data yang dibawa : quantity (data update benih), manggaID (ID mangga yg bakal diupdate)
	// data yang dikirim : update benih

	// get json mangga berdasarkan ID yg dituju
	manggaAsBytes, err := ctx.GetStub().GetState(manggaID)

	if err != nil {
		return nil, fmt.Errorf("Failed to get mangga data. %s", err.Error())
	}

	if manggaAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", manggaID)
	}

	// Create new mangga object
	mangga := new(Mangga)
	// masukan data di manggaAsBytes kedalam mangga
	_ = json.Unmarshal(manggaAsBytes, mangga)

	// update benih quantity
	mangga.KuantitasBenihKg += quantity
	mangga.KuantitasBenihKg = math.Round(mangga.KuantitasBenihKg*100)/100

	if mangga.KuantitasBenihKg > 0 {
		mangga.IsEmpty = false
	}

	manggaAsBytes, err = json.Marshal(mangga)

	return mangga, ctx.GetStub().PutState(mangga.ID, manggaAsBytes)
}

// tahap 3 = proses transaksi antara penangkar dan petani (penangkar)
// Create tx1 with the unique from the benih
// fungsi invoke
func (s *ManggaContract) CreateTrxManggaByPenangkar(ctx contractapi.TransactionContextInterface, manggaData, prevID string) (string, error) {

	if len(manggaData) == 0 {
		return "", fmt.Errorf("Please pass the correct benih data")
	}

	if len(prevID) == 0 {
		return "", fmt.Errorf("Please pass the correct benih aset id")
	}

	//data yang dibawah : manggaData (inputan dari fe), PrevID (data dari block sebelumnya)
	//data yang dikirim : ;nama pengirim, ;kontak pengirim, ;alamat pengirim, ;nama penerima, ;kontak penerima, ;alamat penerima, ;kuantitas benih, ;varietas benih, ;umur benih, ;harga benih (kg), ;harga benih total, ;cara pembayaran, ;createdAt, ;ID, ;BenihID, ;TxID1

	manggaPrev, err := s.GetManggaByID(ctx, prevID) 

	if err != nil {
		return "", fmt.Errorf("Failed while getting benih aset. %s", err.Error())
	}

	manggaNew := new(Mangga)

	// tanggal transaksi
	manggaNew.CreatedAt = time.Now().Unix()

	// Set ID (key)
	manggaNew.ID = ctx.GetStub().GetTxID()
	// Set TxID
	manggaNew.TxID1 = ctx.GetStub().GetTxID()

	// Set asetIDs
	manggaNew.BenihID = manggaPrev.BenihID

	manggaNew.IsAsset = false
	manggaNew.IsConfirmed = false

	// Get penangkar unique value from benih aset
	manggaNew.UmurBenih = manggaPrev.UmurBenih
	manggaNew.VarietasBenih = manggaPrev.VarietasBenih
	manggaNew.KuantitasBenihKg = manggaPrev.KuantitasBenihKg
	// manggaNew.HargaBenihPerKg = manggaPrev.HargaBenihPerKg

	// insert data dari fe : nama pengirim,  nama penerima, harga benih (kg), cara pembayaran. lalu masukan data dari manggaData ke manggaNew
	err = json.Unmarshal([]byte(manggaData), &manggaNew)

	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling benih. %s", err.Error())
	}

	// // check identity
	// alamatPengirim, err := s.checkUserbyUsername(ctx, bawangNew.UsernamePengirim)

	// if err != nil {
	// 	return "", fmt.Errorf("Failed while checking username pengirim. %s", err.Error())
	// }

	// bawangNew.AlamatPengirim = alamatPengirim

	// alamatPenerima, err := s.checkUserbyUsername(ctx, bawangNew.UsernamePenerima)

	// if err != nil {
	// 	return "", fmt.Errorf("Failed while checking username penerima. %s", err.Error())
	// }

	// bawangNew.AlamatPenerima = alamatPenerima

	// calculate benih total price
	manggaNew.HargaBenihTotal = manggaNew.KuantitasBenihKg * manggaNew.HargaBenihPerKg

	// check validasi apakah benih yang dimiliki lebih banyak/sama dengan benih yg dikirim
	if manggaPrev.KuantitasBenihKg - manggaNew.KuantitasBenihKg >= 0 {

		_, err = s.updateKuantitasBenihByID(ctx, prevID, manggaNew.KuantitasBenihKg)

		if err != nil {
			return "", fmt.Errorf("Failed while updating benih aset kuantitas. %s", err.Error())
		}

		manggaAsBytes, err := json.Marshal(manggaNew)

		if err != nil {
			return "", fmt.Errorf("Failed while marshling benih. %s", err.Error())
		}
	
		ctx.GetStub().SetEvent("CreateAsset", manggaAsBytes)
	
		// Put state using key and data
		return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(manggaNew.ID, manggaAsBytes)
	} else {
		return "", fmt.Errorf("Benih tidak mencukupi. %s", err.Error())
	}
}

// tahap 4 = proses menanam benih (petani)
// fungsi invoke
func (s *ManggaContract) TanamBenih(ctx contractapi.TransactionContextInterface, manggaData, prevID string) (string, error) {
	
	if len(manggaData) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga data")
	}

	if len(prevID) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga transaction id")
	}

	// data yang dibawa : manggaData (data dari fe), PrevID (data dari blockchain sebelumnya)
	// data yang dikirim : ;nama petani, ;kuantitas benih, ;tanggal tanam, ;lokasi lahan, ;pupuk, ;ID, ;ManggaID, ;CreatedAt, dan "data sebelumnya" 

	//proses unmarshal data blockchain sebelumnya untuk dipakai sekarang
	manggaPrev, err := s.GetManggaByID(ctx, prevID) 

	if err != nil {
		return "", fmt.Errorf("Failed while getting mangga init. %s", err.Error())
	}

	// set the benih quantity that petani got from penangkar to zero
	// petani plant all of its benih
	// make easier to query which benih that has been planted
	_, err = s.updateKuantitasBenihByID(ctx, prevID, manggaPrev.KuantitasBenihKg)

	manggaNew := new(Mangga)

	manggaNew.CreatedAt = time.Now().Unix()
	
	// Set ID (key)
	manggaNew.ID = ctx.GetStub().GetTxID()

	// Set asetID
	manggaNew.ManggaID = ctx.GetStub().GetTxID()
	manggaNew.BenihID = manggaPrev.BenihID

	manggaNew.TxID1 = manggaPrev.TxID1

	manggaNew.IsAsset = true
	manggaNew.IsConfirmed = false

	// Get penangkar unique field from manggaPrev
	manggaNew.UmurBenih = manggaPrev.UmurBenih
	manggaNew.VarietasBenih = manggaPrev.VarietasBenih
	manggaNew.KuantitasBenihKg = manggaPrev.KuantitasBenihKg
	manggaNew.HargaBenihPerKg = manggaPrev.HargaBenihPerKg
	manggaNew.HargaBenihTotal = manggaPrev.HargaBenihTotal

	// nama petani
	manggaNew.NamaPengirim = manggaPrev.NamaPenerima

	// Insert Tanggal Tanam
	manggaNew.TanggalTanam = time.Now().Unix()

	// menambah data dari fe : pupuk, lokasi lahan
	err = json.Unmarshal([]byte(manggaData), &manggaNew)

	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling mangga. %s", err.Error())
	}

	manggaAsBytes, err := json.Marshal(manggaNew)
	
	if err != nil {
		return "", fmt.Errorf("Failed while marshling mangga. %s", err.Error())
	}

	ctx.GetStub().SetEvent("CreateAsset", manggaAsBytes)

	// Put state using key and data
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(manggaNew.ID, manggaAsBytes)
}

// tahap 5 = proses panen mangga (petani)
// Add other mangga field
// fungsi invoke
func (s *ManggaContract) PanenMangga(ctx contractapi.TransactionContextInterface, manggaData, manggaID string) (string, error) {
	
	if len(manggaData) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga data")
	}

	if len(manggaID) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga id")
	}

	// data yang dibawa : manggaData (data dari fe), manggaID (id mangga yg diproses)
	// data yg dikirim : ;tanggal panen, ;kuantitas mangga, ;ukuran, ;pestisida, ;kadar air, ;perlakuan, ;produktivitas, createdAt

	manggaAsBytes, err := ctx.GetStub().GetState(manggaID)

	if err != nil {
		return "", fmt.Errorf("Failed to get mangga data. %s", err.Error())
	}

	if manggaAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", manggaID)
	}

	// Create new mangga object
	mangga := new(Mangga)
	_ = json.Unmarshal(manggaAsBytes, mangga)

	// insert TanggalPanen
	mangga.TanggalPanen = time.Now().Unix()
	mangga.CreatedAt = time.Now().Unix()

	// data dari fe yaitu : kuantitas mangga, ukuran, pestisida, kadar air, perlakuan, produktivitas
	err = json.Unmarshal([]byte(manggaData), &mangga)

	if err != nil {
		return "", fmt.Errorf("Failed while marshling mangga. %s", err.Error())
	}

	manggaAsBytes, err = json.Marshal(mangga)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling mangga. %s", err.Error())
	}

	// return the previous id and update the state
	return manggaID, ctx.GetStub().PutState(mangga.ID, manggaAsBytes)
}

// tahap 6 = proses transaksi antara petani dan pengumpul (petani) 
// adding unique value from petani and txid2
// fungsi invoke
func (s *ManggaContract) CreateTrxManggaByPetani(ctx contractapi.TransactionContextInterface, manggaData, prevID string) (string, error) {

	if len(manggaData) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga data")
	}

	if len(prevID) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga aset id")
	}

	// data yang dibawa : manggadata (data dari fe), prevID (data dari blockchain sebelumnya)
	// data yang dikirim : ;nama pengirim, ;kontak pengirim, ;alamat pengirim, ;nama penerima, ;kontak penerima, ;alamat penerima, ;harga mangga (kg), harga mangga total, ;tanggal transaksi, cara pembayaran, ;createdAt, ;ID, ;TxId2,  dan "data sebelumnya"

	manggaPrev, err := s.GetManggaByID(ctx, prevID) 

	if err != nil {
		return "", fmt.Errorf("Failed while getting mangga aset. %s", err.Error())
	}

	manggaNew := new(Mangga)

	manggaNew.CreatedAt = time.Now().Unix()
	manggaNew.TanggalTransaksi = time.Now().Unix()

	// Set ID (key)
	manggaNew.ID = ctx.GetStub().GetTxID()
	// Set TxID
	manggaNew.TxID1 = manggaPrev.TxID1
	manggaNew.TxID2 = ctx.GetStub().GetTxID()

	// Set asetID
	manggaNew.ManggaID = manggaPrev.ManggaID
	manggaNew.BenihID = manggaPrev.BenihID

	manggaNew.IsAsset = false
	manggaNew.IsConfirmed = false

	// Get penangkar unique field from mangga aset
	manggaNew.UmurBenih = manggaPrev.UmurBenih
	manggaNew.VarietasBenih = manggaPrev.VarietasBenih
	manggaNew.KuantitasBenihKg = manggaPrev.KuantitasBenihKg
	manggaNew.HargaBenihPerKg = manggaPrev.HargaBenihPerKg
	manggaNew.HargaBenihTotal = manggaPrev.HargaBenihTotal

	// Get petani unique field from mangga aset
	manggaNew.Pupuk = manggaPrev.Pupuk
	manggaNew.TanggalTanam = manggaPrev.TanggalTanam
	manggaNew.LokasiLahan = manggaPrev.LokasiLahan

	manggaNew.KadarAir = manggaPrev.KadarAir
	manggaNew.Ukuran = manggaPrev.Ukuran
	manggaNew.Pestisida = manggaPrev.Pestisida
	manggaNew.Perlakuan = manggaPrev.Perlakuan
	manggaNew.Produktivitas = manggaPrev.Produktivitas
	manggaNew.TanggalPanen = manggaPrev.TanggalPanen
	manggaNew.KuantitasManggaKg = manggaPrev.KuantitasManggaKg

	// data tambahan dari fe : nama pengirim, nama penerima, harga mangga (kg), cara pembayaran
	err = json.Unmarshal([]byte(manggaData), &manggaNew)

	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling mangga. %s", err.Error())
	}

	// calculate mangga total price
	manggaNew.HargaManggaTotal = manggaNew.KuantitasManggaKg * manggaNew.HargaManggaPerKg

	// check apakah mangga yang tersedia lebih banyak/sama dengan mangga yang akan dikirim
	if manggaPrev.KuantitasManggaKg - manggaNew.KuantitasManggaKg >= 0 {

		//update mangga yg tersedia dikurangi mangga yang dikirim
		_, err = s.updateKuantitasManggaByID(ctx, prevID, manggaNew.KuantitasManggaKg)

		if err != nil {
			return "", fmt.Errorf("Failed while updating mangga aset kuantitas. %s", err.Error())
		}

		manggaAsBytes, err := json.Marshal(manggaNew)

		if err != nil {
			return "", fmt.Errorf("Failed while marshling mangga. %s", err.Error())
		}

		ctx.GetStub().SetEvent("CreateAsset", manggaAsBytes)
	
		// Put state using key and data
		return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(manggaNew.ID, manggaAsBytes)
	} else {
		return "", fmt.Errorf("mangga tidak mencukupi. %s", err.Error())
	}
}

// tahap 7 = proses transaksi antara pengumpul dan pedagang (pengumpul) 
// adding unique value from pengumpul and txid3
// fungsi invoke
func (s *ManggaContract) CreateTrxManggaByPengumpul(ctx contractapi.TransactionContextInterface, manggaData, prevID string) (string, error) {

	if len(manggaData) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga data")
	}

	if len(prevID) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga aset id")
	}

	// data yang dibawa : manggaData (data dari fe), prevId (data dari block sebelumnya)
	// data yang dikirim : ;nama pengirim, ;kontak pengirim, ;alamat pengirim, ;nama penerima, ;kontak penerima, ;alamat penerima, ;kuantitas mangga, ;harga mangga (kg), harga mangga total, ;tanggal transaksi, ;teknik sorting, ;metode pengemasan, ;cara pengangkutan, ;cara pembayaran, ;createdAt, ;ID, ;TxID3

	manggaPrev, err := s.GetManggaByID(ctx, prevID) 

	if err != nil {
		return "", fmt.Errorf("Failed while getting mangga aset. %s", err.Error())
	}

	manggaNew := new(Mangga)

	manggaNew.CreatedAt = time.Now().Unix()
	manggaNew.TanggalTransaksi = time.Now().Unix()


	// Set id (key)
	manggaNew.ID = ctx.GetStub().GetTxID()
	// Set txid1
	manggaNew.TxID1 = manggaPrev.TxID1
	// Set txid2
	manggaNew.TxID2 = manggaPrev.TxID2
	// Set txid3
	manggaNew.TxID3 = ctx.GetStub().GetTxID()

	// Set asetID
	manggaNew.ManggaID = manggaPrev.ManggaID
	manggaNew.BenihID = manggaPrev.BenihID

	manggaNew.IsAsset = false
	manggaNew.IsConfirmed = false

	// Get penangkar unique field from mangga aset
	manggaNew.UmurBenih = manggaPrev.UmurBenih
	manggaNew.VarietasBenih = manggaPrev.VarietasBenih
	manggaNew.KuantitasBenihKg = manggaPrev.KuantitasBenihKg
	manggaNew.HargaBenihPerKg = manggaPrev.HargaBenihPerKg
	manggaNew.HargaBenihTotal = manggaPrev.HargaBenihTotal

	// Get petani unique field from mangga aset
	manggaNew.Pupuk = manggaPrev.Pupuk
	manggaNew.TanggalTanam = manggaPrev.TanggalTanam
	manggaNew.LokasiLahan = manggaPrev.LokasiLahan

	manggaNew.KadarAir = manggaPrev.KadarAir
	manggaNew.Ukuran = manggaPrev.Ukuran
	manggaNew.Pestisida = manggaPrev.Pestisida
	manggaNew.Perlakuan = manggaPrev.Perlakuan
	manggaNew.Produktivitas = manggaPrev.Produktivitas
	manggaNew.TanggalPanen = manggaPrev.TanggalPanen

	// tambah data dari fe nama pengirim, nama penerima, kuantitas mangga, harga mangga (kg), teknik sorting, metode pengemasan, cara pengangkutan, cara pembayaran
	err = json.Unmarshal([]byte(manggaData), &manggaNew)

	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling benih. %s", err.Error())
	}

	// calculate mangga total price
	manggaNew.HargaManggaTotal = manggaNew.KuantitasManggaKg * manggaNew.HargaManggaPerKg

	// check quantity
	if manggaPrev.KuantitasManggaKg - manggaNew.KuantitasManggaKg >= 0 {

		_, err = s.updateKuantitasManggaByID(ctx, prevID, manggaNew.KuantitasManggaKg)

		if err != nil {
			return "", fmt.Errorf("Failed while updating mangga aset kuantitas. %s", err.Error())
		}

		manggaAsBytes, err := json.Marshal(manggaNew)

		if err != nil {
			return "", fmt.Errorf("Failed while marshling mangga. %s", err.Error())
		}
	
		// Insert into blockchain
		ctx.GetStub().SetEvent("CreateAsset", manggaAsBytes)
	
		// Put state using key and data
		return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(manggaNew.ID, manggaAsBytes)
	} else {
		return "", fmt.Errorf("mangga tidak mencukupi. %s", err.Error())
	}
}

// tahap 8 = proses transaksi antara pedagang dan konsumen (pedagang) 
// adding unique value from pedagang and txid4
// fungsi invoke
func (s *ManggaContract) CreateTrxManggaByPedagang(ctx contractapi.TransactionContextInterface, manggaData, prevID string) (string, error) {

	if len(manggaData) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga data")
	}

	if len(prevID) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga aset id")
	}

	// data yang dibawa : manggaData (data dari fe), prevId (data dari block sebelumnya)
	// data yang dikirim : ;nama pengirim, ;kontak pengirim, ;alamat pengirim, ;nama penerima, ;kuantitas mangga, ;harga mangga (kg), harga mangga total, ;teknik sorting, ;metode pengemasan, ;cara pengangkutan, ;cara pembayaran, ;createdAt, ;ID, ;TxID3

	manggaPrev, err := s.GetManggaByID(ctx, prevID) 

	if err != nil {
		return "", fmt.Errorf("Failed while getting mangga aset. %s", err.Error())
	}

	manggaNew := new(Mangga)

	manggaNew.CreatedAt = time.Now().Unix()

	// Set id (key)
	manggaNew.ID = ctx.GetStub().GetTxID()
	// Set txid1
	manggaNew.TxID1 = manggaPrev.TxID1
	// Set txid2
	manggaNew.TxID2 = manggaPrev.TxID2
	// Set txid3
	manggaNew.TxID3 = manggaPrev.TxID3
	// Set txid4
	manggaNew.TxID4 = ctx.GetStub().GetTxID()

	// Set asetID
	manggaNew.ManggaID = manggaPrev.ManggaID
	manggaNew.BenihID = manggaPrev.BenihID

	manggaNew.IsAsset = false
	manggaNew.IsConfirmed = false

	// Get penangkar unique field from mangga aset
	manggaNew.UmurBenih = manggaPrev.UmurBenih
	manggaNew.VarietasBenih = manggaPrev.VarietasBenih
	manggaNew.KuantitasBenihKg = manggaPrev.KuantitasBenihKg
	manggaNew.HargaBenihPerKg = manggaPrev.HargaBenihPerKg
	manggaNew.HargaBenihTotal = manggaPrev.HargaBenihTotal

	// Get petani unique field from mangga aset
	manggaNew.Pupuk = manggaPrev.Pupuk
	manggaNew.TanggalTanam = manggaPrev.TanggalTanam
	manggaNew.LokasiLahan = manggaPrev.LokasiLahan

	manggaNew.KadarAir = manggaPrev.KadarAir
	manggaNew.Ukuran = manggaPrev.Ukuran
	manggaNew.Pestisida = manggaPrev.Pestisida
	manggaNew.Perlakuan = manggaPrev.Perlakuan
	manggaNew.Produktivitas = manggaPrev.Produktivitas
	manggaNew.TanggalPanen = manggaPrev.TanggalPanen

	//get pengumpul unique field from mangga aset


	// tambah data dari fe nama pengirim, kontak pengirim, alamat pengirim, nama penerima, kuantitas mangga, harga mangga (kg), teknik sorting, metode pengemasan, cara pengangkutan, cara pembayaran
	err = json.Unmarshal([]byte(manggaData), &manggaNew)

	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling benih. %s", err.Error())
	}

	// calculate mangga total price
	manggaNew.HargaManggaTotal = manggaNew.KuantitasManggaKg * manggaNew.HargaManggaPerKg

	// check quantity
	if manggaPrev.KuantitasManggaKg - manggaNew.KuantitasManggaKg >= 0 {

		_, err = s.updateKuantitasManggaByID(ctx, prevID, manggaNew.KuantitasManggaKg)

		if err != nil {
			return "", fmt.Errorf("Failed while updating mangga aset kuantitas. %s", err.Error())
		}

		manggaAsBytes, err := json.Marshal(manggaNew)

		if err != nil {
			return "", fmt.Errorf("Failed while marshling mangga. %s", err.Error())
		}
	
		// Insert into blockchain
		ctx.GetStub().SetEvent("CreateAsset", manggaAsBytes)
	
		// Put state using key and data
		return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(manggaNew.ID, manggaAsBytes)
	} else {
		return "", fmt.Errorf("mangga tidak mencukupi. %s", err.Error())
	}
}

// tahap 9 = konfirmasi transaksi diterima (petani, pengumpul, pedagang)
// fungsi invoke
func (s *ManggaContract) ConfirmTrxByID(ctx contractapi.TransactionContextInterface, manggaID string) (string, error) {

	if len(manggaID) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga id")
	}

	// data yang dibawa : manggaID (id mangga yg bersangkutan)
	// data yg dikirim : isConfirmed = true

	manggaAsBytes, err := ctx.GetStub().GetState(manggaID)

	if err != nil {
		return "", fmt.Errorf("Failed to get mangga data. %s", err.Error())
	}

	if manggaAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", manggaID)
	}

	// Create new mangga object
	mangga := new(Mangga)
	_ = json.Unmarshal(manggaAsBytes, mangga)

	// Change the confirmed status
	mangga.IsConfirmed = true

	// Change createdAt
	mangga.CreatedAt = time.Now().Unix()

	// Set tanggal masuk when pengumpul confirmed the mangga trx
	if len(mangga.Ukuran) != 0 || len(mangga.Pupuk) != 0 || 
	mangga.TanggalTanam != 0 || mangga.TanggalPanen != 0 {
		mangga.TanggalMasuk = time.Now().Unix()
	}

	manggaAsBytes, err = json.Marshal(mangga)

	if err != nil {
		return "", fmt.Errorf("Failed while marshling mangga. %s", err.Error())
	}

	// Update the state
	return manggaID, ctx.GetStub().PutState(mangga.ID, manggaAsBytes)
}

// tahap 10 = konfirmasi transaksi ditolak (petani, pengumpul, pedagang)
// fungsi invoke
func (s *ManggaContract) RejectTrxByID(ctx contractapi.TransactionContextInterface, manggaIDPrev, manggaIDReject string, kuantitasPrev float64, rejectReason string) (string, error) {
	
	// Create new mangga object for previous mangga
	if len(manggaIDPrev) == 0 {
		return "", fmt.Errorf("Please pass the correct previous mangga id")
	}

	// data yg dibawa : manggaIDPrev (id mangga yg bersangkutan), manggaIDReject, kuantitasPrev (untuk balikin kuantitas mangga ke semula), rejectReason (alasan menolak)

	manggaPrevAsBytes, err := ctx.GetStub().GetState(manggaIDPrev)

	if err != nil {
		return "", fmt.Errorf("Failed to get mangga data. %s", err.Error())
	}

	if manggaPrevAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", manggaIDPrev)
	}
	
	manggaPrev := new(Mangga)
	_ = json.Unmarshal(manggaPrevAsBytes, manggaPrev)

	if len(manggaIDReject) == 0 {
		return "", fmt.Errorf("Please pass the correct rejected mangga id")
	}

	manggaRejectAsBytes, err := ctx.GetStub().GetState(manggaIDReject)

	if err != nil {
		return "", fmt.Errorf("Failed to get mangga data. %s", err.Error())
	}

	if manggaRejectAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", manggaIDReject)
	}

	manggaReject := new(Mangga)
	_ = json.Unmarshal(manggaRejectAsBytes, manggaReject)

	// Change the rejected status
	manggaReject.IsRejected = true

	if manggaPrev.BenihID != "" && manggaPrev.ManggaID == "" {
		manggaPrev.KuantitasBenihKg += kuantitasPrev
	} else {
		manggaPrev.KuantitasManggaKg += kuantitasPrev
	}

	manggaReject.RejectReason = rejectReason

	manggaPrevAsBytes, err = json.Marshal(manggaPrev)
	manggaRejectAsBytes, err = json.Marshal(manggaReject)

	if err != nil {
		return "", fmt.Errorf("Failed while marshling mangga. %s", err.Error())
	}

	// Update the state
	ctx.GetStub().PutState(manggaPrev.ID, manggaPrevAsBytes)

	return manggaIDReject, ctx.GetStub().PutState(manggaReject.ID, manggaRejectAsBytes)
}

// Add mangga asset quantity by ID
// gak dipake?
func (s *ManggaContract) AddManggaKuantitasByID(ctx contractapi.TransactionContextInterface, quantity float64, manggaID string) (*Mangga, error) {
	if len(manggaID) == 0 {
		return nil, fmt.Errorf("Please pass the correct mangga id")
	}

	manggaAsBytes, err := ctx.GetStub().GetState(manggaID)

	if err != nil {
		return nil, fmt.Errorf("Failed to get mangga data. %s", err.Error())
	}

	if manggaAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", manggaID)
	}

	// Create new mangga object
	mangga := new(Mangga)
	_ = json.Unmarshal(manggaAsBytes, mangga)

	// Add mangga quantity
	mangga.KuantitasManggaKg += quantity

	if mangga.KuantitasManggaKg > 0 {
		mangga.IsEmpty = false
	}

	manggaAsBytes, err = json.Marshal(mangga)

	return mangga, ctx.GetStub().PutState(mangga.ID, manggaAsBytes)
}


// CreateUser function to create user and insert it on blockchain
func (s *UserContract) CreateUser(ctx contractapi.TransactionContextInterface, userData string) (string, error) {

	if len(userData) == 0 {
		return "", fmt.Errorf("Please pass the correct mangga data")
	}

	var user User

	user.CreatedAt = time.Now().Unix()

	//create user ID
	user.ID = ctx.GetStub().GetTxID()

	err := json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling user. %s", err.Error())
	}

	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling mangga. %s", err.Error())
	}

	// Insert into blockchain
	ctx.GetStub().SetEvent("CreateAsset", userAsBytes)

	// Put state using key and data
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(user.ID, userAsBytes)
}

// Get mangga object by ID
// mengambil data mangga, lalu melakukan proses unmarshal agar dapat dipakai saat transaksi
func (s *ManggaContract) GetManggaByID(ctx contractapi.TransactionContextInterface, manggaID string) (*Mangga, error) {
	if len(manggaID) == 0 {
		return nil, fmt.Errorf("Please provide correct contract ID")
	}

	// data yang dibawa : manggaID (ID mangga yg dipakai)
	// data yang dikirim : data blockchain yg selesai di unmarshal

	manggaAsBytes, err := ctx.GetStub().GetState(manggaID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if manggagAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", manggaID)
	}

	// create new mangga object and unmarshal manggaAsBytes into mangga
	mangga := new(Mangga)
	_ = json.Unmarshal(manggaAsBytes, mangga)

	return mangga, nil
}

// fungsi dari : benih yang ada dikurangi benih yang dikirim pas transaksi 
// dan fungsi : benih yang diperoleh pas transaksi - benih yang ditanam ==> pasti kosong hasilnya
func (s *ManggaContract) updateKuantitasBenihByID(ctx contractapi.TransactionContextInterface, manggaID string, kuantitasBenihNext float64) (string, error) {
	if len(manggaID) == 0 {
		return "", fmt.Errorf("Please insert the correct benih aset id")
	}

	//data yg dibawa : manggaID (id mangga yg bakal dipakai), kuantitasBenihNext (banyaknya benih)

	manggaAsBytes, err := ctx.GetStub().GetState(manggaID)

	if err != nil {
		return "", fmt.Errorf("Failed to get benih data. %s", err.Error())
	}

	if manggaAsBytes == nil {
		return "", fmt.Errorf("Benih aset with %s id does not exist", manggaID)
	}

	// Create new mangga object
	mangga := new(Mangga)
	_ = json.Unmarshal(manggaAsBytes, mangga)

	mangga.KuantitasBenihKg -= kuantitasBenihNext
	mangga.KuantitasBenihKg = math.Round(mangga.KuantitasBenihKg*100)/100


	if mangga.KuantitasBenihKg == 0 {
		mangga.IsEmpty = true
	}

	manggaAsBytes, err = json.Marshal(mangga)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(mangga.ID, manggaAsBytes)
}

// fungsi dari : mangga yang ada di petani dikurangi mangga yg akan dikirim ke pengumpul
func (s *ManggaContract) updateKuantitasManggaByID(ctx contractapi.TransactionContextInterface, manggaID string, kuantitasManggaNext float64) (string, error) {
	if len(manggaID) == 0 {
		return "", fmt.Errorf("Please insert the correct mangga aset id")
	}

	//data yang dibawa manggaID (id mangga yg sesuai), kuantitasmanggaNext (banyak mangga yg dikirim)

	manggaAsBytes, err := ctx.GetStub().GetState(manggaID)

	if err != nil {
		return "", fmt.Errorf("Failed to get mangga data. %s", err.Error())
	}

	if manggaAsBytes == nil {
		return "", fmt.Errorf("mangga aset with %s id does not exist", manggaID)
	}

	// Create new mangga object
	mangga := new(Mangga)
	_ = json.Unmarshal(manggaAsBytes, mangga)

	mangga.KuantitasManggaKg -= kuantitasManggaNext
	mangga.KuantitasManggaKg = math.Round(mangga.KuantitasManggaKg*100)/100


	if mangga.KuantitasManggaKg == 0 {
		mangga.IsEmpty = true
	}

	manggaAsBytes, err = json.Marshal(mangga)

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(mangga.ID, manggaAsBytes)
}



// =========================================================================================================== //
// GetUserByID get mangga object by ID
// mengambil data user berdasarkan ID nya, lalu lakukan proses unmarshal
// fungsi query
func (s *UserContract) GetUserByID(ctx contractapi.TransactionContextInterface, userID string) (*User, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("Please provide correct contract ID")
		// return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// data yang dibawa userID

	// fetch mangga by its ID
	userAsBytes, err := ctx.GetStub().GetState(userID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", userID)
	}

	// create new mangga object and unmarshal manggaAsBytes into mangga
	user := new(User)
	_ = json.Unmarshal(userAsBytes, user)

	return user, nil
}

// it should be GetManggagForQuery
// fungsi query
func (s *ManggaContract) GetManggaForQuery(ctx contractapi.TransactionContextInterface, queryString string) (string, error) {

	
	queryResults, err := s.getQueryResultForQueryString(ctx, queryString)

	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	return queryResults, nil

}

// GetHistoryForAssetByID to get history for asset by ID
// fungsi query
func (s *ManggaContract) GetHistoryForAssetByID(ctx contractapi.TransactionContextInterface, manggaID string) (string, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(manggaID)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		// buffer.WriteString(", \"Timestamp\":")
		// buffer.WriteString("\"")
		// buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		// buffer.WriteString("\"")

		// buffer.WriteString(", \"IsDelete\":")
		// buffer.WriteString("\"")
		// buffer.WriteString(strconv.FormatBool(response.IsDelete))
		// buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return string(buffer.Bytes()), nil
}


// Check user identity and retrieve its address
// gak dipake?
func (s *ManggaContract) checkUserbyUsername(ctx contractapi.TransactionContextInterface, username string) (string, error) {
	if len(username) == 0 {
		return "", fmt.Errorf("Please insert the correct username")
	}

	queryString := `{"selector":{"username":"`+ username +`"}}`

	res, err := s.GetManggaForQuery(ctx, queryString)

	log.Printf(res)

	type KeyRecord struct {
		Key string `json:"Key"`
		Record User `json:"Record"`
	}

	type Result struct {
		Result []KeyRecord `json:"result"`
		Error error `json:"error"`
		ErrorData error `json:"errorData"`
	}

	if err != nil {
		return "", fmt.Errorf("User not found. %s", err.Error())
	}

	result := new(Result)
	json.Unmarshal([]byte(res), &result)

	orgName := result.Result[0].Record.OrgName

	if orgName == "Penangkar" || orgName == "Petani" {
		return result.Result[0].Record.AlamatLahan, nil
	} else {
		return result.Result[0].Record.AlamatToko, nil
	}

	return "", nil
}



func (s *ManggaContract) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) (string, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
    resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
    defer resultsIterator.Close()
    if err != nil {
        return "", err
    }
    // buffer is a JSON array containing QueryRecords
    var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse,
        err := resultsIterator.Next()
        if err != nil {
            return "", err
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
    fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
    return string(buffer.Bytes()), nil
}

// InitLedger
func (s *ManggaContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// for i, mangga := range manggas {
	// 	manggaAsBytes, _ := json.Marshal(mangga)
	// 	err := ctx.GetStub().PutState("mangga"+strconv.Itoa(i), manggaAsBytes)

	// 	if err != nil {
	// 		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	// 	}
	// }

	return nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(ManggaContract), new(UserContract))
	if err != nil {
		fmt.Printf("Error create mangga chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincodes: %s", err.Error())
	}

}