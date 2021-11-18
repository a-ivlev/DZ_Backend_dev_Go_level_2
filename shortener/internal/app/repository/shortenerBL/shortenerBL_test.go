package shortenerBL_test

import (
	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/repository/shortenerBL"
	"context"
	"github.com/google/uuid"
	"testing"
)

func TestGenerateShortLink(t *testing.T) {

	idtest1, _ := uuid.Parse("17830224-2532-41a5-a389-4fd8244efbaa")
	idtest2, _ := uuid.Parse("9ba87644-294d-4f3e-b5df-51c9d8c330b3")
	idtest3, _ := uuid.Parse("88d3ea27-336e-41e7-b0da-877cc8a041d3")

	cases := map[string]struct {
		ctx context.Context
		id uuid.UUID
		expected string

	}{
		"Test 1": {
			ctx: context.WithValue(context.Background(), "srvHost", "http://localhost"),
			id: idtest1,
			expected: "41A5A389",
		},
		"Test 2": {
			ctx: context.WithValue(context.Background(), "srvHost", "http://localhost"),
			id: idtest2,
			expected: "4F3EB5DF",
		},
		"Test 3": {
			ctx: context.WithValue(context.Background(), "srvHost", "http://localhost"),
			id: idtest3,
			expected: "41E7B0DA",
		},

	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			result := shortenerBL.GenerateShortLink(cs.ctx, cs.id)
			if result != cs.expected {
				t.Errorf("GenerateShortLink returned wrong result: got %v want %v", cs.expected, result)
			}
		})
	}
}

func TestGenerateStatLink(t *testing.T) {

	idtest1, _ := uuid.Parse("17830224-2532-41a5-a389-4fd8244efbaa")
	idtest2, _ := uuid.Parse("9ba87644-294d-4f3e-b5df-51c9d8c330b3")
	idtest3, _ := uuid.Parse("88d3ea27-336e-41e7-b0da-877cc8a041d3")

	cases := map[string]struct {
		ctx context.Context
		id uuid.UUID
		expected string

	}{
		"Test 1": {
			ctx: context.WithValue(context.Background(), "srvHost", "http://localhost"),
			id: idtest1,
			expected: "17830224253241a5a3894fd8244efbaa",
		},
		"Test 2": {
			ctx: context.WithValue(context.Background(), "srvHost", "http://localhost"),
			id: idtest2,
			expected: "9ba87644294d4f3eb5df51c9d8c330b3",
		},
		"Test 3": {
			ctx: context.WithValue(context.Background(), "srvHost", "http://localhost"),
			id: idtest3,
			expected: "88d3ea27336e41e7b0da877cc8a041d3",
		},

	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			result := shortenerBL.GenerateStatLink(cs.ctx, cs.id)
			if result != cs.expected {
				t.Errorf("GenerateShortLink returned wrong result: got %v want %v", cs.expected, result)
			}
		})
	}
}

type MockDB struct {

}

func TestCreateShort(t *testing.T) {

	//cases := map[string]struct {
	//	ctx context.Context
	//	shortener *shortenerBL.Shortener
	//	expected *shortenerBL.Shortener
	//	//errorMessage string
	//}{
	//	"Test 1": {
	//		ctx :  context.Background(),
	//		shortener: &shortenerBL.Shortener{
	//			FullLink: "http://test.local/test",
	//		},
	//		expected: &shortenerBL.Shortener{
	//			FullLink: "http://test.local/test",
	//		},
	//	},
	//	//"Test 2": {
	//	//	srvHost :  "http://127.0.0.1",
	//	//	id: idtest2,
	//	//	expected: "http://127.0.0.1/9ba87644294d4f3eb5df51c9d8c330b3",
	//	//},
	//	//"Test 3": {
	//	//	srvHost :  "http://172.10.0.1",
	//	//	id: idtest3,
	//	//	expected: "http://172.10.0.1/88d3ea27336e41e7b0da877cc8a041d3",
	//	//},
	//
	//}
	//
	//for name, cs := range cases {
	//	t.Run(name, func(t *testing.T) {
	//		test := shortenerBL.NewShotenerBL()
	//		result, _ := test.CreateShort(cs.ctx, *cs.shortener)
	//		if result != cs.expected {
	//			t.Errorf("GenerateShortLink returned wrong result: got %v want %v", cs.expected, result)
	//		}
	//	})
	//}
}

