package telegrambot

//func TestTelegramRobot_SendMessage(t *testing.T) {
//
//	telegramBot := NewTelegramRobot(token, chatID)
//
//	type fields struct {
//		TelegramBot
//	}
//	type args struct {
//		Level      string
//		EventName  string
//		Content    string
//		serverID   int32
//		vendorID   int32
//		vendorName string
//		serverName string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "test_telegram_bot",
//			fields: fields{
//				telegramBot,
//			},
//			args: args{
//				Level:      model.Level_Error,
//				serverID:   3345678,
//				vendorID:   45678,
//				vendorName: "test1",
//				serverName: "test2",
//				EventName:  "RTP",
//				Content:    "testing",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			bot := tt.fields.TelegramBot
//			msg := model.NewServerInfoMessage(tt.args.Level, tt.args.serverID, tt.args.serverName, tt.args.vendorID, tt.args.vendorName)
//
//			if err := bot.SendMessage(msg); (err != nil) != tt.wantErr {
//				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
