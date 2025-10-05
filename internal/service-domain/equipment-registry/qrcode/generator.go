package qrcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
	qrcode "github.com/skip2/go-qrcode"
)

// QRData represents the data encoded in the QR code
type QRData struct {
	URL       string `json:"url"`
	ID        string `json:"id"`
	SerialNo  string `json:"serial"`
	QRCode    string `json:"qr"`
}

// Generator handles QR code generation
type Generator struct {
	baseURL    string
	outputDir  string
	qrSize     int
}

// NewGenerator creates a new QR code generator
func NewGenerator(baseURL, outputDir string) *Generator {
	return &Generator{
		baseURL:   baseURL,
		outputDir: outputDir,
		qrSize:    256, // 256x256 pixels
	}
}

// GenerateQRCode generates a QR code image for equipment
func (g *Generator) GenerateQRCode(equipmentID, serialNumber, qrCodeID string) (string, error) {
	// Create QR data with URL and identifiers
	url := fmt.Sprintf("%s/equipment/%s", g.baseURL, equipmentID)
	
	qrData := QRData{
		URL:      url,
		ID:       equipmentID,
		SerialNo: serialNumber,
		QRCode:   qrCodeID,
	}

	// Encode QR data as JSON
	jsonData, err := json.Marshal(qrData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal QR data: %w", err)
	}

	// Generate QR code image
	qrFilename := fmt.Sprintf("qr_%s.png", equipmentID)
	qrPath := filepath.Join(g.outputDir, qrFilename)

	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate QR code with medium error correction
	err = qrcode.WriteFile(string(jsonData), qrcode.Medium, g.qrSize, qrPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qrPath, nil
}

// GenerateQRCodeBytes generates a QR code as byte array for database storage
func (g *Generator) GenerateQRCodeBytes(equipmentID, serialNumber, qrCodeID string) ([]byte, error) {
	// Create QR data with URL and identifiers
	url := fmt.Sprintf("%s/equipment/%s", g.baseURL, equipmentID)
	
	qrData := QRData{
		URL:      url,
		ID:       equipmentID,
		SerialNo: serialNumber,
		QRCode:   qrCodeID,
	}

	// Encode QR data as JSON
	jsonData, err := json.Marshal(qrData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal QR data: %w", err)
	}

	// Generate QR code as PNG bytes with medium error correction
	qrBytes, err := qrcode.Encode(string(jsonData), qrcode.Medium, g.qrSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qrBytes, nil
}

// GenerateQRLabel generates a printable PDF label with QR code (legacy filesystem version)
func (g *Generator) GenerateQRLabel(equipmentID, equipmentName, serialNumber, manufacturer, qrCodeID, qrImagePath string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Equipment QR Code Label")
	pdf.Ln(12)

	// Equipment details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Equipment:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, equipmentName)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Manufacturer:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, manufacturer)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Serial Number:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, serialNumber)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "QR Code:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, qrCodeID)
	pdf.Ln(12)

	// Add QR code image
	if _, err := os.Stat(qrImagePath); err == nil {
		// Image exists, add it to PDF
		pdf.Image(qrImagePath, 40, pdf.GetY(), 60, 60, false, "", 0, "")
		pdf.Ln(65)
	}

	// Add instructions
	pdf.SetFont("Arial", "I", 10)
	pdf.MultiCell(0, 5, "Scan this QR code with your mobile device to view equipment details and request service.", "", "", false)
	pdf.Ln(5)
	
	// Add URL for reference
	url := fmt.Sprintf("%s/equipment/%s", g.baseURL, equipmentID)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(0, 5, fmt.Sprintf("URL: %s", url))

	// Save PDF
	pdfFilename := fmt.Sprintf("qr_label_%s.pdf", equipmentID)
	pdfPath := filepath.Join(g.outputDir, pdfFilename)

	err := pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF label: %w", err)
	}

	return pdfPath, nil
}

// GenerateQRLabelFromBytes generates a printable PDF label with QR code from byte array
func (g *Generator) GenerateQRLabelFromBytes(equipmentID, equipmentName, serialNumber, manufacturer, qrCodeID string, qrImageBytes []byte) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Equipment QR Code Label")
	pdf.Ln(12)

	// Equipment details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Equipment:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, equipmentName)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Manufacturer:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, manufacturer)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Serial Number:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, serialNumber)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "QR Code:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, qrCodeID)
	pdf.Ln(12)

	// Add QR code image from bytes
	if len(qrImageBytes) > 0 {
		reader := bytes.NewReader(qrImageBytes)
		pdf.RegisterImageReader(fmt.Sprintf("qr_%s", equipmentID), "PNG", reader)
		pdf.Image(fmt.Sprintf("qr_%s", equipmentID), 40, pdf.GetY(), 60, 60, false, "", 0, "")
		pdf.Ln(65)
	}

	// Add instructions
	pdf.SetFont("Arial", "I", 10)
	pdf.MultiCell(0, 5, "Scan this QR code with your mobile device to view equipment details and request service.", "", "", false)
	pdf.Ln(5)
	
	// Add URL for reference
	url := fmt.Sprintf("%s/equipment/%s", g.baseURL, equipmentID)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(0, 5, fmt.Sprintf("URL: %s", url))

	// Output PDF to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF label: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateBatchLabels generates PDF with multiple QR labels (for printing)
func (g *Generator) GenerateBatchLabels(equipmentList []EquipmentInfo) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	
	for _, eq := range equipmentList {
		pdf.AddPage()

		// Title
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "Equipment Service Label")
		pdf.Ln(10)

		// Equipment details
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(40, 7, "Equipment:")
		pdf.SetFont("Arial", "", 11)
		pdf.MultiCell(0, 7, eq.EquipmentName, "", "", false)

		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(40, 7, "Manufacturer:")
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(0, 7, eq.Manufacturer)
		pdf.Ln(7)

		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(40, 7, "Serial Number:")
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(0, 7, eq.SerialNumber)
		pdf.Ln(7)

		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(40, 7, "QR Code:")
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(0, 7, eq.QRCode)
		pdf.Ln(10)

		// Generate QR code on the fly
		qrImagePath, err := g.GenerateQRCode(eq.EquipmentID, eq.SerialNumber, eq.QRCode)
		if err == nil && qrImagePath != "" {
			if _, err := os.Stat(qrImagePath); err == nil {
				pdf.Image(qrImagePath, 40, pdf.GetY(), 50, 50, false, "", 0, "")
				pdf.Ln(55)
			}
		}

		// Instructions
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, "For service requests, scan QR code or WhatsApp the photo to our service number.", "", "", false)
		pdf.Ln(3)
		
		// URL
		url := fmt.Sprintf("%s/equipment/%s", g.baseURL, eq.EquipmentID)
		pdf.SetFont("Arial", "", 8)
		pdf.Cell(0, 4, fmt.Sprintf("Web: %s", url))
	}

	// Save batch PDF
	pdfPath := filepath.Join(g.outputDir, "qr_labels_batch.pdf")
	err := pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate batch labels: %w", err)
	}

	return pdfPath, nil
}

// EquipmentInfo holds equipment information for batch label generation
type EquipmentInfo struct {
	EquipmentID   string
	EquipmentName string
	Manufacturer  string
	SerialNumber  string
	QRCode        string
}

// DecodeQRData decodes JSON data from QR code string
func DecodeQRData(qrString string) (*QRData, error) {
	var qrData QRData
	err := json.Unmarshal([]byte(qrString), &qrData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode QR data: %w", err)
	}
	return &qrData, nil
}

// DecodeQRFromImage decodes QR code data from an image file
func (g *Generator) DecodeQRFromImage(imagePath string) (*QRData, error) {
	// For MVP, we'll use a simplified approach
	// In production, integrate a QR decoding library like github.com/makiuchi-d/gozxing
	
	// For now, return an error indicating this needs to be implemented with proper QR decoder
	return nil, fmt.Errorf("QR decode from image not yet implemented - requires QR decoder library integration")
}
