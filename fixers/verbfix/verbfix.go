package verbfix

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

var fmtExpsRE = regexp.MustCompile(`(fmt\.Sprintf\(\s*` + "`[^`]*`)" + `([^)]*)\)`)

// FixIt takes the whole file to be fixed and chunks it for fixing and assembly.
func FixIt(shebang string) string {
	//flag.Parse()

	fixed := ""

	remainder := shebang

	for {
		if loc := fmtExpsRE.FindStringIndex(remainder); loc != nil {
			fixed += remainder[:loc[0]]
			fixed += stitchTogether(remainder[loc[0]:loc[1]])
			remainder = remainder[loc[1]:]
		} else {
			fixed += remainder
			break
		}
	}

	return fixed
}

func stitchTogether(match string) string {
	r := fixReplacements(match)
	for i, v := range r.oldVerbs {
		r.fmtMatch = strings.Replace(r.fmtMatch, v, r.newVerbs[i], 1)
	}
	return r.fmtMatch + ", " + strings.Join(r.newVars, ", ") + ")"
}

type replacers struct {
	fmtMatch string
	oldVerbs []string
	newVerbs []string
	newVars  []string
}

func fixReplacements(match string) replacers {
	//submatch = capturing group
	varSplitRE := regexp.MustCompile(`\s*,\s*`)
	verbsRE := regexp.MustCompile(`(%(\[[0-9]+\])?[a-zA-Z])`)
	fmtExps := fmtExpsRE.FindAllStringSubmatch(match, -1)

	var r replacers

	//fmtExps[0][1] = fmt.Sprintf(`.*`
	//fmtExps[0][2] = vars

	r.fmtMatch = fmtExps[0][1]

	vars := varSplitRE.Split(fmtExps[0][2], -1)[1:]
	rawVerbs := verbsRE.FindAllStringSubmatch(r.fmtMatch, -1)
	for _, v := range rawVerbs {
		r.oldVerbs = append(r.oldVerbs, v[1])
	}

	glog.Infof("vars: %q\n", vars)
	glog.Infof("r.oldVerbs: %q\n", r.oldVerbs)

	varIndex := 0
	verbIndex := 0
	var keys []string
	varCount := make(map[string]int)
	useIndices := false
	for i := 0; i < max(len(vars), len(r.oldVerbs)); i++ {
		//current
		baseVerb := r.oldVerbs[verbIndex][len(r.oldVerbs[verbIndex])-1:]
		varr := varForVerb(vars, r.oldVerbs, &varIndex, verbIndex)
		key := indexKey(baseVerb, varr)
		keys = append(keys, key)
		verbIndex++

		glog.Infof("verbIndex: %d\tvarIndex %d\tkey: %q\t", verbIndex, varIndex, key)
		if c, ok := varCount[varr]; ok {
			useIndices = true
			varCount[varr] = c + 1
		} else {
			varCount[varr] = 1
			r.newVars = append(r.newVars, varr)
			glog.Info(" not")
		}
		glog.Info(" contained\n")

	}
	glog.Infof("varCount: %q\n", varCount)
	glog.Infof("keys: %q\n", keys)
	glog.Infof("r.newVars: %q\n", r.newVars)

	for _, k := range keys {
		keyParts := strings.Split(k, "-")
		if useIndices {
			r.newVerbs = append(r.newVerbs, fmt.Sprintf("%%[%d]%s", indexOf(r.newVars, keyParts[1])+1, keyParts[0]))
		} else {
			r.newVerbs = append(r.newVerbs, fmt.Sprintf("%%%s", keyParts[0]))
		}
	}
	glog.Infof("r.newVerbs: %q\n", r.newVerbs)

	return r
}

func varForVerb(vars, verbs []string, varIndex *int, verbIndex int) string {
	rawVerb := verbs[verbIndex]
	realVarIndex := *varIndex
	if rawVerb[1:2] == "[" {
		if i, err := strconv.Atoi(rawVerb[2 : len(rawVerb)-2]); err == nil {
			realVarIndex = i - 1

			// still increment on ineffectual index
			if *varIndex == realVarIndex {
				*varIndex++
			}
		} else {
			log.Fatalf("bad verb index: %s", rawVerb)
		}
	} else {
		*varIndex++
	}
	return vars[realVarIndex]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func indexOf(haystack []string, needle string) int {
	for i, hay := range haystack {
		if hay == needle {
			return i
		}
	}
	return -1
}

func indexKey(verb, varName string) string {
	return fmt.Sprintf("%s-%s", verb, varName)
}

func FileContent(filename string) (string, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil || info.IsDir() {
		return "", fmt.Errorf("file does not exist or is a directory")
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func OverwriteFile(filename, content string) error {
	if err := os.Remove(filename); err != nil {
		return err
	}

	return WriteFile(filename, content)
}

func WriteFile(filename, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
