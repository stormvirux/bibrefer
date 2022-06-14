package getdoi

import (
	"fmt"
	"github.com/stormvirux/pdf"
	"math"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

func arxivRegexBuilder() (newRegex *regexp.Regexp, oldRegex *regexp.Regexp) {
	const domain = `astro-ph.GA|astro-p.CO|astro-ph.EP|astro-ph.HE|astro-ph.IM|astro-ph.SR|cond-math.dis-nn|cond-math.mtrl-sci|cond-math.mes-hall|cond-math.other|cond-math.quant-gas|cond-math.soft|cond-math.stat-mech|cond-math.str-el|cond-math.supr-con|gr-qc|hep-ex|hep-lat|hep-ph|hep-th|math-ph|nlin.AO|nlin.CG|nlin.CD|nlin.SI|nlin.PS|nucl-ex|nucl-th|physics.acc-ph|physics.ao-ph|physics.atom-ph|physics.atm-clus|physics.bio-ph|physics.chem-ph|physics.class-ph|physics.comp-ph|physics.data-an|physics.flu-dyn|physics.gen-ph|physics.geo-ph|physics.hist-ph|physics.ins-det|physics.med-ph|physics.optics|physics.ed-ph|physics.soc-ph|physics.plasm-phphysics.pop-ph|physics.space-ph|physics.quant-ph|math.AG|math.AT|math.AP|math.CT|math.CA|math.CO|math.AC|math.CV|math.DG|math.DS|math.FA|math.GM|math.GN|math.GT|math.GR|math.HO|math.IT|math.KT|math.LO|math.MP|math.MG|math.NT|math.NA|math.OA|math.OC|math.PR|math.QA|math.RT|math.RA|math.SP|math.ST|math.SG|cs.AI|cs.CL|cs.CC|cs.CE|cs.CG|cs.GT|cs.CV|cs.CY|cs.CR|cs.DS|cs.DB|cs.DL|cs.DM|cs.DC|cs.ET|cs.FL|cs.GL|cs.GR|cs.AR|cs.HC|cs.IR|cs.IT|cs.LG|cs.LO|cs.MS|cs.MA|cs.MM|cs.|cs.NE|cs.NA|cs.OS|cs.OH|cs.PF|cs.PL|cs.RO|cs.SI|cs.SE|cs.SD|cs.SC|cs.SY|q-bio.BM|q-bio.CB|q-bio.GN|q-bio.MN|q-bio.NC|q-bio.OT|q-bio.PE|q-bio.QM|q-bio.SC|q-bio.TO|q-fin.CP|q-fin.EC|q-fin.GN|q-fin.MF|q-fin.PM|q-fin.PR|q-fin.RM|q-fin.ST|q-fin.TR|stat.AP|stat.CO|stat.ML|stat.ME|stat.OTstat.TH`

	const arxivIDFrom2007 = `\d{4}\.\d{4,5}(v\d+)?`
	const arxivIDBefore2007 = `((` + domain + `)r)/\\d+`

	const regexStringFrom2007 = `(?i)(arxiv:)?(` + arxivIDFrom2007 + `)`
	const regexStringBefore2007 = `(?i)(arxiv:)?(` + arxivIDBefore2007 + `)`

	newRegex = regexp.MustCompile(regexStringFrom2007)
	oldRegex = regexp.MustCompile(regexStringBefore2007)

	return newRegex, oldRegex
}

// Only use to avoid cyclomatic complexity
func verbosePrint(isVerbose bool, message string) {
	if isVerbose {
		fmt.Println(message)
	}
}

func findWords(chars []pdf.Text) (words []pdf.Text) {
	// Sort by Y coordinate and normalize.
	const nudge = 1
	sort.Sort(pdf.TextVertical(chars))
	old := -100000.0
	for i, c := range chars {
		if c.Y != old && math.Abs(old-c.Y) < nudge {
			chars[i].Y = old
		} else {
			old = c.Y
		}
	}

	// Sort by Y coordinate, breaking ties with X.
	// This will bring letters in a single word together.
	sort.Sort(pdf.TextVertical(chars))

	// Loop over chars.
	for i := 0; i < len(chars); {
		// Find all chars in line.
		j := i + 1
		for j < len(chars) && chars[j].Y == chars[i].Y {
			j++
		}
		var end float64
		// Split line into words (really, phrases).
		for k := i; k < j; {
			ck := &chars[k]
			s := ck.S
			end = ck.X + ck.W
			charSpace := ck.FontSize / 6
			wordSpace := ck.FontSize * 2 / 3
			l := k + 1
			for l < j {
				// Grow word.
				cl := &chars[l]
				if sameFont(cl.Font, ck.Font) && math.Abs(cl.FontSize-ck.FontSize) < 0.1 && cl.X <= end+charSpace {
					s += cl.S
					end = cl.X + cl.W
					l++
					continue
				}
				// Add space to phrase before next word.
				if sameFont(cl.Font, ck.Font) && math.Abs(cl.FontSize-ck.FontSize) < 0.1 && cl.X <= end+wordSpace {
					s += " " + cl.S
					end = cl.X + cl.W
					l++
					continue
				}
				break
			}
			f := ck.Font
			// f = strings.TrimSuffix(f, ",Italic")
			// f = strings.TrimSuffix(f, "-Italic")
			words = append(words, pdf.Text{Font: f, FontSize: ck.FontSize, X: ck.X, Y: ck.Y, W: end - ck.X, S: s})
			k = l
		}
		i = j
	}

	return words
}

func sameFont(f1, _ string) bool {
	f1 = strings.TrimSuffix(f1, ",Italic")
	f1 = strings.TrimSuffix(f1, "-Italic")
	_ = strings.TrimSuffix(f1, ",Italic")
	f2 := strings.TrimSuffix(f1, "-Italic")
	return strings.TrimSuffix(f1, ",Italic") == strings.TrimSuffix(f2, ",Italic") || f1 == "Symbol" || f2 == "Symbol" || f1 == "TimesNewRoman" || f2 == "TimesNewRoman"
}

func stripVersion(arXivID string) string {
	r := regexp.MustCompile(`v\d+\z`)
	return r.ReplaceAllString(arXivID, "")
}

func detectGS() (bool, string) {
	cmd := exec.Command("which", "gs")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "gswin*c.exe")
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		// TODO: Log the above
		fmt.Errorf("%w", err)
		return false, ""
	}
	return string(output) != "", string(output)
}
