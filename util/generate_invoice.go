package util

import (
	"fmt"
	"os"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/nibrasmuhamed/cartique/models"
)

func GenerateInvoice(order models.Invoice) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)

	buildHeading(m, order)
	buildFruitList(m, order)
	BuidFooter(m)

	err := m.OutputFileAndClose("public/images/recipt.pdf")
	if err != nil {
		fmt.Println("⚠️  Could not save PDF:", err)
		os.Exit(1)
	}

	fmt.Println("PDF saved successfully")
}

func buildHeading(m pdf.Maroto, order models.Invoice) {
	m.RegisterHeader(func() {
		m.Row(50, func() {
			m.Col(12, func() {
				err := m.FileImage("public/images/my.jpg", props.Rect{
					Center:  true,
					Percent: 75,
				})

				if err != nil {
					fmt.Println("Image file was not loaded 😱 - ", err)
				}
			})
		})
		m.Row(10, func() {
			m.Text(order.CreatedAt,
				props.Text{Align: consts.Left})
		})
		m.Row(5, func() {
			m.Text(order.Name,
				props.Text{Align: consts.Left})
		})
		m.Row(5, func() {
			m.Text(order.Phone,
				props.Text{Align: consts.Left})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Invoice Prepared by Cartique LTD", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
				Color: getDarkPurpleColor(),
			})
		})
	})
}

func buildFruitList(m pdf.Maroto, order models.Invoice) {
	tableHeadings := []string{"Order ID", "Quantity", "Price"}
	lightPurpleColor := getLightPurpleColor()

	m.SetBackgroundColor(getTealColor())
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Products", props.Text{
				Top:    2,
				Size:   13,
				Color:  color.NewWhite(),
				Family: consts.Courier,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())
	contents := [][]string{{order.Product, order.Quantity, order.Price}}
	m.TableList(tableHeadings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 7, 2},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 7, 2},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})

}

func BuidFooter(m pdf.Maroto) {
	m.RegisterFooter(func() {
		m.Row(6, func() {
			m.Text("This is a computer generated bill which do not need signature", props.Text{Align: consts.Center})
		})
		m.Row(10, func() {
			m.Col(6, func() {
				m.Signature("cartique limited.")
			})

		})
	})
}
func getDarkPurpleColor() color.Color {
	return color.Color{
		Red:   88,
		Green: 80,
		Blue:  99,
	}
}

func getLightPurpleColor() color.Color {
	return color.Color{
		Red:   210,
		Green: 200,
		Blue:  230,
	}
}

func getTealColor() color.Color {
	return color.Color{
		Red:   3,
		Green: 166,
		Blue:  166,
	}
}
