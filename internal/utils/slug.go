package utils

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GenerateSlug(title string) string {
	// 1. Chuyển tiếng Việt có dấu thành không dấu (Loại bỏ các ký tự dấu - marks)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	titleWithoutAccents, _, _ := transform.String(t, title)

	// Đặc thù tiếng Việt: Chữ Đ/đ không thuộc loại dấu (mark) ở trên nên phải replace thủ công
	titleWithoutAccents = strings.ReplaceAll(titleWithoutAccents, "đ", "d")
	titleWithoutAccents = strings.ReplaceAll(titleWithoutAccents, "Đ", "D")

	// 2. Chuyển tất cả thành chữ thường
	lower := strings.ToLower(titleWithoutAccents)

	// 3. Thay thế tất cả những gì không phải chữ/số thành dấu gạch ngang
	re := regexp.MustCompile("[^a-z0-9]+")
	slug := strings.Trim(re.ReplaceAllString(lower, "-"), "-")

	// 4. Thêm chuỗi số ngẫu nhiên ở cuối để đảm bảo tính duy nhất
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	slug = slug + "-" + strconv.Itoa(rng.Intn(10000))

	return slug
}
