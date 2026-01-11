package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gocolly/colly"
)

type Teacher struct {
	Index string
	Name  string
	Type  string
	CTI   string
}

func GetTeachersInfo(url string) []Teacher {
	collector := colly.NewCollector(
		colly.AllowedDomains("www.unica.edu.pe"),
	)

	var teachers []Teacher

	collector.OnHTML(".table", func(e *colly.HTMLElement) {
		t := Teacher{}

		e.ForEach("tbody", func(i int, element *colly.HTMLElement) {
			var fields []string
			for _, v := range strings.Split(element.Text, "\n") {
				v = strings.TrimSpace(v)
				if v != "" {
					fields = append(fields, v)
				}
			}
			ctyURL := element.ChildAttr("a", "href")

			if len(fields) == 3 {
				t.Index = fields[0]
				t.Name = fields[1]
				t.Type = fields[2]
			}

			if ctyURL != "" {
				t.CTI = ctyURL
			}
			teachers = append(teachers, t)
		})

	})
	collector.Visit(fmt.Sprintf("https://www.unica.edu.pe/%s/code/docentes.php", url))

	return teachers
}

func main() {

	facus := []any{
		"administracion",
		"agronomia",
		"arquitectura",
		"ciencias",
		"ccbb",
		"cctya",
		"educacion",
		"economia",
		"contabilidad",
		"derecho",
		"farmacia",
		"fias",
		"civil",
		"fimm",
		"sistemas",
		"fime",
		"fipa",
		"quimica",
		"medicinahumana",
		"veterinaria",
		"obstetricia",
		"odontologia",
		"psicologia",
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Println("Seleccione una facultad: (Escribala tal cual)\n")
	i, j := 0, 5
	for ; j < len(facus); i, j = j, j+5 {
		t := facus[i:j]
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", t...)
	}
	f := ""

	t := facus[i:]

	f += strings.Repeat("%s\t", len(t))
	f += "\n"

	fmt.Fprintf(w, f, t...)

	w.Flush()

	var facultad string
	fmt.Print("\n>>> ")
	fmt.Scanln(&facultad)
	facultad = strings.TrimSpace(facultad)

	in := false
	for _, v := range facus {
		s, _ := v.(string)

		if s == facultad {
			in = true
			break
		}
	}

	if !in {
		log.Fatal("Entrada invalida")
	}

	teachers := GetTeachersInfo(facultad)

	fmt.Fprintln(w, "INDICE\tNOMBRE COMPLETO\tCATEGORIA DOCENTE\tCTI VITAE")
	format := "%s\t%s\t%s\t%s\n"
	for _, teacher := range teachers {
		fmt.Fprintf(w, format, teacher.Index, teacher.Name, teacher.Type, teacher.CTI)
	}

	w.Flush()

}
