package main

import (
	_ "bytes"
	"log"
	"os"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	wdStr, err := os.Getwd()
	path := wdStr + "\\pdf.pdf"
	if err != nil {
		log.Fatal(err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	pdf.AddPageFormat("L", gofpdf.SizeType{Wd: 210, Ht: 297})
	pdf.Cell(160, 6, "Bukti Kegiatan Harian")
	pdf.Cell(120, 6, "Nama Pegawai: dr. I Wayan Surya Sedana")
	pdf.Ln(-1)
	pdf.Cell(160, 6, "Pegawai RSUP Sanglah Denpasar")
	pdf.Cell(120, 6, "NIP/Gol: 198702112014121001")
	pdf.Ln(-1)
	pdf.Cell(160, 6, "Bulan: Juni Tahun 2017")
	pdf.Cell(120, 6, "Tempat Tugas: IGD RSUP Sanglah")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 20, "Uraian", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 20, "Poin", "1", 0, "C", false, 0, "")
	pdf.CellFormat(176, 10, "Jumlah Kegiatan Harian", "1", 2, "C", false, 0, "")
	for i := 1; i < 17; i++ {
		pdf.CellFormat(11, 10, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.SetXY(266, 28)
	pdf.CellFormat(25, 20, "Jumlah Poin", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(10, 24, "1", "1", 0, "C", false, 0, "")

	pdf.MultiCell(50, 6, "Melakukan pelayanan medik umum (per pasien : pemeriksaan rawat jalan, IGD, visite rawat inap, tim medis diskusi", "1", "L", false)
	pdf.SetXY(70, 48)
	pdf.CellFormat(20, 24, "0,0032", "1", 0, "C", false, 0, "")
	for i := 1; i < 17; i++ {
		pdf.CellFormat(11, 24, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.CellFormat(25, 24, "xxx", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(10, 12, "2", "1", 0, "C", false, 0, "")
	pdf.MultiCell(50, 6, "Melakukan tindakan medik umum tingkat sederhana (per tindakan)", "1", "L", false)
	pdf.SetXY(70, 72)
	pdf.CellFormat(20, 12, "0,0032", "1", 0, "C", false, 0, "")
	for i := 1; i < 17; i++ {
		pdf.CellFormat(11, 12, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.CellFormat(25, 12, "xxx", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)
	// Baris ke dua
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 20, "Uraian", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 20, "Poin", "1", 0, "C", false, 0, "")
	pdf.CellFormat(176, 10, "Jumlah Kegiatan Harian", "1", 2, "C", false, 0, "")
	for i := 1; i < 16; i++ {
		pdf.CellFormat(11, 10, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.SetFont("Arial", "B", 7)
	pdf.MultiCell(11, 5, "Jumlah Poin", "1", "C", false)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(266, 96)
	pdf.MultiCell(25, 20, "Jumlah X Poin", "1", "C", false)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(10, 24, "1", "1", 0, "C", false, 0, "")

	pdf.MultiCell(50, 6, "Melakukan pelayanan medik umum (per pasien : pemeriksaan rawat jalan, IGD, visite rawat inap, tim medis diskusi", "1", "L", false)
	pdf.SetXY(70, 116)
	pdf.CellFormat(20, 24, "0,0032", "1", 0, "C", false, 0, "")
	for i := 17; i <= 32; i++ {
		pdf.CellFormat(11, 24, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.CellFormat(25, 24, "xxx", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(10, 12, "2", "1", 0, "C", false, 0, "")
	pdf.MultiCell(50, 6, "Melakukan tindakan medik umum tingkat sederhana (per tindakan)", "1", "L", false)
	pdf.SetXY(70, 140)
	pdf.CellFormat(20, 12, "0,0032", "1", 0, "C", false, 0, "")
	for i := 17; i <= 32; i++ {
		pdf.CellFormat(11, 12, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	}
	pdf.CellFormat(25, 12, "xxx", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(256, 6, "Jumlah Point X Volume kegiatan pelayanan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, "xxx", "1", 1, "C", false, 0, "")
	pdf.CellFormat(256, 6, "Target Point kegiatan pelayanan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, "xxx", "1", 1, "C", false, 0, "")

	pdf.AddPage()
	// pdf.SetFont("Arial", "B", 16)
	// wd := pdf.GetStringWidth("Buku Catatan Pribadi")
	// pdf.SetX((210 - wd) / 2)
	// pdf.Cell(wd, 9, "Buku Catatan Pribadi")
	// pdf.Ln(10)
	// pdf.SetFont("Arial", "", 12)
	// pdf.Cell(20, 5, "Nama")
	// pdf.Cell(105, 5, ": dr. I Wayan Surya Sedana")
	// pdf.Ln(-1)
	// pdf.Cell(20, 5, "Bulan")
	// pdf.Cell(105, 5, ": Agustus 2017")
	// pdf.Ln(-1)
	// pdf.Ln(-1)
	// pdf.SetFont("Arial", "", 10)
	// pdf.CellFormat(9, 20, "No", "1", 0, "C", false, 0, "")
	// pdf.CellFormat(15, 20, "Tanggal", "1", 0, "C", false, 0, "")
	// pdf.CellFormat(20, 20, "No. CM", "1", 0, "C", false, 0, "")
	// pdf.CellFormat(60, 20, "Nama", "1", 0, "C", false, 0, "")
	// pdf.CellFormat(40, 20, "Diagnosis", "1", 0, "C", false, 0, "")

	// pdf.MultiCell(20, 5, "Melakukan pelayanan medik umum", "1", "C", false)
	// // fmt.Println(pdf.GetXY())
	// pdf.SetXY(174, 35)
	// pdf.MultiCell(25, 4, "Melakukan tindakan medik umum tingkat sederhana", "1", "C", false)
	// // diag := []string{"a","a","a"}
	// diag := "aaaaa"
	// for i := 1; i < 40; i++ {
	// 	diag = diag + "a"
	// 	if len(diag) > 10 {
	// 		diag = diag[:10]
	// 	}
	// 	num := strconv.Itoa(i)
	// 	pdf.CellFormat(9, 7, num, "1", 0, "C", false, 0, "")
	// 	pdf.CellFormat(15, 7, fmt.Sprintf("tanggal %v", i), "1", 0, "C", false, 0, "")
	// 	pdf.CellFormat(20, 7, "aaaaa", "1", 0, "C", false, 0, "")
	// 	pdf.CellFormat(60, 7, "aaaaa", "1", 0, "C", false, 0, "")
	// 	pdf.CellFormat(40, 7, diag, "1", 0, "C", false, 0, "")
	// 	pdf.CellFormat(20, 7, "aaa", "1", 0, "C", false, 0, "")
	// 	pdf.CellFormat(25, 7, "aaa", "1", 0, "C", false, 0, "")
	// 	pdf.Ln(-1)
	// }

	// pdf.SetHeaderFunc(func() {
	// 	pdf.SetFont("Arial", "B", 15)
	// 	wd := pdf.GetStringWidth("Buku Catatan Pribadi")
	// 	pdf.SetX((210 - wd) / 2)
	// 	pdf.Cell(wd, 9, "Buku Catatan Pribadi")
	// 	pdf.Ln(10)
	// })

	//////////////////////////////////////////////
	// header := []string{"Country", "Capital", "Area (sq km)", "Pop. (thousands)"}
	// for _, str := range header {
	// 	pdf.CellFormat(40, 7, str, "1", 0, "", false, 0, "")
	// }
	// pdf.Ln(-1)
	////////////////////////////////////////////////////
	// //////////////////////////////////////////////////////////
	// 	titleStr := "20000 Leagues Under the Seas"
	// 	pdf.SetTitle(titleStr, false)
	// 	pdf.SetHeaderFunc(func() {
	// 		pdf.SetFont("Arial", "B", 15)
	// 		wd := pdf.GetStringWidth(titleStr) + 6
	// 		pdf.SetX((210 - wd) / 2)
	// 		pdf.SetDrawColor(0, 80, 180)
	// 		pdf.SetFillColor(230, 230, 0)
	// 		pdf.SetTextColor(220, 50, 50)
	// 		pdf.SetLineWidth(1)
	// 		pdf.CellFormat(wd, 9, titleStr, "1", 1, "c", true, 0, "")
	// 		pdf.Ln(10)
	// 	})

	// 	pdf.SetFooterFunc(func() {
	// 		pdf.SetY(-15)
	// 		pdf.SetFont("Arial", "I", 8)
	// 		pdf.SetTextColor(128, 128, 128)
	// 		pdf.CellFormat(0, 10, fmt.Sprintf("Halaman %d", pdf.PageNo()), "", 0, "C", false, 0, "")
	// 	})

	// 	chapterTitle := func(chapNum int, titleStr string) {
	// 		pdf.SetFont("Arial", "", 12)
	// 		pdf.SetFillColor(200, 220, 255)
	// 		pdf.CellFormat(0, 6, fmt.Sprintf("Chapter %d : %s", chapNum, titleStr), "", 1, "L", true, 0, "")
	// 		pdf.Ln(4)
	// 	}

	/////////////////////////////////////////////////////////
	// pdf set header
	// pdf := gofpdf.New("P", "mm", "A4", "")
	// pdf.SetHeaderFunc(func() {
	// 	pdf.SetY(5)
	// 	pdf.SetFont("Arial", "B", 15)
	// 	pdf.Cell(80, 0, "")
	// 	pdf.CellFormat(30, 10, "Title", "1", 0, "C", false, 0, "")
	// 	pdf.Ln(20)
	// })
	// pdf.SetFooterFunc(func() {
	// 	pdf.SetY(-15)
	// 	pdf.SetFont("Arial", "I", 8)
	// 	pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	// })

	// pdf.AliasNbPages("")
	// pdf.AddPage()
	// pdf.SetFont("Times", "", 12)
	// for j := 1; j <= 40; j++ {
	// 	pdf.CellFormat(0, 10, fmt.Sprintf("Printing line number %d", j), "", 1, "", false, 0, "")
	// }

	err = pdf.OutputFileAndClose(path)
	if err != nil {
		log.Fatal(err)
	}

	// wdStr, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Print(wdStr)

	// file, err := os.Create(path)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// berkas, err := os.OpenFile(path, os.O_RDWR, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer berkas.Close()

	// _, err = berkas.WriteString("halo\n")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = berkas.WriteString("Dunia")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = berkas.Sync()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// func makePDF(){
// 	pdf := gofpdf.New("P", "mm", "A4", "")
// 	pdf.AddPage()
// 	pdf.SetFont("Arial", "B", 16)
// 	pdf.Cell(40,10, "Hello World!")
// 	fileStr
// }
