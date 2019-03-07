package pages

// import (
// 	"fmt"
// 		"context"


// 	"github.com/rivo/tview"
// 		"github.com/raedahgroup/godcr/app"

// )

// func CreatWalletPage(tviewApp *tview.Application, ctx context.Context, walletMiddleware app.WalletMiddleware) tview.Primitive {
// var password, confPassword string
// 	form := tview.NewForm().
// 		AddPasswordField("Password", "", 20, '*', func(text string) {
// 			if len(text) == 0 {
// 				fmt.Println("password cannot be less than 4")
// 				return
// 			}
// 			password = text
// 		}).
// 		AddPasswordField("Password", "", 20, '*', func(text string) {
// 			confPassword = text
// 		}).
// 		AddButton("Create", func() {
// 			if password != confPassword {
// 				fmt.Println("password does not match")
// 				return
// 			}
// 			CreateWallet(ctx, password, walletMiddleware)
// 		}).
// 		AddButton("Quit", func() {
// 			tviewApp.Stop()
// 		})
// 	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignCenter).SetRect(30, 10, 40, 10)
// return form
// }