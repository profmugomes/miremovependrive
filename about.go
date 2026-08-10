// Copyright (c) 2024-2026 Murilo Gomes Julio. All Rights Reserved.

// Licensed under the PolyForm Strict License 1.0.0.
// See LICENSE.md for details.

package main

import (
	"fmt"
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/profmugomes/mgsmartflow/v2"
)

func showAbout() {
	w := a.NewWindow("Sobre")
	w.Resize(fyne.NewSize(597, 359))
	w.CenterOnScreen()
	w.SetFixedSize(true)

	flow := mgsmartflow.New()

	lblSoftware := canvas.NewText(fmt.Sprintf("MiRemovePendrive - Version: %s", VERSION_APP), color.Opaque)
	lblSoftware.TextSize = 18
	lblSoftware.TextStyle.Bold = true

	flow.AddRow(lblSoftware)
	flow.Move(lblSoftware, 7, 7)

	lblDesenvolvedor1 := widget.NewLabel("Desenvolvido por:")
	lblDesenvolvedor1.TextStyle = fyne.TextStyle{Bold: true}
	lblDesenvolvedor2 := widget.NewLabel("Murilo Gomes Julio")

	flow.AddColumn(lblDesenvolvedor1, lblDesenvolvedor2)
	flow.Resize(lblDesenvolvedor1, 142, 0)
	lblSite1 := widget.NewLabel("Site:")
	lblSite1.TextStyle = fyne.TextStyle{Bold: true}

	sURL, _ := url.Parse("https://www.profmugomes.com.br")
	lblSite2 := widget.NewHyperlink("https://www.profmugomes.com.br", sURL)

	flow.AddColumn(lblSite1, lblSite2)
	flow.Resize(lblSite1, 34, 0)

	lblCopyright1 := widget.NewLabel("Copyright (c) 2024-2026 Murilo Gomes Julio. All Rights Reserved.")
	lblCopyright1.TextStyle = fyne.TextStyle{Bold: true}
	flow.AddRow(lblCopyright1)

	lblLicense1 := widget.NewLabel("License:")
	lblLicense1.TextStyle = fyne.TextStyle{Bold: true}

	lblLicense2 := widget.NewLabel("PolyForm Strict 1.0.0")

	flow.AddColumn(lblLicense1, lblLicense2)
	flow.Resize(lblLicense1, 62, 0)

	txtLicense := widget.NewRichTextFromMarkdown(`
This project is licensed under the PolyForm Strict License 1.0.0.

### Summary

This software is available for noncommercial use only.

You may:
- ✔ Use the software for noncommercial purposes.
- ✔ Inspect and study the source code.

You may not:
- ✖ Use the software for commercial purposes.
- ✖ Modify the software.
- ✖ Redistribute the software.

See the full license terms at [LICENSE.md](LICENSE.md).

This summary is provided for convenience only.
	`)
	txtLicense.Wrapping = fyne.TextWrapWord

	vBoxLicense := container.NewVScroll(txtLicense)

	flow.AddRow(vBoxLicense)
	flow.Resize(vBoxLicense, w.Canvas().Size().Width, 319)

	w.SetContent(flow.Container)
	w.Show()
}
