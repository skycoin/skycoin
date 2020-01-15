package visor

import "testing"

func TestPage_Cal(t *testing.T) {
	type fields struct {
		Size   uint64
		Number uint64
	}
	type args struct {
		n uint64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantStart uint64
		wantEnd   uint64
		wantPages uint64
		wantErr   error
	}{
		{
			name:      "size=0",
			fields:    fields{},
			wantStart: 0,
			wantEnd:   0,
			wantErr:   ErrZeroPageSize,
		},
		{
			name:      "page num=0",
			fields:    fields{Size: 1},
			wantStart: 0,
			wantEnd:   0,
			wantErr:   ErrZeroPageNum,
		},
		{
			name: "size=1 number=1 input 0",
			fields: fields{
				Size:   1,
				Number: 1,
			},
			args:      args{0},
			wantStart: 0,
			wantEnd:   0,
			wantPages: 0,
			wantErr:   nil,
		},
		{
			name: "size=1 number=1 input=5",
			fields: fields{
				Size:   1,
				Number: 1,
			},
			args:      args{5},
			wantStart: 0,
			wantEnd:   1,
			wantPages: 5,
		},
		{
			name: "size=1 number=2 input=5",
			fields: fields{
				Size:   1,
				Number: 2,
			},
			args:      args{5},
			wantStart: 1,
			wantEnd:   2,
			wantPages: 5,
		},
		{
			name: "size=1 number=3 input=5",
			fields: fields{
				Size:   1,
				Number: 3,
			},
			args:      args{5},
			wantStart: 2,
			wantEnd:   3,
			wantPages: 5,
		},
		{
			name: "size=1 number=4 input=5",
			fields: fields{
				Size:   1,
				Number: 4,
			},
			args:      args{5},
			wantStart: 3,
			wantEnd:   4,
			wantPages: 5,
		},
		{
			name: "size=1 number=5 input=5",
			fields: fields{
				Size:   1,
				Number: 5,
			},
			args:      args{5},
			wantStart: 4,
			wantEnd:   5,
			wantPages: 5,
		},
		{
			name: "size=1 number=6 input=5",
			fields: fields{
				Size:   1,
				Number: 6,
			},
			args:      args{5},
			wantStart: 0,
			wantEnd:   0,
			wantPages: 5,
		},
		{
			name: "size=10 number=1 input=100",
			fields: fields{
				Size:   10,
				Number: 1,
			},
			args:      args{100},
			wantStart: 0,
			wantEnd:   10,
			wantPages: 10,
		},
		{
			name: "size=10 number=9 input=100",
			fields: fields{
				Size:   10,
				Number: 9,
			},
			args:      args{100},
			wantStart: 80,
			wantEnd:   90,
			wantPages: 10,
		},
		{
			name: "size=10 number=10 input=100",
			fields: fields{
				Size:   10,
				Number: 10,
			},
			args:      args{100},
			wantStart: 90,
			wantEnd:   100,
			wantPages: 10,
		},
		{
			name: "size=10 number=11 input=100",
			fields: fields{
				Size:   10,
				Number: 11,
			},
			args:      args{100},
			wantStart: 0,
			wantEnd:   0,
			wantPages: 10,
		},
		{
			name: "size=9 number=10 input=100",
			fields: fields{
				Size:   9,
				Number: 10,
			},
			args:      args{100},
			wantStart: 81,
			wantEnd:   90,
			wantPages: 12,
		},
		{
			name: "size=9 number=11 input=100",
			fields: fields{
				Size:   9,
				Number: 11,
			},
			args:      args{100},
			wantStart: 90,
			wantEnd:   99,
			wantPages: 12,
		},
		{
			name: "size=99 number=2 input=100",
			fields: fields{
				Size:   99,
				Number: 2,
			},
			args:      args{100},
			wantStart: 99,
			wantEnd:   100,
			wantPages: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PageIndex{
				size: tt.fields.Size,
				n:    tt.fields.Number,
			}
			gotStart, gotEnd, gotPages, err := p.Cal(tt.args.n)
			if err != tt.wantErr {
				t.Errorf("Page.Cal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStart != tt.wantStart {
				t.Errorf("Page.Cal() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("Page.Cal() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
			if gotPages != tt.wantPages {
				t.Errorf("Page.Cal() gotPages = %v, want %v", gotPages, tt.wantPages)
			}
		})
	}
}
