package utils_test

import (
	"MediaTools/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChineseToInt(t *testing.T) {
	tests := []struct {
		input     string
		expectInt int
		exceptErr error
	}{
		{"零", 0, nil},
		{"一", 1, nil},
		{"十", 10, nil},
		{"十二", 12, nil},
		{"二十", 20, nil},
		{"一百", 100, nil},
		{"一百零一", 101, nil},
		{"二百五", 250, nil}, // 二百五十的简写
		{"一千零一", 1001, nil},
		{"一千二百三十四", 1234, nil},
		{"一万", 10000, utils.ErrNumberTooLarge},
		{"一万零一", 10001, utils.ErrNumberTooLarge},
		{"十万", 100000, utils.ErrNumberTooLarge},
		{"一百二十三万四千五百六十七", 1234567, utils.ErrNumberTooLarge},
		{"一亿", 100000000, utils.ErrNumberTooLarge},
		{"三亿二千万", 320000000, utils.ErrNumberTooLarge},
		{"五亿三", 530000000, utils.ErrNumberTooLarge}, // 五亿三千万的简写
		{"二十一", 21, nil},
		{"一百一十", 110, nil},
		{"一百一十一", 111, nil},
		{"一千一百一十一", 1111, nil},
		{"一万一千一百一十一", 11111, utils.ErrNumberTooLarge},
		{"一百万", 1000000, utils.ErrNumberTooLarge},
		{"一千零二十", 1020, nil},
		{"一千二百零三", 1203, nil},
		{"一千零二", 1002, nil},
		{"一万零三百", 10300, utils.ErrNumberTooLarge},
		{"一万零三", 10003, utils.ErrNumberTooLarge},
		{"一亿零一", 100000001, utils.ErrNumberTooLarge},
		{"一亿零一万零一", 100010001, utils.ErrNumberTooLarge},
		{"十亿", 1000000000, utils.ErrNumberTooLarge},
		{"一千二百三十四万五千六百七十八", 12345678, utils.ErrNumberTooLarge},
		{"一千二百三十四亿五千六百七十八万九千零一十二", 123456789012, utils.ErrNumberTooLarge},
		{"一百零一", 101, nil},
		{"一百零十", 110, nil},
		{"一百零", 100, nil},
		{"十万零三百", 100300, utils.ErrNumberTooLarge},
		{"十万零三", 100003, utils.ErrNumberTooLarge},
		{"一千二百三", 1203, nil},
		{"一千二", 1200, nil},
		{"一万三", 13000, utils.ErrNumberTooLarge},
		{"一万零三百五", 10305, utils.ErrNumberTooLarge},
		{"一千零三十", 1030, nil},
		{"一千零三", 1003, nil},
		{"一千三", 1300, nil},
		{"一百三", 130, nil},
		{"一千二百零三", 1203, nil},
		{"一千二百三", 1203, nil},
		{"一千零二十", 1020, nil},
		{"一千零二", 1002, nil},
		{"一万零三百", 10300, utils.ErrNumberTooLarge},
		{"一万零三", 10003, utils.ErrNumberTooLarge},
		{"一亿零一", 100000001, utils.ErrNumberTooLarge},
		{"一亿零一万零一", 100010001, utils.ErrNumberTooLarge},
		{"十亿", 1000000000, utils.ErrNumberTooLarge},
		// 错误用例
		{"", 0, utils.ErrInvalidCharacter},
		{"abc", 0, utils.ErrInvalidCharacter},
		{"一百二十x", 0, utils.ErrInvalidCharacter},
	}

	for _, test := range tests {
		result, err := utils.ChineseToInt(test.input)
		require.Equal(t, test.exceptErr, err, "输入: %s", test.input)
		if test.exceptErr == nil {
			require.Equal(t, test.expectInt, result, "输入: %s", test.input)
		}
	}
}
