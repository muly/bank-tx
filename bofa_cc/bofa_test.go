package bofa_cc

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseStatementPeriod(t *testing.T) {
	type args struct {
		periodStr string
	}
	tests := []struct {
		name      string
		args      args
		beginDate time.Time
		endDate   time.Time
		wantErr   bool
	}{
		{
			name:      "case 1: same year",
			args:      args{periodStr: "September 12 - October 11, 2024"},
			beginDate: time.Date(2024, 9, 12, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2024, 10, 11, 0, 0, 0, 0, time.UTC),
			wantErr:   false,
		},
		{
			name:      "case 2: different years",
			args:      args{periodStr: "December 12 - January 11, 2024"},
			beginDate: time.Date(2023, 12, 12, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			wantErr:   false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseStatementPeriod(tt.args.periodStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStatementPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.beginDate) {
				t.Errorf("parseStatementPeriod() got = %v, want %v", got, tt.beginDate)
			}
			if !reflect.DeepEqual(got1, tt.endDate) {
				t.Errorf("parseStatementPeriod() got1 = %v, want %v", got1, tt.endDate)
			}
		})
	}
}

func Test_addYearToDate(t *testing.T) {
	type args struct {
		dateStr     string
		startPeriod time.Time
		endPeriod   time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "case 1: same year",
			args:    args{
				dateStr:     "09/20",
				startPeriod: time.Date(2024, 9, 12, 0, 0, 0, 0, time.UTC),
				endPeriod:   time.Date(2024, 10, 11, 0, 0, 0, 0, time.UTC),
			},
			want:    time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "case 2: different year, 2 month range",
			args:    args{
				dateStr:     "12/20",
				startPeriod: time.Date(2023, 12, 12, 0, 0, 0, 0, time.UTC),
				endPeriod:   time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			},
			want:    time.Date(2023, 12, 20, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "case 3: different year, 2 month range",
			args:    args{
				dateStr:     "01/10",
				startPeriod: time.Date(2023, 12, 12, 0, 0, 0, 0, time.UTC),
				endPeriod:   time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			},
			want:    time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "case 4: different year, 3 month range",
			args:    args{
				dateStr:     "12/10",
				startPeriod: time.Date(2023, 11, 12, 0, 0, 0, 0, time.UTC),
				endPeriod:   time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
			},
			want:    time.Date(2023, 12, 10, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addYearToDate(tt.args.dateStr, tt.args.startPeriod, tt.args.endPeriod)
			if (err != nil) != tt.wantErr {
				t.Errorf("addYearToDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addYearToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
