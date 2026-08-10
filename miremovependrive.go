// Copyright (c) 2024-2026 Murilo Gomes Julio. All Rights Reserved.

// Licensed under the PolyForm Strict License 1.0.0.
// See LICENSE.md for details.

package main

import (
	"fmt"
	"net/url"
	"os/user"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/profmugomes/mgcolumnview/v2"
	"github.com/profmugomes/mgdialogbox/v2"
	"github.com/profmugomes/mgrun/v2"
	"github.com/profmugomes/mgsmartflow/v2"
)

const VERSION_APP string = "4.0.0"

var a fyne.App

func convertBytes(value float64) string {
	var base float64
	var tipo string

	if value < 1024 {
		base = value
		tipo = "bytes"
	} else if value < 1048576 {
		base = value / 1024
		tipo = "KB"
	} else if value < 1073741824 {
		base = value / 1024 / 1024
		tipo = "MB"
	} else if value < 1099511627776 {
		base = value / 1024 / 1024 / 1024
		tipo = "GB"
	} else if value < 1125899906842624 {
		base = value / 1024 / 1024 / 1024 / 1024
		tipo = "TB"
	}

	return fmt.Sprintf("%.2f %s", base, tipo)
}

func getRun(command string) string {
	sRun := mgrun.New(command)

	var sLine string
	sRun.OnStdout(func(line string) {
		sLine = line
	})

	if err := sRun.Run(); err != nil {
		fmt.Printf("Erro na execução: %v\n", err)
	}

	return sLine
}

func getDrivers() []string {
	userCurrent, _ := user.Current()
	sResult := getRun(fmt.Sprintf("ls %s", strings.Join([]string{"/run/media/", userCurrent.Name, "/"}, "")))

	list := strings.Split(sResult, "\n")

	itens := []string{}
	for _, row := range list {
		if row != "" {
			itens = append(itens, row)
		}
	}

	return itens
}

func main() {
	sIcon := fyne.NewStaticResource("miremovependrive.png", resourceAppIconPngData)

	a = app.NewWithID("br.com.mugomes.miremovependrive")
	a.Settings().SetTheme(&myDarkTheme{})
	a.SetIcon(sIcon)

	window := a.NewWindow("MiRemovePendrive")
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(530, 275))

	mnuAbout := fyne.NewMenu("Sobre",
		fyne.NewMenuItem("Verificar Atualização", func() {
			url, _ := url.Parse("https://github.com/mugomes/miremovependrive/releases")
			a.OpenURL(url)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Apoie MiRemovePendrive", func() {
			url, _ := url.Parse("https://mugomes.github.io/apoie.html")
			a.OpenURL(url)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sobre MiRemovePendrive", func() {
			showAbout()
		}),
	)

	window.SetMainMenu(fyne.NewMainMenu(mnuAbout))

	flow := mgsmartflow.New()

	lblDisp := widget.NewLabel("Selecione o dispositivo que deseja remover:")
	lblDisp.TextStyle = fyne.TextStyle{Bold: true}
	flow.AddRow(lblDisp)

	cboDisp := widget.NewSelect(nil, nil)
	cboDisp.PlaceHolder = " "
	cboDisp.Options = getDrivers()

	if len(cboDisp.Options) > 0 {
		cboDisp.Selected = cboDisp.Options[0]
	}

	flow.AddRow(cboDisp)
	flow.Gap(cboDisp, 0, 14)

	lblInfo := widget.NewLabel("")

	header := []string{"Disponível", "Usado", "Total"}
	wHeader := []float32{150, 170, 150}
	cv := mgcolumnview.NewColumnView(header, wHeader, false)

	userCurrent, _ := user.Current()

	var btnRemover *widget.Button
	var btnInfo *widget.Button

	checkDisable := func() {
		if len(cboDisp.Options) > 0 {
			btnRemover.Enable()
			btnInfo.Enable()
		} else {
			btnRemover.Disable()
			btnInfo.Disable()
		}
	}

	btnRemover = widget.NewButton("Remover", func() {
		go func() {
			fyne.Do(func() {
				flow.AddRow(lblInfo)
				flow.Container.Remove(cv)

				lblInfo.SetText(fmt.Sprintf("Obtendo dispositivo %s...", cboDisp.Selected))
			})

			sLSBLK := fmt.Sprintf("lsblk -l | grep \"/run/media/%s/%s\"", userCurrent.Name, cboDisp.Selected)
			time.Sleep(1 * time.Second)
			sDriver := strings.Split(getRun(sLSBLK), " ")[0]
			time.Sleep(1 * time.Second)

			fyne.Do(func() {
				lblInfo.SetText(fmt.Sprintf("Desmontando dispositivo %s...", sDriver))
			})

			uDisks1 := fmt.Sprintf("udisksctl unmount -b /dev/%s", sDriver)
			time.Sleep(1 * time.Second)
			sInfo := getRun(uDisks1)
			time.Sleep(1 * time.Second)

			fyne.Do(func() {
				lblInfo.SetText(fmt.Sprintf("Desligando dispositivo %s...", sDriver))
			})

			if sInfo != "" {
				uDisks2 := fmt.Sprintf("udisksctl power-off -b /dev/%s", sDriver)
				time.Sleep(1 * time.Second)
				sInfo = getRun(uDisks2)
				time.Sleep(1 * time.Second)

				fyne.Do(func() {
					flow.Container.Remove(lblInfo)
					flow.Container.Refresh()

					if sInfo == "" {
						mgdialogbox.NewAlert(a, "MiRemovePendrive", "Dispositivo removido com sucesso!", false, "Ok", nil)

						cboDisp.ClearSelected()
						cboDisp.Options = getDrivers()
						cboDisp.Refresh()

						if len(cboDisp.Options) > 0 {
							cboDisp.Selected = cboDisp.Options[0]
						} else {
							checkDisable()
						}
					} else {
						mgdialogbox.NewAlert(a, "MiRemovePendrive", "Dispositivo ocupado, não é possível remover!", true, "Ok", nil)
					}
				})
			} else {
				fyne.Do(func() {
					flow.Container.Remove(lblInfo)
					mgdialogbox.NewAlert(a, "MiRemovePendrive", "Dispositivo ocupado, não é possível desmontá-lo!", true, "Ok", nil)
				})
			}
		}()
	})

	var btnAtualizar *widget.Button

	btnInfo = widget.NewButton("Informação", func() {
		flow.Container.Remove(cv)
		flow.Gap(btnAtualizar, 0, 19)
		flow.AddRow(cv)
		flow.Resize(cv, window.Canvas().Size().Width, 84)

		userMedia := fmt.Sprintf("df -B1 | grep \"/run/media/%s/%s\"", userCurrent.Name, cboDisp.Selected)
		sResult := getRun(userMedia)
		sSize := strings.Split(strings.TrimSpace(sResult), " ")

		aItens := []string{}
		for _, row := range sSize {
			if row != "" {
				aItens = append(aItens, row)
			}
		}

		sAvailable, _ := strconv.ParseFloat(aItens[3], 64)
		sUsed, _ := strconv.ParseFloat(aItens[2], 64)
		sTextUsed := fmt.Sprintf("%s %s", convertBytes(sUsed), aItens[4])
		sTotal, _ := strconv.ParseFloat(aItens[1], 64)

		cv.RemoveAll()
		cv.AddRow([]string{convertBytes(sAvailable), sTextUsed, convertBytes(sTotal)})
	})

	btnAtualizar = widget.NewButton("Atualizar", func() {
		cboDisp.Options = getDrivers()

		if len(cboDisp.Options) > 0 {
			cboDisp.Selected = cboDisp.Options[0]
		}

		checkDisable()
	})

	checkDisable()

	flow.AddColumn(btnRemover, btnInfo, btnAtualizar)

	window.SetContent(flow.Container)
	window.ShowAndRun()
}
