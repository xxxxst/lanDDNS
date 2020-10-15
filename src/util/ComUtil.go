package util

import (
	// "io"
	"regexp"
	"strings"
	"path"
	"os"
	"io/ioutil"
)

// 比较两个[]byte是否相同
func Equal(a []byte, b []byte) bool {
	// If one is nil, the other must also be nil.
    if (a == nil) != (b == nil) { 
        return false; 
    }

    if len(a) != len(b) {
        return false;
    }

    for i := range a {
        if a[i] != b[i] {
            return false;
        }
    }

    return true;
}

// func removeOneBom(bytes []byte, bom []byte) ([]byte, bool){
// 	isOk := Equal(bytes[:3], []byte{0xef, 0xbb, 0xbf});
// 	if(isOk){
// 		bytes = bytes[3:];
// 	}

// 	return bytes, isOk;
// }

// 去除文件bom
func removeBom(data []byte) []byte {
	arrBom := [][]byte {
		{0xef, 0xbb, 0xbf},			//utf-8
		{0xfe, 0xff},				//utf-16大端
		{0xff, 0xfe},				//utf-16小端
		{0x00, 0x00, 0xfe, 0xff},	//utf-32大端
		{0xff, 0xfe, 0x00, 0x00},	//utf-32小端
	}
	rst := data;
	isOk := false;
	for i:=0; i < len(arrBom); i++ {
		if len(data) < len(arrBom[i]) {
			continue;
		}
		
		isOk = Equal(data[:len(arrBom[i])], arrBom[i]);
		if(isOk){
			rst = data[len(arrBom[i]):];
			break;
		}
	}

	return rst;
}

// 读文件-去除bom
func ReadFile(path string) []byte {
	bytes,_ := ioutil.ReadFile(path);
	bytes = removeBom(bytes);
	return bytes;
}

// 读文件-去除bom
func ReadFileString(path string) string {
	bytes,_ := ioutil.ReadFile(path);
	bytes = removeBom(bytes);
	return string(bytes);
}

// 写文件
func SaveFileString(path string, text string) bool {
	fs, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, os.ModePerm);
	if(err != nil) {
		return false;
	}
	_, err = fs.WriteString(text);
	fs.Close();
	
	if(err != nil) {
		return false;
	}
	return true;
}

// 复制文件
func CopyFile(src string, dst string) {
	if(src == dst) {
		return;
	}

	fsSrc, err := os.Open(src);
	if(err != nil) {
		return;
	}
	defer fsSrc.Close();

	fsDst, err1 := os.Create(dst);
	if(err1 != nil){
		return;
	}
	defer fsDst.Close();

	bufSize := 1024 * 1024;

	buf := make([]byte, bufSize);
	for {
		n, err2 := fsSrc.Read(buf);
		if(n <= 0) {
			break;
		}
		fsDst.Write(buf[:n]);
		if err2 != nil {
			break;
		}
	}
}

func Min(x, y int) int {
    if x <= y {
        return x
    }
    return y
}

func Max(x, y int) int {
    if x >= y {
        return x
    }
    return y
}

func SplitStr(text string, what string) []string {
	runeText := []rune(text);
	arrRst := []string{};

	arr := SplitRune(runeText, what);
	for _,val := range arr {
		arrRst = append(arrRst, string(val));
	}

	return arrRst;
}

func SplitRune(text []rune, what string) [][]rune {
	arrRst := [][]rune{};
	whatRunes := []rune(what);
	lenWaht := len(whatRunes);

	if len(text) == 0 {
		arrRst = append(arrRst, []rune{});
		return arrRst;
	}

	if lenWaht == 0 {
		arrRst = append(arrRst, text);
		return arrRst;
	}

	idx := 0;
    for i:=0; i<len(text); i++ {
        found := true
        for j := range whatRunes {
			if i+j >= len(text) {
				found = false;
				break;
			}
            if text[i+j] != whatRunes[j] {
                found = false
                break
            }
        }
        if found {
			// fmt.Println("aaa:", idx, i, text[idx:i]);
			arrRst = append(arrRst, text[idx:i]);
			i += lenWaht;
			idx = i;
        }
	}
	if idx < len(text) {
		arrRst = append(arrRst, text[idx:]);
	} else if idx == len(text) && idx != 0 {
		arrRst = append(arrRst, []rune{});
	}
    return arrRst;
}

func SearchStr(text string, what string) int {
	runeText := []rune(text);

	return SearchRune(runeText, what);
}

func SearchRune(text []rune, what string) int {
    whatRunes := []rune(what)

    for i := range text {
        found := true
        for j := range whatRunes {
			if i+j >= len(text) {
				return -1;
			}
            if text[i+j] != whatRunes[j] {
                found = false
                break
            }
        }
        if found {
            return i
        }
    }
    return -1
}

func regReplaceOne(reg *regexp.Regexp, str string, strReplace string) string {
	found := reg.FindString(str);
    if found != "" {
        return strings.Replace(str, found, strReplace, 1);
	}
	return str;
}

// func FormatPath(path string) string {
var FormatPath = (func() (func(path string) string) {
	reg1 := regexp.MustCompile("[\\/\\\\]+");
	reg2 := regexp.MustCompile("(/|^)(?:\\s*\\.\\s*/)+");
	reg4 := regexp.MustCompile("(([^/]+/)|^)\\s*\\.\\.\\s*/");
	reg3 := regexp.MustCompile("^/");

	return func(path string) string {
		//aa/\\/bb\\c  =>  a/b/c
		path = strings.TrimSpace(reg1.ReplaceAllString(path, "/"));

		// ./a././b./c  =>  a/b/c
		isRelative := !((len(path)>=1&&path[0] == '/') || (len(path)>1&&path[1]==':'));
		path = reg2.ReplaceAllString(path, "/");
		if(isRelative) {
			path = reg3.ReplaceAllString(path, "");
		}

		// a/b/../c  =>  a/c
		length := 0;
		for length != len(path) {
			length = len(path);
			path = regReplaceOne(reg4, path, "");
		}

		return path;
	}

})();

func FileExists(path string) bool {
	st, err := os.Stat(path);
	if err != nil {
		return false
	}

	return !st.IsDir();
}

func DirectoryExists(path string) bool {
	st, err := os.Stat(path);
	if err != nil {
		return false
	}

	return st.IsDir();
}

func IsDocType(strPath string) bool {
	str := strings.ToLower(path.Ext(strPath));
	return (str == ".md" || str == ".html");
}

func IsTextType(strPath string) bool {
	str := strings.ToLower(path.Ext(strPath));
	return (str == ".md" || str == ".html" || str == ".txt");
}

func NeedEncodeType(strPath string) bool {
	str := strings.ToLower(path.Ext(strPath));
	return (str == ".md" || str == ".html" || str == ".txt");
}

func IsMdType(strPath string) bool {
	str := strings.ToLower(path.Ext(strPath));
	return (str == ".md");
}

func GetFileName(strPath string) string {
	reg := regexp.MustCompile("^.*[\\\\/]");
	return reg.ReplaceAllString(strPath, "");
}

func GetFileNameWithoutExtension(strPath string) string {
	reg := regexp.MustCompile("\\.[^\\.]*$");
	str := GetFileName(strPath);
	return reg.ReplaceAllString(str, "");
}