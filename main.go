package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"cryptopepe.io/cryptopepe-svg/builder"
	"math/big"
	"cryptopepe.io/cryptopepe-reader/pepe"
	"math/rand"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"cryptopepe.io/cryptopepe-svg/builder/look"
	"strings"
	"sort"
)

func main() {

	previewer := NewPreviewer()
	previewer.Start("localhost:8080")
}

type Previewer struct {
	router *mux.Router

	svgBuilder *builder.SVGBuilder

	indexTmpl *template.Template

	parts []string
}

func NewPreviewer() *Previewer {

	previewer := new(Previewer)

	previewer.svgBuilder = new(builder.SVGBuilder)
	previewer.svgBuilder.Load()

	previewer.indexTmpl = template.Must(template.ParseFiles("index.tmpl"))

	previewer.LoadPepeParts()

	previewer.router = mux.NewRouter()
	previewer.router.HandleFunc("/svg/{pepeId}", previewer.GetPepePartCombi).Methods("GET")
	previewer.router.HandleFunc("/look/{pepeId}", previewer.GetPepeLook).Methods("GET")
	previewer.router.HandleFunc("/", previewer.GetIndex).Methods("GET")

	return previewer
}

func (previewer *Previewer) Start(address string) {
	http.Handle("/", previewer.router)
	log.Println("Listening on " + address + " ...")
	log.Fatal(http.ListenAndServe(address, nil))
}

func parsePepeId(w http.ResponseWriter, r *http.Request) (*big.Int, bool) {
	vars := mux.Vars(r)
	idParam := vars["pepeId"]
	pepeId := new(big.Int)

	//base 0: base is determined by prefix (e.g. 0x.... or 0b....)
	pepeId, ok := pepeId.SetString(idParam, 0)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id!"))
		return nil, false
	}

	return pepeId, true
}

var sh256bit *big.Int
func init() {
	sh256bit = big.NewInt(1)
	sh256bit.Lsh(sh256bit, 256)
}
func getPepeDna(pepeId int64) *pepe.PepeDNA {
	rngSrc := rand.NewSource(pepeId)
	rng := rand.New(rngSrc)

	dna := pepe.PepeDNA{}
	dna[0] = new(big.Int).Rand(rng, sh256bit)
	dna[1] = new(big.Int).Rand(rng, sh256bit)
	return &dna
}

type PreviewPepe struct {
	PepeId int64
}

type PreviewData struct {
	Pepes []*PreviewPepe
}

func (previewer *Previewer) GetIndex(w http.ResponseWriter, r *http.Request) {

	prevData := new(PreviewData)
	pepeCount := len(previewer.parts) /// 64*64
	prevData.Pepes = make([]*PreviewPepe, pepeCount)
	for i := 0; i < pepeCount; i++ {
		prevData.Pepes[i] = &PreviewPepe{PepeId: int64(i)}
	}

	previewer.indexTmpl.Execute(w, prevData)
}

func (previewer *Previewer) GetPepePartCombi(w http.ResponseWriter, r *http.Request) {

	pepeId, ok := parsePepeId(w, r)
	if !ok {
		return
	}
	var pepeLook = look.PepeLook{
		Skin: look.Skin{
			Color: "#389945",
		},
		Head: look.Head{
			Hair: look.Hair{
				HairColor: "#ab7e2b",
				HatColor:  "#c97225",
				HatColor2: "#cf1d32",
				HairType:  "none",
			},
			Eyes: look.Eyes{
				EyeColor: "#477b64",
				EyeType:  "eyes>colored_eyes",
			},
			Mouth: "mouth>basic_lips",
		},
		Body: look.Body{
			Neck: "none",
			Shirt: look.Shirt{
				ShirtColor: "#1ca479",
				ShirtType:  "shirt>basic_shirt",
			},
		},
		Extra: look.Extra{
			Glasses: look.Glasses{
				PrimaryColor:   "#0d0606",
				SecondaryColor: "#00c3c2",
				GlassesType:    "none",
			},
		},
		BackgroundColor: "#dbdefb",
	}

	index := pepeId.Int64()
	index = index % int64(len(previewer.parts))
	part := previewer.parts[index]
	if strings.HasPrefix(part, "shirt") {
		pepeLook.Body.Shirt.ShirtType = part
	}
	if strings.HasPrefix(part, "hair") {
		pepeLook.Head.Hair.HairType = part
	}
	if strings.HasPrefix(part, "glasses") {
		pepeLook.Extra.Glasses.GlassesType = part
	}
	if strings.HasPrefix(part, "mouth") {
		pepeLook.Head.Mouth = part
	}
	if strings.HasPrefix(part, "neck") {
		pepeLook.Body.Neck = part
	}
	if strings.HasPrefix(part, "eyes") {
		pepeLook.Head.Eyes.EyeType = part
	}
	pepe.ResolveLookConflicts(&pepeLook)

	h := w.Header()
	h.Set("Content-Type", "image/svg+xml")
	h.Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	err := previewer.svgBuilder.ConvertToSVG(w, &pepeLook)
	if err != nil {
		log.Fatalf("Error! Failed to convert pepe %s data to SVG! %v\n", pepeId.Text(10), err)
	}
}

func (previewer *Previewer) GetPepeSVG(w http.ResponseWriter, r *http.Request) {

	pepeId, ok := parsePepeId(w, r)
	if !ok {
		return
	}
	dna := getPepeDna(pepeId.Int64())
	pepeLook := dna.ParsePepeDNA()

	h := w.Header()
	h.Set("Content-Type", "image/svg+xml")
	h.Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	err := previewer.svgBuilder.ConvertToSVG(w, pepeLook)
	if err != nil {
		log.Fatalf("Error! Failed to convert pepe %s data to SVG! %v\n", pepeId.Text(10), err)
	}
}

func (previewer *Previewer) GetPepeLook(w http.ResponseWriter, r *http.Request) {
	pepeId, ok := parsePepeId(w, r)
	if !ok {
		return
	}
	dna := getPepeDna(pepeId.Int64())
	pepeLook := dna.ParsePepeDNA()

	h := w.Header()
	h.Set("Content-Type", "application/json")
	h.Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(pepeLook)
	if err != nil {
		log.Fatalf("Error! Failed to convert pepe %s look to JSON! %v\n", pepeId.Text(10), err)
	}
}

func (previewer *Previewer) LoadPepeParts() {
	dat, err := ioutil.ReadFile("vendor/cryptopepe.io/cryptopepe-svg/builder/tmpl/mappings.json")
	if err != nil {
		panic(err)
	}

	var parsed map[string]map[string]string

	if err := json.Unmarshal(dat, &parsed); err != nil {
		panic(err)
	}

	idList := parsed["id2name"]

	previewer.parts = make([]string, 0, len(idList))
	for k := range idList {
		previewer.parts = append(previewer.parts, k)
	}
	sort.Strings(previewer.parts)
}
