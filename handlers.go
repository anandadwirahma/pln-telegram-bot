package main

import (
	"github.com/yanzay/tbot"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (a *application) greatingHandler(m *tbot.Message) {
	buttons := makeMenuButtons()

	msg := "Halo Electrizen\n\n" +
		"Tekan <b>Baca Meteran</b> untuk Baca Meter Mandiri pemakaian listrik (pascabayar) dan cek tagihan"

	a.client.SendMessage(m.Chat.ID, msg, tbot.OptParseModeHTML, tbot.OptInlineKeyboardMarkup(buttons))
}

func (a *application) commonHandler(m *tbot.Message) {
	a.checkState(m)
}

func (a *application) callbackHandler(cq *tbot.CallbackQuery) {
	humanMove := cq.Data
	msg := draw(humanMove)
	a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	a.client.SendMessage(cq.Message.Chat.ID, msg)
}

func (a *application) checkState(m *tbot.Message) {
	switch tmpState {
	case "mainOne":
		tmpMsg = m.Text
		idPelanggan = m.Text
		tmpState = "uploadPhoto"

		msg := "Kirim Photo kWh meter dan pastikan terlihat dengan jelas , seperti gambar berikut"

		a.client.SendPhotoFile(m.Chat.ID, "kwh.jpg", tbot.OptCaption(msg))
	case "uploadPhoto":

		if m.Document != nil {
			photo, err := a.client.GetFile(m.Document.FileID)
			if err != nil {
				log.Print(err)
				return
			}
			url := a.client.FileURL(photo)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			out, err := os.Create(m.Document.FileUniqueID + ".png")
			if err != nil {
				log.Println(err)
				return
			}
			defer out.Close()
			io.Copy(out, resp.Body)

			err = getOcrData(m.Document.FileUniqueID + ".png")
			if err != nil {
				log.Println(err)
				return
			}
		}

		msg := "Baik, pembacaan meteran akan kami proses terlebih dahulu"
		a.client.SendMessage(m.Chat.ID, msg)
		a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
		time.Sleep(3)

		billingMsg := "<b>Berikut informasi billing anda</b> \n" +
			"<b>==============================</b> \n\n" +
			"ID Pelanggan : " + idPelanggan + "\n" +
			"Nama : " + "Ananda Dwi Rahma\n" +
			"Pemakaian KWH : " + strconv.Itoa(UsageCurrentMonth) + "kWh \n" +
			"Daya Listrik : " + dayaPelanggan + " \n" +
			"Biaya : Rp. " + strconv.Itoa(billingAmount) + " \n" +
			"<b>Silahkan lakukan pembayaran sebelum jatuh tempo ya electrizen.</b>"
		a.client.SendMessage(m.Chat.ID, billingMsg, tbot.OptParseModeHTML)

		idPelanggan = ""
		tmpMsg = ""
		tmpState = ""
	}
}
