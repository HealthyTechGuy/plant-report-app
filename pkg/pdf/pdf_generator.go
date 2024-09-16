package pdf

// PDFGenerator defines the methods for generating PDF reports
type PDFGenerator interface {
	GeneratePDF(plantInfo *PlantInfo) ([]byte, error)
}
