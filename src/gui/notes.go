package gui

// NotesRPC note rpc
// type NotesRPC struct {
// 	Notes           wallet.Notes
// 	WalletDirectory string
// }

// // Ng global note
// var Ng *NotesRPC

// // InitWalletRPC init wallet rpc
// func InitWalletRPC(walletDir string, options ...wallet.Option) {
// 	Ng = NewNotesRPC(walletDir)
// }

// // NewNotesRPC new notes rpc
// func NewNotesRPC(walletDir string) *NotesRPC {
// 	rpc := &NotesRPC{}
// 	if err := os.MkdirAll(walletDir, os.FileMode(0700)); err != nil {
// 		logger.Panicf("Failed to create notes directory %s: %v", walletDir, err)
// 	}
// 	rpc.WalletDirectory = walletDir
// 	w, err := wallet.LoadNotes(rpc.WalletDirectory)
// 	if err != nil {
// 		logger.Panicf("Failed to load all notes: %v", err)
// 	}
// 	wallet.CreateNoteFileIfNotExist(walletDir)
// 	rpc.Notes = w
// 	return rpc
// }

// GetNotesReadable returns readable notes
// func (nt *NotesRPC) GetNotesReadable() wallet.ReadableNotes {
// 	return nt.Notes.ToReadable()
// }

// // Create a wallet Name is set by creation date
// func notesCreate(gateway *daemon.Gateway) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		logger.Info("API request made to create a note")
// 		note := r.FormValue("note")
// 		transactionID := r.FormValue("transaction_id")
// 		newNote := wallet.Note{
// 			TxID:  transactionID,
// 			Value: note,
// 		}
// 		Ng.Notes.SaveNote(Ng.WalletDirectory, newNote)
// 		rlt := Ng.GetNotesReadable()
// 		wh.SendOr500(w, rlt)
// 	}
// }

// // Returns a wallet by ID if GET.  Creates or updates a wallet if POST.
// func notesHandler(gateway *daemon.Gateway) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		//ret := wallet.Wallets.ToPublicReadable()
// 		ret := Ng.GetNotesReadable()
// 		wh.SendOr404(w, ret)
// 	}
// }

// mux.Handle("/notes", notesHandler(gateway))

// mux.Handle("/notes/create", notesCreate(gateway))
